package crawlingrepository

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

func DistinctBankIdAndBankNameStmt() string {
	s := `SELECT
				DISTINCT bankId,
				bankName
			FROM
				Banks
			WHERE
				UserId = ? and OfficeName = ?;`
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

func DistinctCardIdAndCardNameStmt() string {
	s := `SELECT
				DISTINCT CardId,
				CardName
			FROM
				Cards
			WHERE
				UserId = ? and OfficeName = ?`
	return s
}
