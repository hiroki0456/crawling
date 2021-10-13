package crawlingrepository

import (
	"strconv"

	"cloud.google.com/go/spanner"
)

func OfficeReadStmt(userId string) spanner.Statement {
	s := spanner.Statement{
		SQL: `SELECT
					OfficeName 
				FROM
					Banks
				WHERE
					UserId = @UserId
				GROUP BY
					OfficeName`,
		Params: map[string]interface{}{
			"UserId": userId,
		},
	}
	return s
}

func DistinctBankNameCountStmt(userId string, officeName string) spanner.Statement {
	s := spanner.Statement{
		SQL: `SELECT
					count(distinct BankName)
				FROM
					banks
				WHERE
					UserId = @UserId and OfficeName = @OfficeName`,
		Params: map[string]interface{}{
			"UserId":     userId,
			"OfficeName": officeName,
		},
	}
	return s
}

func SumAmountOnBanksStmt(bankCount int64, officeName string, userId string) spanner.Statement {
	s := spanner.Statement{
		SQL: `SELECT
					sum(amount)
				FROM
					(SELECT
							amount
						FROM
							banks
						WHERE
							UserId = @UserId and OfficeName = @OfficeName
						ORDER BY
							updatedAt desc
						LIMIT
							` + strconv.Itoa(int(bankCount)) + `)`,
		Params: map[string]interface{}{
			"UserId":     userId,
			"OfficeName": officeName,
		},
	}
	return s
}

func BankNameAndBankAmountStmt(bankCount int64, officeName string, userId string) spanner.Statement {
	s := spanner.Statement{

		SQL: `SELECT
					BankId, BankName, Amount
				FROM
					(SELECT
							*
						FROM
							banks
						WHERE
							UserId = @UserId and OfficeName = @OfficeName
						ORDER BY
							updatedAt desc
						LIMIT
							` + strconv.Itoa(int(bankCount)) + `)`,
		Params: map[string]interface{}{
			"UserId":     userId,
			"OfficeName": officeName,
		},
	}
	return s
}

func DetailStmt(userId string, bankId string) spanner.Statement {
	s := spanner.Statement{
		SQL: `SELECT
					BankName,
					TradingContent,
					(Payment + Withdrawal) as Payment,
					UpdatedDate,
					GettingDate
				FROM
					details
				WHERE
					UserId = @UserId and bankId = @BankId`,
		Params: map[string]interface{}{
			"UserId": userId,
			"BankId": bankId,
		},
	}
	return s
}

func CardCountStmt(userId string, officeName string) spanner.Statement {
	s := spanner.Statement{
		SQL: `SELECT
					count(distinct CardName)
				FROM
					cards
				WHERE
					UserId = @UserId and OfficeName = @OfficeName`,
		Params: map[string]interface{}{
			"UserId":     userId,
			"OfficeName": officeName,
		},
	}
	return s
}

func CardSumStmt(cardCount int64, userId string, officeName string) spanner.Statement {
	s := spanner.Statement{
		SQL: `SELECT
					sum(amount)
				FROM
					(SELECT
							amount
						FROM
							cards
						WHERE
							UserId = @UserId and OfficeName = @OfficeName
						ORDER BY
							updatedAt desc
						LIMIT
							` + strconv.Itoa(int(cardCount)) +
			`)`,
		Params: map[string]interface{}{
			"UserId":     userId,
			"OfficeName": officeName,
		},
	}
	return s
}

func CardInfoStmt(cardCount int64, userId string, officeName string) spanner.Statement {
	s := spanner.Statement{
		SQL: `SELECT
					CardId, CardName, Amount
				FROM
					(SELECT
							*
						FROM
							cards
						WHERE
							UserId = @UserId and OfficeName = @OfficeName
						ORDER BY
							updatedAt desc
						LIMIT
							` + strconv.Itoa(int(cardCount)) + `)`,
		Params: map[string]interface{}{
			"UserId":     userId,
			"OfficeName": officeName,
		},
	}
	return s
}