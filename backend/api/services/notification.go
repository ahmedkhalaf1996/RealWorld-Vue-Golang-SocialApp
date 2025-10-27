package services

import (
	"Server/database"
	"Server/models"
	"Server/realtime"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SendNotification(notification models.Notification) error {
	userCol := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(notification.UserID)
	if err != nil {
		log.Printf("Invalid user Id : %v", err)
		return err
	}

	var user models.UserModel
	err = userCol.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&user)
	if err != nil {
		log.Printf("faild to fetch user data : %v", err)
		user.Name = "Unknown user"
		user.ImageUrl = ""
	}

	realTimeNotification := realtime.Notification{
		ID:        notification.ID.Hex(),
		Details:   notification.Deatils,
		MainUID:   notification.MainUID,
		TargetID:  notification.TargetID,
		IsReaded:  notification.IsReaded,
		CraetedAt: notification.CreatedAt,
		User: realtime.User{
			Name:     user.Name,
			ImageUrl: user.ImageUrl,
		},
	}

	notificationManager := realtime.GetNotificationManager()
	err = notificationManager.SendNotificatonToUser(notification.MainUID, realTimeNotification)

	if err != nil {
		log.Printf("Faild to send real time notification : %v", err)
		return err
	}

	log.Printf("Notification send succueffully to user %s", notification.MainUID)
	return nil
}
