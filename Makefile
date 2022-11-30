t:  # test
	go test -v -race -covermode atomic -coverprofile coverage.out && go tool cover -html coverage.out -o coverage.html

f:  # format
	go fmt .

b:  # benchmark
	go test -bench . -benchmem -cpu 1

r:  # report benchmark
	go test -cpuprofile cpu.prof -memprofile mem.prof -bench . -cpu 1

c:  # cpu profile
	go tool pprof cpu.prof

m:  # memory profile
	go tool pprof mem.prof
