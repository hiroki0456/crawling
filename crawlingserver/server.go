package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	pb "upsider.crawling/crawlingproto"
	"upsider.crawling/crawlingrepository"
	"upsider.crawling/healthcheck"
)

type server struct{}

func (*server) UserHandler(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {

	newC := crawlingrepository.NewCrawling()
	today := time.Now()

	err := newC.Crawling(req.Pass, req.UserInput)
	if err != nil {
		return &pb.UserResponse{
			IsSuccess: false,
		}, err
	}

	db := crawlingrepository.NewDatabase()
	if err := db.UserCreate(crawlingrepository.Users, req.UserInput.UserId, &today); err != nil {
		return &pb.UserResponse{
			IsSuccess: false,
		}, err
	}
	if err := db.BankCreate(req.UserInput.UserId, crawlingrepository.Banks, &today); err != nil {
		return &pb.UserResponse{
			IsSuccess: false,
		}, err
	}

	if err := db.DetailCreate(req.UserInput.UserId, crawlingrepository.Details, &today); err != nil {
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

	return &pb.FreeeResponse{
		Office: offices,
	}, nil
}

func (*server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (res *pb.HealthCheckResponse, err error) {
	h := healthcheck.NewHealthCheck()
	err = h.AccessCheck(req)
	if err != nil {
		healthcheck.NoticeSlack(err)
		return nil, err
	}

	err = h.LoginCheck(req)
	if err != nil {
		healthcheck.NoticeSlack(err)
		return nil, err
	}

	err = h.PageTransitionCheck(req)
	if err != nil {
		healthcheck.NoticeSlack(err)
		return nil, err
	}

	return &pb.HealthCheckResponse{
		IsSuccess: true,
	}, nil
}
