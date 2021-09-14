package crawlinginterface

// import (
// 	"upsider.crawling/crawlingproto"
// 	"upsider.crawling/crawlingrepository"
// )

// type Crawling interface {
// 	Crawling(siteKind int32, pass string, input *crawlingproto.UserInput) (dom string, err error)
// }

// type crawling struct {
// 	cr crawlingrepository.CrawlingRepository
// }

// func NewCrawling(cr crawlingrepository.CrawlingRepository) Crawling {
// 	return &crawling{cr}
// }

// func (c *crawling) Crawling(siteKind int32, pass string, input *crawlingproto.UserInput) (dom string, err error) {
// 	if siteKind == 1 {
// 		dom, err = c.cr.FreeeCrawling(pass, input)
// 		if err != nil {
// 			return "", err
// 		}
// 	} else if siteKind == 2 {
// 		dom, err = c.cr.FreeeCrawling(pass, input)
// 		if err != nil {
// 			return "", err
// 		}
// 	}
// 	return dom, nil
// }

type CrawlingInterface interface {
	Crawling() (dom string, err error)
}

type crawlingRepository struct {
	rp CrawlingInterface
}

func NewCrawling(rp CrawlingInterface) *crawlingRepository {
	crawlinginterface := crawlingRepository{
		rp: rp,
	}
	return &crawlinginterface
}

func (c *crawlingRepository) Exec() {
	c.rp.Crawling()
}
