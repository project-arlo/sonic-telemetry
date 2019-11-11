To make go packages form proto:

protoc --proto_path=. --go_out=plugins=grpc:. *.proto
