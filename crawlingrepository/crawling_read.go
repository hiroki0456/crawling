package crawlingrepository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "upsider.crawling/crawlingproto"
)

type CrawlingReadInterface interface {
	OfficeRead(ctx context.Context, req *pb.FreeeRequest) (offices []*pb.Office, err error)
	BankRead(ctx context.Context, req *pb.FreeeRequest, officeName string, startDay string, lastDay string) (err error)
	CardRead(ctx context.Context, req *pb.FreeeRequest, officeName string, startDay string, lastDay string) (err error)
	detailRead(ctx context.Context, userId string, bankId string, startDay string, lastDay string) (details []*pb.Detail, err error)
}

type CrawlingRead struct {
	client *spanner.Client
}

type Date struct {
	detailDate time.Time
	createdAt  time.Time
	updatedAt  time.Time
}

var PbBanks *pb.Banks
var PbCards *pb.Cards

func NewCrawlingRead(db *spanner.Client) CrawlingReadInterface {
	return &CrawlingRead{db}
}

func (c *CrawlingRead) OfficeRead(ctx context.Context, req *pb.FreeeRequest) (offices []*pb.Office, err error) {
	ro := c.client.ReadOnlyTransaction()
	defer ro.Close()
	stmt := OfficeReadStmt(req.UserInput.UserId)
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("銀行口座数と残高のクエリ取得に失敗しました: %s", err)
		}
		office := &pb.Office{}
		date := &Date{}
		if err := row.Columns(&office.OfficeName, &date.updatedAt); err != nil {
			return nil, fmt.Errorf("銀行口座数と残高の取得に失敗しました: %s", err)
		}
		office.Crawling = timestamppb.New(date.updatedAt)

		offices = append(offices, office)
	}

	return offices, nil
}

func (c *CrawlingRead) BankRead(ctx context.Context, req *pb.FreeeRequest, officeName string, startDay string, lastDay string) (err error) {
	PbBanks = &pb.Banks{}
	err = getBankCount(c, ctx, req.UserInput.UserId, officeName)
	if err != nil {
		return err
	}

	if PbBanks.BankCount == 0 {
		return nil
	}

	err = getBankSum(c, ctx, req.UserInput.UserId, officeName)
	if err != nil {
		return err
	}

	err = getBankDetails(c, ctx, req.UserInput.UserId, officeName, startDay, lastDay)
	if err != nil {
		return err
	}

	return nil
}

func (c *CrawlingRead) CardRead(ctx context.Context, req *pb.FreeeRequest, officeName string, startDay string, lastDay string) (err error) {
	PbCards = &pb.Cards{}

	err = getCardCount(c, ctx, req.UserInput.UserId, officeName)
	if err != nil {
		return err
	}

	if PbCards.CardCount == 0 {
		return nil
	}

	err = getCardSum(c, ctx, req.UserInput.UserId, officeName)
	if err != nil {
		return err
	}

	err = getCardDetails(c, ctx, req.UserInput.UserId, officeName, startDay, lastDay)
	if err != nil {
		return err
	}

	return nil
}

func (c *CrawlingRead) detailRead(ctx context.Context, userId string, bankId string, startDay string, lastDay string) (details []*pb.Detail, err error) {
	ro := c.client.ReadOnlyTransaction()
	defer ro.Close()

	details = []*pb.Detail{}

	stmt := DetailStmt(userId, bankId, startDay, lastDay)
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	for {
		detail := &pb.Detail{}
		date := &Date{}
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("明細のクエリ取得に失敗しました: %s", err)
		}

		if err := row.Columns(&detail.DetailName, &detail.Contents, &detail.Amount, &detail.Balance, &date.detailDate, &date.createdAt, &date.updatedAt); err != nil {
			return nil, fmt.Errorf("明細の取得に失敗しました: %s", err)
		}
		detail.DetailDate = timestamppb.New(date.detailDate)
		detail.CreatedAt = timestamppb.New(date.createdAt)
		detail.UpdatedAt = timestamppb.New(date.updatedAt)

		details = append(details, detail)
	}

	return details, nil
}

func getBankCount(c *CrawlingRead, ctx context.Context, userId string, officeName string) error {
	ro := c.client.ReadOnlyTransaction()
	defer ro.Close()
	stmt := DistinctBankNameCountStmt(userId, officeName)
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("銀行口座数のクエリ取得に失敗しました: %s", err)
		}

		if err := row.Columns(&PbBanks.BankCount); err != nil {
			return fmt.Errorf("銀行口座数の取得に失敗しました: %s", err)
		}
	}
	return nil
}

func getBankSum(c *CrawlingRead, ctx context.Context, userId string, officeName string) error {
	ro := c.client.ReadOnlyTransaction()
	defer ro.Close()
	stmt := SumAmountOnBanksStmt(PbBanks.BankCount, officeName, userId)
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("銀行口座残高のクエリ取得に失敗しました: %s", err)
		}
		if err := row.Columns(&PbBanks.BankSum); err != nil {
			return fmt.Errorf("銀行口残高の取得に失敗しました: %s", err)
		}
	}
	return nil
}

func getBankDetails(c *CrawlingRead, ctx context.Context, userId string, officeName string, startDay string, lastDay string) error {
	ro := c.client.ReadOnlyTransaction()
	defer ro.Close()
	stmt := BankNameAndBankAmountStmt(PbBanks.BankCount, officeName, userId)
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	bankList := []*pb.Bank{}

	for {
		bank := &pb.Bank{}
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("各銀行名と各銀行の残高のクエリ取得に失敗しました: %s", err)
		}

		if err := row.Columns(&bank.BankId, &bank.BankName, &bank.BankAmount); err != nil {
			return fmt.Errorf("各銀行名と各銀行の残高の取得に失敗しました: %s", err)
		}

		details, err := c.detailRead(ctx, userId, bank.BankId, startDay, lastDay)
		if err != nil {
			return fmt.Errorf("%sの明細取得に失敗しました: %s", bank.BankName, err)
		}
		bank.Detail = details
		bankList = append(bankList, bank)
		PbBanks.DetailCount = int64(len(details))
	}
	PbBanks.Bank = bankList

	return nil
}

func getCardCount(c *CrawlingRead, ctx context.Context, userId string, officeName string) error {
	ro := c.client.ReadOnlyTransaction()
	defer ro.Close()
	stmt := CardCountStmt(userId, officeName)
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("クレジットカード数のクエリ取得に失敗しました: %s", err)
		}

		if err := row.Columns(&PbCards.CardCount); err != nil {
			return fmt.Errorf("クレジットカード数の取得に失敗しました: %s", err)
		}
	}
	return nil
}

func getCardSum(c *CrawlingRead, ctx context.Context, userId string, officeName string) error {
	ro := c.client.ReadOnlyTransaction()
	defer ro.Close()
	stmt := CardSumStmt(PbCards.CardCount, userId, officeName)
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("クレジットカード残高のクエリ取得に失敗しました: %s", err)
		}

		if err := row.Columns(&PbCards.CardSum); err != nil {
			return fmt.Errorf("クレジットカード残高の取得に失敗しました: %s", err)
		}
	}
	return nil
}

func getCardDetails(c *CrawlingRead, ctx context.Context, userId string, officeName string, startDay string, lastDay string) error {
	ro := c.client.ReadOnlyTransaction()
	defer ro.Close()
	stmt := CardInfoStmt(PbCards.CardCount, userId, officeName)
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	cardList := []*pb.Card{}

	for {
		card := &pb.Card{}
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("各クレジットカード名と各クレジットカード残高のクエリ取得に失敗しました: %s", err)
		}
		if err := row.Columns(&card.CardId, &card.CardName, &card.CardAmount); err != nil {
			return fmt.Errorf("各クレジットカード名と各クレジットカードの残高の取得に失敗しました: %s", err)
		}

		details, err := c.detailRead(ctx, userId, card.CardId, startDay, lastDay)
		if err != nil {
			return fmt.Errorf("%sの明細取得に失敗しました: %s", card.CardName, err)
		}

		card.Detail = details

		cardList = append(cardList, card)
		PbCards.DetailCount = int64(len(details))
	}
	PbCards.Card = cardList

	return nil
}

func GetLastId(userId string) (string, error) {
	ctx := context.Background()
	client, _ := spanner.NewClient(ctx, "projects/test-project/instances/test-instance/databases/test-database")
	ro := client.ReadOnlyTransaction()
	defer ro.Close()
	stmt := spanner.Statement{SQL: `select LastId from Users where UserIdOfficeName = @UserId`, Params: map[string]interface{}{"UserId": userId}}
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	var lastId string

	for {

		row, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return "", fmt.Errorf("明細のクエリ取得に失敗しました: %s", err)
		}

		if err := row.Columns(&lastId); err != nil {
			return "", fmt.Errorf("明細の取得に失敗しました: %s", err)
		}
	}
	return lastId, nil
}
