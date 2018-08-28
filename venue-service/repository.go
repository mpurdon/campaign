package main

import (
	"errors"
	pb "github.com/mpurdon/gomicro-example/venue-service/proto/venue"
)

type Repository interface {
	FindAvailable(specification *pb.VenueSpecification) (*pb.Venue, error)
}

type VenueRepository struct {
	venues []*pb.Venue
}

func (repo *VenueRepository) FindAvailable(spec *pb.VenueSpecification) (*pb.Venue, error) {

	Logger.Infof("Attempting to find venue in %s with capacity of at least %d\n", spec.Location, spec.Capacity)

	for _, venue := range repo.venues {
		if spec.Capacity <= venue.Capacity && spec.Location == venue.Location {
			return venue, nil
		}
	}

	return nil, errors.New("no venue matches the given specifications")
}
