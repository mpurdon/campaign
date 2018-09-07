package main

import (
	"fmt"
	"github.com/jinzhu/gorm"

	pb "github.com/mpurdon/gomicro-example/user-service/proto/user"
	"reflect"
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
	orm *gorm.DB
}

/*
func toUserMessage(source UserModel) *pb.User {

	model := pb.User{}

	modelValue := reflect.ValueOf(model)
	sourceValue := reflect.ValueOf(source)

	modelType := modelValue.Type().Elem()
	modelValue = modelValue.Elem()

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelType.FieldByIndex([]int{i})
		fieldValue := sourceValue.FieldByName(field.Name)

		reflect.ValueOf(model).Elem().FieldByName(field.Name).Set(fieldValue)
	}

	return &model
}
*/

func toUserModel(message *pb.User) UserModel {

	model := UserModel{}
	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()

	fmt.Printf("Reflected message: %+v\n", message)

	for i := 0; i < modelValue.NumField(); i++ {
		f := modelType.FieldByIndex([]int{i})
		fmt.Printf("\nHandling field %+v (%s)\n", f.Name, f.Type)

		messageField := reflect.ValueOf(message).Elem().FieldByName(f.Name)
		if !messageField.IsValid() {
			fmt.Printf("Message does not have the field")
			continue
		}
		messageFieldValue := reflect.ValueOf(messageField)
		fmt.Printf("message value for '%s': %+v\n", f.Name, messageFieldValue)

		modelField := reflect.ValueOf(&model).Elem().FieldByName(f.Name)
		modelField.Set(messageField)
	}

	return model
}

/**
 * Create a database record for the given account
 */
func (repo *UserRepository) Create(user *pb.User) error {
	Logger.Infof("inserting account %s into database", user.Name)

	userModel := toUserModel(user)

	if err := repo.orm.Create(&userModel).Error; err != nil {
		Logger.Errorf("query error adding account: %v", err)
		return err
	}

	Logger.Infof("added account: %v", userModel)
	return nil
}

/**
 * Get a account
 */
func (repo *UserRepository) Get(user *pb.User) (*pb.User, error) {
	Logger.Infof("Getting a account from the database.")

	model := toUserModel(user)

	if err := repo.orm.First(model).Error; err != nil {
		Logger.Errorf("query error getting account: %s", err)
		return nil, err
	}

	return user, nil
}

/**
 * Get a account
 */
func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	Logger.Infof("Getting a account from the database by email %s.", email)

	user := &pb.User{}
	if err := repo.orm.Where("email = ?", email).First(&user).Error; err != nil {
		Logger.Errorf("query error getting account: %s", err)
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

	if err := repo.orm.Find(&users).Error; err != nil {
		Logger.Errorf("query error getting all users: %s", err)
		return nil, err
	}

	return users, nil
}

/**
 * Authenticate a account
 */
func (repo *UserRepository) Auth() ([]*pb.User, error) {
	Logger.Infof("Authenticating a account.")

	var users []*pb.User

	if err := repo.orm.Find(&users).Error; err != nil {
		Logger.Errorf("query error getting all users: %s", err)
		return nil, err
	}

	return users, nil
}
