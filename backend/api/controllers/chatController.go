package controllers

import (
	"Server/database"
	"Server/models"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SendMessage
// @Summary send message to friend user
// @Description SendMessage form one user to another
// @Tags Chat
// @Accept json
// @Produce json
// @Param message body models.SendMessageM true "user SendMessage deatils"
// @Success 201 {object} models.Message
// @Failure 400 {object} map[string]interface{}
// @Router /chat/sendmessage [post]
func SendMessage(c *fiber.Ctx) error {

	var MessageSchema = database.DB.Collection("messages")
	var UnReadedMsgSchema = database.DB.Collection("unReadedmessages")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.SendMessageM
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"deatils": err.Error(),
		})
	}

	var msg models.Message
	c.BodyParser(&msg)
	// save the message to db
	result, err := MessageSchema.InsertOne(ctx, &msg)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "faild to save msg",
			"deatils": err.Error(),
		})
	}

	// update or create the unreaded message count and is readed
	var unRreadedMsg models.UnReadedMsg
	filtter := bson.M{"mainUserid": msg.Recever, "otherUserid": msg.Sender}
	update := bson.M{"$inc": bson.M{"numOfUnreadedMessages": 1}, "$set": bson.M{"isReaded": false}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	err = UnReadedMsgSchema.FindOneAndUpdate(ctx, filtter, update, opts).Decode(&unRreadedMsg)
	if err != nil && err != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "Faild to update unareded message count",
			"deatils": err.Error(),
		})
	}

	// Return the created message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Message sent Successfully",
		"result":  result.InsertedID,
	})

}

// GetMsgsByNums
// @Summary get message by pagenation
// @Description GetMsgsByNumbetween two users by pagenation
// @Tags Chat
// @Accept json
// @Produce json
// @Param from query int true "Staring point page num"
// @Param firstuid query string true "first user id"
// @Param seconduid query string true "second user id"
// @Success 201 {object} []models.Message
// @Failure 400 {object} map[string]interface{}
// @Router /chat/getmsgsbynums [get]
func GetMsgsByNums(c *fiber.Ctx) error {

	var MessageSchema = database.DB.Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	from, err := strconv.Atoi(c.Query("from"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid value form from",
			"err":     err.Error(),
		})
	}

	firstuid := c.Query("firstuid")
	seconduid := c.Query("seconduid")

	// construct the filer
	senderFilter := bson.M{"sender": firstuid, "recever": seconduid}
	receiverFilter := bson.M{"sender": seconduid, "recever": firstuid}
	filter := bson.M{"$or": []bson.M{senderFilter, receiverFilter}}

	// pagenation options
	options := options.Find()
	options.SetSort(bson.D{{Key: "_id", Value: -1}})
	options.SetSkip(int64(from * 2))
	options.SetLimit(2)

	// query the db
	cursor, err := MessageSchema.Find(ctx, filter, options)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to retrieve messages",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	// iterate over the cursor and build the res array
	var messages []models.Message
	for cursor.Next(ctx) {
		var msg models.Message
		err := cursor.Decode(&msg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Faild to decode messages",
				"error":   err.Error(),
			})
		}
		messages = append(messages, msg)
	}

	// reverce the message array

	for i := 0; i < len(messages)/2; i++ {
		j := len(messages) - 1 - i
		messages[i], messages[j] = messages[j], messages[i]
	}

	if len(messages) == 0 {
		messages = []models.Message{}
	}

	// Return the created message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msgs": messages,
	})

}

// GetuserUnreadedmessage
// @Summary Get unreaded message count & recodes for user
// @Description Get unreaded message count & recodes for user
// @Tags Chat
// @Accept json
// @Produce json
// @Param userid query string true "user id"
// @Failure 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /chat/get-user-unreadedmsg [get]
func GetUserUnreadedMsg(c *fiber.Ctx) error {

	var UnReadedMsgSchema = database.DB.Collection("unReadedmessages")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userid := c.Query("userid")
	if userid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user id query param is requied!",
		})
	}
	// filter
	filter := bson.M{"mainUserid": userid, "isReaded": false}

	// query the db
	cursor, err := UnReadedMsgSchema.Find(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to retrieve unareded messages",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	// iterate over the cursor and build the res array
	var urms []models.UnReadedMsg
	totalUnreadMessageCount := 0

	for cursor.Next(ctx) {
		var urm models.UnReadedMsg
		err := cursor.Decode(&urm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Faild to decode unreded message ",
				"error":   err.Error(),
			})
		}
		if !urm.IsReaded {
			urms = append(urms, urm)
		}
		totalUnreadMessageCount += urm.NumOfUnreadedMessages
	}

	if len(urms) == 0 {
		urms = []models.UnReadedMsg{}
	}

	// Return the created message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": urms,
		"total":    totalUnreadMessageCount,
	})

}

// MarkMsgAsReaded
// @Summary mark messages as read for user
// @Description mark messages as read for user uupate the recoded make is read true num 0
// @Tags Chat
// @Accept json
// @Produce json
// @Param mainuid query string true "main user id"
// @Param otheruid query string true "ohter user id"
// @Failure 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /chat/mark-msg-asreaded [get]
func MarkMsgAsReaded(c *fiber.Ctx) error {

	var UnReadedMsgSchema = database.DB.Collection("unReadedmessages")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mainuid := c.Query("mainuid")
	otheruid := c.Query("otheruid")
	if mainuid == "" || otheruid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "mainuid and other userid query params is requied!",
		})
	}
	// filter
	filter := bson.M{"mainUserid": mainuid, "otherUserid": otheruid}
	update := bson.M{"$set": bson.M{"isReaded": true, "numOfUnreadedMessages": 0}}

	// update the docoument
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	result := UnReadedMsgSchema.FindOneAndUpdate(ctx, filter, update, options)

	if result.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messages": "Faild to mark message as readed",
			"error":    result.Err().Error(),
		})

	}

	// check
	var updateDoc bson.M
	if err := result.Decode(&updateDoc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messages": "Faild to decode update docoument",
			"error":    result.Err().Error(),
		})

	}

	isMarked := updateDoc != nil

	// Return the created message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isMarked": isMarked,
	})

}
