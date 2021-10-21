package config

const (
	UsersTable string = `
    CREATE TABLE Users (
        Id STRING(36) not null,
        UserIdOfficeName STRING(MAX) not null,
        UserId STRING(MAX) not null,
        OfficeName STRING(MAX) not null,
        LastId STRING(MAX) not null,
        updatedAt TIMESTAMP not null,
    ) PRIMARY KEY (UserIdOfficeName);
    `

	BanksTable string = `
    CREATE TABLE Banks (
        Id STRING(36) not null,
        UserId STRING(MAX) not null,
        BankId STRING(MAX) not null,
        LastCommit STRING(MAX),
        OfficeName STRING(MAX) not null,
        BankName STRING(MAX) not null,
        Amount int64,
        updatedAt TIMESTAMP not null,
    ) PRIMARY KEY (Id);
    `

	CardsTable string = `
    CREATE TABLE Cards (
        Id STRING(36) not null,
        UserId STRING(MAX) not null,
        CardId STRING(MAX) not null,
        LastCommit STRING(MAX),
        OfficeName STRING(MAX) not null,
        CardName STRING(MAX) not null,
        Amount int64,
        updatedAt TIMESTAMP not null,
    ) PRIMARY KEY (Id);
    `

	OthersTable string = `
    CREATE TABLE Others (
        Id STRING(36) not null,
        UserId STRING(MAX) not null,
        OtherId STRING(MAX) not null,
        LastCommit STRING(MAX),
        OfficeName STRING(MAX) not null,
        OtherName STRING(MAX) not null,
        Amount int64,
        updatedAt TIMESTAMP not null,
    ) PRIMARY KEY (Id);
    `

	DetailsTable = `
    CREATE TABLE Details (
        Id STRING(36) not null,
        UserId STRING(MAX) not null,
        BankId STRING(MAX),
        OfficeName STRING(MAX) not null,
        BankName STRING(MAX),
        TradingDate TIMESTAMP,
        TradingContent STRING(MAX),
        Payment int64,
        Withdrawal int64,
        Balance int64,
        UpdatedDate TIMESTAMP,
        GettingDate TIMESTAMP,
        crawling TIMESTAMP
    ) PRIMARY KEY (Id);
    `
)
