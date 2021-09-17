package crawlingrepository

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
)

type DB interface {
	UserCreate(userId string, updatedAt string) error
	BunkCreate(userId string, bunks []*Bunk) error
	DetailCreate(userId string, details []*Detail) error
}

type db struct {
	Client *spanner.Client
}

type User struct {
	Id        string `spanner:"Id"`
	UserId    string `spanner:"UserId"`
	UpdatedAt string `spanner:"updatedAt"`
}

func NewDatabase() DB {
	ctx := context.Background()
	client, _ := spanner.NewClient(ctx, "projects/test-project/instances/test-instance/databases/test-database")
	return &db{
		Client: client,
	}
}

func (d *db) UserCreate(userId string, updatedAt string) (err error) {
	id, _ := uuid.NewRandom()
	u := User{Id: id.String(), UserId: userId, UpdatedAt: updatedAt}
	m, err := spanner.InsertStruct("Users", u)
	if err != nil {
		return err
	}
	_, err = d.Client.Apply(context.Background(), []*spanner.Mutation{m})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (d *db) BunkCreate(userId string, bunks []*Bunk) (err error) {
	for _, v := range bunks {
		v.UserId = userId
		m, err := spanner.InsertOrUpdateStruct("Bunks", v)
		if err != nil {
			return err
		}
		_, err = d.Client.Apply(context.Background(), []*spanner.Mutation{m})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return err
}

func (d *db) DetailCreate(userId string, details []*Detail) (err error) {
	for _, v := range details {
		v.UserId = userId
		m, err := spanner.InsertOrUpdateStruct("Details", v)
		if err != nil {
			return err
		}
		_, err = d.Client.Apply(context.Background(), []*spanner.Mutation{m})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return err
}
