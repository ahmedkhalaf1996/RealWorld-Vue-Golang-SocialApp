package realtime

import (
	"Server/database"
	"Server/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserFriends(userID string) <-chan []string {
	ch := make(chan []string)
	go func() {
		defer close(ch)

		// call grp client fun
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		userFriends, err := getFollowingFollowersClient(ctx, userID)

		if err != nil {
			fmt.Printf("Erro On Frinds Method Realtime")
			ch <- []string{}
			return
		}

		ch <- userFriends
	}()
	return ch

}

func getFollowingFollowersClient(ctx context.Context, userID string) ([]string, error) {
	UserSchema := database.DB.Collection("users")

	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}

	var user models.UserModel

	uid, _ := primitive.ObjectIDFromHex(userID)

	err := UserSchema.FindOne(ctx, bson.M{"_id": uid}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	friendsMap := make(map[string]bool)
	for _, id := range user.Following {
		friendsMap[id] = true
	}
	for _, id := range user.Followers {
		friendsMap[id] = true
	}

	var friends []string
	for id := range friendsMap {
		friends = append(friends, id)
	}

	log.Printf("Found %d frids for user %s", len(friends), userID)
	return friends, nil
}
