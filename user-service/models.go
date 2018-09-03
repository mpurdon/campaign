package main

/**
 * Models for the User Service
 *
 * @see http://gorm.io/docs/models.html
 * @see https://github.com/gothinkster/golang-gin-realworld-example-app/blob/master/users/models.go
 */

import (
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&UserModel{})
}

type UserModel struct {
	gorm.Model
	Name     string
	Email    string `gorm:"type:varchar(100);unique_index"`
	Password string `gorm:"column:password;not null"`
}

func (UserModel) TableName() string {
	return "users"
}

func (u *UserModel) setPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}

	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.Password = string(passwordHash)

	return nil
}

// Database will only save the hashed string, you should check it by util function.
func (u *UserModel) checkPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)

	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

/*
// You could input the conditions and it will return an UserModel in database with error info.
// 	userModel, err := FindOneUser(&UserModel{Username: "username0"})
func FindOneUser(condition interface{}) (UserModel, error) {
	var model UserModel
	err := db.Where(condition).First(&model).Error

	return model, err
}

// You could input an UserModel which will be saved in database returning with error info
// 	if err := SaveOne(&userModel); err != nil { ... }
func SaveOne(data interface{}) error {
	err := db.Save(data).Error

	return err
}

// You could update properties of an UserModel to database returning with error info.
//  err := db.Model(userModel).Update(UserModel{Username: "wangzitian0"}).Error
func (model *UserModel) Update(data interface{}) error {
	err := db.Model(model).Update(data).Error
	return err
}

// You could add a following relationship as userModel1 following userModel2
// 	err = userModel1.following(userModel2)
func (u UserModel) following(v UserModel) error {
	var follow FollowModel
	err := db.FirstOrCreate(&follow, &FollowModel{
		FollowingID:  v.ID,
		FollowedByID: u.ID,
	}).Error

	return err
}

// You could check whether  userModel1 following userModel2
// 	followingBool = myUserModel.isFollowing(self.UserModel)
func (u UserModel) isFollowing(v UserModel) bool {
	var follow FollowModel
	db.Where(FollowModel{
		FollowingID:  v.ID,
		FollowedByID: u.ID,
	}).First(&follow)

	return follow.ID != 0
}

// You could delete a following relationship as userModel1 following userModel2
// 	err = userModel1.unFollowing(userModel2)
func (u UserModel) unFollowing(v UserModel) error {
	err := db.Where(FollowModel{
		FollowingID:  v.ID,
		FollowedByID: u.ID,
	}).Delete(FollowModel{}).Error

	return err
}

// You could get a following list of userModel
// 	followings := userModel.GetFollowings()
func (u UserModel) GetFollowings() []UserModel {
	tx := db.Begin()
	var follows []FollowModel
	var followings []UserModel
	tx.Where(FollowModel{
		FollowedByID: u.ID,
	}).Find(&follows)
	for _, follow := range follows {
		var userModel UserModel
		tx.Model(&follow).Related(&userModel, "Following")
		followings = append(followings, userModel)
	}
	tx.Commit()

	return followings
}
*/
