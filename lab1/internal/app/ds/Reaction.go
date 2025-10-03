package ds

type Reaction struct {
	ID               int    `gorm:"primaryKey"`
	Title            string `gorm:"not null"`
	Src              string
	SrcUr            string
	Details          string
	IsDelete         bool    `gorm:"type:boolean not null;default:false"`
	StartingMaterial string  `gorm:"not null"`
	DensitySM        float32 `gorm:"not null"`
	//VolumeSM         float32
	MolarMassSM    int     `gorm:"not null"`
	ResultMaterial string  `gorm:"not null"`
	DensityRM      float32 `gorm:"not null"`
	//VolumeRM         float32
	MolarMassRM int `gorm:"not null"`
}
