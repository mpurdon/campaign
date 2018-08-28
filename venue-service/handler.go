package main

import (
	"context"
	pb "github.com/mpurdon/gomicro-example/venue-service/proto/venue"
)

type service struct {
	repo Repository
}

func (s *service) FindAvailable(ctx context.Context, req *pb.VenueSpecification, res *pb.Response) error {
	venue, err := s.repo.FindAvailable(req)
	if err != nil {
		Logger.Error(err)
		return err
	}

	res.Venue = venue
	return nil
}
