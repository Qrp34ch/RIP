//package repository
//
//import (
//	"fmt"
//	"strings"
//)
//
//type Repository struct {
//}
//
//func NewRepository() (*Repository, error) {
//	return &Repository{}, nil
//}
//
//type Step struct { // вот наша новая структура
//	ID               int    // поля структур, которые передаются в шаблон
//	Title            string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
//	Src              string
//	SrcUr            string
//	Details          string
//	Cart             bool
//	StartingMaterial string
//	DensitySM        float32
//	VolumeSM         float32
//	MolarMassSM      int
//	ResultMaterial   string
//	DensityRM        float32
//	VolumeRM         float32
//	MolarMassRM      int
//}
//
//type Reagent struct { // вот наша новая структура
//	ID        int    // поля структур, которые передаются в шаблон
//	Title     string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
//	Density   float32
//	Volume    float32
//	MolarMass int
//}
//
//func (r *Repository) GetSteps() ([]Step, error) {
//	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
//	steps := []Step{ // массив элементов из наших структур
//		{
//			ID:               1,
//			Title:            "Хлорирование толуола",
//			Src:              "http://localhost:9000/aspirinimages/img/toluol.jpg",
//			Details:          "Толуол хлорируют в присутствии катализатора хлорид алюминия",
//			SrcUr:            "http://localhost:9000/aspirinimages/img/urToluol.png",
//			Cart:             false,
//			StartingMaterial: "Толуол",
//			DensitySM:        0.87,
//			MolarMassSM:      92,
//			ResultMaterial:   "Пара-хлорметилбензол",
//			DensityRM:        1.1,
//			MolarMassRM:      127,
//		},
//		{
//			ID:               2,
//			Title:            "Окисление аддукта",
//			Src:              "http://localhost:9000/aspirinimages/img/adduct.png",
//			Details:          "Аддукт окисляют атомарным кислородом (озоном) при температуре t=0-5°С в водной эмульсии",
//			SrcUr:            "http://localhost:9000/aspirinimages/img/urAdduct.png",
//			Cart:             false,
//			StartingMaterial: "Пара-хлорметилбензол",
//			DensitySM:        1.1,
//			MolarMassSM:      127,
//			ResultMaterial:   "О-хлорбензойная кислота",
//			DensityRM:        1.25,
//			MolarMassRM:      157,
//		},
//		{
//			ID:               3,
//			Title:            "Омыление о-хлорбензойной кислоты",
//			Src:              "http://localhost:9000/aspirinimages/img/ohlorbenzol.jpg",
//			Details:          "О-хлорбензойную кислоту омыляют 30% водным раствором гидроксида натрия",
//			SrcUr:            "http://localhost:9000/aspirinimages/img/urOhlor.png",
//			Cart:             false,
//			StartingMaterial: "О-хлорбензойная кислота",
//			DensitySM:        1.25,
//			MolarMassSM:      157,
//			ResultMaterial:   "Натрия салицилат",
//			DensityRM:        1.7,
//			MolarMassRM:      160,
//		},
//		{
//			ID:               4,
//			Title:            "Свободная салициловая кислота",
//			Src:              "http://localhost:9000/aspirinimages/img/salicylicacid.jpg",
//			Details:          "Солевую форму салициловой кислоты переводят в свободную кислоту",
//			SrcUr:            "http://localhost:9000/aspirinimages/img/urSalAcid.png",
//			Cart:             true,
//			StartingMaterial: "Натрия салицилат",
//			DensitySM:        1.7,
//			MolarMassSM:      160,
//			ResultMaterial:   "Салициловая кислота",
//			DensityRM:        1.44,
//			MolarMassRM:      138,
//		},
//		{
//			ID:               5,
//			Title:            "Получение аспирина",
//			Src:              "http://localhost:9000/aspirinimages/img/aspirin.png",
//			Details:          "Реакция салициловой кислоты и уксусной кислоты с катализатором - серной кислотой. Для расчета необходимо указать объем салициловой и уксусной кислот.",
//			SrcUr:            "http://localhost:9000/aspirinimages/img/urAspirin.png",
//			Cart:             true,
//			StartingMaterial: "Салициловая кислота",
//			DensitySM:        1.44,
//			MolarMassSM:      138,
//			ResultMaterial:   "Аспирин",
//			DensityRM:        1.44,
//			MolarMassRM:      180,
//		},
//	}
//	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
//	// тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
//	if len(steps) == 0 {
//		return nil, fmt.Errorf("массив пустой")
//	}
//
//	return steps, nil
//}
//
//func (r *Repository) GetReagents() ([]Reagent, error) {
//	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
//	reagents := []Reagent{ // массив элементов из наших структур
//		{
//			ID:        1,
//			Title:     "Толуол",
//			Density:   0.87,
//			MolarMass: 92,
//		},
//		{
//			ID:        2,
//			Title:     "Пара-хлорметилбензол",
//			Density:   1.1,
//			MolarMass: 127,
//		},
//		{
//			ID:        3,
//			Title:     "О-хлорбензойная кислота",
//			Density:   1.25,
//			MolarMass: 157,
//		},
//		{
//			ID:        4,
//			Title:     "Натрия салицилат",
//			Density:   1.7,
//			MolarMass: 160,
//		},
//		{
//			ID:        5,
//			Title:     "Салициловая кислота",
//			Density:   1.44,
//			MolarMass: 138,
//		},
//		{
//			ID:        6,
//			Title:     "Аспирин",
//			Density:   1.44,
//			MolarMass: 180,
//		},
//	}
//	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
//	// тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
//	if len(reagents) == 0 {
//		return nil, fmt.Errorf("массив пустой")
//	}
//
//	return reagents, nil
//}
//
//func (r *Repository) GetStep(id int) (Step, error) {
//	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
//	steps, err := r.GetSteps()
//	if err != nil {
//		return Step{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
//	}
//
//	for _, step := range steps {
//		if step.ID == id {
//			return step, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
//		}
//	}
//	return Step{}, fmt.Errorf("этап не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
//}
//
//func (r *Repository) GetStepsByTitle(title string) ([]Step, error) {
//	steps, err := r.GetSteps()
//	if err != nil {
//		return []Step{}, err
//	}
//
//	var result []Step
//	for _, step := range steps {
//		if strings.Contains(strings.ToLower(step.Title), strings.ToLower(title)) {
//			result = append(result, step)
//		}
//	}
//
//	return result, nil
//}
//
//func (r *Repository) GetStepsInCart() ([]Step, error) {
//	steps, err := r.GetSteps()
//	if err != nil {
//		return []Step{}, err
//	}
//
//	var result []Step
//	for _, step := range steps {
//		if step.Cart {
//			result = append(result, step)
//		}
//	}
//
//	return result, nil
//}

package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	// Возвращаем объект Repository с подключенной базой данных
	return &Repository{
		db: db,
	}, nil
}
