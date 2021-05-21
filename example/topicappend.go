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
	address = "localhost:10000"
)

func main() {
	// Set up a connection.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	grainite := proto.NewGrainiteClient(conn)

	// Send 5 append requests. Each appends 100 events to topic 'mytopic',
	// using 20 keys.
	req := &proto.AppendRequest{}
	req.App = "myapp"
	req.TopicName = "mytopic"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := 0; i < 5; i++ {
		for j := 0; j < 100; j++ {
			var key string = fmt.Sprintf("key_%d", j%20)
			var payload string = fmt.Sprintf("payload_%d", (i*100)+j)
			event := &proto.Event{Key: []byte(key), Payload: []byte(payload)}
			req.Events = append(req.Events, event)
		}
		resp, err := grainite.TopicAppend(ctx, req)
		if err != nil {
			log.Fatalf("Could not Append: %v", err)
		}
		// print out the response
		log.Printf("Append result: %v", resp)
	}
}
