


build_client:
	go build -o UploadClient cmd/cli/main.go


build_server:
	go build -o grpcUploadServer cmd/server/main.go