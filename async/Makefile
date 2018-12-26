.PHONY: build buildPlugin clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/func1 func1/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/func2 func2/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
