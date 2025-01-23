package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/maksimUlitin/internal/helpers"
	"github.com/maksimUlitin/internal/lib"
	"github.com/maksimUlitin/internal/models"
	"github.com/maksimUlitin/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	userCollection *mongo.Collection = storage.OpenCollection(*storage.Client, "user")
	validate                         = validator.New()
	ctx, cancel                      = context.WithTimeout(context.Background(), 100*time.Second)
	user           models.User
	foundUser      models.User
	allUsers       []bson.M
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		logger.Error("Failed to hash password", "error", err)
		return ""
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return false, "email or password is incorrect"
	}
	return true, ""
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.BindJSON(&user); err != nil {
			logger.Warn("Invalid JSON in signup request", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			logger.Warn("Validation error in signup request", "error", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"Error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			logger.Error("Error checking for existing email", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
			return
		}

		if count > 0 {
			logger.Warn("Attempted signup with existing email", "email", user.Email)
			c.JSON(http.StatusConflict, gin.H{"error": "this email already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			logger.Error("Error checking for existing phone number", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the phone number"})
			return
		}

		if count > 0 {
			logger.Warn("Attempted signup with existing phone number", "phone", user.Phone)
			c.JSON(http.StatusConflict, gin.H{"error": "this phone number already exists"})
			return
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, *&user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			logger.Error("Failed to insert new user", "error", insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
			return
		}
		defer cancel()

		logger.Info("New user signed up successfully", "user_id", user.UserId)
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.BindJSON(&user); err != nil {
			logger.Warn("Invalid JSON in login request", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			logger.Warn("Login attempt with non-existent email", "email", user.Email)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			logger.Warn("Login attempt with incorrect password", "user_id", foundUser.UserId)
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, foundUser.UserId)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.UserId)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.UserId}).Decode(&foundUser)
		if err != nil {
			logger.Error("Error fetching user after token update", "error", err, "user_id", foundUser.UserId)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.Info("User logged in successfully", "user_id", foundUser.UserId)
		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			logger.Warn("Unauthorized access attempt to GetUsers", "error", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			logger.Error("Error occurred while listing user items", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
			return
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			logger.Error("Error decoding user list result", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding user items"})
			return
		}

		logger.Info("Users list retrieved successfully", "count", len(allUsers))
		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			logger.Warn("Unauthorized access attempt to GetUser", "error", err, "requested_user_id", userId)
			c.JSON(http.StatusForbidden, gin.H{"Error": err.Error()})
			return
		}

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			logger.Error("Error fetching user", "error", err, "user_id", userId)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the user"})
			return
		}

		logger.Info("User retrieved successfully", "user_id", userId)
		c.JSON(http.StatusOK, user)
	}
}
