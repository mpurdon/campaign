package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/micro/go-micro/cmd"
	"golang.org/x/net/context"

	microclient "github.com/micro/go-micro/client"
	pb "github.com/mpurdon/gomicro-example/campaign-service/proto/campaign"
)

const (
	defaultFilename = "campaign.json"
)

/**
 * Parse the JSON file.
 */
func parseFile(file string) ([]*pb.Campaign, error) {
	var campaigns []*pb.Campaign

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &campaigns)

	return campaigns, err
}

func main() {

	cmd.Init()

	client := pb.NewCampaignServiceClient("go.micro.srv.campaign", microclient.DefaultClient)

	// Contact the server and print out its response.
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	campaigns, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	for _, campaign := range campaigns {
		r, err := client.CreateCampaign(context.TODO(), campaign)
		if err != nil {
			log.Fatalf("Could not create: %v", err)
		}
		log.Printf("Created: %t", r.Created)
	}

	getAll, err := client.GetCampaigns(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list campaigns: %v", err)
	}

	for _, v := range getAll.Campaigns {
		log.Println(v)
	}
}
