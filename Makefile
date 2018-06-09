build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/players/read players/read/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/players/write players/write/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/players/create players/create/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/players/delete players/delete/main.go

	env GOOS=linux go build -ldflags="-s -w" -o bin/matches/read matches/read/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/matches/write matches/write/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/matches/delete matches/delete/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/matches/create matches/create/main.go

	env GOOS=linux go build -ldflags="-s -w" -o bin/groups/read groups/read/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/groups/write groups/write/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/groups/delete groups/delete/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/groups/create groups/create/main.go

	env GOOS=linux go build -ldflags="-s -w" -o bin/bookings/read bookings/read/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/bookings/write bookings/write/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/bookings/delete bookings/delete/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/bookings/create bookings/create/main.go

	sls deploy
