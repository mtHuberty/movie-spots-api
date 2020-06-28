run:
	go run cmd/movie-spots-api/main.go

run-dev:
	CompileDaemon -build="go build -o ./bin ./..." -command="./bin/movie-spots-api" -color=true -log-prefix=false -exclude-dir=.git

build:
	scripts/build.sh

test:
	GO_ENV=test go test ./...

vet:
	go vet ./...

test-only:
	if [[ -n '$(strip $(path))' ]]; then GO_ENV=test go test -tags test -p 1 ./$(path)/... -v -cover; else echo "\033[1;31mNeed to specify the 'path' argument (i.e. make test-only path=internal/middleware/movement)\033[m"; fi