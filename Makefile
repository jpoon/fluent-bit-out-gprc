all:
	go build -buildmode=c-shared -o out_grpc.so ./out-grpc

protoc:
	protoc -I api/ --go_out=plugins=grpc:api api/api.proto

clean:
	rm -rf *.so *.h *~
