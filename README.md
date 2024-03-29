# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

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

### starting schema (supposedly will not change)
CREATE TYPE metric_type AS ENUM ('gauge', 'counter');
create table metrics
(
name varchar(255) primary key ,
type metric_type not null default 'gauge',
value double precision default null,
delta numeric default null
);

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

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

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
go test -coverprofile cover.out ./...
go tool cover -html=cover.out

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
./metricstest -test.v -test.run=^TestIteration*$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server

### 3
./metricstest -test.v -test.run=^TestIteration3$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server

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
### 8