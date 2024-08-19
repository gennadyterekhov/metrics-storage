package handlers

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/gennadyterekhov/metrics-storage/internal/common/protobuf"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
)

var _ pb.MetricsServer = &Server{}

type Server struct {
	pb.UnimplementedMetricsServer
	GetMetricService  *services.GetMetricService
	SaveMetricService *services.SaveMetricService
	PingService       *services.PingService
}

func (s *Server) Ping(_ context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	if s.PingService.Repository == nil {
		return nil, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
	}

	return &pb.PingResponse{
		Message: "ok",
	}, nil
}

func (s *Server) GetAllMetrics(ctx context.Context, _ *pb.GetAllMetricsRequest) (*pb.GetAllMetricsResponse, error) {
	htmlPage := s.GetMetricService.GetMetricsListAsHTML(ctx)
	resp := &pb.GetAllMetricsResponse{
		Html: htmlPage,
	}
	return resp, nil
}

func (s *Server) GetMetric(ctx context.Context, request *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	response := &pb.GetMetricResponse{
		Type: request.Type,
	}

	if request.Type == pb.MetricType_COUNTER {
		counter, err := s.GetMetricService.Repository.GetCounter(ctx, request.Name)
		if err != nil {
			return nil, err
		}
		response.Counter = counter
	} else {
		gauge, err := s.GetMetricService.Repository.GetGauge(ctx, request.Name)
		if err != nil {
			return nil, err
		}
		response.Gauge = gauge
	}

	return response, nil
}

func (s *Server) SaveMetricList(ctx context.Context, request *pb.SaveMetricListRequest) (*pb.SaveMetricListResponse, error) {
	resp := &pb.SaveMetricListResponse{
		Message: "ok",
	}
	for _, v := range request.Request {
		_, err := s.SaveMetric(ctx, v)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func (s *Server) SaveMetric(ctx context.Context, request *pb.SaveMetricRequest) (*pb.SaveMetricResponse, error) {
	if request.Type == pb.MetricType_COUNTER {
		s.SaveMetricService.Repository.AddCounter(ctx, request.Name, request.Counter)
	} else {
		s.SaveMetricService.Repository.SetGauge(ctx, request.Name, request.Gauge)
	}

	if s.SaveMetricService.Config.StoreInterval == 0 && s.SaveMetricService.Config.FileStorage != "" {
		s.SaveMetricService.SaveToDisk(ctx)
	}

	return &pb.SaveMetricResponse{
		Message: "ok",
	}, nil
}
