package controllers

import (
	"Server/database"
	"Server/models"
	"context"
	"math"
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create Post
// @Summary create  a new post
// @Description create new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body models.CreateOrUpdatePost true "post create  deatils"
// @Success 201 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts [post]
func CraetePost(c *fiber.Ctx) error {

	var UserSchema = database.DB.Collection("users")
	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.CreateOrUpdatePost
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"deatils": err.Error(),
		})
	}

	// start set data
	var post models.PostModel
	post.Creator = c.Locals("userId").(string)
	post.Likes = make([]string, 0)
	post.Comments = make([]string, 0)
	post.CreatedAt = time.Now()
	post.Title = body.Title
	post.Message = body.Message
	post.SelectedFile = body.SelectedFile
	//

	var user models.UserModel
	objId, _ := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	err := UserSchema.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//
	post.Name = user.Name
	// set data end
	// craete post
	result, err := PostSchema.InsertOne(ctx, &post)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	} else {
		var createdPost *models.PostModel
		query := bson.M{"_id": result.InsertedID}

		PostSchema.FindOne(ctx, query).Decode(&createdPost)
		return c.Status(fiber.StatusCreated).JSON(createdPost)
	}

}

// Get Post
// @Summary Get  a new post
// @Description Get new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post id"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Router /posts/{id} [get]
func GetPost(c *fiber.Ctx) error {

	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "post id is required",
		})
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var post *models.PostModel
	query := bson.M{"_id": objID}

	err = PostSchema.FindOne(ctx, query).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "post Not Found",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"post": post,
		})

}

// Update Post
// @Summary Update  post
// @Description Update post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post Id"
// @Param post body models.CreateOrUpdatePost true "update post  deatils"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts/{id} [patch]
func UpdatePost(c *fiber.Ctx) error {

	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var newData models.CreateOrUpdatePost
	if err := c.BodyParser(&newData); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"deatils": err.Error(),
		})
	}

	// authorization start
	var authPost models.PostModel
	primID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	PostSchema.FindOne(ctx, bson.M{"_id": primID}).Decode(&authPost)

	if authPost.Creator != c.Locals("userId").(string) {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "You Are Not authorized to update this post.",
		})
	}

	// set data end
	authPost.Title = newData.Title
	authPost.Message = newData.Message
	authPost.SelectedFile = newData.SelectedFile

	// craete post
	// update := bson.M{"title": newData.Title, "message":newData.Message, "selectedFile": newData.SelectedFile}
	_, err = PostSchema.UpdateOne(ctx, bson.M{"_id": authPost.ID}, bson.M{"$set": authPost})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	} else {

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": authPost})
	}

}

// GetAllPosts Post
// @Summary Get All Posts
// @Description GetAllPosts with pagination
// @Tags Posts
// @Accept json
// @Produce json
// @Param page query int false "page number"
// @Param id query string true "user id"
// @Success 200 {object} []models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts [get]
func GetAllPosts(c *fiber.Ctx) error {

	var PostSchema = database.DB.Collection("posts")
	var userSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.UserModel
	var posts []models.PostModel

	userid := c.Query("id")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// get user follwoing list ides and add our user id to it
	MainUserid, _ := primitive.ObjectIDFromHex(userid)
	userSchema.FindOne(ctx, bson.M{"_id": MainUserid}).Decode(&user)
	user.Following = append(user.Following, userid)
	///

	var LIMIT = 2

	findOptions := options.Find()
	filter := bson.M{"creator": bson.M{"$in": user.Following}}

	// get total num of posts
	total, err := PostSchema.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No Posts",
		})
	}

	findOptions.SetSkip((int64(page) - 1) * int64(LIMIT))
	findOptions.SetLimit(int64(LIMIT))
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	// start get psots
	cursor, err := PostSchema.Find(ctx, filter, findOptions)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No Posts",
		})
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var post models.PostModel
		cursor.Decode(&post)
		posts = append(posts, post)
	}

	if posts == nil {
		posts = make([]models.PostModel, 0)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":          posts,
		"currentPage":   page,
		"numberOfPages": math.Ceil(float64(total) / float64(LIMIT)),
	})

}

// GetPostsUsersBySearch Post
// @Summary Get Posts users by serach query
// @Description get posts adnd users matching the search query
// @Tags Posts
// @Accept json
// @Produce json
// @Param searchQuery query string true "Search query"
// @Failure 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts/search [get]
func GetPostsUsersBySearch(c *fiber.Ctx) error {

	var PostSchema = database.DB.Collection("posts")
	var userSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var users []models.UserModel
	var posts []models.PostModel

	//
	filterPost := bson.M{}
	filterUser := bson.M{}

	//
	findOptionsPost := options.Find()
	findOptionsUser := options.Find()

	if search := c.Query("searchQuery"); search != "" {
		// post
		filterPost = bson.M{
			"$or": []bson.M{
				{
					"title": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
				{
					"description": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
			},
		}
		//
		filterUser = bson.M{
			"$or": []bson.M{
				{
					"name": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
				{
					"email": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
			},
		}
	}
	// end
	cursorPosts, _ := PostSchema.Find(ctx, filterPost, findOptionsPost)
	defer cursorPosts.Close(ctx)

	cursorUsers, _ := userSchema.Find(ctx, filterUser, findOptionsUser)
	defer cursorUsers.Close(ctx)
	//

	for cursorUsers.Next(ctx) {
		var user models.UserModel
		cursorUsers.Decode(&user)
		users = append(users, user)
	}

	for cursorPosts.Next(ctx) {
		var post models.PostModel
		cursorPosts.Decode(&post)
		posts = append(posts, post)
	}

	return c.JSON(fiber.Map{
		"user":  users,
		"posts": posts,
	})
}

// Comment Post
// @Summary comment  post
// @Description comment post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post Id"
// @Param post body models.ComnmentPost true "comment value"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts/{id}/commentPost [post]
func CommentPost(c *fiber.Ctx) error {

	var PostSchema = database.DB.Collection("posts")
	var UserSchema = database.DB.Collection("users")
	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var b models.ComnmentPost
	if err := c.BodyParser(&b); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"deatils": err.Error(),
		})
	}

	var post models.PostModel
	postid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	err = PostSchema.FindOne(ctx, bson.M{"_id": postid}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}
	//
	newComment := bson.M{"comments": append(post.Comments, b.Value)}
	_, err = PostSchema.UpdateOne(ctx, bson.M{"_id": postid}, bson.M{"$set": newComment})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}
	err = PostSchema.FindOne(ctx, bson.M{"_id": postid}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}
	// CREATE NOtification start
	userID := c.Locals("userId").(string)
	objId, _ := primitive.ObjectIDFromHex(userID)
	var user models.UserModel

	// get nuser data
	userResult := UserSchema.FindOne(ctx, bson.M{"_id": objId})
	if userResult.Err() != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"success": false,
			"message": "User Not found",
		})
	}

	userResult.Decode(&user)
	// Create Notification
	notification := models.Notification{
		MainUID:   post.Creator,
		TargetID:  postid.Hex(),
		Deatils:   user.Name + " Commented on your Post",
		User:      models.User{Name: user.Name, Avatart: user.ImageUrl},
		CreatedAt: time.Now(),
	}
	_, err = NotificationSchema.InsertOne(ctx, notification)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Faild to create notification",
			"error":   err.Error(),
		})
	}
	// end
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": post,
	})

}

// like Post
// @Summary like or unkike a post
// @Description Like or un like a post  by it's id
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post Id"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts/{id}/likePost [patch]
func LikePost(c *fiber.Ctx) error {

	var PostSchema = database.DB.Collection("posts")
	var UserSchema = database.DB.Collection("users")
	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var post models.PostModel
	postid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	err = PostSchema.FindOne(ctx, bson.M{"_id": postid}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}
	// after getting post
	userID, errb := c.Locals("userId").(string)
	if !errb {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"deatils": "you are not authorized",
		})
	}

	// check
	if slices.Contains(post.Likes, userID) {
		i := sort.SearchStrings(post.Likes, userID)
		post.Likes = slices.Delete(post.Likes, i, i+1)
	} else {
		post.Likes = append(post.Likes, userID)
		//  START craete Notification
		objId, _ := primitive.ObjectIDFromHex(userID)
		var user models.UserModel

		// get nuser data
		userResult := UserSchema.FindOne(ctx, bson.M{"_id": objId})
		if userResult.Err() != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"success": false,
				"message": "User Not found",
			})
		}

		userResult.Decode(&user)
		// Create Notification
		notification := models.Notification{
			MainUID:   post.Creator,
			TargetID:  post.ID.Hex(),
			Deatils:   user.Name + " Liked your Post",
			User:      models.User{Name: user.Name, Avatart: user.ImageUrl},
			CreatedAt: time.Now(),
		}
		_, err = NotificationSchema.InsertOne(ctx, notification)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Faild to create notification",
				"error":   err.Error(),
			})
		}
		// End create notfication
	}

	// update post
	updateLikelist := bson.M{"likes": post.Likes}
	_, err = PostSchema.UpdateOne(ctx, bson.M{"_id": postid}, bson.M{"$set": updateLikelist})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}
	err = PostSchema.FindOne(ctx, bson.M{"_id": postid}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"post": post,
	})
}

// Delete Post
// @Summary Delete  post by id
// @Description Delete post by post id need to prvided auth token for post craetor
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post Id"
// @Failure 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts/{id} [delete]
func DeletePost(c *fiber.Ctx) error {

	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// authorization start
	var authPost models.PostModel
	primID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	PostSchema.FindOne(ctx, bson.M{"_id": primID}).Decode(&authPost)

	if authPost.Creator != c.Locals("userId").(string) {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "You Are Not authorized to delete this post.",
		})
	}

	//
	result, err := PostSchema.DeleteOne(ctx, bson.M{"_id": primID})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	if result.DeletedCount == 1 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Post Deleted Successfully!",
		})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "can't Delete Post!",
		})
	}

}
