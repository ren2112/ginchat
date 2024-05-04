package models

import (
	"ginchat/utils"
	"gorm.io/gorm"
)

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int //对应owner和target的关系类型 1好友 2群友 3
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	utils.DB.Where("owner_id=? and type=1", userId).Find(&contacts)
	objIds := make([]uint64, 0)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}

func AddFriend(userId uint, targetId uint) int {
	user := UserBasic{}
	if targetId != 0 {
		user = FindByID(targetId)
		if user.Salt != "" {
			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = targetId
			contact.Type = 1
			utils.DB.Create(&contact)
			return 0
		}
		return -1
	}
	return -1
}
