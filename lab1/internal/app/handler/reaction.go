package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab1/internal/app/ds"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) GetReactions(ctx *gin.Context) {
	var reactions []ds.Reaction
	var err error

	searchQuery := ctx.Query("query") // получаем значение из нашего поля
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		reactions, err = h.Repository.GetReactions()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		reactions, err = h.Repository.GetReactionsByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"synthesisCount": h.Repository.GetReactionsInSynthesis(),
		"reactions":      reactions,
		"id":             h.Repository.FindUserSynthesis(1),
		"query":          searchQuery, // передаем введенный запрос обратно на страницу
		// в ином случае оно будет очищаться при нажатии на кнопку
	})
}

func (h *Handler) GetReaction(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error("Invalid ID format:", err)
		ctx.Redirect(http.StatusFound, "/reaction")
		return
	}

	reaction, err := h.Repository.GetReaction(id)
	//if err != nil {
	//	logrus.Error(err)
	//}
	if err != nil {
		logrus.Warnf("Reaction %d not found or deleted: %v", id, err)
		// Показываем страницу с ошибкой или перенаправляем
		ctx.Redirect(http.StatusFound, "/reaction")
		return
	}

	ctx.HTML(http.StatusOK, "reaction.html", gin.H{
		"reaction": reaction,
	})
}

func (h *Handler) AddReactionInSynthesis(ctx *gin.Context) {
	// считываем значение из формы, которую мы добавим в наш шаблон
	strId := ctx.PostForm("reaction_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = h.Repository.GetReaction(id)
	if err != nil {
		logrus.Warnf("Cannot add deleted reaction %d to synthesis", id)
		ctx.Redirect(http.StatusFound, "/reaction")
		return
	}

	err = h.Repository.AddReactionInSynthesis(uint(id))
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return
	}

	// после вызова сразу произойдет обновление страницы
	ctx.Redirect(http.StatusFound, "/reaction")
}

func (h *Handler) GetSynthesis(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	synthesisStatus, err := h.Repository.SynthesisStatusById(uint(id))
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/reaction")
	}

	// если заявка по которой переходим удалена, то перенаправляем на главную
	if synthesisStatus == "удалён" {
		ctx.Redirect(http.StatusFound, "/reaction")
	}
	var reactions []ds.Reaction
	reactions, err = h.Repository.GetSynthesis(uint(id))
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "synthesis.html", gin.H{
		"reactions": reactions,
		"id":        id,
		"user":      h.Repository.GetUserNameByID(1),
		"date":      h.Repository.GetDateUpdate(uint(id)),
	})
}

func (h *Handler) RemoveSynthesis(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	err = h.Repository.RemoveSynthesis(uint(id))
	ctx.Redirect(http.StatusFound, "/reaction")
}
func (h *Handler) GetReactionsAPI(ctx *gin.Context) {
	var reactions []ds.Reaction
	var err error

	searchQuery := ctx.Query("query") // получаем значение из нашего поля
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		reactions, err = h.Repository.GetReactions()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		reactions, err = h.Repository.GetReactionsByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reactions": reactions,
		"query":     searchQuery,
	})
}
func (h *Handler) GetReactionAPI(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	//if err != nil {
	//	logrus.Error(err)
	//}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	reaction, err := h.Repository.GetReaction(id)
	//if err != nil {
	//	logrus.Error(err)
	//}

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Reaction not found or deleted"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reaction": reaction,
	})
}

func (h *Handler) CreateReactionAPI(ctx *gin.Context) {
	var reactionInput struct {
		Title string `json:"title,omitempty"`
		//Src              string  json:"src,omitempty"
		//SrcUr            string  json:"src_ur,omitempty"
		Details          string  `json:"details,omitempty"`
		IsDelete         bool    `json:"is_delete,omitempty"`
		StartingMaterial string  `json:"starting_material,omitempty"`
		DensitySM        float32 `json:"density_sm,omitempty"`
		VolumeSM         float32 `json:"volume_sm,omitempty"`
		MolarMassSM      int     `json:"molar_mass_sm,omitempty"`
		ResultMaterial   string  `json:"result_material,omitempty"`
		DensityRM        float32 `json:"density_rm,omitempty"`
		VolumeRM         float32 `json:"volume_rm,omitempty"`
		MolarMassRM      int     `json:"molar_mass_rm,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&reactionInput); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	newReaction := ds.Reaction{
		Title: reactionInput.Title,
		//Src              string  json:"src,omitempty"
		//SrcUr            string  json:"src_ur,omitempty"
		Details:          reactionInput.Details,
		IsDelete:         reactionInput.IsDelete,
		StartingMaterial: reactionInput.StartingMaterial,
		DensitySM:        reactionInput.DensitySM,
		VolumeSM:         reactionInput.VolumeSM,
		MolarMassSM:      reactionInput.MolarMassSM,
		ResultMaterial:   reactionInput.ResultMaterial,
		DensityRM:        reactionInput.DensityRM,
		VolumeRM:         reactionInput.VolumeRM,
		MolarMassRM:      reactionInput.MolarMassRM,
	}

	err := h.Repository.AddReaction(&newReaction)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"data":    newReaction,
		"message": "Реакция успешно создана",
	})
}

func (h *Handler) ChangeReactionAPI(ctx *gin.Context) {
	idReactionStr := ctx.Param("id")
	id, err := strconv.Atoi(idReactionStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var reactionInput struct {
		Title            string  `json:"title,omitempty"`
		Src              string  `json:"src,omitempty"`
		SrcUr            string  `json:"src_ur,omitempty"`
		Details          string  `json:"details,omitempty"`
		IsDelete         bool    `json:"is_delete,omitempty"`
		StartingMaterial string  `json:"starting_material,omitempty"`
		DensitySM        float32 `json:"density_sm,omitempty"`
		VolumeSM         float32 `json:"volume_sm,omitempty"`
		MolarMassSM      int     `json:"molar_mass_sm,omitempty"`
		ResultMaterial   string  `json:"result_material,omitempty"`
		DensityRM        float32 `json:"density_rm,omitempty"`
		VolumeRM         float32 `json:"volume_rm,omitempty"`
		MolarMassRM      int     `json:"molar_mass_rm,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&reactionInput); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	changeReaction := ds.Reaction{
		Title:            reactionInput.Title,
		Src:              reactionInput.Src,
		SrcUr:            reactionInput.SrcUr,
		Details:          reactionInput.Details,
		IsDelete:         reactionInput.IsDelete,
		StartingMaterial: reactionInput.StartingMaterial,
		DensitySM:        reactionInput.DensitySM,
		VolumeSM:         reactionInput.VolumeSM,
		MolarMassSM:      reactionInput.MolarMassSM,
		ResultMaterial:   reactionInput.ResultMaterial,
		DensityRM:        reactionInput.DensityRM,
		VolumeRM:         reactionInput.VolumeRM,
		MolarMassRM:      reactionInput.MolarMassRM,
	}
	err = h.Repository.ChangeReaction(uint(id), &changeReaction)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	updatedReaction, err := h.Repository.GetReaction(int(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    updatedReaction,
		"message": "Реакция успешно обновлена",
	})
}

func (h *Handler) DeleteReactionAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Проверяем существование записи перед удалением
	_, err = h.Repository.GetReaction(int(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// Используем ваш существующий метод DeleteFuel (мягкое удаление)
	err = h.Repository.DeleteReaction(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Реакция успешно удалена",
	})
}

func (h *Handler) AddReactionInSynthesisAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	err = h.Repository.AddReactionInSynthesis(uint(id))
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Реакция добавлена в заявку",
	})
}

func (h *Handler) UploadReactionImageAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Получаем файл из формы
	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("файл изображения обязателен"))
		return
	}

	// Загружаем изображение
	err = h.Repository.UploadReactionImage(uint(id), file)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Получаем обновленные данные услуги
	updatedReaction, err := h.Repository.GetReaction(int(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    updatedReaction,
		"message": "Изображение успешно загружено",
	})
}
