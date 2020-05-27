package model

type RecruitList struct {
	Makerid string	`gorm:"primary_key;column:makerid" json:"makerid"`
	Gamer1  string	`gorm:"column:gamer1" json:"gamer1"`
	Gamer2  string	`gorm:"column:gamer2" json:"gamer2"`
	Gamer3  string	`gorm:"column:gamer3" json:"gamer3"`
}

func (RecruitList) TableName() string { return "recruitList" }
