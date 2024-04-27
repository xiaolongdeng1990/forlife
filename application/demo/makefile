
packagename=rpcx_demo


.PHONY: build
build: clean
	go build -o ${packagename} main.go

.PHONY: run
run:
	go run -race main.go

.PHONY: clean
clean:
	go clean
