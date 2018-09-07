package main

import (
	"errors"
	"golang.org/x/net/context"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"

	accountService "github.com/mpurdon/gomicro-example/account-service/proto/account"
	pb "github.com/mpurdon/gomicro-example/campaign-service/proto/campaign"
	venueProto "github.com/mpurdon/gomicro-example/venue-service/proto/venue"
)

/**
 * Wraps a handler function to provide authentication
 */
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}
		token := meta["Token"]
		Logger.Infof("Authenticating with token: %s", token)

		// Auth here
		authClient := accountService.NewAccountClient("fc.account", client.DefaultClient)
		authResp, err := authClient.ValidateToken(ctx, &accountService.Token{
			Token: token,
		})
		Logger.Infof("Auth response: %s", authResp)
		Logger.Infof("Error: %v", err)
		if err != nil {
			return err
		}

		err = fn(ctx, req, resp)
		return err
	}
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo        Repository
	venueClient venueProto.VenueServiceClient
}

/**
 * Create campaigns handler
 */
func (s *service) CreateCampaign(ctx context.Context, req *pb.Campaign, res *pb.Response) error {

	Logger.Infof("Attempting to create campaign %s\n", req.Name)

	var capacity int32
	if len(req.Rewards) > 0 {
		capacity = int32(req.Rewards[0].Available)
	}

	response, err := s.venueClient.FindAvailable(context.Background(), &venueProto.VenueSpecification{
		Location: req.Location,
		Capacity: capacity,
	})

	if err != nil {
		Logger.Warnf("Could not find a venue for the campaign: %s", err)
		return err
	}

	Logger.Infof("Found venue: %s, setting campaign venue to %s \n", response.Venue.Name, response.Venue.Id)

	// We set the VesselId as the vessel we got back from our
	// vessel service
	req.VenueId = response.Venue.Id

	// Save our campaign
	campaign, err := s.repo.Create(req)
	if err != nil {
		Logger.Errorf("Could not create campaign: %s", err)
		return err
	}

	res.Created = true
	res.Campaign = campaign

	return nil
}

/**
 * Get campaigns handler
 */
func (s *service) GetCampaigns(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	campaigns, err := s.repo.GetAll()

	if err != nil {
		return err
	}
	res.Campaigns = campaigns

	return nil
}
