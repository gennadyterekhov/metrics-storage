package agent

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/gennadyterekhov/metrics-storage/internal/agent/config"

	"github.com/gennadyterekhov/metrics-storage/internal/agent/client"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/poller"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/sender"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

// Agent instance with dependencies
type Agent struct {
	Config               *config.Config
	Poller               *poller.Poller
	Sender               *sender.MetricsSender
	MetricsStorageClient *client.MetricsStorageClient
}

func New() *Agent {
	conf := config.GetConfig()
	_, err := fmt.Printf("Agent started with server addr %v\n", conf.Addr)
	if err != nil {
		panic(err)
	}

	metricsStorageClient := client.New(conf)
	inst := &Agent{
		Config:               conf,
		Poller:               poller.New(conf.PollInterval),
		MetricsStorageClient: metricsStorageClient,
		Sender:               sender.New(metricsStorageClient, conf),
	}

	return inst
}

func (ag *Agent) Start() error {
	rootContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go gracefulShutdown(rootContext)

	ag.RunAgent(rootContext)

	return nil
}

func (ag *Agent) RunAgent(ctx context.Context) {
	// we only need to send the latest metrics,
	// it's no use if the poller puts 10 metrics into the channel, we need only the latest
	// there are 3 possible cases
	// (1) config.ReportInterval == config.PollInterval
	// in this case we can use a buffered channel with cap 1
	//
	// (2) config.ReportInterval > config.PollInterval
	// in this case we can also use a buffered channel with cap 1,
	// but we need to take-and-put new metrics when polling
	//
	// (3) config.ReportInterval < config.PollInterval
	// this means we report the same metrics multiple times
	// it means that we need to take-and-put the same metrics when reporting
	metricsChannel := make(chan metric.MetricsSet, 1)

	go ag.pollingRoutine(ctx, metricsChannel)
	go ag.reportingRoutine(ctx, metricsChannel)

	<-ctx.Done()
	close(metricsChannel)
}

func (ag *Agent) pollingRoutine(ctx context.Context, metricsChannel chan metric.MetricsSet) {
	logger.Custom.Infoln("polling started")

	for {
		select {
		case <-ctx.Done():
			logger.Custom.Infoln("poll context finished")
			return
		default:
			if !ag.Poller.IsRunning {
				if len(metricsChannel) == 0 {
					// if empty, just store in channel
					metricsChannel <- *ag.Poller.Poll()
				} else {
					// take latest and replace it with a new poll
					<-metricsChannel
					metricsChannel <- *ag.Poller.Poll()
				}
			}

			time.Sleep(time.Duration(ag.Config.PollInterval) * time.Second)
		}
	}
}

func (ag *Agent) reportingRoutine(ctx context.Context, metricsChannel chan metric.MetricsSet) {
	logger.Custom.Infoln("reporting started")
	var metricsSet metric.MetricsSet
	for {
		select {
		case <-ctx.Done():
			logger.Custom.Infoln("report context finished")
			return
		default:

			if !ag.Sender.IsRunning {
				if len(metricsChannel) == 0 {
					// nothing to report yet, need to wait for poller
					continue
				}
				metricsSet = <-metricsChannel
				metricsChannel <- metricsSet

				ag.Sender.Report(&metricsSet)
			}

			time.Sleep(time.Duration(ag.Config.ReportInterval) * time.Second)
		}
	}
}

// gracefulShutdown - this code runs if app gets any of (syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
func gracefulShutdown(ctx context.Context) {
	<-ctx.Done()
	logger.Custom.Infoln("graceful shutdown. waiting a little")
	time.Sleep(time.Second)
}
