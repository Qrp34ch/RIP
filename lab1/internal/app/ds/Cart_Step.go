package ds

type CartStep struct {
	ID uint `gorm:"primaryKey"`
	// здесь создаем Unique key, указывая общий uniqueIndex
	CartID uint `gorm:"not null;uniqueIndex:idx_cart_step"`
	StepID uint `gorm:"not null;uniqueIndex:idx_cart_step"`

	Cart Cart `gorm:"foreignKey:CartID"`
	Step Step `gorm:"foreignKey:StepID"`
}
