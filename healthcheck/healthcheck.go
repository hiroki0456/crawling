package healthcheck

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"upsider.crawling/config"
	pb "upsider.crawling/crawlingproto"
)

type HealthCheck interface {
	AccessCheck(req *pb.HealthCheckRequest) error
}

type healthCheck struct{}

func NewHealthCheck() HealthCheck {
	return &healthCheck{}
}

func (*healthCheck) AccessCheck(req *pb.HealthCheckRequest) error {

	ctx := config.NewChromedpContext()
	loginURL := "https://accounts.secure.freee.co.jp/login/accounting"

	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.Navigate(loginURL).Do(ctx)
			var illegalCheck string
			chromedp.Location(&illegalCheck).Do(ctx)
			fmt.Println(illegalCheck)
			if illegalCheck != "https://accounts.secure.freee.co.jp/login/accounting" {
				return fmt.Errorf("topページに遷移できません: %s", illegalCheck)
			}
			return nil
		}),
	)
	return err
}
