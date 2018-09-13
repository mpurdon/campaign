package main

import (
	"errors"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"

	"fmt"
	microerrors "github.com/micro/go-micro/errors"
	pb "github.com/mpurdon/gomicro-example/account-service/proto/account"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo         Repository
	tokenService Authable
	Publisher    micro.Publisher
}

/*
 * Create a account.
 */
func (s *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {

	Logger.Infof("Creating account: %v", req)

	// Generates a hashed version of our password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(fmt.Sprintf("error hashing password: %v", err))
	}

	req.Password = string(hashedPass)

	if err := s.repo.Create(req); err != nil {
		return errors.New(fmt.Sprintf("error creating account: %v", err))
	}

	token, err := s.tokenService.Encode(req)
	if err != nil {
		return err
	}

	res.User = req
	res.Token = &pb.Token{Token: token}

	// Publish message to broker
	if err := s.Publisher.Publish(ctx, req); err != nil {
		Logger.Errorf("publishing account creation failed: %+v", err)
		return err
	}

	return nil
}

/**
 * Get account handler
 */
func (s *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {

	user, err := s.repo.Get(req)
	if err != nil {
		Logger.Errorf("query error getting account: %v", err)
		return err
	}

	res.User = user

	return nil
}

/**
 * Get all users handler
 */
func (s *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := s.repo.GetAll()

	if err != nil {
		return err
	}
	res.Users = users

	return nil
}

/**
 * Get a campaign by its GUID
func (repo *CampaignRepository) GetCampaign(guid string) (*pb.Campaign, error) {
    Logger.Infof("Getting campaign with GUID %s from the database.", guid)

    var campaign *pb.Campaign
    campaign.Guid = guid

    if err := repo.orm.First(&campaign).Error; err != nil {
        Logger.Errorf("query error getting campaign: %s", err)
        return nil, err
    }

    return campaign, nil
}
*/

/**
 * Handle Authentication
 */
func (s *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {

	// tracing
	tr := trace.New("api.v1", "Hotel.Rates")
	defer tr.Finish()

	Logger.Infof("ctx: %+v", ctx)
	Logger.Infof("req: %+v", req)
	Logger.Infof("Logging in with: %s|%s", req.Email, req.Password)

	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return microerrors.BadRequest("Account.Auth", err.Error())
	}

	// Compares our given password against the hashed password
	// stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return microerrors.Forbidden("Account.Auth", err.Error())
	}

	token, err := s.tokenService.Encode(user)
	if err != nil {
		return microerrors.InternalServerError("Account.Auth", err.Error())
	}
	res.Token = token

	return nil
}

/**
 * Validate a given token
 */
func (s *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {

	// Decode token
	claims, err := s.tokenService.Decode(req.Token)

	if err != nil {
		return err
	}

	if claims.User.Id == "" {
		return errors.New("invalid account")
	}

	res.Valid = true

	return nil
}
