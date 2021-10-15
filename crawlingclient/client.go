package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	pb "upsider.crawling/crawlingproto"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := pb.NewCrawlingServiceClient(cc)

	// crawlingWrite(c)
	crawlingRead(c)
	// //
}

func crawlingRead(c pb.CrawlingServiceClient) {
	req := &pb.FreeeRequest{
		UserInput: &pb.UserInput{
			UserId: "volleyball0456@gmail.com",
		},
		Pass:     "volleyball0456",
		StartDay: "2021-10-16",
		// LastDay: "2021-10-",
	}

	res, err := c.FreeeRead(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Println(res)
}

func crawlingWrite(c pb.CrawlingServiceClient) {
	req := &pb.UserRequest{
		UserInput: &pb.UserInput{
			UserId: "volleyball0456@gmail.com",
		},
		Pass:     "volleyball0456",
		SiteKind: 1,
	}

	res, err := c.UserHandler(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Println(res)
}
