.DEFAULT_GOAL := run 
.PHONY : hello fmt vet  build  run 
hello : 
	@echo "BUilding and running the file "
fmt : 
	go fmt ./...
vet : fmt 
	go vet ./...
build: vet
	mkdir bin
	go build -o ./bin 
run: build
	./bin/ratelimit
exec : 
	./bin/ratelimit