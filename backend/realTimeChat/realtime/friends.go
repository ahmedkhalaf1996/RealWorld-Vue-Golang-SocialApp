package realtime

import (
	"fmt"
	"realTimeChat/servegrpc"
)

func GetUserFriends(userID string) <-chan []string {
	ch := make(chan []string)
	go func() {
		defer close(ch)

		// call grp client fun
		userFriends, err := servegrpc.GetFollowingFollowersClient(userID)

		if err != nil {
			fmt.Printf("Erro On Frinds Method Realtime")
			ch <- []string{}
			return
		}

		var friends []string
		for _, userIDsList := range userFriends {
			friends = append(friends, userIDsList.UserIdsList...)
		}

		ch <- friends
	}()
	return ch

}

// func GetUserFriends(userID string) <-chan []string {
// 	ch := make(chan []string)
// 	go func() {
// 		defer close(ch)
// 		switch userID {
// 		case "1":
// 			ch <- []string{"2", "3", "4"}
// 		case "2":
// 			ch <- []string{"1", "3", "4"}
// 		case "3":
// 			ch <- []string{"1", "2", "4"}
// 		case "4":
// 			ch <- []string{"1"}
// 		default:
// 			ch <- []string{}
// 		}
// 	}()
// 	return ch

// }
