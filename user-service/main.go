package main

import (
	"os"

	"github.com/micro/go-micro"

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
	defer db.Close()

	repo := &UserRepository{}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags.
	srv.Init()

	// Register handler
	pb.RegisterUserServiceHandler(srv.Server(), &service{repo})

	// Run the server
	if err := srv.Run(); err != nil {
		Logger.Fatal(err)
	}
}
