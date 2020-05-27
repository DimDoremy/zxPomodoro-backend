package model

type UserPiece struct {
	Openid         string `gorm:"primary_key;column:openid" json:"openid"`
	XuanYuanJian   int	`gorm:"column:xuanyuanjian" json:"xuanyuanjian"`
	DongHuangZhong int	`gorm:"column:donghuangzhong" json:"donghuangzhong"`
	PanGuFu        int	`gorm:"column:pangufu" json:"pangufu"`
	LianYaoHu      int	`gorm:"column:lianyaohu" json:"lianyaohu"`
	HaoTianTa      int	`gorm:"column:haotianta" json:"haotianta"`
	FuXiQin        int	`gorm:"column:fuxiqin" json:"fuxiqin"`
	ShenNongDing   int	`gorm:"column:shennongding" json:"shennongding"`
	KongTongYin    int	`gorm:"column:kongtongyin" json:"kongtongyin"`
	KunLunJing     int	`gorm:"column:kunlunjing" json:"kunlunjing"`
	NvWaShi        int	`gorm:"column:nvwashi" json:"nvwashi"`
}

func (UserPiece)TableName() string {
	return "userPiece"
}