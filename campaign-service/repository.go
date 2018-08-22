package main

import (
	"github.com/jinzhu/gorm"
	pb "github.com/mpurdon/gomicro-example/campaign-service/proto/campaign"
)

const (
	dbName             = "fc"
	campaignCollection = "campaigns"
)

type Repository interface {
	Create(*pb.Campaign) (*pb.Campaign, error)
	GetAll() ([]*pb.Campaign, error)
}

type CampaignRepository struct {
	db *gorm.DB
}

/**
 * Create a database record for the given campaign
 */
func (repo *CampaignRepository) Create(campaign *pb.Campaign) (*pb.Campaign, error) {
	Logger.Infof("Creating campaign %s", campaign.Name)

	if err := repo.db.Create(campaign).Error; err != nil {
		Logger.Errorf("query error adding campaign: %s", err)
		return nil, err
	}

	Logger.Infof("Added campaign with id: %d", campaign.Id)
	return campaign, nil
}

/**
 * Get all campaigns
 */
func (repo *CampaignRepository) GetAll() ([]*pb.Campaign, error) {
	Logger.Infof("Getting all campaigns from the database.")

	var campaigns []*pb.Campaign

	if err := repo.db.Find(&campaigns).Error; err != nil {
		Logger.Errorf("query error getting all campaigns: %s", err)
		return nil, err
	}

	return campaigns, nil
}
