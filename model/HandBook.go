package model

type HandBook struct {
	Openid   string `gorm:"primary_key;column:openid" json:"openid"`
	Monster1 int	`gorm:"column:monster1" json:"monster1"`
	Monster2 int	`gorm:"column:monster2" json:"monster2"`
	Monster3 int	`gorm:"column:monster3" json:"monster3"`
	Monster4 int	`gorm:"column:monster4" json:"monster4"`
	Monster5 int	`gorm:"column:monster5" json:"monster5"`
	Monster6 int	`gorm:"column:monster6" json:"monster6"`
	Monster7 int	`gorm:"column:monster7" json:"monster7"`
}

func (HandBook)TableName() string {
	return "handBook"
}