# Crypto Tracker App

This is a Crypto Tracker App that allows users to track the price of cryptocurrencies to selected currency. The data is fetched from the Coincap API (https://docs.coincap.io/). The project is fully integrated with JWT (JSON Web Token) for user authentication and authorization.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Endpoints](#endpoints)

## Installation

```
git clone https://github.com/michaelwongycn-work/crypto-tracker.git
cd crypto-tracker
mv application_config_example.json application_config.json
go mod download
```

## Usage

```
go run main.go
```

The server will start running on http://localhost:2000.

## Endpoint

The following endpoints are available:

GET /ping
Check if the server is running.

POST /login
Login with email & password.

POST /register
Register with email, password, & password confirmation.

POST /logout
Logout the current user.

POST /refresh-token
Refresh the current user's token.

GET /crypto
Retrieve the user's cryptocurrency assets.

POST /crypto
Insert a new cryptocurrency asset for the user.

DELETE /crypto
Delete a cryptocurrency asset for the user.

All endpoints require authentication except for /login, /register, and /refresh-token.
