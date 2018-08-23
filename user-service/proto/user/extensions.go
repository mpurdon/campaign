package go_micro_srv_user

import (
	"github.com/jinzhu/gorm"
	//"github.com/google/uuid"
	"github.com/satori/go.uuid"
)

func (model *User) BeforeCreate(scope *gorm.Scope) error {
	guid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	return scope.SetColumn("Guid", guid.String())
}
