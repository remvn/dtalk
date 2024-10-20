# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main ./cmd/main

run:
	@go run ./cmd/main

run-lk:
	@livekit-server --dev

docker-build:
	@docker build -t remvn/dtalk:latest . 
	@docker build -t remvn/dtalk-web:latest ./web 

docker-publish:
	@docker push remvn/dtalk:latest
	@docker push remvn/dtalk-web:latest

gen-deploy:
	@go run ./cmd/gen-deploy

compose-up:
	@docker compose --file ./deploy/docker-compose.yml up 

compose-down:
	@docker compose --file ./deploy/docker-compose.yml down 

test:
	@echo "Testing..."

clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean
