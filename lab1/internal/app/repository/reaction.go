package repository

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"lab1/internal/app/ds"
)

func (r *Repository) GetReactions() ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Find(&reactions).Error
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
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
	err := r.db.Where("id = ?", id).First(&reaction).Error
	if err != nil {
		return ds.Reaction{}, err
	}
	return reaction, nil
}

func (r *Repository) GetReactionsByTitle(title string) ([]ds.Reaction, error) {
	var reactions []ds.Reaction
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&reactions).Error
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (r *Repository) GetReactionsInSynthesis() int64 {
	var synthesisID uint
	var count int64
	creatorID := 1
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

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
			return nil, err
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

	newSynthesisReaction := ds.SynthesisReaction{
		SynthesisID: synthesisID,
		ReactionID:  id,
	}

	err = r.db.Create(&newSynthesisReaction).Error
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
	//не каждый человек знает песни Газманова так что
	//а я я ясные дни забираю себе а я хмурые дни возвращаю судьбе <3 ВЛАДИК ЧИКАНЧИ ПРИШЕЛ ЗА ШЛЯПОЙ В МАГАЗИН (ОНА ЕМУ НЕ ПОДОШЛА)
	r.db.Model(&ds.Synthesis{}).Where("id = ?", synthesisID).Select("date_update").First(&dateUpdateTime)
	dateUpdate = dateUpdateTime.Format("02.01.2006 15:04:05")
	return dateUpdate
}
