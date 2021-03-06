package main

import (
	"context"
	"log"
	"time"

	"../api"
	"github.com/fluent/fluent-bit-go/output"
	"google.golang.org/grpc"
)
import (
	"C"
	"fmt"
	"unsafe"
)

var clientConn *grpc.ClientConn

//export FLBPluginRegister
func FLBPluginRegister(ctx unsafe.Pointer) int {
	return output.FLBPluginRegister(ctx, "grpc", "GRPC output plugin")
}

//export FLBPluginInit
func FLBPluginInit(ctx unsafe.Pointer) int {
	// Example to retrieve an optional configuration parameter
	param := output.FLBPluginConfigKey(ctx, "param")
	fmt.Printf("[out-grpc] plugin parameter = '%s'\n", param)

	var err error
	clientConn, err = grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	var (
		ret    int
		ts     interface{}
		record map[interface{}]interface{}
	)

	var count int

	// Create Fluent Bit decoder
	dec := output.NewDecoder(data, int(length))

	c := api.NewEventClient(clientConn)

	// Iterate Records
	count = 0
	for {
		// Record
		ret, ts, record = output.GetRecord(dec)
		if ret != 0 {
			break
		}

		rec := make(map[string]string)
		for k, v := range record {
			strKey := fmt.Sprintf("%s", k)
			strValue := fmt.Sprintf("%s", v)
			rec[strKey] = strValue
		}

		// Timestamp
		var timestamp time.Time
		switch tts := ts.(type) {
		case output.FLBTime:
			timestamp = tts.Time
		case uint64:
			// From our observation, when ts is of type uint64 it appears to
			// be the amount of seconds since unix epoch.
			timestamp = time.Unix(int64(tts), 0)
		default:
			timestamp = time.Now()
		}

		log.Printf("[%d] %s: %s %s\n", count, timestamp.String(), C.GoString(tag), record)
		response, err := c.RecordEvents(context.Background(), &api.Record{Tag: C.GoString(tag), Record: rec})
		if err != nil {
			log.Fatalf("Error calling RecordEvents: %s", err)
		}
		log.Printf("Response from server: %d", response.EventCount)
		count++
	}

	// Return options:
	//
	// output.FLB_OK    = data have been processed.
	// output.FLB_ERROR = unrecoverable error, do not try this again.
	// output.FLB_RETRY = retry to flush later.
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	clientConn.Close()
	return output.FLB_OK
}

func main() {
}
