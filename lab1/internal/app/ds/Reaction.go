package ds

type Reaction struct {
	ID               int `gorm:"primaryKey"`
	Title            string
	Src              string
	SrcUr            string
	Details          string
	IsDelete         bool `gorm:"type:boolean not null;default:false"`
	StartingMaterial string
	DensitySM        float32
	VolumeSM         float32
	MolarMassSM      int
	ResultMaterial   string
	DensityRM        float32
	VolumeRM         float32
	MolarMassRM      int
}
