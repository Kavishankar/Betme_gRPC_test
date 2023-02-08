package main

import (
	"betme_test/proto"
	"betme_test/utils"
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// gRPC server address
var addr string = "localhost:50051"

func main() {

	// Open new connection
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}

	// Close connection before main exits
	defer conn.Close()

	// Create New TestServiceClient and send request
	c := proto.NewTestServiceClient(conn)
	req := &proto.Dates{
		Dates: []*proto.Date{
			{
				Year:  2018,
				Month: 11,
				Day:   11,
			},
			{
				Year:  2018,
				Month: 11,
			},
		},
	}
	stream, err := c.GetFiles(context.Background(), req)
	if err != nil {
		log.Fatalf("Error getting files: %v\n", err)
	}

	// For each File response from stream, print filename, date and feed in log
	// TODO: Write the data received in clientdata dir?
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}

		// Parsing the date from fileData is different for each Feed
		if msg.Feed == proto.Feed_FEED_X {
			log.Printf("Received file %s with date %v from feed %v\n", msg.Path, utils.GetDateFromFeedXData([]byte(msg.Data)), msg.Feed)
		} else if msg.Feed == proto.Feed_FEED_Y {
			log.Printf("Received file %s with date %v from feed %v\n", msg.Path, utils.GetDateFromFeedYData([]byte(msg.Data)), msg.Feed)
		}
	}

}
