package controllers

import (
	"Server/database"
	"Server/models"
	"Server/services"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUserBy ID
// @Summary Get User By ID
// @Description GetUser Deatils By ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param page query int false "page number"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} map[string]interface{}
// @Router /user/getUser/{id} [get]
func GetUserByID(c *fiber.Ctx) error {

	var UserSchema = database.DB.Collection("users")
	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.UserModel
	var posts []bson.M

	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// Create Cache key for user profile & psots
	cacheKey := fmt.Sprintf("user:profile:%s:page:%d", c.Params("id"), page)
	// try to get form cache first .. from redis
	cachedData, err := database.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedRes models.CachedGetUserResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedRes); err == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"user":          cachedRes.User,
				"posts":         cachedRes.Posts,
				"currentPage":   cachedRes.CurrentPage,
				"numberOfPages": cachedRes.NumberOfPages,
				"cached":        true,
			})
		}
	} else {
		log.Printf("Cache miss for user profile %s: %s", c.Params("id"), err)
	}

	var LIMIT = 3

	userResult := UserSchema.FindOne(ctx, bson.M{"_id": objId})
	if userResult.Err() != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not found",
			"sucess":  false,
		})
	}
	userResult.Decode(&user)

	filter := bson.M{"creator": objId}

	// get total num of user posts
	total, err := PostSchema.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No Posts",
		})
	}

	/// Aggregation pipleline for posts with comments and user data

	pipeline := []bson.M{
		{"$match": bson.M{"creator": objId}},
		{"$sort": bson.M{"_id": -1}},
		{"$skip": int64((page - 1) * LIMIT)},
		{"$limit": int64(LIMIT)},
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "creator",
			"foreignField": "_id",
			"as":           "user",
		}},
		{"$unwind": "$user"},
		{"$lookup": bson.M{
			"from": "comments",
			"let":  bson.M{"postId": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$postId", "$$postId"}}}},
				{"$sort": bson.M{"createdAt": -1}},
				{"$lookup": bson.M{
					"from": "users",
					"let":  bson.M{"uid": "$userId"},
					"pipeline": []bson.M{
						{"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$_id", "$$uid"}}}},
						{"$project": bson.M{"name": 1, "imageUrl": 1}},
					},
					"as": "user",
				}},
				{"$unwind": bson.M{"path": "$user", "preserveNullAndEmptyArrays": true}},
				{"$project": bson.M{"_id": 1, "value": 1, "createdAt": 1, "userId": 1, "user": 1}},
			},
			"as": "comments",
		}},
		{"$project": bson.M{
			"_id":           1,
			"creator":       1,
			"title":         1,
			"message":       1,
			"name":          1,
			"selectedFile":  1,
			"likes":         1,
			"createdAt":     1,
			"comments":      1,
			"user.name":     1,
			"user.imageUrl": 1,
		}},
	}

	cursor, err := PostSchema.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "posts aggregation fiald",
			"details": err.Error(),
		})
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &posts); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "fiald to read aggregation posts",
				"details": err.Error()})
	}

	if posts == nil {
		posts = make([]bson.M, 0)
	}

	// Prepare response for caching

	response := models.CachedGetUserResponse{
		User:          user,
		Posts:         posts,
		CurrentPage:   page,
		NumberOfPages: math.Ceil(float64(total) / float64(LIMIT)),
	}

	responseJSON, err := json.Marshal(response)
	if err == nil {
		database.RedisClient.Set(ctx, cacheKey, responseJSON, 10*time.Second)
	} else {
		log.Printf("Fiald to marshal response for caching :%s ", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":          user,
		"posts":         posts,
		"currentPage":   page,
		"numberOfPages": math.Ceil(float64(total) / float64(LIMIT)),
	})
}

// UpdateUser
// @Summary update user data
// @Description update user deatils
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUser true "deatils "
// @Success 201 {object} models.UserModel
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /user/Update/{id} [patch]
func UpdateUser(c *fiber.Ctx) error {

	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//
	extUid := c.Locals("userId").(string)

	if extUid != c.Params("id") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "You Are Not Authroized to Update This Profile",
		})
	}

	userid, _ := primitive.ObjectIDFromHex(c.Params("id"))

	var user models.UpdateUser
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"deatils": err.Error(),
		})
	}

	update := bson.M{"name": user.Name, "imageUrl": user.ImageUrl, "bio": user.Bio}

	result, err := UserSchema.UpdateOne(ctx, bson.M{"_id": userid}, bson.M{"$set": update})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "cannot update the user data",
			"deatils": err.Error(),
		})
	}
	//
	var updateUsser models.UserModel
	if result.MatchedCount == 1 {
		err := UserSchema.FindOne(ctx, bson.M{"_id": userid}).Decode(&updateUsser)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"deatils": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": updateUsser})

}

// Following Users
// @Summary Follow/UnFollow User
// @Description follow or  un follow a user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /user/{id}/following [patch]
func FollowingUser(c *fiber.Ctx) error {

	var UserSchema = database.DB.Collection("users")
	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var FirstUser models.UserModel
	var SecondUser models.UserModel

	FirstUserID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	SecondUserID, _ := primitive.ObjectIDFromHex(c.Locals("userId").(string))

	err := UserSchema.FindOne(ctx, bson.M{"_id": FirstUserID}).Decode(&FirstUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": SecondUserID}).Decode(&SecondUser)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	fuid := c.Params("id")
	suid := c.Locals("userId").(string)

	if slices.Contains(FirstUser.Followers, suid) {
		i := sort.SearchStrings(FirstUser.Followers, suid)
		FirstUser.Followers = slices.Delete(FirstUser.Followers, i, i+1)
		// remove form the following list on second user
		index := sort.SearchStrings(SecondUser.Following, fuid)
		SecondUser.Following = slices.Delete(SecondUser.Following, index, index+1)
	} else {
		FirstUser.Followers = append(FirstUser.Followers, suid)
		SecondUser.Following = append(SecondUser.Following, fuid)

		// Create Notification
		notification := models.Notification{
			MainUID:   FirstUser.ID.Hex(),
			TargetID:  SecondUser.ID.Hex(),
			Deatils:   SecondUser.Name + " Start Following You!",
			UserID:    SecondUser.ID.Hex(),
			CreatedAt: time.Now(),
		}
		res, err := NotificationSchema.InsertOne(ctx, notification)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Faild to create notification",
				"error":   err.Error(),
			})
		}

		// set the id fiald of the notficato object
		notification.ID = res.InsertedID.(primitive.ObjectID)
		// call grpc
		services.SendNotification(notification)
	}

	updateFirst := bson.M{"followers": FirstUser.Followers}
	updateSecond := bson.M{"following": SecondUser.Following}

	_, err = UserSchema.UpdateOne(ctx, bson.M{"_id": FirstUserID}, bson.M{"$set": updateFirst})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}
	_, err = UserSchema.UpdateOne(ctx, bson.M{"_id": SecondUserID}, bson.M{"$set": updateSecond})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": FirstUserID}).Decode(&FirstUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": SecondUserID}).Decode(&SecondUser)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"SecondUser": SecondUser, "FirstUser": FirstUser})

}

// GetSugUser Users
// @Summary Get Suggersted users
// @Description get suggested userses based on the current user's following list
// @Tags Users
// @Accept json
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /user/getSug [get]
func GetSugUser(c *fiber.Ctx) error {

	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var MainUser models.UserModel
	var SugListId []string
	var AllSugUsers []models.UserModel

	MainUserID, err := primitive.ObjectIDFromHex(c.Query("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": MainUserID}).Decode(&MainUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	// Get SugUsers id then put them in suglistid
	for _, FID := range MainUser.Following {
		var singleUser models.UserModel
		convFID, _ := primitive.ObjectIDFromHex(FID)
		err = UserSchema.FindOne(ctx, bson.M{"_id": convFID}).Decode(&singleUser)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"deatils": err.Error(),
			})
		}

		// following
		for _, id := range singleUser.Following {
			if slices.Contains(SugListId, id) || id != c.Query("id") {
				SugListId = append(SugListId, id)
			}
		}

		// Followers
		for _, id := range singleUser.Followers {
			if slices.Contains(SugListId, id) || id != c.Query("id") {
				SugListId = append(SugListId, id)
			}
		}

	}

	// Gest Sug Users by id .
	if len(SugListId) > 0 {

		var objectides []primitive.ObjectID
		for _, id := range SugListId {
			objid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				continue
			}
			objectides = append(objectides, objid)
		}

		// fetch all users in one qeery using $in operator
		cursor, err := UserSchema.Find(ctx, bson.M{
			"_id": bson.M{"$in": objectides},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"deatils": err.Error(),
			})
		}

		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &AllSugUsers); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"deatils": err.Error(),
			})
		}
	}

	if AllSugUsers == nil {
		AllSugUsers = make([]models.UserModel, 0)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"users": AllSugUsers})
}

// DeleteUser
// @Summary delete user
// @Description delete user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object}  map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @security BearerAuth
// @Router /user/delete/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//
	extUid := c.Locals("userId").(string)

	if extUid != c.Params("id") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "You Are Not Authroized to Delete This User",
		})
	}

	userID, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid User id",
		})
	}

	result, err := UserSchema.DeleteOne(ctx, bson.M{"_id": userID})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "faild to delete user",
			"error":   err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}
	// sucuss
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User Deleted Successfully!",
	})
}
