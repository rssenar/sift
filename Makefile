test:
	go test -test.v . -cover

test-race-cond:
	go test -test.v -tags race_cond -race .
