package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Reaction struct { // вот наша новая структура
	ID               int    // поля структур, которые передаются в шаблон
	Title            string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
	Src              string
	SrcUr            string
	Details          string
	Synthesis        bool
	StartingMaterial string
	DensitySM        float32
	VolumeSM         float32
	MolarMassSM      int
	ResultMaterial   string
	DensityRM        float32
	VolumeRM         float32
	MolarMassRM      int
}

func (r *Repository) GetReactions() ([]Reaction, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	reactions := []Reaction{ // массив элементов из наших структур
		{
			ID:               1,
			Title:            "Хлорирование толуола",
			Src:              "http://localhost:9000/aspirinimages/img/toluol.jpg",
			Details:          "Толуол хлорируют в присутствии катализатора хлорид алюминия",
			SrcUr:            "http://localhost:9000/aspirinimages/img/urToluol.png",
			Synthesis:        false,
			StartingMaterial: "Толуол",
			DensitySM:        0.87,
			MolarMassSM:      92,
			ResultMaterial:   "Пара-хлорметилбензол",
			DensityRM:        1.1,
			MolarMassRM:      127,
		},
		{
			ID:               2,
			Title:            "Окисление аддукта",
			Src:              "http://localhost:9000/aspirinimages/img/adduct.png",
			Details:          "Аддукт окисляют атомарным кислородом (озоном) при температуре t=0-5°С в водной эмульсии",
			SrcUr:            "http://localhost:9000/aspirinimages/img/urAdduct.png",
			Synthesis:        false,
			StartingMaterial: "Пара-хлорметилбензол",
			DensitySM:        1.1,
			MolarMassSM:      127,
			ResultMaterial:   "О-хлорбензойная кислота",
			DensityRM:        1.25,
			MolarMassRM:      157,
		},
		{
			ID:               3,
			Title:            "Омыление о-хлорбензойной кислоты",
			Src:              "http://localhost:9000/aspirinimages/img/ohlorbenzol.jpg",
			Details:          "О-хлорбензойную кислоту омыляют 30% водным раствором гидроксида натрия",
			SrcUr:            "http://localhost:9000/aspirinimages/img/urOhlor.png",
			Synthesis:        false,
			StartingMaterial: "О-хлорбензойная кислота",
			DensitySM:        1.25,
			MolarMassSM:      157,
			ResultMaterial:   "Натрия салицилат",
			DensityRM:        1.7,
			MolarMassRM:      160,
		},
		{
			ID:               4,
			Title:            "Свободная салициловая кислота",
			Src:              "http://localhost:9000/aspirinimages/img/salicylicacid.jpg",
			Details:          "Солевую форму салициловой кислоты переводят в свободную кислоту",
			SrcUr:            "http://localhost:9000/aspirinimages/img/urSalAcid.png",
			Synthesis:        true,
			StartingMaterial: "Натрия салицилат",
			DensitySM:        1.7,
			MolarMassSM:      160,
			ResultMaterial:   "Салициловая кислота",
			DensityRM:        1.44,
			MolarMassRM:      138,
		},
		{
			ID:               5,
			Title:            "Получение аспирина",
			Src:              "http://localhost:9000/aspirinimages/img/aspirin.png",
			Details:          "Реакция салициловой кислоты и уксусной кислоты с катализатором - серной кислотой. Для расчета необходимо указать объем салициловой и уксусной кислот.",
			SrcUr:            "http://localhost:9000/aspirinimages/img/urAspirin.png",
			Synthesis:        true,
			StartingMaterial: "Салициловая кислота",
			DensitySM:        1.44,
			MolarMassSM:      138,
			ResultMaterial:   "Аспирин",
			DensityRM:        1.44,
			MolarMassRM:      180,
		},
	}
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	// тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
	if len(reactions) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return reactions, nil
}

func (r *Repository) GetReaction(id int) (Reaction, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	reactions, err := r.GetReactions()
	if err != nil {
		return Reaction{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, reaction := range reactions {
		if reaction.ID == id {
			return reaction, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Reaction{}, fmt.Errorf("этап не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}

func (r *Repository) GetReactionsByTitle(title string) ([]Reaction, error) {
	reactions, err := r.GetReactions()
	if err != nil {
		return []Reaction{}, err
	}

	var result []Reaction
	for _, reaction := range reactions {
		if strings.Contains(strings.ToLower(reaction.Title), strings.ToLower(title)) {
			result = append(result, reaction)
		}
	}

	return result, nil
}

func (r *Repository) GetReactionsInSynthesis() ([]Reaction, error) {
	reactions, err := r.GetReactions()
	if err != nil {
		return []Reaction{}, err
	}

	var result []Reaction
	for _, reaction := range reactions {
		if reaction.Synthesis {
			result = append(result, reaction)
		}
	}

	return result, nil
}
