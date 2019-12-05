dep:
	dep ensure

compile:
	env GOOS=linux go build -ldflags="-s -w" -o bin/footie main.go

build: dep compile

run: dep
	go build && ./footie http

deploy: build
	sls deploy
