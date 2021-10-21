package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	pb "upsider.crawling/crawlingproto"
	"upsider.crawling/crawlingrepository"
)

type server struct{}

func (*server) UserHandler(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {

	newC := crawlingrepository.NewCrawling()
	today := time.Now()

	err := newC.Crawling(req.Pass, req.UserInput)
	if err != nil {
		return &pb.UserResponse{
			IsSuccess:     false,
			IllegalUrl:    false,
			IllegalAccess: false,
			IllegalLogin:  false,
			IllegalGet:    false,
		}, err
	}

	db := crawlingrepository.NewDatabase()
	if err := db.UserCreate(crawlingrepository.Users, req.UserInput.UserId, &today); err != nil {
		return &pb.UserResponse{
			IsSuccess: false,
		}, err
	}
	if err := db.BankCreate(req.UserInput.UserId, crawlingrepository.Banks, &today); <-err != nil {
		err := <-err
		return &pb.UserResponse{
			IsSuccess: false,
		}, err
	}

	if err := db.DetailCreate(req.UserInput.UserId, crawlingrepository.Details, &today); <-err != nil {
		err := <-err
		return &pb.UserResponse{
			IsSuccess: false,
		}, err
	}

	return &pb.UserResponse{
		IsSuccess: true,
	}, nil
}

func (*server) FreeeRead(ctx context.Context, req *pb.FreeeRequest) (*pb.FreeeResponse, error) {
	client, err := sql.Open("mysql", "root@/freee?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		log.Fatal(err)
	}
	cr := crawlingrepository.NewCrawlingRead(client)
	// 銀行口座取得

	offices, err := cr.OfficeRead(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, office := range offices {
		err := cr.BankRead(ctx, req, office.OfficeName, req.StartDay, req.LastDay)
		if err != nil {
			return nil, err
		}

		err = cr.CardRead(ctx, req, office.OfficeName, req.StartDay, req.LastDay)
		if err != nil {
			return nil, err
		}
		office.Banks = crawlingrepository.PbBanks
		office.Cards = crawlingrepository.PbCards
	}

	fmt.Println(len(offices[0].Cards.Card[1].Detail))

	return &pb.FreeeResponse{
		Office: offices,
	}, nil
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
