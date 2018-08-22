package main

import (
	"golang.org/x/net/context"

	pb "github.com/mpurdon/gomicro-example/user-service/proto/user"
	"google.golang.org/grpc/grpclog"
)

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo         Repository
	tokenService Authable
}

/**
 * Get user handler
 */
func (s *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {

	if err := s.repo.db.First(&req).Error; err != nil {
		grpclog.Errorf("query error getting user: %v", err)
		return err
	}

	return nil
}

/**
 * Get a campaign by its GUID
func (repo *CampaignRepository) GetCampaign(guid string) (*pb.Campaign, error) {
    grpclog.Infof("Getting campaign with GUID %s from the database.", guid)

    var campaign *pb.Campaign
    campaign.Guid = guid

    if err := repo.db.First(&campaign).Error; err != nil {
        grpclog.Errorf("query error getting campaign: %s", err)
        return nil, err
    }

    return campaign, nil
}
*/
