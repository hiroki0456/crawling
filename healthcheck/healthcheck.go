package healthcheck

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"upsider.crawling/config"
	pb "upsider.crawling/crawlingproto"
)

type HealthCheck interface {
	AccessCheck(req *pb.HealthCheckRequest) error
	LoginCheck(req *pb.HealthCheckRequest) error
	PageTransitionCheck(req *pb.HealthCheckRequest) error
}

type healthCheck struct{}

func NewHealthCheck() HealthCheck {
	return &healthCheck{}
}

func (*healthCheck) AccessCheck(req *pb.HealthCheckRequest) error {

	ctx, cancel := context.WithTimeout(config.NewChromedpContext(), 3*time.Minute)
	defer cancel()
	loginURL := "https://accounts.secure.freee.co.jp/login/accounting"

	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.Navigate(loginURL).Do(ctx)
			var illegalCheck string
			chromedp.Location(&illegalCheck).Do(ctx)
			if illegalCheck != "https://accounts.secure.freee.co.jp/login/accounting" {
				return fmt.Errorf("URLアクセス不可: %s", illegalCheck)
			}
			return nil
		}),
	)
	return err
}

func (*healthCheck) LoginCheck(req *pb.HealthCheckRequest) error {

	ctx, cancel := context.WithTimeout(config.NewChromedpContext(), 3*time.Minute)
	defer cancel()
	loginURL := "https://accounts.secure.freee.co.jp/login/accounting"

	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.Navigate(loginURL).Do(ctx)
			chromedp.WaitVisible(`//input[@name="email"]`).Do(ctx)
			chromedp.WaitVisible(`//input[@name="password"]`).Do(ctx)
			chromedp.SendKeys(`//input[@name="email"]`, req.UserId).Do(ctx)
			chromedp.SendKeys(`//input[@name="password"]`, req.Pass).Do(ctx)
			chromedp.Submit(`//input[@name="password"]`).Do(ctx)

			var illegalCheck string
			chromedp.Location(&illegalCheck).Do(ctx)
			if illegalCheck != "https://secure.freee.co.jp/" {
				return fmt.Errorf("ログイン不可: %s", illegalCheck)
			}
			return nil
		}),
	)
	return err
}

func (*healthCheck) PageTransitionCheck(req *pb.HealthCheckRequest) error {

	ctx, cancel := context.WithTimeout(config.NewChromedpContext(), 3*time.Minute)
	defer cancel()
	officeURL := "https://secure.freee.co.jp/user/show_companies"
	detailURL := "https://secure.freee.co.jp/wallet_txns"
	topURL := "https://secure.freee.co.jp/"
	loginURL := "https://accounts.secure.freee.co.jp/login/accounting"

	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.Navigate(loginURL).Do(ctx)
			chromedp.WaitVisible(`//input[@name="email"]`).Do(ctx)
			chromedp.WaitVisible(`//input[@name="password"]`).Do(ctx)
			chromedp.SendKeys(`//input[@name="email"]`, req.UserId).Do(ctx)
			chromedp.SendKeys(`//input[@name="password"]`, req.Pass).Do(ctx)
			chromedp.Submit(`//input[@name="password"]`).Do(ctx)

			var illegalCheck string
			chromedp.Location(&illegalCheck).Do(ctx)

			if illegalCheck != topURL {
				return fmt.Errorf("画面遷移エラー【topページ】: %s", illegalCheck)
			}
			chromedp.WaitVisible(`#footer`).Do(ctx)

			chromedp.Navigate(detailURL).Do(ctx)
			chromedp.Location(&illegalCheck).Do(ctx)
			if illegalCheck != detailURL {
				return fmt.Errorf("画面遷移エラー【口座詳細ページ】: %s", illegalCheck)
			}
			chromedp.WaitVisible(`#footer`).Do(ctx)

			chromedp.Navigate(officeURL).Do(ctx)
			chromedp.Location(&illegalCheck).Do(ctx)
			if illegalCheck != officeURL {
				return fmt.Errorf("面遷移エラー【事業所詳細ページ】: %s", illegalCheck)
			}
			return nil

		}),
	)
	return err
}
