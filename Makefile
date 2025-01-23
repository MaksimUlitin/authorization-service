.PHONY: up run

up:
	docker-compose up -d

run:
	go run cmd/main.go

all: up run
