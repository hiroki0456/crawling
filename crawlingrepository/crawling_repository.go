package crawlingrepository

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/sclevine/agouti"
	"upsider.crawling/crawlingproto"
)

type CrawlingSite struct {
	Pass  string
	Input *crawlingproto.UserInput
}

type Freee struct {
	CrawlingSite
}

type Bunk struct {
	Id     string `spanner:"Id"`
	UserId string `spanner:"UserId"`
	BunkId string `spanner:"BunkName"`
	Amount int64  `spanner:"Amount"`
	Kind   string `spanner:"Kind"`
}

type Detail struct {
	Id             string `spanner:"Id"`
	UserId         string `spanner:"UserId"`
	BunkName       string `spanner:"BunkName"`
	TradingDate    string `spanner:"TradingDate"`
	TradingContent string `spanner:"TradingContent"`
	Payment        int64  `spanner:"Payment"`
	Withdrawal     int64  `spanner:"Withdrawal"`
	Balance        int64  `spanner:"Balance"`
	UpdatedDate    string `spanner:"UpdatedDate"`
	GettingDate    string `spanner:"GettingDate"`
}

func (f *Freee) Crawling() ([]*Bunk, []*Detail, error) {
	driver := agouti.ChromeDriver()
	defer driver.Stop()

	if err := driver.Start(); err != nil {
		log.Printf("Failed to start driver: %v", err)
	}

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Printf("Failed to open page: %v", err)
	}

	// URLをDBに保存
	err = page.Navigate("https://accounts.secure.freee.co.jp/login/accounting?a=false&e=0&o=true&_gl=1*y6dhsy*_ga*ODkyNTExMDEzLjE2MzA5ODk2OTc.*_ga_9998VV0FMT*MTYzMTM2NDIwMS41LjAuMTYzMTM2NDIwMS42MA..")
	if err != nil {
		log.Printf("Failed to navigate: %v", err)
	}

	page.Find("#user_email").Fill(f.CrawlingSite.Input.GetUserID())
	page.FindByXPath("/html/body/div[3]/div/div[1]/form/div/div[3]/input").Fill(f.CrawlingSite.Pass)
	page.FindByXPath("/html/body/div[3]/div/div[1]/form/div/div[5]/input").Click()

	dom, err := page.HTML()
	if err != nil {
		return nil, nil, err
	}

	bunks, err := bunkAndSum(dom)
	if err != nil {
		return nil, nil, err
	}

	err = page.Navigate("https://secure.freee.co.jp/wallet_txns")
	if err != nil {
		return nil, nil, err
	}

	page.SetImplicitWait(10)

	pagenationDOM := page.FindByXPath("/html/body/div[2]/div/div/div[3]/div/ul").All("li")
	count, _ := pagenationDOM.Count()
	count -= 2
	detailList := []*Detail{}
	for i := 0; i < count; i++ {
		pagenationDOM.At(i + 1).Find("a").Click()
		page.SetImplicitWait(10)
		dom, _ = page.HTML()
		detail, err := details(dom)
		if err != nil {
			return nil, nil, err
		}
		detailList = append(detailList, detail...)
	}
	// dom, err = page.HTML()
	// if err != nil {
	// 	return nil, err
	// }

	return bunks, detailList, nil

	// fmt.Println(sumAndBunk)
}

func bunkAndSum(dom string) ([]*Bunk, error) {
	readerCurContents := strings.NewReader(dom)
	contentsDom, err := goquery.NewDocumentFromReader(readerCurContents)
	if err != nil {
		return nil, err
	}

	bunks := []*Bunk{}
	contentsDom.Find("section.vb-contentsBase").Each(func(i int, v1 *goquery.Selection) {
		v1.Find("div.walletable___StyledDiv-sc-3etvmj-0").Each(func(i int, v2 *goquery.Selection) {
			strAmount := v2.Find("div.walletable___StyledDiv8-sc-3etvmj-8").Text()
			strAmount = strings.Replace(strAmount, ",", "", -1)
			amount, _ := strconv.ParseInt(strAmount, 10, 64)
			id, _ := uuid.NewRandom()
			bunks = append(bunks, &Bunk{
				Id:     id.String(),
				BunkId: v2.Find("a.walletable___StyledA-sc-3etvmj-3").Text(),
				Amount: amount,
				Kind:   v1.Find("h2.vb-sectionTitle").Text(),
			})
		})
	})

	return bunks, nil
}

func details(dom string) ([]*Detail, error) {
	readerCurContents := strings.NewReader(dom)
	contentsDom, err := goquery.NewDocumentFromReader(readerCurContents)
	if err != nil {
		return nil, err
	}

	details := []*Detail{}
	contentsDom.Find("tr.line").Each(func(i int, v *goquery.Selection) {
		strP := v.Find("td.number").Eq(0).Text()
		strW := v.Find("td.number").Eq(1).Text()
		strB := v.Find("td.number").Eq(2).Text()
		strP = strings.Replace(strP, ",", "", -1)
		strW = strings.Replace(strW, ",", "", -1)
		strB = strings.Replace(strB, ",", "", -1)
		payment, _ := strconv.ParseInt(strP, 10, 64)
		withdrawal, _ := strconv.ParseInt(strW, 10, 64)
		balance, _ := strconv.ParseInt(strB, 10, 64)

		id, _ := uuid.NewRandom()
		details = append(details, &Detail{
			Id:             id.String(),
			BunkName:       v.Find("td.walletable-name").Text(),
			TradingDate:    v.Find("td.date-cell").Eq(0).Text(),
			TradingContent: v.Find("td.description").Text(),
			Payment:        payment,
			Withdrawal:     withdrawal,
			Balance:        balance,
			UpdatedDate:    v.Find("td.date-cell").Eq(1).Text(),
			GettingDate:    v.Find("td.date-cell").Eq(2).Text(),
		})

	})

	return details, nil
}
