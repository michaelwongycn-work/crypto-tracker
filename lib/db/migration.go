package db

const (
	usersTable            = "users"
	usersTableSchema      = `CREATE TABLE users (ID INTEGER PRIMARY KEY, email TEXT UNIQUE, password TEXT)`
	userAssetsTable       = "user_assets"
	userAssetsTableSchema = `CREATE TABLE user_assets (ID INTEGER PRIMARY KEY, userId INTEGER, assetId INTEGER, FOREIGN KEY (userId) REFERENCES users(ID), CONSTRAINT unique_user_crypto UNIQUE (userId, assetId))`
	userTokensTable       = "user_tokens"
	userTokensTableSchema = `CREATE TABLE user_tokens (userId INTEGER PRIMARY KEY, accessToken TEXT, refreshToken TEXT, expirationTime INTEGER)`
)
