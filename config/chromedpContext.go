package config

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

func NewChromedpContext() (ctx context.Context) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1920, 1080),
		chromedp.Flag("remote-debugging-port", "9222"),
	)

	chromedp.WithBrowserOption()
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ = chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)

	return ctx

}
