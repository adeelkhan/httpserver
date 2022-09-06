.PHONY: build clean

build:
	go build -o ./build/ httpserver.go 

run:
	./build/httpserver

clean:
	rm -rf ./build/