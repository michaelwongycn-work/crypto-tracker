package cryptoDB

import (
	"context"
	"database/sql"
	"time"

	"github.com/michaelwongycn/crypto-tracker/domain/model"
	"github.com/michaelwongycn/crypto-tracker/lib/log"
)

const (
	noRowsFoundErrorMsg      = "no rows found for the query"
	errorScanningRowErrorMsg = "error when scanning row"
	errorQueryingSQLErrorMsg = "error when querying SQL"
)

type cryptoDBImpl struct {
	db      *sql.DB
	timeout time.Duration
}

func NewCryptoDBImpl(timeout time.Duration, db *sql.DB) CryptoDBInterface {
	return &cryptoDBImpl{
		db:      db,
		timeout: timeout * time.Second,
	}
}

func (d *cryptoDBImpl) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*model.User, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	var data model.User
	row := d.db.QueryRowContext(ctx, getUserByEmailAndPasswordQuery, email, password)

	err := row.Scan(&data.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.PrintLogErr(ctx, noRowsFoundErrorMsg, err)
			return nil, err
		} else {
			log.PrintLogErr(ctx, errorScanningRowErrorMsg, err)
			return nil, err
		}
	}
	return &data, nil
}

func (d *cryptoDBImpl) InsertUser(ctx context.Context, email, password string) error {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	_, err := d.db.ExecContext(ctx, insertUserQuery, email, password)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return err
	}

	return nil
}

func (d *cryptoDBImpl) GetUserToken(ctx context.Context, userId int) (*model.UserToken, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	UserToken := model.UserToken{
		UserId: userId,
	}
	row := d.db.QueryRowContext(ctx, getUserTokenQuery, userId)

	err := row.Scan(&UserToken.AccessToken, &UserToken.RefreshToken, &UserToken.ExpirationTime)
	if err != nil {
		if err == sql.ErrNoRows {
			log.PrintLogErr(ctx, noRowsFoundErrorMsg, err)
			return nil, err
		} else {
			log.PrintLogErr(ctx, errorScanningRowErrorMsg, err)
			return nil, err
		}
	}
	return &UserToken, nil
}

func (d *cryptoDBImpl) InsertUserToken(ctx context.Context, userId int, accessToken, refreshToken string, expirationTime int64) error {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	_, err := d.db.ExecContext(ctx, insertUserTokenQuery, userId, accessToken, refreshToken, expirationTime)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return err
	}

	return nil
}

func (d *cryptoDBImpl) DeleteUserToken(ctx context.Context, userId int) error {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	_, err := d.db.ExecContext(ctx, deleteUserTokenQuery, userId)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return err
	}

	return nil
}

func (d *cryptoDBImpl) GetUserAssetsByUserId(ctx context.Context, userId int) (*[]model.UserAsset, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	rows, err := d.db.QueryContext(ctx, getUserAssetsByUserIdQuery, userId)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return nil, err
	}

	var data []model.UserAsset
	for rows.Next() {
		var userAsset model.UserAsset
		err := rows.Scan(&userAsset.ID, &userAsset.UserId, &userAsset.AssetId)
		if err != nil {
			if err == sql.ErrNoRows {
				log.PrintLogErr(ctx, noRowsFoundErrorMsg, err)
				return nil, err
			} else {
				log.PrintLogErr(ctx, errorScanningRowErrorMsg, err)

				return nil, err
			}
		}
		data = append(data, userAsset)
	}
	return &data, nil
}

func (d *cryptoDBImpl) InsertUserAsset(ctx context.Context, userId int, assetId string) error {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	_, err := d.db.ExecContext(ctx, insertUserAssetQuery, userId, assetId)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return err
	}

	return nil
}

func (d *cryptoDBImpl) DeleteUserAsset(ctx context.Context, userId int, assetId string) error {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	_, err := d.db.ExecContext(ctx, deleteUserAssetQuery, userId, assetId)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return err
	}

	return nil
}
