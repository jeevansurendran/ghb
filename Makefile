protoc:
	protoc --proto_path=api --go_out=api --go_opt=paths=source_relative --go-grpc_out=api --go-grpc_opt=paths=source_relative api/http.proto

tidy:
	go mod tidy
