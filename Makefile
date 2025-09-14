protoc:
	protoc --proto_path=api --go_out=api --go_opt=paths=source_relative --go-grpc_out=api --go-grpc_opt=paths=source_relative api/http.proto
	protoc --proto_path=test --go_out=test --go_opt=paths=source_relative --go-grpc_out=test --go-grpc_opt=paths=source_relative test/test.proto

tidy:
	go mod tidy
