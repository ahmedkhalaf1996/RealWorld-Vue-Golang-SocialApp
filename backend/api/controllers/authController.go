package controllers

import (
	"Server/database"
	"Server/models"
	"context"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Register
// @Summary Gegister a new user
// @Description Register an ew user by providing email, password , first name , last name
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "user register deatils"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} map[string]interface{}
// @Router /user/signup [post]
func Register(c *fiber.Ctx) error {

	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.CreateUser
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"deatils": err.Error(),
		})
	}

	CheckUser := UserSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&body)

	if CheckUser == nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "user with email" + body.Email + "Alraedy Exist!",
		})
	}

	// hashing password
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	newUser := models.UserModel{
		Name:      body.FirstName + " " + body.LastName,
		Email:     body.Email,
		Password:  string(hashPassword),
		Followers: make([]string, 0),
		Following: make([]string, 0),
	}

	result, err := UserSchema.InsertOne(ctx, &newUser)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(err)
	}

	// get the new user
	var createdUser *models.UserModel
	query := bson.M{"_id": result.InsertedID}

	UserSchema.FindOne(ctx, query).Decode(&createdUser)
	// create the token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    createdUser.ID.Hex(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	JwtSecret := os.Getenv("JWT_SECRET")

	token, _ := claims.SignedString([]byte(JwtSecret))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": createdUser,
		"token":  token,
	})
}

// Login
// @Summary login a  user
// @Description Login an user by providing email, password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.LoginUser true "user Login deatils"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} map[string]interface{}
// @Router /user/signin [post]
func Login(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.LoginUser
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"deatils": err.Error(),
		})
	}

	var user models.UserModel
	CheckEmail := UserSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&user)

	// check if user with prvided email exist or not
	if CheckEmail != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "Invalid User With Email" + body.Email,
		})
	}

	// check if we have the same pass or not
	checkPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if checkPass != nil {
		return c.Status(fiber.StatusBadGateway).JSON(string(checkPass.Error()))
	}

	// create the token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID.Hex(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	JwtSecret := os.Getenv("JWT_SECRET")

	token, _ := claims.SignedString([]byte(JwtSecret))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": user,
		"token":  token,
	})
}

// Refresh Userdata
// @Summary refresh user data and token
// @Description  refresh user data and return issue a new token
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /user/refresh [get]
func RefreshUser(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing Authorizaiton header",
		})
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)

	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthenticated",
		})
	}

	userID, err := primitive.ObjectIDFromHex(claims.Issuer)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user Id in token",
		})
	}

	var user models.UserModel
	err = UserSchema.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// c new claims
	// create the token
	newclaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID.Hex(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	newToken, _ := newclaims.SignedString([]byte(jwtSecret))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": user,
		"token":  newToken,
	})

}
