test:
	curl -i 'http://localhost:3031/count'
	curl -i 'http://localhost:3032/count'
	curl -i 'http://localhost:3033/count'
	curl -XPOST -i 'http://localhost:3033/inc'
	sleep 10
	curl -i 'http://localhost:3031/count'
	curl -i 'http://localhost:3032/count'
	curl -i 'http://localhost:3033/count'