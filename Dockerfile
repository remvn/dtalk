# 1: build stage
FROM golang:1.23.2-bookworm AS build

WORKDIR /app 

COPY go.mod go.sum ./ 
RUN go mod download 

COPY ./cmd ./cmd
COPY ./internal ./internal 

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main/

# 2: final stage 
FROM scratch

WORKDIR /app 

COPY --from=build /app/main .

ENTRYPOINT ["./main"]
