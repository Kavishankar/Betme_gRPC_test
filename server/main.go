package main

import (
	"betme_test/proto"
	"betme_test/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	proto.TestServiceServer
}

// Address for gRPC server
var addr string = "localhost:50051"

var (
	DataMap    map[proto.Feed]string
	HandlerMap map[proto.Feed]func(*proto.Dates, proto.TestService_GetFilesServer)
)

func main() {

	// Populate DataMap
	DataMap = map[proto.Feed]string{
		proto.Feed_FEED_X: "data/feed_x",
		proto.Feed_FEED_Y: "data/feed_y",
	}

	// Populate HandlerMap
	HandlerMap = map[proto.Feed]func(*proto.Dates, proto.TestService_GetFilesServer){
		proto.Feed_FEED_X: FeedXHandler,
		proto.Feed_FEED_Y: FeedYHandler,
	}

	// Create a listener
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}

	log.Printf("Listening on %s...\n", addr)

	serv := grpc.NewServer()
	proto.RegisterTestServiceServer(serv, &Server{})

	// Serve using gRPC server
	if err = serv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}

}

// GetFiles - GetFiles rpc implementation
func (s *Server) GetFiles(dates *proto.Dates, stream proto.TestService_GetFilesServer) error {
	start := time.Now()
	log.Printf("Received request: %v\n", dates)
	wg := &sync.WaitGroup{}
	for _, handler := range HandlerMap {
		wg.Add(1)
		go func(handler func(*proto.Dates, proto.TestService_GetFilesServer)) {
			handler(dates, stream)
			wg.Done()
		}(handler)
	}
	wg.Wait()
	log.Printf("Finished Processing request in %v\n", time.Since(start))
	return nil
}

// FeedXHandler - Handles feed_x fir in data (dir walk, check date and stream data if needed)
func FeedXHandler(dates *proto.Dates, stream proto.TestService_GetFilesServer) {
	feedName := proto.Feed_FEED_X
	feedDir := DataMap[feedName]
	// Get all season_id dirs
	seasons, err := ioutil.ReadDir(feedDir)
	if err != nil {
		log.Fatalf("Error reading directory: %v\n", err)
	}
	// Loop through season_id dirs
	for _, dirs := range seasons {
		// Only dirs
		if !dirs.IsDir() {
			continue
		}
		fixturesDir := fmt.Sprintf("%s/%s/fixtures", feedDir, dirs.Name())
		// Get all fixtures dir inside season_id/fixtures dir
		fixtures, err := ioutil.ReadDir(fixturesDir)
		if err != nil {
			log.Fatalf("Error reading directory: %v\n", err)
		}
		// Loop through fixture dirs
		for _, fixture := range fixtures {
			// Only files
			if fixture.IsDir() {
				continue
			}
			// Generate full path for each file, check if date in file content matches
			// with any of the dates in request and stream the File data
			path := fmt.Sprintf("%s/%s", fixturesDir, fixture.Name())
			fileData, err := os.ReadFile(path)
			if err != nil {
				log.Fatalf("Error reading file: %v\n", err)
			}
			if utils.DataMatchesAnyDate(feedName, fileData, dates) {
				payload := proto.File{
					Feed: feedName,
					Path: path,
					Data: string(fileData),
				}
				stream.Send(&payload)
			}
		}
	}
}

// FeedYHandler - Handles feed_y fir in data (dir walk, check date and stream data if needed)
func FeedYHandler(dates *proto.Dates, stream proto.TestService_GetFilesServer) {
	feedName := proto.Feed_FEED_Y
	feedDir := DataMap[feedName]
	// Get all year dirs
	years, err := ioutil.ReadDir(feedDir)
	if err != nil {
		log.Fatalf("Error reading directory: %v\n", err)
	}
	// Loop through year dirs
	for _, year := range years {
		yearDir := fmt.Sprintf("%s/%s/fixtures", feedDir, year.Name())
		// Get all fixture dirs from year/fixtures dir
		fixtures, err := ioutil.ReadDir(yearDir)
		if err != nil {
			log.Fatalf("Error reading directory: %v\n", err)
		}
		// Loop through fixture dirs
		for _, fixture := range fixtures {
			if !fixture.IsDir() {
				continue
			}
			// Generate full path for each fixture.json, check if date in file content matches
			// with any of the dates in request and stream the File data
			path := fmt.Sprintf("%s/%s/fixture.json", yearDir, fixture.Name())
			fileData, err := os.ReadFile(path)
			if err != nil {
				log.Fatalf("Error reading file: %v\n", err)
			}
			if utils.DataMatchesAnyDate(feedName, fileData, dates) {
				payload := proto.File{
					Feed: feedName,
					Path: path,
					Data: string(fileData),
				}
				stream.Send(&payload)
			}
		}
	}
}
