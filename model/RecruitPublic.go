package model

type EnumType int8

type RecruitPublic struct {
	Mission   string `gorm:"primary_key;column:mission" json:"mission"`
	IsOpen    bool   `gorm:"column:isopen" json:"isopen"`
	IsSpecial bool   `gorm:"column:isspecial" json:"isspecial"`
	Rank      int    `gorm:"column:rank" json:"rank"`
}

func (RecruitPublic) TableName() string { return "recruitPublic" }
