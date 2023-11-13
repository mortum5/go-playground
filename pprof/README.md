# Profiling golang

Based on articles https://www.matoski.com/article/golang-profiling-flamegraphs/

## Guide for using

1. Include `_ "net/http/pprof"`
2. Check web server `hostname:port/debug/pprof`
3. Download profile for 30s in binary format
```sh
PPROF_TMPDIR=$(PWD) go tool pprof http://localhost:9090/debug/pprof/{profile_name}?seconds=30
```

4. Run interactively in web
```sh 
go tool pprof -http=localhost:5050 pprof.samples.cpu.001.pb.gz
```

5. Get trace for 30 seconds 
```sh
curl -o trace.out "http://localhost:6060/debug/pprof/trace?seconds=30"
```

6. Interactive watch of trace and goroutines analyzes `go tool trace trace.out`