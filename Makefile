gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb
clean:
	rm pb/*.go
server:
	go run main.go