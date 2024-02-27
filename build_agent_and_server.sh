
go build -o cmd/server/server cmd/server/*.go
go build -o cmd/agent/agent cmd/agent/*.go
#./cmd/server/server 2>&1 & ./cmd/agent/agent 2>&1
#./cmd/server/server > logs/serv.log & ./cmd/agent/agent > logs/agent.log