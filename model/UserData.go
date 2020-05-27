package model

type UserData struct {
	Openid     string `gorm:"primary_key;column:openid" json:"openid"`
	Nickname   string `gorm:"column:nickname" json:"nickname"`
	AvatarUrl  string `gorm:"column:avatarurl" json:"avatarurl"`
	Time       int	`json:"time"`
	Soul       int	`json:"soul"`
	LifePoint  int	`gorm:"column:lifepoint" json:"lifepoint"`
	Experience int	`json:"experience"`
}

func (UserData)TableName() string {
	return "userData"
}