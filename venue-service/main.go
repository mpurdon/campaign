package main

import (
	"github.com/micro/go-micro"
	pb "github.com/mpurdon/gomicro-example/venue-service/proto/venue"
)

func main() {
	// Ensure that all log messages are written on shutdown
	defer Logger.Sync()

	venues := []*pb.Venue{
		&pb.Venue{Id: "venue_001", Name: "First Venue", Location: "Toronto", Capacity: 10, UserId: "user_0001"},
		&pb.Venue{Id: "venue_002", Name: "Second Venue", Location: "Toronto", Capacity: 200, UserId: "user_0001"},
	}

	repo := &VenueRepository{venues: venues}

	srv := micro.NewService(
		micro.Name("go.micro.srv.venue"),
		micro.Version("latest"),
	)

	srv.Init()

	pb.RegisterVenueServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		Logger.Error(err)
	}
}
