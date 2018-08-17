package main

import (
	"github.com/micro/go-micro"
	pb "github.com/mpurdon/gomicro-example/campaign-service/proto/campaign"
	"golang.org/x/net/context"

	"fmt"
	venueProto "github.com/mpurdon/gomicro-example/venue-service/proto/venue"
	//venueProto "../venue-service/proto/venue"
	"log"
	"os"
)

const (
	port = ":50051"
)

type Repository interface {
	Create(*pb.Campaign) (*pb.Campaign, error)
	GetAll() []*pb.Campaign
}

// CampaignRepository - Dummy repository, this simulates the use of a data store
// of some kind. We'll replace this with a real implementation later on.
type CampaignRepository struct {
	campaigns []*pb.Campaign
}

func (repo *CampaignRepository) Create(campaign *pb.Campaign) (*pb.Campaign, error) {
	updated := append(repo.campaigns, campaign)
	repo.campaigns = updated
	return campaign, nil
}

func (repo *CampaignRepository) GetAll() []*pb.Campaign {
	return repo.campaigns
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo        Repository
	venueClient venueProto.VenueServiceClient
}

func (s *service) CreateCampaign(ctx context.Context, req *pb.Campaign, res *pb.Response) error {

	response, err := s.venueClient.FindAvailable(context.Background(), &venueProto.VenueSpecification{
		Location: req.Location,
		Capacity: int32(req.Ca),
	})

	log.Printf("Found vessel: %s \n", response.Venue.Name)
	if err != nil {
		return err
	}

	// We set the VesselId as the vessel we got back from our
	// vessel service
	req.VenueId = response.Venue.Id

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

	fmt.Println("Running with environment...")
	for _, e := range os.Environ() {
		//pair := strings.Split(e, "=")
		fmt.Println(e)
	}

	repo := &CampaignRepository{}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.campaign"),
		micro.Version("latest"),
	)

	venueClient := venueProto.NewVenueServiceClient("go.micro.srv.venue", srv.Client())

	// Init will parse the command line flags.
	srv.Init()

	// Register handler
	pb.RegisterCampaignServiceHandler(srv.Server(), &service{repo, venueClient})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
