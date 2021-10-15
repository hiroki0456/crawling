package crawlingrepository

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
	"github.com/google/uuid"
	"upsider.crawling/crawlingproto"
)

type CrawlingInterface interface {
	Crawling(pass string, input *crawlingproto.UserInput) error
}

type crawlingRepository struct{}

func NewCrawling() CrawlingInterface {
	return &crawlingRepository{}
}

type User struct {
	Id               string     `spanner:"Id"`
	UserIdOfficeName string     `spanner:"UserIdOfficeName"`
	UserId           string     `spanner:"UserId"`
	OfficeName       string     `spanner:"OfficeName"`
	LastId           string     `spanner:"LastId"`
	UpdatedAt        *time.Time `spanner:"updatedAt"`
}

type Bank struct {
	Id         string `spanner:"Id"`
	UserId     string `spanner:"UserId"`
	BankId     string `spanner:"BankId"`
	lastCommit string `spanner:"LastCommit"`
	OfficeName string `spanner:"OfficeName"`
	BankName   string `spanner:"BankName"`
	Amount     int64  `spanner:"Amount"`
	Kind       string `spanner:"Kind"`
}

type Detail struct {
	Id             string    `spanner:"Id"`
	UserId         string    `spanner:"UserId"`
	BankId         string    `spanner:"BankId"`
	OfficeName     string    `spanner:"OfficeName"`
	BankName       string    `spanner:"BankName"`
	TradingDate    time.Time `spanner:"TradingDate"`
	TradingContent string    `spanner:"TradingContent"`
	Payment        int64     `spanner:"Payment"`
	Withdrawal     int64     `spanner:"Withdrawal"`
	Balance        int64     `spanner:"Balance"`
	UpdatedDate    time.Time `spanner:"UpdatedDate"`
	GettingDate    time.Time `spanner:"GettingDate"`
}

// スクレイピング済みの口座情報と明細情報
var Users []*User
var Banks []*Bank
var Details []*Detail

// スクレイピング時に必要な事業所名、銀行口座名、銀行口座ID
var bankNameAndId []map[string]string
var officeName string

func (c *crawlingRepository) Crawling(pass string, input *crawlingproto.UserInput) error {

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), chromedp.Flag("headless", false))
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
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

		chromedp.SetValue(loginIdSel, input.UserId, chromedp.BySearch).Do(ctx)
		chromedp.SetValue(loginPassSel, pass, chromedp.BySearch).Do(ctx)
		chromedp.Click(loginButtonSel).Do(ctx)
		chromedp.Location(&illegalCheck).Do(ctx)
		fmt.Println(illegalCheck)
		chromedp.WaitVisible(`.walletable_group___StyledDiv5-sc-1uncx9n-4.kQEfxP`, chromedp.ByQuery).Do(ctx)

		return nil
	})

	getBanksActionFunc := chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.Navigate(topURL).Do(ctx)
		chromedp.Location(&illegalCheck).Do(ctx)
		if illegalCheck == "chrome-error://chromewebdata/" {
			return fmt.Errorf("URLの遷移に失敗しました: %s", illegalCheck)
		}
		chromedp.WaitVisible(`.walletable_controls___StyledSpan-sc-11p3ona-0`, chromedp.ByQuery).Do(ctx)
		lastCommitNodes := []*cdp.Node{}
		chromedp.Nodes(`.walletable_controls___StyledSpan-sc-11p3ona-0`, &lastCommitNodes, chromedp.ByQueryAll).Do(ctx)
		for _, n := range lastCommitNodes {
			chromedp.Click(n.FullXPath(), chromedp.NodeVisible).Do(ctx)
			chromedp.WaitVisible(`.walletable_controls___StyledMdExpandLess-sc-11p3ona-1`, chromedp.ByQuery).Do(ctx)
		}
		bankNode := []*cdp.Node{}
		chromedp.Nodes(`.walletable_group___StyledDiv-sc-1uncx9n-0.dHyIIm`, &bankNode, chromedp.ByQueryAll).Do(ctx)

		// var lastCommit string
		// chromedp.Text(`.sync_all_walletables___StyledDiv2-tf1121-1`, &lastCommit, chromedp.ByQuery).Do(ctx)
		var wg sync.WaitGroup
		wg.Add(len(bankNode))
		for _, n := range bankNode {
			go func(n *cdp.Node) {
				defer wg.Done()
				res, err := dom.GetOuterHTML().WithNodeID(n.NodeID).Do(ctx)
				if err != nil {
					fmt.Printf("error %s", err)
				}
				scrapingOfBanks(res)
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
		chromedp.ScrollIntoView(`#footer`, chromedp.ByQuery).Do(ctx)
		detailBankNode := []*cdp.Node{}
		chromedp.Nodes(`select#walletable > option`, &detailBankNode, chromedp.ByQueryAll).Do(ctx)
		detailBankNode = append(detailBankNode[:0], detailBankNode[1:]...)

		for _, bankIdNode := range detailBankNode {
			res, _ := dom.GetOuterHTML().WithNodeID(bankIdNode.NodeID).Do(ctx)
			bankName, err := scrapingDetailBankName(res)
			if err != nil {
				return err
			}
			var bankId string
			chromedp.Value(bankIdNode.FullXPath(), &bankId, chromedp.BySearch).Do(ctx)

			for _, bank := range Banks {
				if bank.BankName == bankName && bank.BankId == "" {
					bank.BankId = bankId
				}
			}
			bankNameAndId = append(bankNameAndId, map[string]string{"bankName": bankName, "bankId": bankId})
		}

		paginationNode := []*cdp.Node{}
		chromedp.Nodes(`.sw-pagination > ul > li > a`, &paginationNode, chromedp.ByQueryAll).Do(ctx)

		detailsDom := []string{}
		for i := 0; i < len(paginationNode)-2; i++ {
			if i+1 == 1 {
				pagingDetailContent := ""
				chromedp.OuterHTML(`.wallet-txn-list-table`, &pagingDetailContent, chromedp.ByQueryAll).Do(ctx)
				detailsDom = append(detailsDom, pagingDetailContent)
				continue
			}
			chromedp.Click(paginationNode[i+1].FullXPath(), chromedp.NodeVisible).Do(ctx)
			chromedp.WaitVisible(paginationNode[i+1].Parent.FullXPath()+`[@class='active sw-active']`, chromedp.BySearch).Do(ctx)
			pagingDetailContent := ""
			chromedp.OuterHTML(`.wallet-txn-list-table`, &pagingDetailContent, chromedp.ByQueryAll).Do(ctx)

			detailsDom = append(detailsDom, pagingDetailContent)
		}

		scrapingDetails(detailsDom, input.UserId)
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
		chromedp.Nodes(`table.list-table > tbody > tr`, &officeNode, chromedp.ByQueryAll).Do(ctx)

		for i := range officeNode {
			chromedp.Nodes(`table.list-table > tbody > tr`, &officeNode, chromedp.ByQueryAll).Do(ctx)
			officeButtonNode, _ := dom.QuerySelector(officeNode[i].NodeID, ".btn.btn-primary").Do(ctx)

			buttonDom, _ := dom.GetOuterHTML().WithNodeID(officeButtonNode).Do(ctx)
			if buttonDom == "" {
				chromedp.Text(officeNode[i].FullXPath()+`/td[1]`, &officeName, chromedp.BySearch).Do(ctx)
				chromedp.Navigate(topURL).Do(ctx)
				chromedp.WaitVisible(`.walletable_controls___StyledSpan-sc-11p3ona-0`, chromedp.ByQuery).Do(ctx)
				chromedp.Run(ctx, getBanksActionFunc, getDetailActionFunc)
				chromedp.Navigate(officeURL).Do(ctx)
				chromedp.WaitVisible(`#footer`).Do(ctx)
				continue
			}

			chromedp.Text(officeNode[i].FullXPath()+`/td[1]`, &officeName, chromedp.BySearch).Do(ctx)
			chromedp.Click(officeNode[i].FullXPath()+`/td[5]/a`, chromedp.BySearch).Do(ctx)
			chromedp.WaitVisible(`.walletable_controls___StyledSpan-sc-11p3ona-0`, chromedp.ByQuery).Do(ctx)
			chromedp.Run(ctx, getBanksActionFunc, getDetailActionFunc)
			chromedp.Navigate(officeURL).Do(ctx)
			chromedp.WaitVisible(`#footer`).Do(ctx)
		}
		return nil
	})

	// crawling開始
	err := chromedp.Run(ctx,
		loginActionFunc,
		crawlingActionFunc,
	)
	if err != nil {
		return err
	}

	return nil
}

func scrapingOfBanks(d string) error {

	readerCurContents := strings.NewReader(d)
	contentsDom, err := goquery.NewDocumentFromReader(readerCurContents)
	if err != nil {
		return err
	}

	// lastCommit = strings.Replace(lastCommit, "最終同期日", "", -1)

	contentsDom.Find(`div.walletable___StyledDiv-sc-3etvmj-0`).Each(func(i int, v *goquery.Selection) {
		strAmount := v.Find("div.walletable___StyledDiv8-sc-3etvmj-8").Text()
		strAmount = strings.Replace(strAmount, ",", "", -1)
		amount, err := strconv.ParseInt(strAmount, 10, 64)
		if err != nil {
			log.Printf("intへのconvertに失敗しました:　%v", err)
		}

		id := uuid.New()
		Banks = append(Banks, &Bank{
			Id:         id.String(),
			OfficeName: officeName,
			// lastCommit: lastCommit,
			BankName: v.Find("a.walletable___StyledA-sc-3etvmj-3").Text(),
			Amount:   amount,
			Kind:     v.Parent().Parent().Find("h2.vb-sectionTitle").Text(),
		})
	})

	return err
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
	var val string
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
			if val == "" {
				val, _ = v.Find("input.checkbox").Attr("value")
				Users = append(Users, &User{OfficeName: officeName, LastId: val})
			}

			lastId, err := GetLastId(userId + "_" + officeName)
			if err != nil {
				log.Printf("最終更新日の取得に失敗しました:　%v", err)
			}

			preLastId, _ := v.Find("input.checkbox").Attr("value")

			if preLastId == lastId {
				cancel = false
				return false
			}
			strP := v.Find("td.number").Eq(0).Text()
			strW := "-" + v.Find("td.number").Eq(1).Text()
			strB := v.Find("td.number").Eq(2).Text()
			strP = strings.Replace(strP, ",", "", -1)
			strW = strings.Replace(strW, ",", "", -1)
			strB = strings.Replace(strB, ",", "", -1)
			payment, err := strconv.ParseInt(strP, 10, 64)
			if err != nil {
				log.Printf("intへのconvertに失敗しました:　%v", err)
			}
			withdrawal, err := strconv.ParseInt(strW, 10, 64)
			if err != nil {
				log.Printf("intへのconvertに失敗しました:　%v", err)
			}
			balance, err := strconv.ParseInt(strB, 10, 64)
			if err != nil {
				log.Printf("intへのconvertに失敗しました:　%v", err)
			}

			// laioutについてはgolangが指定している
			layout := "2006/01/02"
			strTradingDate := strings.Replace(v.Find("td.date-cell").Eq(0).Text(), "-", "/", -1)
			strUpdatedDate := strings.Replace(v.Find("td.date-cell").Eq(1).Text(), "-", "/", -1)
			strGettingDate := strings.Replace(v.Find("td.date-cell").Eq(2).Text(), "-", "/", -1)

			tradingDate, err := time.Parse(layout, strTradingDate)
			if err != nil {
				log.Printf("dateへのconvertに失敗しました:　%v", err)
			}
			updatedDate, err := time.Parse(layout, strUpdatedDate)
			if err != nil {
				log.Printf("dateへのconvertに失敗しました:　%v", err)
			}
			gettingDate, err := time.Parse(layout, strGettingDate)
			if err != nil {
				log.Printf("dateへのconvertに失敗しました:　%v", err)
			}

			var bankId string
			for _, bank := range bankNameAndId {
				if bank["bankName"] == v.Find("td.walletable-name").Text() {
					bankId = bank["bankId"]
					break
				}
			}

			id := uuid.New()
			Details = append(Details, &Detail{
				Id:             id.String(),
				BankId:         bankId,
				OfficeName:     officeName,
				BankName:       v.Find("td.walletable-name").Text(),
				TradingDate:    tradingDate,
				TradingContent: v.Find("td.description").Text(),
				Payment:        payment,
				Withdrawal:     withdrawal,
				Balance:        balance,
				UpdatedDate:    updatedDate,
				GettingDate:    gettingDate,
			})
			return true
		})
	}

	return nil
}
