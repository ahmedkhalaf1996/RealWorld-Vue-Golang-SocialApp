package servergrpc

import (
	"Server/models"
	"Server/protos"
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client protos.NotificationGrpcServiceClient
}

func NewClient() (*Client, error) {
	// conn, err := grpc.NewClient(":8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// DevOps docker Compose usage
	conn, err := grpc.NewClient("GolangNotifyService:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: protos.NewNotificationGrpcServiceClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) SendGrpcNotification(ctx context.Context, xId, details, mainUID, targetID string,
	isReaded bool, craetedAt time.Time, userName, userAvatar string) error {
	// Prepare the request
	request := &protos.NotificationGrpcRequest{
		XId:       xId,
		Deatils:   details,
		Mainuid:   mainUID,
		Targetid:  targetID,
		Isreded:   isReaded,
		CreatedAt: &timestamp.Timestamp{Seconds: craetedAt.Unix()},
		User: &protos.Usergrpc{
			Name:   userName,
			Avatar: userAvatar,
		},
	}

	// call the grpc client func
	_, err := c.client.SendGrpcNotification(ctx, request)
	if err != nil {
		log.Printf("Faild To send notification : %v", err)
	}
	return err
}

func SendNotification(notification models.Notification) error {
	// client, err := NewClient()
	// if err != nil {
	// 	log.Printf("Faild to create grpc client %v", err)
	// 	return err
	// }
	// defer client.Close()

	// ctx := context.Background()
	// err = client.SendGrpcNotification(ctx,
	// 	notification.ID.Hex(),
	// 	notification.Deatils,
	// 	notification.MainUID,
	// 	notification.TargetID,
	// 	notification.IsReaded,
	// 	notification.CreatedAt,
	// 	notification.User.Name,
	// 	notification.User.Avatart)
	// if err != nil {
	// 	log.Printf("Faild to send notification : %v", err)
	// 	return err
	// }
	return nil
}
