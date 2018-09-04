package main

import (
	"encoding/json"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"

	pb "github.com/mpurdon/gomicro-example/user-service/proto/user"

	_ "github.com/micro/go-plugins/broker/nats"
)

const topic = "user.created"

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.srv.email"),
		micro.Version("latest"),
	)

	srv.Init()

	pubsub := srv.Server().Options().Broker
	if err := pubsub.Connect(); err != nil {
		Logger.Fatalf("could not connect to broker: %+v", err)
	}

	// Subscribe to messages from the Broker
	_, err := pubsub.Subscribe(topic, func(p broker.Publication) error {
		var user *pb.User
		if err := json.Unmarshal(p.Message().Body, &user); err != nil {
			return err
		}

		Logger.Info("subscribing to %s events", topic)

		go sendEmail(user)
		return nil
	})

	// Run the server
	if err := srv.Run(); err != nil {
		Logger.Errorf("could not run the server: %+v", err)
	}

	if err != nil {
		Logger.Errorf("Could not subscribe to %s, %+v", topic, err)
	}

}

/**
 * Send email to a given user
 */
func sendEmail(user *pb.User) error {
	Logger.Infof("Sending email to: %s", user.Name)
	return nil
}
