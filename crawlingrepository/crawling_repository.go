package crawlingrepository

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
	"github.com/sclevine/agouti"
	"upsider.crawling/crawlingproto"
)

// type CrawlingRepository interface {
// 	FreeeCrawling(pass string, input *crawlingproto.UserInput) (dom string, err error)
// }

// type crawlingRepository struct{}

// func NewCrawling() CrawlingRepository {
// 	return &crawlingRepository{}
// }
type sumAndBunk struct {
	sum      string
	bunkName []string
}

type CrawlingSite struct {
	Pass  string
	Input *crawlingproto.UserInput
}

type Freee struct {
	CrawlingSite
}

func (f *Freee) Crawling() (dom string, err error) {

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

	time.Sleep(10 * time.Second)

	dom, err = page.HTML()
	if err != nil {
		return "", err
	}

	sumAndBunk, err := bunkAndSum(dom)
	if err != nil {
		return "", err
	}

	fmt.Println(sumAndBunk)

	page.Navigate("https://secure.freee.co.jp/wallet_txns")
	time.Sleep(10 * time.Second)
	return dom, nil
}

func bunkAndSum(dom string) (*sumAndBunk, error) {
	readerCurContents := strings.NewReader(dom)
	contentsDom, err := goquery.NewDocumentFromReader(readerCurContents)
	if err != nil {
		return nil, err
	}
	sum := 0
	contentsDom.Find("div.walletable_group___StyledDiv6-sc-1uncx9n-5").Each(func(i int, v *goquery.Selection) {
		strs := strings.Split(v.Text(), ",")
		strNum := ""
		for _, v := range strs {
			strNum += v
		}
		intNum, err := strconv.Atoi(strNum)
		if err != nil {
			return
		}
		sum += intNum
	})
	commaNum := humanize.Comma(int64(sum))

	bunkNames := []string{}
	contentsDom.Find("a.walletable___StyledA-sc-3etvmj-3").Each(func(i int, v *goquery.Selection) {
		bunkNames = append(bunkNames, v.Text())
	})

	fmt.Println(bunkNames)

	return &sumAndBunk{sum: commaNum, bunkName: bunkNames}, nil
}

// func meisaiIndex(dom string) (*sumAndBunk, error) {
// 	readerCurContents := strings.NewReader(dom)
// 	contentsDom, err := goquery.NewDocumentFromReader(readerCurContents)
// 	if err != nil {
// 		return nil, err
// 	}

// 	contentsDom.Find("ul.fr-dropdown-menu").First().click()

// }
