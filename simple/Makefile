.PHONY: build buildPlugin clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/service1 service1/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/service2 service2/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/service3 service3/main.go


clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
