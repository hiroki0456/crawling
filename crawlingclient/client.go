package main

import (
	"context"
	"fmt"
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
	// crawlingHealthCheck(c)
	// //
}

func crawlingRead(c pb.CrawlingServiceClient) {
	req := &pb.FreeeRequest{
		UserInput: &pb.UserInput{
			UserId: "volleyball0456@gmail.com",
		},
		// StartDay: "2021-10-14",
		// LastDay:  "2021-10-18",
	}

	res, err := c.FreeeRead(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	fmt.Println(res.Office[1].OfficeName)
	// log.Println(res.Office[0])
	// log.Println(res.Office[1])
	// log.Println(res.Office[1])
	// log.Println(res.Office[0].Banks.Bank[1].Detail)
	// log.Println(res.Office[1].Banks.Bank[0].Detail)

}

func crawlingWrite(c pb.CrawlingServiceClient) {
	req := &pb.UserRequest{
		UserInput: &pb.UserInput{
			UserId: "volleyball0456@gmail.com",
		},
		Pass:     "hiro0456",
		SiteKind: 1,
	}

	res, err := c.UserHandler(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Println(res)
}

func crawlingHealthCheck(c pb.CrawlingServiceClient) {
	req := &pb.HealthCheckRequest{
		UserId: "volleyball0456@gmail.com",
		Pass:   "hiro0456",
	}

	res, err := c.HealthCheck(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
