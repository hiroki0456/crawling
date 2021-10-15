package crawlingrepository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
)

type DB interface {
	UserCreate(users []*User, userId string, updatedAt *time.Time) error
	BankCreate(userId string, banks []*Bank, today *time.Time) <-chan error
	DetailCreate(userId string, details []*Detail) <-chan error
}

type db struct {
	Client *spanner.Client
}

type CreateBank struct {
	Id         string     `spanner:"Id"`
	UserId     string     `spanner:"UserId"`
	BankId     string     `spanner:"BankId"`
	lastCommit string     `spanner:"LastCommit"`
	OfficeName string     `spanner:"OfficeName"`
	BankName   string     `spanner:"BankName"`
	Amount     int64      `spanner:"Amount"`
	UpdatedAt  *time.Time `spanner:"updatedAt"`
}

type CreateCard struct {
	Id         string     `spanner:"Id"`
	UserId     string     `spanner:"UserId"`
	CardId     string     `spanner:"CardId"`
	lastCommit string     `spanner:"LastCommit"`
	OfficeName string     `spanner:"OfficeName"`
	CardName   string     `spanner:"CardName"`
	Amount     int64      `spanner:"Amount"`
	UpdatedAt  *time.Time `spanner:"updatedAt"`
}

type Other struct {
	Id         string     `spanner:"Id"`
	UserId     string     `spanner:"UserId"`
	OtherId    string     `spanner:"OtherId"`
	lastCommit string     `spanner:"LastCommit"`
	OfficeName string     `spanner:"OfficeName"`
	OtherName  string     `spanner:"OtherName"`
	Amount     int64      `spanner:"Amount"`
	UpdatedAt  *time.Time `spanner:"updatedAt"`
}

func NewDatabase() DB {
	ctx := context.Background()
	client, _ := spanner.NewClient(ctx, "projects/test-project/instances/test-instance/databases/test-database")
	return &db{
		Client: client,
	}
}

func (d *db) UserCreate(users []*User, userId string, updatedAt *time.Time) (err error) {
	for _, v := range users {
		id, _ := uuid.NewRandom()
		v.Id = id.String()
		v.UserId = userId
		v.UpdatedAt = updatedAt
		v.UserIdOfficeName = v.UserId + "_" + v.OfficeName

		m, err := spanner.InsertOrUpdateStruct("Users", v)
		if err != nil {
			return err
		}
		_, err = d.Client.Apply(context.Background(), []*spanner.Mutation{m})
		if err != nil {
			return err
		}
	}

	return err
}

func (d *db) BankCreate(userId string, banks []*Bank, today *time.Time) <-chan error {
	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(banks))
	for _, v := range banks {
		go func(v *Bank) {
			defer wg.Done()

			if v.Kind == "銀行口座" {

				v.UserId = userId
				m, err := spanner.InsertStruct("Banks", &CreateBank{
					Id:         v.Id,
					UserId:     v.UserId,
					BankId:     v.BankId,
					lastCommit: v.lastCommit,
					OfficeName: v.OfficeName,
					BankName:   v.BankName,
					Amount:     v.Amount,
					UpdatedAt:  today,
				})
				if err != nil {
					errCh <- err
				}
				_, err = d.Client.Apply(context.Background(), []*spanner.Mutation{m})
				if err != nil {
					fmt.Println(err)
					errCh <- err
				}
			} else if v.Kind == "クレジットカード" {
				v.UserId = userId
				m, err := spanner.InsertStruct("Cards", &CreateCard{
					Id:         v.Id,
					UserId:     v.UserId,
					CardId:     v.BankId,
					lastCommit: v.lastCommit,
					OfficeName: v.OfficeName,
					CardName:   v.BankName,
					Amount:     v.Amount,
					UpdatedAt:  today,
				})
				if err != nil {
					errCh <- err
				}
				_, err = d.Client.Apply(context.Background(), []*spanner.Mutation{m})
				if err != nil {
					fmt.Println(err)
					errCh <- err
				}
			} else {
				v.UserId = userId
				m, err := spanner.InsertStruct("Others", &Other{
					Id:         v.Id,
					UserId:     v.UserId,
					OtherId:    v.BankId,
					lastCommit: v.lastCommit,
					OfficeName: v.OfficeName,
					OtherName:  v.BankName,
					Amount:     v.Amount,
					UpdatedAt:  today,
				})
				if err != nil {
					errCh <- err
				}
				_, err = d.Client.Apply(context.Background(), []*spanner.Mutation{m})
				if err != nil {
					fmt.Println(err)
					errCh <- err
				}
			}

		}(v)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	return errCh
}

func (d *db) DetailCreate(userId string, details []*Detail) <-chan error {
	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(details))
	for _, v := range details {
		go func(v *Detail) {
			defer wg.Done()
			v.UserId = userId
			m, err := spanner.InsertStruct("Details", v)
			if err != nil {
				errCh <- err
			}
			_, err = d.Client.Apply(context.Background(), []*spanner.Mutation{m})
			if err != nil {
				errCh <- err
			}
		}(v)
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()
	return errCh
}
