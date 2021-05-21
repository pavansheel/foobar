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
	// set up a connection
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	grainite := proto.NewGrainiteClient(conn)

	// append 100 events to topic 'mytopic', using 20 keys
	req := &proto.AppendRequest{}
	req.App = "myapp"
	req.TopicName = "mytopic"
	for i := 0; i < 100; i++ {
		var key string = fmt.Sprintf("key_%d", i%20)
		var payload string = fmt.Sprintf("payload_%d", i)
		event := &proto.Event{Key: []byte(key), Payload: []byte(payload)}
		req.Events = append(req.Events, event)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := grainite.TopicAppend(ctx, req)
	if err != nil {
		log.Fatalf("could not Append: %v", err)
	}
	// print out the response
	log.Printf("append successful: %v", resp)
}
