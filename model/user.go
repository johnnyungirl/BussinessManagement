package model

import "time"

// import "gorm.io/gorm"

//User -> model for users table
type User struct {
	ID          uint       `gorm:"primarykey;autoIncrement:true"`
	CreateDate  *time.Time `gorm:"type:timestamp without time zone;default:null"`
	UpdateDate  *time.Time `gorm:"type:timestamp without time zone;default:null"`
	IsDelete    bool       `gorm:"type:boolean"`
	Name        string     `gorm:"type:varchar(100)" json:"name" `
	PhoneNumber string     `gorm:"type:varchar(100)" json:"phonenumber"`
	Address     string     `gorm:"type:varchar(100)" json:"address"`
	Website     string     `gorm:"type:varchar(100);" json:"website"`
	State       string     `gorm:"type:varchar(100)" json:"state"`
	Role        string     `gorm:"type:varchar(100)" json:"role"`
	Password    string     `gorm:"type:varchar(500)" json:"password"`
	Email       string     `gorm:"type:varchar(100);unique" json:"email"`
}

//TableName --> Table for Product Model
func (User) TableName() string {
	return "users"
}
