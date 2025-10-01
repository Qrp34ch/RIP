package ds

type SynthesisReaction struct {
	ID uint `gorm:"primaryKey"`
	// здесь создаем Unique key, указывая общий uniqueIndex
	SynthesisID uint `gorm:"not null;uniqueIndex:idx_synthesis_reaction"`
	ReactionID  uint `gorm:"not null;uniqueIndex:idx_synthesis_reaction"`

	Synthesis Synthesis `gorm:"foreignKey:SynthesisID"`
	Reaction  Reaction  `gorm:"foreignKey:ReactionID"`
}
