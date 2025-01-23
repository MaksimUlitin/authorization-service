package helpers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/maksimUlitin/internal/lib"
	"github.com/maksimUlitin/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	UserType  string
	jwt.StandardClaims
}

var (
	userCollection *mongo.Collection = storage.OpenCollection(*storage.Client, "user")
	SEKRET_KEY     string            = os.Getenv("SEKRET_KEY")
	ctx, cancel                      = context.WithTimeout(context.Background(), 100*time.Second)
	updateObj      primitive.D
)

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Uid:       uid,
		UserType:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(SEKRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims).SignedString([]byte(SEKRET_KEY))
	if err != nil {
		logger.Error("Failed to generate tokens", "error", err, "user_id", uid)
		return
	}
	logger.Info("Tokens generated successfully", "user_id", uid)
	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) { return []byte(SEKRET_KEY), nil })
	if err != nil {
		msg = err.Error()
		logger.Warn("Token validation failed", "error", msg)
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		logger.Warn("Invalid token claims", "error", msg)
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		logger.Warn("Token expired", "user_id", claims.Uid)
		return
	}

	logger.Info("Token validated successfully", "user_id", claims.Uid)
	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

	usert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &usert,
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
	defer cancel()
	if err != nil {
		logger.Error("Failed to update tokens in database", "error", err, "user_id", userId)
		return
	}
	logger.Info("Tokens updated successfully in database", "user_id", userId)
	return
}
