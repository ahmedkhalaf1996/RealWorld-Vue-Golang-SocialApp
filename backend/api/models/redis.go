package models

import "go.mongodb.org/mongo-driver/bson"

// chachedrespose for all posts
type CachedGetAllPostResponse struct {
	Data          []bson.M `json:"data"`
	CurrentPage   int      `json:"currentPage"`
	NumberOfPages float64  `json:"numberOfPages"`
}

// cachedrespose for user profile
type CachedGetUserResponse struct {
	User          UserModel `json:"user"`
	Posts         []bson.M  `json:"posts"`
	CurrentPage   int       `json:"currentPage"`
	NumberOfPages float64   `json:"numberOfPages"`
}
