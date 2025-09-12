run:
	go run main.go

build:
	env GOOS=linux GOARCH=amd64 go build -o tickets *.go

deploy: build
	serverless deploy --stage prod

clean:
	rm -rf ./bin ./vendor Gopkg.lock ./serverless