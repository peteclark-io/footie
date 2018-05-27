build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/players/read players/read/main.go
