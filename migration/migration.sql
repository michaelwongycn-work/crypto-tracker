CREATE DATABASE cryptoTracker;
USE cryptoTracker;

CREATE TABLE users (
    ID INTEGER PRIMARY KEY,
    email TEXT UNIQUE,
    password TEXT
);

CREATE TABLE user_assets (
    ID INTEGER PRIMARY KEY,
    userId INTEGER,
    assetId INTEGER,
    FOREIGN KEY (userId) REFERENCES users(ID),
    CONSTRAINT unique_user_crypto UNIQUE (userId, assetId)
);

CREATE TABLE user_tokens (
	userId INTEGER PRIMARY KEY,
    accessToken TEXT,
	refreshToken TEXT,
	expirationTime INTEGER
);