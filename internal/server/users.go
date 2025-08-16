package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	domainErrors "github.com/Dorrrke/library0706/internal/domain/errors"
	"github.com/Dorrrke/library0706/internal/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (api *LibraryApi) register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Println("failed to bind body: ", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	api.log.Debug().Str("user email", user.Email).Str("user pass", user.Pass).Msg("incoming user from request")

	valid := validator.New()

	err := valid.Struct(user)
	if err != nil {
		api.log.Error().Err(err).Msg("failed validation")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.UID = uuid.New().String()
	api.log.Debug().Msgf("user uid: %s", user.UID)

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to hash password")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Pass = string(hash)

	err = api.db.SaveUser(user)
	if err != nil {
		if errors.Is(err, domainErrors.ErrUserAlredyExist) {
			api.log.Error().Err(err).Send()
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		api.log.Error().Err(err).Msg("failed to save user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(200, "OK")
}

func (api *LibraryApi) login(ctx *gin.Context) {
	var user models.UserLogin
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Println("failed to bind body: ", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	api.log.Debug().Str("user email", user.Email).Str("user pass", user.Pass).Msg("incoming user from request")

	valid := validator.New()

	err := valid.Struct(user)
	if err != nil {
		api.log.Error().Err(err).Msg("failed validation")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbUser, err := api.db.GetUser(user.Email)
	if err != nil {
		if errors.Is(err, domainErrors.ErrIvalidCreds) {
			api.log.Error().Err(err).Send()
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		api.log.Error().Err(err).Msg("failed to get user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Pass), []byte(user.Pass))
	if err != nil {
		api.log.Error().Err(err).Msg("failed to compare password")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// ---- Работа с JWT ----
	accessToken, _, err := createToken(dbUser.UID, accessTokenDuration)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to create access token")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	refreshToken, refreshTokenID, err := createToken(dbUser.UID, refreshTokenDuration)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to create refresh token")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = api.db.SaveRefreshToken(refreshToken, refreshTokenID, dbUser.UID)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to save refresh token")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		int(refreshTokenDuration/time.Second),
		"/",
		"localhost",
		true,
		true,
	)

	ctx.Header("Authorization", "Bearer "+accessToken)

	ctx.JSON(http.StatusOK, dbUser)
}

func (api *LibraryApi) refreshHandler(ctx *gin.Context) {
	token, err := ctx.Cookie("refresh_token")
	if err != nil {
		api.log.Error().Err(err).Msg("failed to get refresh token")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, err := parseToken(token)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to parse token")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// TODO: проверка наличия refreshToken в БД
	valid, err := api.db.CheckRefreshToken(claims.ID)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to check refresh token")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !valid {
		api.log.Error().Err(err).Msg("invalid refresh token")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "bad refresh token"})
		return
	}

	newAccessToken, _, err := createToken(claims.UserID, accessTokenDuration)
	if err != nil {
		api.log.Error().Err(err).Msg("failed to create access token")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Authorization", "Bearer "+newAccessToken)
	ctx.JSON(http.StatusOK, "OK")
}
