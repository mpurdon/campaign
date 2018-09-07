package main

import (
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"

	microclient "github.com/micro/go-micro/client"
	pb "github.com/mpurdon/gomicro-example/user-service/proto/user"
)

func main() {

	// Create new account client
	client := pb.NewUserServiceClient(
		"fc.account",
		microclient.DefaultClient,
	)

	// Define our flags
	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:  "name",
				Usage: "User full name",
			},
			cli.StringFlag{
				Name:  "email",
				Usage: "User email address",
			},
			cli.StringFlag{
				Name:  "password",
				Usage: "User password",
			},
		),
	)

	// Start as service
	service.Init(

		micro.Action(func(c *cli.Context) {

			/*
				name := c.String("name")
				email := c.String("email")
				password := c.String("password")
			*/

			name := "Matthew Purdon"
			email := "mdjpurdon@gmail.com"
			password := "password"

			Logger.Infof("creating account from arguments '%s<%s>'", name, email)

			user := &pb.User{
				Name:     name,
				Email:    email,
				Password: password,
			}
			Logger.Infof("attempting to create account: %v", user)

			// Call our account service
			r, err := client.Create(context.TODO(), user)

			if err != nil {
				Logger.Fatalf("Could not create account: %v", err)
			}

			Logger.Infof("Created account with id: %d", r.User.Id)

			Logger.Info("attempting to get all users:")
			response, err := client.GetAll(context.Background(), &pb.Request{})
			if err != nil {
				Logger.Fatalf("Could not list users: %v", err)
			}

			for _, v := range response.Users {
				Logger.Info(v)
			}

			login := &pb.User{
				Email:    email,
				Password: password,
			}
			Logger.Infof("attempting to log in %v", login)

			authResponse, err := client.Auth(context.TODO(), login)
			if err != nil {
				Logger.Fatalf("could not authenticate account: %s, error: %v\n", email, err)
			}

			Logger.Infof("your access token is: %s \n", authResponse.Token)

			os.Exit(0)
		}),
	)

	// Run the service
	if err := service.Run(); err != nil {
		Logger.Error(err)
	}

	os.Exit(0)
}
