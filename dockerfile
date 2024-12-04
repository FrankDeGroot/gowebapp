FROM golang AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN env CGO_ENABLED=0 env GOOS=linux GOARCH=amd64 go build -o myapp
FROM scratch
COPY --from=builder /app/myapp /myapp
ENTRYPOINT ["/myapp"]
