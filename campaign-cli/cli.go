package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/micro/go-micro/cmd"
	"golang.org/x/net/context"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	pb "github.com/mpurdon/gomicro-example/campaign-service/proto/campaign"
	"os"
)

var campaignFilenames = [...]string{
	"campaign_001.json",
	"campaign_001.json",
}

/**
 * Parse the JSON file.
 */
func parseFile(file string) (*pb.Campaign, error) {
	var campaign *pb.Campaign

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &campaign)

	return campaign, err
}

func main() {

	cmd.Init()

	// Create new JWT token based client
	client := pb.NewCampaignServiceClient("fc.campaign", microclient.DefaultClient)

	if len(os.Args) < 2 {
		Logger.Fatal(errors.New("not enough arguments, expecting token"))
	}

	var token string
	token = os.Args[1]

	// Create a new context which contains our given token.
	// This same context will be passed into both the calls we make
	// to our campaign-service.
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

	dataFile := ""

	for _, campaignFile := range campaignFilenames {

		dataFile = "data/" + campaignFile
		Logger.Infof("Loading campaign data from file: %s", dataFile)
		campaign, err := parseFile(dataFile)

		if err != nil {
			Logger.Fatalf("Could not parse data file: %v", err)
		}

		r, err := client.CreateCampaign(ctx, campaign)
		if err != nil {
			Logger.Fatalf("Could not create campaign: %v", err)
		}
		Logger.Infof("Created campaign: %t", r.Created)
	}

	response, err := client.GetCampaigns(ctx, &pb.GetRequest{})
	if err != nil {
		Logger.Fatalf("Could not list campaigns: %v", err)
	}
	for _, v := range response.Campaigns {
		Logger.Info(v)
	}
}
