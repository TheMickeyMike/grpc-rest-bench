module github.com/TheMickeyMike/grpc-rest-bench/benchmarks

go 1.17

replace github.com/TheMickeyMike/grpc-rest-bench/warehouse => ../warehouse

require (
	github.com/TheMickeyMike/grpc-rest-bench/warehouse v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
)

require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.20.0 // indirect
	golang.org/x/text v0.3.3 // indirect
)