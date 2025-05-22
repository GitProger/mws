.PHONY: all install clean build generate

all: clean install build

clean:
	rm -rf ./gen_*

build: generate
	go build -o server main.go
	cd client && go build -o client main.go

generate:
	go generate ./...
