package crawlingrepository

import (
	"strconv"
)

func OfficeReadStmt() string {
	s := `SELECT
			OfficeName, updatedAt
		FROM
			Users
		WHERE
			UserId = ?
	`
	return s
}

func DistinctBankNameCountStmt() string {
	s := `SELECT
			count(distinct BankName)
				FROM
					banks
				WHERE
					UserId = ? and OfficeName = ?`
	return s
}

func SumAmountOnBanksStmt(bankCount int64) string {
	s := `SELECT
					sum(amount)
				FROM
					(SELECT
							amount
						FROM
							banks
						WHERE
							UserId = ? and OfficeName = ?
						ORDER BY
							updatedAt desc
						LIMIT
							` + strconv.Itoa(int(bankCount)) + `) as bank`

	return s
}

func BankNameAndBankAmountStmt(bankCount int64) string {
	s := `SELECT
			BankId, BankName, Amount
				FROM
					(SELECT
							*
						FROM
							banks
						WHERE
							UserId = ? and OfficeName = ?
						ORDER BY
							updatedAt desc
						LIMIT
							` + strconv.Itoa(int(bankCount)) + `) as bank`
	return s
}

func DetailStmt(startDay string, lastDay string) string {
	if startDay == "" {
		startDay = "2010-01-01"
	}
	if lastDay == "" {
		lastDay = "2030-12-31"
	}
	s := `SELECT
			BankName,
			TradingContent,
			(Payment + Withdrawal) as Payment,
			Balance,
			TradingDate,
			GettingDate,
			UpdatedDate
				FROM
					details
				WHERE
					UserId = ? and bankId = ? and officeName = ? and TradingDate > '` + startDay + `' and TradingDate < '` + lastDay + `'`
	return s
}

func CardCountStmt() string {
	s := `SELECT
			count(distinct CardName)
				FROM
					cards
				WHERE
					UserId = ? and OfficeName = ?`
	return s
}

func CardSumStmt(cardCount int64) string {
	s := `SELECT
			sum(amount)
				FROM
					(SELECT
							amount
						FROM
							cards
						WHERE
							UserId = ? and OfficeName = ?
						ORDER BY
							updatedAt desc
						LIMIT
							` + strconv.Itoa(int(cardCount)) +
		`) as card`
	return s
}

func CardInfoStmt(cardCount int64) string {
	s := `SELECT
			CardId, CardName, Amount
				FROM
					(SELECT
							*
						FROM
							cards
						WHERE
							UserId = ? and OfficeName = ?
						ORDER BY
							updatedAt desc
						LIMIT
							` + strconv.Itoa(int(cardCount)) + `) as card`

	return s
}
