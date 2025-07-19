package servegrpc

import (
	"context"
	"log"
	"realTimeChat/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetFollowingFollowersClient(id string) ([]*protos.UserIDsList, error) {
	conn, err := grpc.NewClient(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connected:%v", err)
		return nil, err
	}
	defer conn.Close()

	client := protos.NewRealtimeChatServiceClient(conn)
	// call grpc method
	ctx := context.Background()
	req := &protos.UserID{Userid: id}
	resp, err := client.GetUserFollowingFollowers(ctx, req)
	if err != nil {
		log.Printf("error with caliing getuserfollwingfollowers grpc: %v", err)
		return nil, err
	}

	return resp.GetUserIDsLists(), nil

}

func SendMessageClient(sender, receiver, content string) error {
	conn, err := grpc.NewClient(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connected:%v", err)
		return err
	}
	defer conn.Close()
	client := protos.NewRealtimeChatServiceClient(conn)
	ctx := context.Background()

	req := &protos.MessageRequest{
		Sender:   sender,
		Receiver: receiver,
		Content:  content,
	}

	_, err = client.SendMessage(ctx, req)
	if err != nil {
		log.Printf("error while calling sendmessage grpc %v", err)
		return err
	}

	return nil

}
