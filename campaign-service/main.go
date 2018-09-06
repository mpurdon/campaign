package main

import (
	"os"

	"github.com/micro/go-micro"

	pb "github.com/mpurdon/gomicro-example/campaign-service/proto/campaign"
	venueProto "github.com/mpurdon/gomicro-example/venue-service/proto/venue"
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

	// Automatically migrates the campaign struct
	// into database columns/types etc. This will
	// check for changes and migrate them each time
	// this service is restarted.
	db.AutoMigrate(&pb.Campaign{})

	repo := &CampaignRepository{
		db: db,
	}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("fc.campaign"),
		micro.Version("latest"),
		micro.WrapHandler(AuthWrapper),
	)

	venueClient := venueProto.NewVenueServiceClient("fc.venue", srv.Client())

	// Init will parse the command line flags.
	srv.Init()

	// Register handler
	pb.RegisterCampaignServiceHandler(srv.Server(), &service{repo, venueClient})

	// Run the server
	if err := srv.Run(); err != nil {
		Logger.Fatal(err)
	}
}
