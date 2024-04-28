package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/michaelwongycn/crypto-tracker/domain/request"
	"github.com/michaelwongycn/crypto-tracker/domain/response"
	"github.com/michaelwongycn/crypto-tracker/lib/auth"
	"github.com/michaelwongycn/crypto-tracker/usecase/user"
)

const (
	invalidCredentialsErrorMsg      = "Invalid Credentials"
	passwordNotMatchErrorMsg        = "Password doesn't match"
	emailAlreadyRegisteredErrorMsg  = "Email already registered"
	unableToParseTokenErrorMsg      = "Unable to parse token"
	assetNotFoundErrorMsg           = "Asset not found"
	assetAlreadyRegisteredErrorMsg  = "Asset already registered"
	unableToGetAssetDataErrorMsg    = "Unable to get asset data"
	failedToAddUserToDBErrorMsg     = "Failed to add user to the database"
	failedToAddAssetToDBErrorMsg    = "Failed to add asset to the database"
	failedToDeleteAssetToDBErrorMsg = "Failed to delete asset from the database"
	internalServerErrorMsg          = "Internal Server Error"
)

type controllerImpl struct {
	userUsecase user.UserUsecase
}

func NewControllerImpl(userUsecase user.UserUsecase) Controller {
	return &controllerImpl{
		userUsecase: userUsecase,
	}
}

func setResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (c *controllerImpl) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong!"))
}

func (c *controllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	var credentials request.UserAuthRequest
	authResponse := response.AuthResponse{}
	response := response.ReadResponse{}
	response.Time = requestTime

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		response.Message = err.Error()
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	accessToken, refreshToken, err := c.userUsecase.Login(ctx, credentials.Email, credentials.Password)
	if err != nil {
		response.Message = invalidCredentialsErrorMsg
		setResponse(w, http.StatusOK, response)
		return
	}

	authResponse.AccessToken = *accessToken
	authResponse.RefreshToken = *refreshToken

	response.Message = ""
	response.Data = authResponse
	setResponse(w, http.StatusOK, response)
}

func (c *controllerImpl) Register(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	var credentials request.UserRegisterRequest
	response := response.WriteResponse{}
	response.Time = requestTime

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		response.Message = err.Error()
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	if credentials.Password != credentials.PasswordConfirmation {
		response.Message = passwordNotMatchErrorMsg
		setResponse(w, http.StatusOK, response)
		return
	}

	err := c.userUsecase.Register(ctx, credentials.Email, credentials.Password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			response.Message = emailAlreadyRegisteredErrorMsg
			setResponse(w, http.StatusConflict, response)
			return
		}
		response.Message = failedToAddUserToDBErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = ""
	setResponse(w, http.StatusOK, response)
}

func (c *controllerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	response := response.WriteResponse{}
	response.Time = requestTime

	accessToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	claims, err := auth.ParseToken(accessToken)
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	userId := int(claims["sub"].(float64))
	err = c.userUsecase.Logout(ctx, accessToken, userId)
	if err != nil {
		response.Message = internalServerErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = ""
	setResponse(w, http.StatusOK, response)
}

func (c *controllerImpl) RefreshToken(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	credentials := request.UserRefreshTokenRequest{}
	authResponse := response.AuthResponse{}
	response := response.ReadResponse{}
	response.Time = requestTime

	accessToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	claims, err := auth.ParseToken(accessToken)
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}
	accessTokenUserId := int(claims["sub"].(float64))

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		response.Message = err.Error()
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	claims, err = auth.ParseToken(credentials.RefreshToken)
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}
	refreshTokenUserId := int(claims["sub"].(float64))

	if refreshTokenUserId != accessTokenUserId {
		response.Message = invalidCredentialsErrorMsg
		setResponse(w, http.StatusOK, response)
		return
	}

	newAccessToken, newRefreshToken, err := c.userUsecase.RefreshToken(ctx, credentials.RefreshToken, refreshTokenUserId)
	if err != nil {
		response.Message = invalidCredentialsErrorMsg
		setResponse(w, http.StatusOK, response)
		return
	}

	authResponse.AccessToken = *newAccessToken
	authResponse.RefreshToken = *newRefreshToken

	response.Message = ""
	response.Data = authResponse
	setResponse(w, http.StatusOK, response)
}

func (c *controllerImpl) ShowUserAsset(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	response := response.ReadResponse{}
	response.Time = requestTime

	accessToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	claims, err := auth.ParseToken(accessToken)
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	userId := int(claims["sub"].(float64))

	assets, err := c.userUsecase.GetUserAssetsByUserId(ctx, userId)
	if err != nil {
		response.Message = unableToGetAssetDataErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = ""
	response.Data = assets
	setResponse(w, http.StatusOK, response)
}

func (c *controllerImpl) InsertUserAsset(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	var credentials request.UserInsertAssetRequest
	response := response.WriteResponse{}
	response.Time = requestTime

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		response.Message = err.Error()
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	accessToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	claims, err := auth.ParseToken(accessToken)
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	userId := int(claims["sub"].(float64))
	err = c.userUsecase.InsertUserAsset(ctx, userId, credentials.AssetID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: user_assets.userId, user_assets.assetId") {
			response.Message = assetAlreadyRegisteredErrorMsg
			setResponse(w, http.StatusConflict, response)
			return
		} else if strings.Contains(err.Error(), "API request failed with status code: 404") {
			response.Message = assetNotFoundErrorMsg
			setResponse(w, http.StatusNotFound, response)
			return
		}
		response.Message = failedToAddAssetToDBErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}
	response.Message = ""
	setResponse(w, http.StatusOK, response)
}

func (c *controllerImpl) DeleteUserAsset(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	var credentials request.UserInsertAssetRequest
	response := response.WriteResponse{}
	response.Time = requestTime

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		response.Message = err.Error()
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	accessToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	claims, err := auth.ParseToken(accessToken)
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	userId := int(claims["sub"].(float64))
	err = c.userUsecase.DeleteUserAsset(ctx, userId, credentials.AssetID)
	if err != nil {
		if strings.Contains(err.Error(), "API request failed with status code: 404") {
			response.Message = assetNotFoundErrorMsg
			setResponse(w, http.StatusNotFound, response)
			return
		}
		response.Message = failedToDeleteAssetToDBErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = ""
	setResponse(w, http.StatusOK, response)
}
