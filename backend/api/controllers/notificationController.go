package controllers

import (
	"Server/database"
	"Server/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationResponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty"`
	Deatils   string             `json:"deatils"`
	MainUID   string             `json:"mainuid"`
	TargetID  string             `json:"targetid"`
	Isreded   bool               `json:"isreded"`
	CreatedAt time.Time          `json:"createdAt"`
	User      UserData           `json:"user"`
}

type UserData struct {
	Name     string `json:"name"`
	ImageUrl string `json:"iamgeUrl"`
}

// helper for pupulte user data
func populateNotificationsWithUserData(ctx context.Context, notificatons []models.Notification) ([]NotificationResponse, error) {
	var response []NotificationResponse
	userSchema := database.DB.Collection("users")

	for _, notificaton := range notificatons {
		userObjID, err := primitive.ObjectIDFromHex(notificaton.UserID)
		if err != nil {
			continue // skip invalid user ides

		}
		var user models.UserModel
		err = userSchema.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&user)
		if err != nil {
			user.Name = "UnKnown User"
			user.ImageUrl = ""
		}
		notificatonResp := NotificationResponse{
			ID:        notificaton.ID,
			Deatils:   notificaton.Deatils,
			MainUID:   notificaton.MainUID,
			TargetID:  notificaton.TargetID,
			Isreded:   notificaton.IsReaded,
			CreatedAt: notificaton.CreatedAt,
			User: UserData{
				Name:     user.Name,
				ImageUrl: user.ImageUrl,
			},
		}
		response = append(response, notificatonResp)
	}

	return response, nil
}

// MarkNotAsReaded Post
// @Summary Mark Notfication AsReaded  for a user
// @Description MarkNotAsReaded
// @Tags Notifications
// @Accept json
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /notification/mark-notification-asreaded [get]
func MarknotAsReaded(c *fiber.Ctx) error {

	// parse query paramter
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "id in query is Requierd",
		})
	}

	// construct the filter and update
	filter := bson.M{"mainuid": bson.M{"$regex": id, "$options": "i"}}
	update := bson.M{"$set": bson.M{"isreded": true}}

	// update
	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := NotificationSchema.UpdateMany(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to mark notifications as read",
			"error":   err.Error(),
		})
	}

	// retreive the udpated not
	options := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := NotificationSchema.Find(ctx, filter, options)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild retruieve the udpated notifications",
			"error":   err.Error(),
		})
	}

	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to decoded notifications",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"notifications": notifications,
	})
}

// GetUserNotification Post
// @Summary Get user notifications
// @Description GetUserNotification
// @Tags Notifications
// @Accept json
// @Produce json
// @Param userid path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /notification/{userid} [get]
func GetUserNotification(c *fiber.Ctx) error {

	// parse query paramter
	id := c.Params("userid")
	if id == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "id in Params is Requierd",
		})
	}

	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// construct the filter and update
	filter := bson.M{"mainuid": bson.M{"$regex": id, "$options": "i"}}
	options := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	// retreive the udpated not
	cursor, err := NotificationSchema.Find(ctx, filter, options)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild retruieve the udpated notifications",
			"error":   err.Error(),
		})
	}

	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to decoded notifications",
			"error":   err.Error(),
		})
	}

	if len(notifications) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"notifications": []NotificationResponse{},
		})
	}

	// Populate with fresh user data
	response, err := populateNotificationsWithUserData(ctx, notifications)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Failed to populate user data",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"notifications": response,
	})
}
