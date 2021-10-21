package crawlingrepository

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB interface {
	UserCreate(users []*User, userId string, updatedAt *time.Time) error
	BankCreate(userId string, banks []*Bank, today *time.Time) <-chan error
	DetailCreate(userId string, details []*Detail, today *time.Time) <-chan error
}

type db struct {
	Client *sql.DB
}

func NewDatabase() DB {
	client, err := sql.Open("mysql", "root@/freee")
	if err != nil {
		log.Fatal(err)
	}
	return &db{
		Client: client,
	}
}

func (d *db) UserCreate(users []*User, userId string, updatedAt *time.Time) (err error) {
	for _, v := range users {
		v.UserId = userId
		v.UserIdOfficeName = v.UserId + "_" + v.OfficeName

		updateStmt, err := d.Client.Prepare("UPDATE Users set userIdOfficeName = ?,userId = ?,officeName = ?,lastId = ?,updatedAt = ? where officeName = ?")
		if err != nil {
			return err
		}
		result, err := updateStmt.Exec(v.UserIdOfficeName, v.UserId, v.OfficeName, v.LastId, updatedAt, v.OfficeName)
		if err != nil {
			return err
		}

		rowsAffect, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffect == 0 {
			insertStmt, err := d.Client.Prepare("INSERT INTO Users(userIdOfficeName,userId,officeName,lastId,updatedAt) VALUES(?, ?, ?, ?, ?)")
			if err != nil {
				return err
			}
			_, err = insertStmt.Exec(v.UserIdOfficeName, v.UserId, v.OfficeName, v.LastId, updatedAt)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (d *db) BankCreate(userId string, banks []*Bank, today *time.Time) <-chan error {
	errCh := make(chan error)
	for _, v := range banks {
		if v.Kind == "銀行口座" {

			insertStmt, err := d.Client.Prepare("INSERT INTO Banks(userId,bankId,LastCommitDate,officeName,bankName,amount,updatedAt) VALUES(?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				errCh <- err
			}
			_, err = insertStmt.Exec(userId, v.BankId, v.LastCommit, v.OfficeName, v.BankName, v.Amount, today)
			if err != nil {
				errCh <- err
			}

		} else if v.Kind == "クレジットカード" {
			insertStmt, err := d.Client.Prepare("INSERT INTO Cards(userId,cardId,LastCommitDate,officeName,cardName,amount,updatedAt) VALUES(?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				errCh <- err
			}
			_, err = insertStmt.Exec(userId, v.BankId, v.LastCommit, v.OfficeName, v.BankName, v.Amount, today)
			if err != nil {
				errCh <- err
			}
		} else {
			insertStmt, err := d.Client.Prepare("INSERT INTO Others(userId,otherId,LastCommitDate,officeName,otherName,amount,updatedAt) VALUES(?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				errCh <- err
			}
			_, err = insertStmt.Exec(userId, v.BankId, v.LastCommit, v.OfficeName, v.BankName, v.Amount, today)
			if err != nil {
				errCh <- err
			}
		}
	}

	go func() {
		// wg.Wait()

		close(errCh)
	}()

	return errCh
}

func (d *db) DetailCreate(userId string, details []*Detail, today *time.Time) <-chan error {
	errCh := make(chan error)

	for _, v := range details {
		insertStmt, err := d.Client.Prepare("INSERT INTO Details(userId, bankId, officeName, bankName, tradingDate, tradingContent, payment, withdrawal, balance, UpdatedDate, GettingDate, crawling) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			errCh <- err
		}
		_, err = insertStmt.Exec(userId, v.BankId, v.OfficeName, v.BankName, v.TradingDate, v.TradingContent, v.Payment, v.Withdrawal, v.Balance, v.UpdatedDate, v.GettingDate, today)
		if err != nil {
			errCh <- err
		}
	}
	go func() {
		// wg.Wait()
		close(errCh)
	}()
	return errCh
}
