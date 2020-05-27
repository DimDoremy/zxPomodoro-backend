package model

type RecruitBlock struct {
	MakerId  string `gorm:"primary_key;column:makerid" json:"makerid"`
	BossId   int	`gorm:"column:bossid" json:"bossid"`
	BossLife int	`gorm:"column:bosslife" json:"bosslife"`
}

func (RecruitBlock)TableName() string {
	return "recruitBlock"
}