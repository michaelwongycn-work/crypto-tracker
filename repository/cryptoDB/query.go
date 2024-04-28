package cryptoDB

const (
	getUserByEmailAndPasswordQuery = "SELECT id FROM users WHERE email = ? AND password = ?"
	insertUserQuery                = "INSERT INTO users (email, password) VALUES (?, ?)"
	getUserTokenQuery              = "SELECT accessToken,refreshToken, expirationTime FROM user_tokens WHERE userId = ?"
	insertUserTokenQuery           = "INSERT OR REPLACE INTO user_tokens (userId,accessToken, refreshToken, expirationTime) VALUES (?, ?, ?, ?)"
	deleteUserTokenQuery           = "DELETE FROM user_tokens WHERE userId = ?"

	getUserAssetsByUserIdQuery = "SELECT * FROM user_assets WHERE userId = ?"
	insertUserAssetQuery       = "INSERT INTO user_assets (userId, assetId) VALUES (?, ?)"
	deleteUserAssetQuery       = "DELETE FROM user_assets WHERE userId = ? AND assetId = ?"
)
