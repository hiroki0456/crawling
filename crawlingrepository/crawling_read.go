package crawlingrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	pb "upsider.crawling/crawlingproto"
)

type CrawlingReadInterface interface {
	OfficeRead(ctx context.Context, req *pb.FreeeRequest) (offices []*pb.Office, err error)
	BankRead(ctx context.Context, req *pb.FreeeRequest, officeName string, startDay string, lastDay string) (err error)
	CardRead(ctx context.Context, req *pb.FreeeRequest, officeName string, startDay string, lastDay string) (err error)
	detailRead(ctx context.Context, userId string, bankId string, officeName string, startDay string, lastDay string) (details []*pb.Detail, err error)
}

type CrawlingRead struct {
	client *sql.DB
}

type Date struct {
	detailDate time.Time
	createdAt  time.Time
	updatedAt  time.Time
}

var PbBanks *pb.Banks
var PbCards *pb.Cards

func NewCrawlingRead(db *sql.DB) CrawlingReadInterface {
	return &CrawlingRead{db}
}

func (c *CrawlingRead) OfficeRead(ctx context.Context, req *pb.FreeeRequest) (offices []*pb.Office, err error) {
	rows, err := c.client.Query(OfficeReadStmt(), req.UserInput.UserId)
	if err != nil {
		return nil, fmt.Errorf("事業所名のクエリ取得に失敗しました: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		office := &pb.Office{}
		date := &Date{}
		err := rows.Scan(&office.OfficeName, &date.updatedAt)
		if err != nil {
			return nil, fmt.Errorf("事業所名の取得に失敗しました: %s", err)
		}
		office.Crawling = timestamppb.New(date.updatedAt)

		offices = append(offices, office)
	}

	err = rows.Err()
	if err != nil {
		panic(err.Error())
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

func (c *CrawlingRead) detailRead(ctx context.Context, userId string, bankId string, officeName string, startDay string, lastDay string) (details []*pb.Detail, err error) {
	rows, err := c.client.Query(DetailStmt(startDay, lastDay), userId, bankId, officeName)
	if err != nil {
		return nil, fmt.Errorf("明細のクエリ取得に失敗しました: %s", err)
	}
	details = []*pb.Detail{}

	for rows.Next() {
		detail := &pb.Detail{}
		date := &Date{}

		if err := rows.Scan(&detail.DetailName, &detail.Contents, &detail.Amount, &detail.Balance, &date.detailDate, &date.createdAt, &date.updatedAt); err != nil {
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
	rows, err := c.client.Query(DistinctBankNameCountStmt(), userId, officeName)
	if err != nil {
		return fmt.Errorf("銀行口座数のクエリ取得に失敗しました: %s", err)
	}

	for rows.Next() {
		err := rows.Scan(&PbBanks.BankCount)
		if err != nil {
			return fmt.Errorf("銀行口座数の取得に失敗しました: %s", err)
		}
	}

	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func getBankSum(c *CrawlingRead, ctx context.Context, userId string, officeName string) error {
	rows, err := c.client.Query(SumAmountOnBanksStmt(PbBanks.BankCount), userId, officeName)
	if err != nil {
		return fmt.Errorf("銀行口座残高のクエリ取得に失敗しました: %s", err)
	}

	for rows.Next() {
		err := rows.Scan(&PbBanks.BankSum)
		if err != nil {
			return fmt.Errorf("銀行口残高の取得に失敗しました: %s", err)
		}
	}

	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func getBankDetails(c *CrawlingRead, ctx context.Context, userId string, officeName string, startDay string, lastDay string) error {
	rows, err := c.client.Query(BankNameAndBankAmountStmt(PbBanks.BankCount), userId, officeName)
	if err != nil {
		return fmt.Errorf("各銀行名と各銀行の残高のクエリ取得に失敗しました: %s", err)
	}
	bankList := []*pb.Bank{}

	for rows.Next() {
		bank := &pb.Bank{}
		if err := rows.Scan(&bank.BankId, &bank.BankName, &bank.BankAmount); err != nil {
			return fmt.Errorf("各銀行名と各銀行の残高の取得に失敗しました: %s", err)
		}

		details, err := c.detailRead(ctx, userId, bank.BankId, officeName, startDay, lastDay)
		if err != nil {
			return fmt.Errorf("%sの明細取得に失敗しました: %s", bank.BankName, err)
		}
		bank.Detail = details
		bank.DetailCount = int64(len(details))
		bankList = append(bankList, bank)
	}
	PbBanks.Bank = bankList

	return nil
}

func getCardCount(c *CrawlingRead, ctx context.Context, userId string, officeName string) error {
	rows, err := c.client.Query(CardCountStmt(), userId, officeName)
	if err != nil {
		return fmt.Errorf("クレジットカード数のクエリ取得に失敗しました: %s", err)
	}

	for rows.Next() {
		if err := rows.Scan(&PbCards.CardCount); err != nil {
			return fmt.Errorf("クレジットカード数の取得に失敗しました: %s", err)
		}
	}
	return nil
}

func getCardSum(c *CrawlingRead, ctx context.Context, userId string, officeName string) error {
	rows, err := c.client.Query(CardSumStmt(PbCards.CardCount), userId, officeName)
	if err != nil {
		return fmt.Errorf("クレジットカード残高のクエリ取得に失敗しました: %s", err)
	}

	for rows.Next() {
		if err := rows.Scan(&PbCards.CardSum); err != nil {
			return fmt.Errorf("クレジットカード残高の取得に失敗しました: %s", err)
		}
	}
	return nil
}

func getCardDetails(c *CrawlingRead, ctx context.Context, userId string, officeName string, startDay string, lastDay string) error {
	rows, err := c.client.Query(CardInfoStmt(PbCards.CardCount), userId, officeName)
	if err != nil {
		return fmt.Errorf("各クレジットカード名と各クレジットカード残高のクエリ取得に失敗しました: %s", err)
	}

	cardList := []*pb.Card{}

	for rows.Next() {
		card := &pb.Card{}

		if err := rows.Scan(&card.CardId, &card.CardName, &card.CardAmount); err != nil {
			return fmt.Errorf("各クレジットカード名と各クレジットカードの残高の取得に失敗しました: %s", err)
		}

		details, err := c.detailRead(ctx, userId, card.CardId, officeName, startDay, lastDay)
		if err != nil {
			return fmt.Errorf("%sの明細取得に失敗しました: %s", card.CardName, err)
		}

		card.Detail = details
		card.DetailCount = int64(len(details))
		cardList = append(cardList, card)

	}
	PbCards.Card = cardList

	return nil
}

func GetLastId(userId string) (string, error) {
	client, err := sql.Open("mysql", "root@/freee")
	if err != nil {
		return "", err
	}
	defer client.Close()
	var lastId string
	client.QueryRow("SELECT lastId FROM Users where UserIdOfficeName = ?", userId).Scan(&lastId)

	return lastId, nil
}
