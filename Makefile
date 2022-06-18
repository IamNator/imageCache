
client:
	go build -o UploadClient cmd/cli/main.go

server:
	go build -o grpcUploadServer cmd/server/main.go
