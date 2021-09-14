package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"upsider.crawling/crawlinginterface"
	pb "upsider.crawling/crawlingproto"
	"upsider.crawling/crawlingrepository"
)

type server struct {
}

func (*server) UserHandler(ctx context.Context, req *pb.UserRequest) (*pb.ApiResponse, error) {
	// pass := req.Pass
	r := &crawlingrepository.Freee{
		CrawlingSite: crawlingrepository.CrawlingSite{
			Pass:  req.Pass,
			Input: req.UserInput,
		},
	}

	newC := crawlinginterface.NewCrawling(r)

	newC.Exec()

	return nil, nil
}
func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCrawlingServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
