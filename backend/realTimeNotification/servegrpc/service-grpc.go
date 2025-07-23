package servegrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "realTimeNotification/protos"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type notificationServer struct {
	pb.UnimplementedNotificationGrpcServiceServer
	wsMu *sync.Mutex
	ws   map[string]*websocket.Conn
}

type Notification struct {
	ID        string    `json:"_id"`
	Details   string    `json:"details"`
	MainUID   string    `json:"mainuid"`
	TargetID  string    `json:"targetid"`
	IsReaded  bool      `json:"isreded"`
	CraetedAt time.Time `json:"createdAt"`
	User      User      `json:"user"`
}

type User struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (s *notificationServer) SendGrpcNotification(ctx context.Context, req *pb.NotificationGrpcRequest) (*empty.Empty, error) {

	fmt.Printf("Sending notifjcation to user %s : %s\n", req.Mainuid, req.Deatils)

	// send the notification to websocket server
	s.wsMu.Lock()
	defer s.wsMu.Unlock()

	if conn, ok := s.ws[req.Mainuid]; ok {
		notification := Notification{
			ID:        req.XId,
			MainUID:   req.Mainuid,
			Details:   req.Deatils,
			TargetID:  req.Targetid,
			IsReaded:  req.Isreded,
			CraetedAt: time.Unix(req.CreatedAt.Seconds, 0),
			User: User{
				Name:   req.User.Name,
				Avatar: req.User.Avatar,
			},
		}
		err := conn.WriteJSON(notification)
		if err != nil {
			log.Printf("Error sending notifncation to websocket server: %v", err)
		}
	}
	return &empty.Empty{}, nil
}

func StartGRPCServer(ws map[string]*websocket.Conn, wsMu *sync.Mutex) error {
	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		return fmt.Errorf("Faild to listen on port 8090: %v", err)
	}

	grpcServer := grpc.NewServer()
	notificationService := &notificationServer{ws: ws, wsMu: wsMu}

	pb.RegisterNotificationGrpcServiceServer(grpcServer, notificationService)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("faild to server gRPC server : %v", err)
		}
	}()

	return nil
}
