package healthcheck

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"upsider.crawling/config"
	pb "upsider.crawling/crawlingproto"
)

type HealthCheck interface {
	AccessCheck(req *pb.HealthCheckRequest) error
	LoginCheck(req *pb.HealthCheckRequest) error
	PageTransitionCheck(req *pb.HealthCheckRequest) error
	CrawlingCheck(req *pb.HealthCheckRequest) error
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

	loginURL := "https://accounts.secure.freee.co.jp/login/accounting"
	topURL := "https://secure.freee.co.jp/"
	detailURL := "https://secure.freee.co.jp/wallet_txns"
	officeURL := "https://secure.freee.co.jp/user/show_companies"

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
				return fmt.Errorf("画面遷移エラー【Topページ】: %s", illegalCheck)
			}

			chromedp.WaitVisible(`#footer`).Do(ctx)
			chromedp.Navigate(detailURL).Do(ctx)
			chromedp.Location(&illegalCheck).Do(ctx)
			var resText string
			findText := &resText
			if illegalCheck != detailURL {
				chromedp.Text(`/html/body/div[2]/div/div[1]/h1`, findText).Do(ctx)
				if resText == "ページが見つかりませんでした" {
					return fmt.Errorf("画面遷移エラー【口座明細一覧ページ】: %s", illegalCheck)
				}
			}

			chromedp.WaitVisible(`#footer`).Do(ctx)
			chromedp.Navigate(officeURL).Do(ctx)
			chromedp.Location(&illegalCheck).Do(ctx)
			if illegalCheck != detailURL {
				chromedp.Text(`/html/body/div[2]/div/div[1]/h1`, findText).Do(ctx)
				if resText == "ページが見つかりませんでした" {
					return fmt.Errorf("画面遷移エラー【事業所切り替えページ】: %s", illegalCheck)
				}
			}
			return nil
		}),
	)
	return err
}

func (*healthCheck) CrawlingCheck(req *pb.HealthCheckRequest) error {

	var bankNameAndId []map[string]string
	var officeName string

	ctx, cancel := context.WithTimeout(config.NewChromedpContext(), 3*time.Minute)
	defer cancel()

	loginURL := "https://accounts.secure.freee.co.jp/login/accounting"
	officeURL := "https://secure.freee.co.jp/user/show_companies"
	detailURL := "https://secure.freee.co.jp/wallet_txns"
	topURL := "https://secure.freee.co.jp/"
	illegalCheck := ""

	loginIdSel := `/html/body/div[3]/div/div[1]/form/div/div[2]/input`
	loginPassSel := `/html/body/div[3]/div/div[1]/form/div/div[3]/input`
	loginButtonSel := `.btn.btn-primary.login-page-button.login-button.transition`

	loginActionFunc := chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.Navigate(loginURL).Do(ctx)
		chromedp.Location(&illegalCheck).Do(ctx)
		if illegalCheck == "chrome-error://chromewebdata/" {
			return fmt.Errorf("URLの遷移に失敗しました: %s", illegalCheck)
		}

		chromedp.ScrollIntoView(`body`).Do(ctx)

		chromedp.SetValue(loginIdSel, req.UserId, chromedp.BySearch).Do(ctx)
		chromedp.SetValue(loginPassSel, req.Pass, chromedp.BySearch).Do(ctx)
		chromedp.Click(loginButtonSel).Do(ctx)
		chromedp.Location(&illegalCheck).Do(ctx)
		if illegalCheck == "https://accounts.secure.freee.co.jp/login/accounting?a=false&e=0&o=false" {
			return fmt.Errorf("ログインに失敗しました: %s", illegalCheck)
		}
		chromedp.WaitVisible(`.walletable_group___StyledDiv5-sc-1uncx9n-4.kQEfxP`, chromedp.ByQuery).Do(ctx)

		return nil
	})

	getBanksActionFunc := chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.Navigate(topURL).Do(ctx)
		chromedp.Location(&illegalCheck).Do(ctx)
		if illegalCheck == "chrome-error://chromewebdata/" {
			return fmt.Errorf("画面遷移エラー【Topページ】:  %s", illegalCheck)
		}

		err := chromedp.Run(
			ctx,
			RunWithTimeOut(&ctx, 3, chromedp.Tasks{
				chromedp.WaitVisible(`.walletable_controls___StyledSpan-sc-11p3ona-0`, chromedp.ByQuery),
			}),
		)
		if err != nil {
			return fmt.Errorf("項目取得エラー【口座詳細を開くボタン】: %s", illegalCheck)
		}

		bankNode := []*cdp.Node{}

		err = chromedp.Run(
			ctx,
			RunWithTimeOut(&ctx, 3, chromedp.Tasks{
				chromedp.Nodes(`.walletable_group___StyledDiv-sc-1uncx9n-0.dHyIIm`, &bankNode, chromedp.ByQueryAll),
			}),
		)
		if err != nil {
			return fmt.Errorf("項目取得エラー【銀行、カードの有無】: %s", illegalCheck)
		}

		if len(bankNode) == 0 {
			return fmt.Errorf("銀行、並びにカード情報が取れませんでした。")
		}

		var lastCommit string
		err = chromedp.Run(
			ctx,
			RunWithTimeOut(&ctx, 3, chromedp.Tasks{
				chromedp.Text(`.sync_all_walletables___StyledDiv2-tf1121-1`, &lastCommit, chromedp.ByQuery),
			}),
		)
		if err != nil {
			return fmt.Errorf("項目取得エラー【最終同期日時】: %s", illegalCheck)
		}

		var wg sync.WaitGroup
		wg.Add(len(bankNode))
		for _, n := range bankNode {
			go func(n *cdp.Node) {
				defer wg.Done()
				res, err := dom.GetOuterHTML().WithNodeID(n.NodeID).Do(ctx)
				if err != nil {
					fmt.Printf("error %s", err)
				}
				scrapingOfBanks(res, lastCommit)
			}(n)
		}
		return nil
	})

	getDetailActionFunc := chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.Navigate(detailURL).Do(ctx)
		chromedp.Location(&illegalCheck).Do(ctx)
		if illegalCheck == "chrome-error://chromewebdata/" {
			return fmt.Errorf("URLの遷移に失敗しました: %s", illegalCheck)
		}

		err := chromedp.Run(
			ctx,
			RunWithTimeOut(&ctx, 3, chromedp.Tasks{
				chromedp.ScrollIntoView(`.active.sw-active`, chromedp.ByQuery),
			}),
		)
		if err != nil {
			return fmt.Errorf("項目取得エラー【明細一覧のスクロールバー】: %s", illegalCheck)
		}

		detailBankNode := []*cdp.Node{}
		err = chromedp.Run(
			ctx,
			RunWithTimeOut(&ctx, 3, chromedp.Tasks{
				chromedp.Nodes(`select#walletable > option`, &detailBankNode, chromedp.ByQueryAll),
			}),
		)
		if err != nil {
			return fmt.Errorf("項目取得エラー【明細一覧】: %s", illegalCheck)
		}
		detailBankNode = append(detailBankNode[:0], detailBankNode[1:]...)

		for _, bankIdNode := range detailBankNode {
			res, _ := dom.GetOuterHTML().WithNodeID(bankIdNode.NodeID).Do(ctx)
			bankName, err := scrapingDetailBankName(res)
			if err != nil {
				return err
			}
			var bankId string
			chromedp.Value(bankIdNode.FullXPath(), &bankId, chromedp.BySearch).Do(ctx)

			bankNameAndId = append(bankNameAndId, map[string]string{"officeName": officeName, "bankName": bankName, "bankId": bankId})
		}

		paginationNode := []*cdp.Node{}
		err = chromedp.Run(
			ctx,
			RunWithTimeOut(&ctx, 3, chromedp.Tasks{
				chromedp.Nodes(`.sw-pagination > ul > li > a`, &paginationNode, chromedp.ByQueryAll),
			}),
		)
		if err != nil {
			return fmt.Errorf("項目取得エラー【明細一覧ページ選択】: %s", illegalCheck)
		}

		var loopNumString string
		chromedp.Text(paginationNode[len(paginationNode)-2].FullXPath(), &loopNumString, chromedp.BySearch).Do(ctx)
		loopNum, _ := strconv.Atoi(loopNumString)
		fmt.Println(loopNum)

		detailsDom := []string{}
		sum := 7
		for i := 0; i < loopNum; i++ {
			if i+1 == 1 {
				pagingDetailContent := ""
				err = chromedp.Run(
					ctx,
					RunWithTimeOut(&ctx, 3, chromedp.Tasks{
						chromedp.OuterHTML(`.wallet-txn-list-table`, &pagingDetailContent, chromedp.ByQueryAll),
					}),
				)
				if err != nil {
					return fmt.Errorf("項目取得エラー【明細一覧表】: %s", illegalCheck)
				}
				detailsDom = append(detailsDom, pagingDetailContent)
				continue
			}
			if i+1 < 6 {
				err = chromedp.Run(
					ctx,
					RunWithTimeOut(&ctx, 3, chromedp.Tasks{
						chromedp.Click(`/html/body/div[2]/div/div/div[3]/div/ul/li[`+strconv.Itoa(i+2)+`]/a`, chromedp.NodeVisible),
					}),
				)
				if err != nil {
					return fmt.Errorf("項目取得エラー【明細一覧】: %s", illegalCheck)
				}
				chromedp.WaitVisible(paginationNode[i+1].Parent.FullXPath()+`[@class='active sw-active']`, chromedp.BySearch).Do(ctx)
				pagingDetailContent := ""
				err = chromedp.Run(
					ctx,
					RunWithTimeOut(&ctx, 3, chromedp.Tasks{
						chromedp.OuterHTML(`.wallet-txn-list-table`, &pagingDetailContent, chromedp.ByQueryAll),
					}),
				)
				if err != nil {
					return fmt.Errorf("項目取得エラー【明細一覧表】: %s", illegalCheck)
				}

				detailsDom = append(detailsDom, pagingDetailContent)

				continue
			}
			if i+1 == 6 {
				chromedp.Click(paginationNode[i+1].FullXPath(), chromedp.NodeVisible).Do(ctx)
				chromedp.WaitVisible(paginationNode[i+1].FullXPath()+`[@data-num='7']`, chromedp.BySearch).Do(ctx)
				pagingDetailContent := ""
				err = chromedp.Run(
					ctx,
					RunWithTimeOut(&ctx, 3, chromedp.Tasks{
						chromedp.OuterHTML(`.wallet-txn-list-table`, &pagingDetailContent, chromedp.ByQueryAll),
					}),
				)
				if err != nil {
					return fmt.Errorf("項目取得エラー【明細一覧表】: %s", illegalCheck)
				}
				detailsDom = append(detailsDom, pagingDetailContent)
				continue
			}

			if i+1 < loopNum-3 {
				err = chromedp.Run(
					ctx,
					RunWithTimeOut(&ctx, 3, chromedp.Tasks{
						chromedp.Click(`/html/body/div[2]/div/div/div[3]/div/ul/li[7]/a`, chromedp.NodeVisible),
					}),
				)
				if err != nil {
					return fmt.Errorf("項目取得エラー【明細一覧ページ切り替え】: %s", illegalCheck)
				}
				err = chromedp.Run(
					ctx,
					RunWithTimeOut(&ctx, 3, chromedp.Tasks{
						chromedp.WaitVisible(`/html/body/div[2]/div/div/div[3]/div/ul/li[6]/a[@data-num='`+strconv.Itoa(i+1)+`']`, chromedp.BySearch),
					}),
				)
				if err != nil {
					return fmt.Errorf("項目取得エラー【明細一覧ページ切り替え】: %s", illegalCheck)
				}
				pagingDetailContent := ""
				chromedp.OuterHTML(`.wallet-txn-list-table`, &pagingDetailContent, chromedp.ByQueryAll).Do(ctx)
				detailsDom = append(detailsDom, pagingDetailContent)
				continue
			}

			if i+1 <= loopNum {
				err = chromedp.Run(
					ctx,
					RunWithTimeOut(&ctx, 3, chromedp.Tasks{
						chromedp.Click(`/html/body/div[2]/div/div/div[3]/div/ul/li[`+strconv.Itoa(sum)+`]/a`, chromedp.BySearch),
					}),
				)
				if err != nil {
					return fmt.Errorf("項目取得エラー【明細一覧ページ切り替え】: %s", illegalCheck)
				}
				err = chromedp.Run(
					ctx,
					RunWithTimeOut(&ctx, 3, chromedp.Tasks{
						chromedp.WaitVisible(`/html/body/div[2]/div/div/div[3]/div/ul/li[`+strconv.Itoa(sum)+`]`+`[@class='active sw-active']`, chromedp.BySearch),
					}),
				)
				if err != nil {
					return fmt.Errorf("項目取得エラー【明細一覧ページ切り替え】: %s", illegalCheck)
				}
				sum += 1
				pagingDetailContent := ""
				chromedp.OuterHTML(`.wallet-txn-list-table`, &pagingDetailContent, chromedp.ByQueryAll).Do(ctx)
				detailsDom = append(detailsDom, pagingDetailContent)
				continue
			}
		}

		scrapingDetails(detailsDom, req.UserId)
		return nil
	})

	crawlingActionFunc := chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.Navigate(officeURL).Do(ctx)
		chromedp.Location(&illegalCheck).Do(ctx)
		if illegalCheck == "chrome-error://chromewebdata/" {
			return fmt.Errorf("URLの遷移に失敗しました: %s", illegalCheck)
		}
		chromedp.WaitVisible(`#footer`).Do(ctx)
		officeNode := []*cdp.Node{}
		err := chromedp.Run(
			ctx,
			RunWithTimeOut(&ctx, 3, chromedp.Tasks{
				chromedp.Nodes(`table.list-table > tbody > tr`, &officeNode, chromedp.ByQueryAll),
			}),
		)
		if err != nil {
			return fmt.Errorf("項目取得エラー【事業所の選択】: %s", illegalCheck)
		}

		for i := range officeNode {
			chromedp.Nodes(`table.list-table > tbody > tr`, &officeNode, chromedp.ByQueryAll).Do(ctx)
			officeButtonNode, _ := dom.QuerySelector(officeNode[i].NodeID, ".btn.btn-primary").Do(ctx)

			buttonDom, _ := dom.GetOuterHTML().WithNodeID(officeButtonNode).Do(ctx)
			if buttonDom == "" {
				fmt.Println("test")
				chromedp.Text(officeNode[i].FullXPath()+`/td[1]`, &officeName, chromedp.BySearch).Do(ctx)
				chromedp.Navigate(topURL).Do(ctx)
				chromedp.WaitVisible(`.walletable_controls___StyledSpan-sc-11p3ona-0`, chromedp.ByQuery).Do(ctx)
				err := chromedp.Run(ctx, getBanksActionFunc, getDetailActionFunc)
				if err != nil {
					return err
				}
				chromedp.Navigate(officeURL).Do(ctx)
				chromedp.WaitVisible(`#footer`, chromedp.ByID).Do(ctx)
				continue
			}
			fmt.Println("test2")
			chromedp.Text(officeNode[i].FullXPath()+`/td[1]`, &officeName, chromedp.BySearch).Do(ctx)
			chromedp.Click(officeNode[i].FullXPath()+`/td[5]/a`, chromedp.BySearch).Do(ctx)
			chromedp.WaitVisible(`.walletable_controls___StyledSpan-sc-11p3ona-0`, chromedp.ByQuery).Do(ctx)
			err := chromedp.Run(ctx, getBanksActionFunc, getDetailActionFunc)
			if err != nil {
				return err
			}
			chromedp.Navigate(officeURL).Do(ctx)
			chromedp.WaitVisible(`#footer`, chromedp.ByID).Do(ctx)
		}
		return nil
	})

	err := chromedp.Run(ctx,
		loginActionFunc,
		crawlingActionFunc,
	)
	if err != nil {
		return err
	}

	return nil
}

func scrapingOfBanks(d string, lastCommit string) error {

	readerCurContents := strings.NewReader(d)
	contentsDom, err := goquery.NewDocumentFromReader(readerCurContents)
	if err != nil {
		return err
	}

	contentsDom.Find(`div.walletable___StyledDiv-sc-3etvmj-0`).Each(func(i int, v *goquery.Selection) {
		if err != nil {
			log.Fatalf("intへのconvertに失敗しました: %v", err)
			return
		}
	})

	return nil
}

func scrapingDetailBankName(d string) (string, error) {
	readerCurContents := strings.NewReader(d)
	contentsDom, err := goquery.NewDocumentFromReader(readerCurContents)
	if err != nil {
		return "", err
	}

	return contentsDom.Find(`option`).Text(), nil
}

func scrapingDetails(dl []string, userId string) error {
	// var val string
	cancel := true

	for _, d := range dl {
		if !cancel {
			break
		}
		readerCurContents := strings.NewReader(d)
		contentsDom, err := goquery.NewDocumentFromReader(readerCurContents)
		if err != nil {
			return err
		}

		contentsDom.Find("tr.line").EachWithBreak(func(i int, v *goquery.Selection) bool {
			if err != nil {
				log.Printf("最終更新日の取得に失敗しました: %v", err)
			}
			return true
		})
	}
	return nil
}

func RunWithTimeOut(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}
