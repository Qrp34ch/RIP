package repository

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"lab1/internal/app/ds"
)

func (r *Repository) GetReactions() ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Where("is_delete = ?", false).Find(&reactions).Error
	if err != nil {
		return nil, err
	}
	if len(reactions) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return reactions, nil
}

func (r *Repository) GetReaction(id int) (ds.Reaction, error) {
	reaction := ds.Reaction{}
	err := r.db.Where("id = ? AND is_delete = ?", id, false).First(&reaction).Error
	if err != nil {
		return ds.Reaction{}, err
	}
	return reaction, nil
}

func (r *Repository) GetReactionsByTitle(title string) ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Where("title ILIKE ? AND is_delete = ?", "%"+title+"%", false).Find(&reactions).Error
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (r *Repository) GetReactionsInSynthesis() int64 {
	var synthesisID uint
	var count int64
	creatorID := 1
	err := r.db.Model(&ds.Synthesis{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("id").First(&synthesisID).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.SynthesisReaction{}).Where("synthesis_id = ?", synthesisID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_reactions:", err)
	}

	return count
}

func (r *Repository) GetSynthesis(synthesisID uint) ([]ds.Reaction, error) {
	var reactionID uint
	var reaction ds.Reaction
	var synthesisReaction []ds.SynthesisReaction
	var result []ds.Reaction
	err := r.db.Where("synthesis_id = ?", synthesisID).Find(&synthesisReaction).Error
	if err != nil {
		return nil, err
	}

	for _, mm := range synthesisReaction {
		reactionID = mm.ReactionID
		reaction, err = r.GetReaction(int(reactionID))

		if err != nil {
			logrus.Warnf("Reaction %d not found or deleted, skipping", reactionID)
			continue
		}
		result = append(result, reaction)
	}
	return result, nil
}

func (r *Repository) FindUserSynthesis(userID uint) uint {
	var synthesisID uint
	err := r.db.Model(&ds.Synthesis{}).Where("creator_id = ? AND status = ?", userID, "черновик").Select("id").First(&synthesisID).Error
	if err != nil {
		return 0
	}
	return synthesisID
}

func (r *Repository) AddReactionInSynthesis(id uint) error {
	userID := 1
	moderatorID := 2
	var synthesisID uint
	var count int64

	err := r.db.Model(&ds.Synthesis{}).Where("creator_id = ? AND status = ?", userID, "черновик").Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		newSynthesis := ds.Synthesis{
			Status:      "черновик",
			DateCreate:  time.Now(),
			DateUpdate:  time.Now(),
			CreatorID:   uint(userID),
			ModeratorID: uint(moderatorID),
		}
		err := r.db.Create(&newSynthesis).Error
		if err != nil {
			return err
		}
	}

	err = r.db.Model(&ds.Synthesis{}).Where("creator_id = ? AND status = ?", userID, "черновик").Select("id").First(&synthesisID).Error
	if err != nil {
		return err
	}

	var existingSynthesisReaction ds.SynthesisReaction
	err = r.db.Where("synthesis_id = ? AND reaction_id = ?", synthesisID, id).First(&existingSynthesisReaction).Error

	if err == nil {
		existingSynthesisReaction.Count++
		err = r.db.Save(&existingSynthesisReaction).Error
	} else {
		newSynthesisReaction := ds.SynthesisReaction{
			SynthesisID: synthesisID,
			ReactionID:  id,
			Count:       1,
			VolumeSM:    10.0,
			VolumeRM:    0,
		}
		err = r.db.Create(&newSynthesisReaction).Error
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) RemoveSynthesis(id uint) error {
	deleteQuery := "UPDATE syntheses SET status = $1, date_finish = $2, date_update = $3 WHERE id = $4"
	r.db.Exec(deleteQuery, "удалён", time.Now(), time.Now(), id)
	return nil
}

func (r *Repository) GetUserNameByID(userID uint) string {
	var userName string
	r.db.Model(&ds.Users{}).Where("id = ?", userID).Select("fio").First(&userName)
	return userName
}

func (r *Repository) GetDateUpdate(synthesisID uint) string {
	var dateUpdateTime time.Time
	var dateUpdate string
	r.db.Model(&ds.Synthesis{}).Where("id = ?", synthesisID).Select("date_update").First(&dateUpdateTime)
	dateUpdate = dateUpdateTime.Format("02.01.2006 15:04:05")
	return dateUpdate
}

func (r *Repository) GetPurity(synthesisID uint) float32 {
	var purity float32
	r.db.Model(&ds.Synthesis{}).Where("id = ?", synthesisID).Select("purity").First(&purity)
	return purity
}

func (r *Repository) GetReactionCount(synthesisID uint, reactionID uint) uint {
	var synthesisReaction ds.SynthesisReaction
	err := r.db.Where("synthesis_id = ? AND reaction_id = ?", synthesisID, reactionID).First(&synthesisReaction).Error

	if err != nil {
		return 0
	}
	return synthesisReaction.Count
}

func (r *Repository) SynthesisStatusById(synthesisID uint) (string, error) {
	var SynthesisStatus string
	err := r.db.Model(&ds.Synthesis{}).Where("id = ?", synthesisID).Select("status").First(&SynthesisStatus).Error
	if err != nil {
		return "", err
	}
	return SynthesisStatus, err
}

type SynthesisReactionWithCount struct {
	ds.Reaction
	VolumeSM float32 `json:"volume_sm"`
	VolumeRM float32 `json:"volume_rm,omitempty"`
	//SynthesisReactionID uint    `json:"synthesis_reaction_id"`
	Count uint
}

func (r *Repository) GetSynthesisWithCounts(synthesisID uint) ([]SynthesisReactionWithCount, error) {
	var synthesisReactions []ds.SynthesisReaction
	var result []SynthesisReactionWithCount

	err := r.db.Where("synthesis_id = ?", synthesisID).Find(&synthesisReactions).Error
	if err != nil {
		return nil, err
	}

	for _, sr := range synthesisReactions {
		reaction, err := r.GetReaction(int(sr.ReactionID))
		if err != nil {
			logrus.Warnf("Reaction %d not found or deleted, skipping", sr.ReactionID)
			continue
		}

		result = append(result, SynthesisReactionWithCount{
			Reaction: reaction,
			Count:    sr.Count,
			VolumeSM: sr.VolumeSM,
			VolumeRM: sr.VolumeRM,
		})
	}
	return result, nil
}

func (r *Repository) GetSynthesisCount(creatorID uint) int64 {
	var synthesisID uint
	var count int64
	//creatorID := 1
	err := r.db.Model(&ds.Synthesis{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("id").First(&synthesisID).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.SynthesisReaction{}).Where("synthesis_id = ?", synthesisID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in list_chats:", err)
	}

	return count
}

func (r *Repository) GetSynthesisID(userID uint) int {
	var synthesisID int
	err := r.db.Model(&ds.Synthesis{}).Where("creator_id = ? AND status = ?", userID, "черновик").Select("id").First(&synthesisID).Error
	if err != nil {
		return 0
	}
	return synthesisID
}

func (r *Repository) GetSynthesisByID(synthesisID uint) (*ds.Synthesis, []ds.Reaction, error) {
	var synthesis ds.Synthesis

	err := r.db.
		Preload("Creator").
		Preload("Moderator").
		Where("id = ?", synthesisID).
		First(&synthesis).Error

	if err != nil {
		return nil, nil, fmt.Errorf("синтез с ID %d не найдена", synthesisID)
	}

	var reactions []ds.Reaction
	err = r.db.
		Table("reactions").
		Joins("JOIN synthesis_reactions ON reactions.id = synthesis_reactions.reaction_id").
		Where("synthesis_reactions.synthesis_id = ?", synthesisID).
		Find(&reactions).Error

	if err != nil {
		return nil, nil, fmt.Errorf("ошибка загрузки реакции: %w", err)
	}

	return &synthesis, reactions, nil
}

func (r *Repository) DeleteSynthesis(synthesisID uint) error {
	var synthesis ds.Synthesis
	err := r.db.Where("id = ?", synthesisID).First(&synthesis).Error
	if err != nil {
		return fmt.Errorf("синтез с ID %d не найдена", synthesisID)
	}

	err = r.db.Model(&ds.Synthesis{}).Where("id = ?", synthesisID).Updates(map[string]interface{}{
		"status":      "удалён",
		"date_update": time.Now(),
	}).Error

	if err != nil {
		return fmt.Errorf("ошибка при удалении синтеза: %w", err)
	}

	return nil
}
