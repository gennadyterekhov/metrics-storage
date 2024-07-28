# metrics-storage

## linting

revive -config revive_config.toml -formatter friendly ./... > revive_result.txt


## db installation (one-time use)

      sudo -i -u postgres
      psql -U postgres
      postgres=# create database metrics_db;
      postgres=# create database metrics_db_test;
      postgres=# create user metrics_user with encrypted password 'metrics_pass';
      postgres=# grant all privileges on database metrics_db to metrics_user;
      postgres=# grant all privileges on database metrics_db_test to metrics_user;
      alter database metrics_db owner to metrics_user;
      alter database metrics_db_test owner to metrics_user;
      alter schema public owner to metrics_user;

after that, use this to connect to db in cli

      psql -U metrics_user -d metrics_db

or

      psql -U metrics_user -d metrics_db_test



## test locally

cd internal
go test ./...

### clear test cache

go clean -testcache

### specific package

go test github.com/gennadyterekhov/metrics-storage/internal/agent/client

### specific test

go test github.com/gennadyterekhov/metrics-storage/internal/agent/client -run TestCanSendCounterValue -test.v
go test -run=TestCanSendCounterValue ./...

### mocks

mockgen -destination=mocks/mock_store.go -package=mocks project/store Store

## test coverage

      go clean -testcache
      
      # create coverage file
      go test -v -coverpkg=./... -coverprofile coverage/coverage.out ./...

      # see results in html page
      go tool cover -html=coverage/coverage.out
      
      # see results in CLI
      go tool cover -func coverage/coverage.out


## benchmarks
создать журнал профилирования с помощью бенчмарка
go test -bench=. -cpuprofile=cpu.out  
или  
go test -bench=. -memprofile=mem.out  

анализировать журнал профилирования   
go tool pprof -http=":9090" <test exe> cpu.out  
например  
go tool pprof -http=":9090" BenchmarkSaveMetricService_cpu.test BenchmarkSaveMetricService_cpu.out  

## documentation
### godoc
start godoc server

      godoc -play -http=:8080 -goroot="/Users/gena/code/yandex/practicum/golang_advanced/metrics-storage"  

after that, visit  
- [including 3rd party](http://localhost:8080/pkg/?m=all)  
- [local project doc](http://localhost:8080/pkg/github.com/gennadyterekhov/metrics-storage/?m=all)  

or download html doc using:

      wget -r -np -N -E -p -k http://localhost:8080/pkg/github.com/gennadyterekhov/metrics-storage


### swagger
[comment format doc](https://github.com/swaggo/swag#declarative-comments-format)

update doc files:

      swag init --generalInfo controllers.go --dir ./internal/server/httpui/handlers/handlers --output ./swagger/

visit [swagger editor](https://editor.swagger.io/) and paste `swagger/swagger.yaml`



## CI tests locally

cd cmd/agent
go build -o agent *.go
cd ../server
go build -o server *.go

build:
go build -o cmd/agent/agent cmd/agent/main.go && go build -o cmd/server/server cmd/server/main.go
go build -o cmd/agent/agent cmd/agent/*.go && go build -o cmd/server/server cmd/server/*.go

./metricstest -test.v -test.run=^TestIteration1$ -agent-binary-path=cmd/agent/agent

### all

./metricstest -test.v -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server

### 1

./metricstest -test.v -test.run=^TestIteration1$ -binary-path=cmd/server/server

### 2

./metricstest -test.v -test.run=^TestIteration2[AB]*$ -source-path=. -agent-binary-path=cmd/agent/agent
./metricstest -test.v -source-path=. -agent-binary-path=cmd/agent/agent
./metricstest -test.v -test.run=^TestIteration*$ -source-path=. -agent-binary-path=cmd/agent/agent
-binary-path=cmd/server/server

### 3

./metricstest -test.v -test.run=^TestIteration3$ -source-path=. -agent-binary-path=cmd/agent/agent
-binary-path=cmd/server/server

./metricstest -test.v \
-source-path=. \
-agent-binary-path=cmd/agent/agent \
-binary-path=cmd/server/server

### 7

какая-то фигня тут
использует параметры не так как предполагается
TEMP_FILE=tmp
./metricstest -test.v -test.run=^TestIteration7$ \
-agent-binary-path=cmd/agent/agent \
-binary-path=cmd/server/server \
-server-port=8080 \
-source-path=. \
-file-storage-path=tmp \
-database-dsn=t

### iter8

          SERVER_PORT=8080
          ADDRESS="localhost:8080"
          TEMP_FILE=storage-for-ci-tests.json
          ./metricstest -test.v -test.run=^TestIteration8$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=8080 \
            -source-path=. \
            -file-storage-path=storage-for-ci-tests.json \
            -database-dsn=null

### iter9

          SERVER_PORT=8080
          ADDRESS="localhost:8080"
          TEMP_FILE=storage-for-ci-tests.json
          ./metricstest -test.v -test.run=^TestIteration9$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -file-storage-path=storage-for-ci-tests.json \
            -server-port=8080 \
            -source-path=. \
            -database-dsn=null \
            -key=a

# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без
   префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например,
в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.