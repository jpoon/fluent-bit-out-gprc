# fluent-bit-out-grpc

Example of a [Fluent Bit GoLang output plugin](https://github.com/fluent/fluent-bit/blob/master/GOLANG_OUTPUT_PLUGIN.md) that pushes events to a GRPC server.

## Usage

* Server

	```
	$ go run server/server.go
	```

* Fluent-Bit

	```
	$ make protoc
	$ make all
	$ fluent-bit -i dummy -e out_grpc.so -o grpc
	```
