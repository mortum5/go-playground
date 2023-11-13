profile:
	go run main.go &
	sleep 1 
	PPROF_TMPDIR=$(PWD) PPROF_BINARY_PATH=$(PWD) go tool pprof ./pprof http://127.0.0.1:6060/debug/pprof/profile  &
	sleep 1
	ab -k -c 8 -n 100000 "http://127.0.0.1:6060/concat/?str=test&count=50" &
	ab -k -c 8 -n 100000 "http://127.0.0.1:6060/fib/?n=50&type=recursive" &
	ab -k -c 8 -n 100000 "http://127.0.0.1:6060/fib/?n=50&type=iterative" &

memory:
	go run main.go &
	sleep 1
	PPROF_TMPDIR=$(PWD) go tool pprof pprof http://127.0.0.1:6060/debug/pprof/heap &
	sleep 1
	ab -k -c 80 -n 1000000 "http://127.0.0.1:6060/concat/?str=test&count=5000"

bench:
	go test -bench=.
	
webP:
	go tool pprof pprof -http=localhost:7532 pprof.pprof.samples.cpu.001.pb.gz

webM:
	go tool pprof -http=localhost:7532 pprof.pprof.alloc_objects.alloc_space.inuse_objects.inuse_space.001.pb.gz

clean:
	rm pprof*