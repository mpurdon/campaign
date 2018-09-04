package main

import (
	"github.com/micro/go-micro"
	"os"

	pb "github.com/mpurdon/gomicro-example/user-service/proto/user"
)

/**
 * Main
 */
func main() {

	// Ensure that all log messages are written on shutdown
	defer Logger.Sync()

	Logger.Info("Running with environment...")
	for _, e := range os.Environ() {
		Logger.Info(e)
	}

	db := createConnection()
	Migrate(db)
	defer db.Close()

	repo := &UserRepository{
		orm: db,
	}

	tokenService := &TokenService{repo}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags.
	srv.Init()

	// Get instance of the broker using our defaults
	pubsub := srv.Server().Options().Broker

	// Register handler
	pb.RegisterUserServiceHandler(srv.Server(), &service{repo, tokenService, pubsub})

	// Run the server
	if err := srv.Run(); err != nil {
		Logger.Fatal("could not start the server: %+v\n", err)
	}
}
