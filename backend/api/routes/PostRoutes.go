package routes

import (
	// "Server/controllers"
	// "Server/validation"

	"Server/controllers"
	"Server/middleware"
	"Server/validation"

	"github.com/gofiber/fiber/v2"
)

func SetupPostRoutes(app *fiber.App) {
	// create
	app.Post("/posts", middleware.AuthMiddleware, validation.ValidatePost, controllers.CraetePost)
	// getall
	app.Get("/posts", controllers.GetAllPosts)
	// search
	app.Get("/posts/search", controllers.GetPostsUsersBySearch)
	// get one
	app.Get("/posts/:id", controllers.GetPost)
	// update
	app.Patch("/posts/:id", middleware.AuthMiddleware, validation.ValidatePost, controllers.UpdatePost)
	// commnet
	app.Post("/posts/:id/commentPost", middleware.AuthMiddleware, controllers.CommentPost)
	// like
	app.Patch("/posts/:id/likePost", middleware.AuthMiddleware, controllers.LikePost)

	// delete
	app.Delete("/posts/:id", middleware.AuthMiddleware, controllers.DeletePost)
}
