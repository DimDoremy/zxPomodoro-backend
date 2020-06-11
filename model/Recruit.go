package model

type EnumType int8

type RecruitPublic struct {
	Title     string `gorm:"column:title" json:"title"`
	Describe  string `gorm:"column:description" json:"describe"`
	Mission   string `gorm:"primary_key;column:mission" json:"mission"`
	IsOpen    bool   `gorm:"column:isopen" json:"isopen"`
	IsSpecial bool   `gorm:"column:isspecial" json:"isspecial"`
	Rank      int    `gorm:"column:rank" json:"rank"`
}

func (RecruitPublic) TableName() string { return "recruitPublic" }

// 个人招募结构体
type RequestPersonal struct {
	Title    string `json:"title"`
	Describe string `json:"describe"`
	Url      string `json:"url"`
	Openid   string `json:"openid"`
	Makerid  string `json:"makerid"`
	Time     int64  `json:"time"`
}

// 公共招募结构体
type RequestPublic struct {
	Openid    string `json:"openid"`
	Makerid   string `json:"makerid"`
	Score     int64  `json:"score"`
	IsSpecial bool   `json:"is_special"`
}
