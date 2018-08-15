package main

import (
	pb "./proto/campaign"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"

	"fmt"
)

const (
	port = ":50051"
)

type IRepository interface {
	Create(*pb.Campaign) (*pb.Campaign, error)
	GetAll() []*pb.Campaign
}

// Repository - Dummy repository, this simulates the use of a data store
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	campaigns []*pb.Campaign
}

func (repo *Repository) Create(campaign *pb.Campaign) (*pb.Campaign, error) {
	updated := append(repo.campaigns, campaign)
	repo.campaigns = updated
	return campaign, nil
}

func (repo *Repository) GetAll() []*pb.Campaign {
	return repo.campaigns
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo IRepository
}

func (s *service) CreateCampaign(ctx context.Context, req *pb.Campaign, res *pb.Response) error {

	// Save our campaign
	campaign, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Campaign = campaign

	return nil
}

func (s *service) GetCampaigns(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	campaigns := s.repo.GetAll()
	res.Campaigns = campaigns

	return nil
}

/**
 * Main
 */
func main() {

	repo := &Repository{}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.campaign"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags.
	srv.Init()

	// Register handler
	pb.RegisterCampaignServiceHandler(srv.Server(), &service{repo})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
