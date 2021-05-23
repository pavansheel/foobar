package main

import (
	"context"
	"fmt"
	"log"
	"time"

	proto "github.com/pavansheel/foobar/grpcclient"

	"google.golang.org/grpc"
)

const (
	address = "localhost:5066"
)

func main() {
	// Set up a connection.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	grainite := proto.NewGrainiteClient(conn)

	// Send 5 append requests to topic 'mytopic', containing 100 events each.
	for i := 0; i < 5; i++ {
		req := &proto.AppendRequest{}
		req.App = "myapp"
		req.TopicName = "mytopic"
		// Add 100 events, over 20 keys
		for j := 0; j < 100; j++ {
			var key string = fmt.Sprintf("key_%d", j%20)
			var payload string = fmt.Sprintf("payload_%d", (i*100)+j)
			event := &proto.Event{Key: []byte(key), Payload: []byte(payload)}
			req.Events = append(req.Events, event)
		}
		// Make the append request
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := grainite.TopicAppend(ctx, req)
		if err != nil {
			log.Fatalf("Could not Append: %v", err)
		}
		// Print the response
		var ok_cnt int = 0
		var failed_cnt int = 0
		for _, v := range resp.Status {
			if v.Status.Error == 0 {
				ok_cnt++
			} else {
				failed_cnt++
			}
		}
		log.Printf("Append: %v success %v failed", ok_cnt, failed_cnt)
	}
}
