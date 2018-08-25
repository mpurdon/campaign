package main

import (
	"github.com/jinzhu/gorm"

	pb "github.com/mpurdon/gomicro-example/user-service/proto/user"
)

const (
	dbName         = "fc"
	userCollection = "users"
)

type Repository interface {
	Create(*pb.User) error
	Get(*pb.User) (*pb.User, error)
	GetByEmail(email string) (*pb.User, error)
	GetAll() ([]*pb.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

/**
 * Create a database record for the given user
 */
func (repo *UserRepository) Create(user *pb.User) error {
	Logger.Infof("Creating user %s", user.Name)

	if err := repo.db.Create(user).Error; err != nil {
		Logger.Errorf("query error adding user: %s", err)
		return err
	}

	Logger.Infof("Added user with id: %d", user.Id)
	return nil
}

/**
 * Get a user
 */
func (repo *UserRepository) Get(user *pb.User) (*pb.User, error) {
	Logger.Infof("Getting a user from the database.")

	if err := repo.db.First(&user).Error; err != nil {
		Logger.Errorf("query error getting user: %s", err)
		return nil, err
	}

	return user, nil
}

/**
 * Get a user
 */
func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	Logger.Infof("Getting a user from the database by email %s.", email)

	user := &pb.User{}
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		Logger.Errorf("query error getting user: %s", err)
		return nil, err
	}

	return user, nil
}

/**
 * Get all users
 */
func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	Logger.Infof("Getting all users from the database.")

	var users []*pb.User

	if err := repo.db.Find(&users).Error; err != nil {
		Logger.Errorf("query error getting all users: %s", err)
		return nil, err
	}

	return users, nil
}

/**
 * Authenticate a user
 */
func (repo *UserRepository) Auth() ([]*pb.User, error) {
	Logger.Infof("Authenticating a user.")

	var users []*pb.User

	if err := repo.db.Find(&users).Error; err != nil {
		Logger.Errorf("query error getting all users: %s", err)
		return nil, err
	}

	return users, nil
}
