dev:
	go run cmd/main/main.go

build:
	cd cmd/main && go build -o ../../bin/main main.go

start:
	bin/main
