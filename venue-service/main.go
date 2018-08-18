package main

import (
	//pb "./proto/venue"
	"context"
	"errors"
	"fmt"
	"github.com/micro/go-micro"
	pb "github.com/mpurdon/gomicro-example/venue-service/proto/venue"
)

type Repository interface {
	FindAvailable(specification *pb.VenueSpecification) (*pb.Venue, error)
}

type VenueRepository struct {
	venues []*pb.Venue
}

func (repo *VenueRepository) FindAvailable(spec *pb.VenueSpecification) (*pb.Venue, error) {

	fmt.Printf("Attempting to find venue in %s with capacity of at least %d\n", spec.Location, spec.Capacity)

	for _, venue := range repo.venues {
		if spec.Capacity <= venue.Capacity && spec.Location == venue.Location {
			return venue, nil
		}
	}

	return nil, errors.New("no venue matches the given specifications")
}

type service struct {
	repo Repository
}

func (s *service) FindAvailable(ctx context.Context, req *pb.VenueSpecification, res *pb.Response) error {
	venue, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}

	res.Venue = venue
	return nil
}

func main() {
	venues := []*pb.Venue{
		&pb.Venue{Id: "venue_001", Name: "First Venue", Location: "Toronto", Capacity: 8, UserId: "user_0001"},
	}

	repo := &VenueRepository{venues: venues}

	srv := micro.NewService(
		micro.Name("go.micro.srv.venue"),
		micro.Version("latest"),
	)

	srv.Init()

	pb.RegisterVenueServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
