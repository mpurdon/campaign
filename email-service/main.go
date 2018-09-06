package main

import (
	"github.com/micro/go-micro"
	"golang.org/x/net/context"

	pb "github.com/mpurdon/gomicro-example/user-service/proto/user"
)

const topic = "user.created"

type Subscriber struct{}

func (sub *Subscriber) Process(ctx context.Context, user *pb.User) error {
	Logger.Info("picked up a new message...")
	Logger.Infof("Sending email to:", user.Name)

	return nil
}

func main() {
	srv := micro.NewService(
		micro.Name("fc.email"),
		micro.Version("latest"),
	)

	srv.Init()

	// Subscribe to messages from the Broker
	micro.RegisterSubscriber(topic, srv.Server(), new(Subscriber))

	// Run the server
	if err := srv.Run(); err != nil {
		Logger.Errorf("could not run the server: %+v", err)
	}

}
