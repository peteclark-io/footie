build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/players/read players/read/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/players/write players/write/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/players/delete players/delete/main.go
