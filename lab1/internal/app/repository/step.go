package repository

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"lab1/internal/app/ds"
)

func (r *Repository) GetSteps() ([]ds.Step, error) {
	var steps []ds.Step
	err := r.db.Find(&steps).Error
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	if err != nil {
		return nil, err
	}
	if len(steps) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return steps, nil
}

func (r *Repository) GetStep(id int) (ds.Step, error) {
	step := ds.Step{}
	err := r.db.Where("id = ?", id).First(&step).Error
	if err != nil {
		return ds.Step{}, err
	}
	return step, nil
}

func (r *Repository) GetStepsByTitle(title string) ([]ds.Step, error) {
	var steps []ds.Step
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&steps).Error
	if err != nil {
		return nil, err
	}
	return steps, nil
}

func (r *Repository) GetStepsInCart() int64 {
	var cartID uint
	var count int64
	creatorID := 1
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	err := r.db.Model(&ds.Cart{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("id").First(&cartID).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.CartStep{}).Where("cart_id = ?", cartID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_steps:", err)
	}

	return count
}

func (r *Repository) GetCart(cartID uint) ([]ds.Step, error) {
	//userID := 1
	//var cartID uint
	//cartID = r.FindUserCart(userID)
	var stepID uint
	var step ds.Step
	var cartStep []ds.CartStep
	var result []ds.Step
	//err := r.db.Model(&ds.Cart{}).Where("creator_id = ? AND status = ?", userID, "черновик").Select("id").First(&cartID).Error
	//if err != nil {
	//	return nil, err
	//}
	err := r.db.Where("cart_id = ?", cartID).Find(&cartStep).Error
	if err != nil {
		return nil, err
	}

	for _, mm := range cartStep {
		stepID = mm.StepID
		step, err = r.GetStep(int(stepID))
		if err != nil {
			return nil, err
		}
		result = append(result, step)
	}

	return result, nil
}

func (r *Repository) FindUserCart(userID uint) uint {
	var cartID uint
	err := r.db.Model(&ds.Cart{}).Where("creator_id = ? AND status = ?", userID, "черновик").Select("id").First(&cartID).Error
	if err != nil {
		return 0
	}
	return cartID
}

func (r *Repository) AddStepInCart(id uint) error {
	userID := 1
	moderatorID := 2
	var cartID uint
	var count int64

	err := r.db.Model(&ds.Cart{}).Where("creator_id = ? AND status = ?", userID, "черновик").Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		newCart := ds.Cart{
			Status:      "черновик",
			DateCreate:  time.Now(),
			DateUpdate:  time.Now(),
			CreatorID:   uint(userID),
			ModeratorID: uint(moderatorID),
		}
		err := r.db.Create(&newCart).Error
		if err != nil {
			return err
		}
	}

	err = r.db.Model(&ds.Cart{}).Where("creator_id = ? AND status = ?", userID, "черновик").Select("id").First(&cartID).Error
	if err != nil {
		return err
	}

	newCartStep := ds.CartStep{
		CartID: cartID,
		StepID: id,
	}

	err = r.db.Create(&newCartStep).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *Repository) RemoveCart(id uint) error {
	deleteQuery := "UPDATE carts SET status = $1, date_finish = $2, date_update = $3 WHERE id = $4"
	r.db.Exec(deleteQuery, "удалён", time.Now(), time.Now(), id)
	return nil
}
