.PHONY:api

all: api

api:
	go build -o ./bin/api ./cmd/api.go

clean:
	rm -rf ./bin/*