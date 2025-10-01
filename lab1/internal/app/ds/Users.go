package ds

type Users struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Login       string `gorm:"type:varchar(25);unique;not null" json:"login"`
	Password    string `gorm:"type:varchar(100);not null" json:"-"`
	IsModerator bool   `gorm:"type:boolean;default:false" json:"is_moderator"`
	FIO         string `gorm:"type:varchar(200);not null" json:"fio"`
}
