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

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		reactions, err = h.Repository.GetReactions()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		reactions, err = h.Repository.GetReactionsByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"synthesisCount": h.Repository.GetReactionsInSynthesis(),
		"reactions":      reactions,
		"id":             h.Repository.FindUserSynthesis(1),
		"query":          searchQuery,
	})
}

func (h *Handler) GetReaction(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("Invalid ID format:", err)
		ctx.Redirect(http.StatusFound, "/reaction")
		return
	}

	reaction, err := h.Repository.GetReaction(id)
	if err != nil {
		logrus.Warnf("Reaction %d not found or deleted: %v", id, err)
		ctx.Redirect(http.StatusFound, "/reaction")
		return
	}

	ctx.HTML(http.StatusOK, "reaction.html", gin.H{
		"reaction": reaction,
	})
}

func (h *Handler) AddReactionInSynthesis(ctx *gin.Context) {
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
	ctx.Redirect(http.StatusFound, "/reaction")
}

func (h *Handler) GetSynthesis(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/reaction")
		return
	}

	synthesisStatus, err := h.Repository.SynthesisStatusById(uint(id))
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/reaction")
		return
	}

	if synthesisStatus == "удалён" {
		ctx.Redirect(http.StatusFound, "/reaction")
		return
	}

	synthesisReactions, err := h.Repository.GetSynthesisWithCounts(uint(id))
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/reaction")
		return
	}

	ctx.HTML(http.StatusOK, "synthesis.html", gin.H{
		"synthesisReactions": synthesisReactions,
		"id":                 id,
		"user":               h.Repository.GetUserNameByID(1),
		"date":               h.Repository.GetDateUpdate(uint(id)),
		"purity":             h.Repository.GetPurity(uint(id)),
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

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		reactions, err = h.Repository.GetReactions()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		reactions, err = h.Repository.GetReactionsByTitle(searchQuery)
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
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	reaction, err := h.Repository.GetReaction(id)

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

	_, err = h.Repository.GetReaction(int(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

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

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("файл изображения обязателен"))
		return
	}

	err = h.Repository.UploadReactionImage(uint(id), file)
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
		"message": "Изображение успешно загружено",
	})
}

func (h *Handler) GetSynthesisIconAPI(ctx *gin.Context) {

	synthesisID := h.Repository.GetSynthesisID(1)
	synthesisCount := h.Repository.GetSynthesisCount(1)

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"id_synthesis": synthesisID,
		"items_count":  synthesisCount,
	})
}

func (h *Handler) GetSynthesesAPI(ctx *gin.Context) {
	var filter struct {
		Status    string `form:"status"`
		StartDate string `form:"start_date"`
		EndDate   string `form:"end_date"`
	}

	type SynthesisWithLogin struct {
		ID             uint    `form:"id"`
		Status         string  `form:"status"`
		DateCreate     string  `form:"date_create"`
		DateUpdate     string  `form:"date_update"`
		DateFinish     string  `form:"date_finish"`
		CreatorID      uint    `form:"creator_id"`
		ModeratorID    uint    `form:"moderator_id"`
		Purity         float32 `form:"purity"`
		CreatorLogin   string  `form:"creator_login"`
		ModeratorLogin string  `form:"moderator_login"`
	}

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	syntheses, err := h.Repository.GetSyntheses(filter.Status, filter.StartDate, filter.EndDate)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	response := make([]SynthesisWithLogin, len(syntheses))
	for i, calc := range syntheses {
		response[i] = SynthesisWithLogin{
			ID:           calc.ID,
			Status:       calc.Status,
			DateCreate:   calc.DateCreate.Format("02.01.2006"),
			DateUpdate:   calc.DateUpdate.Format("02.01.2006"),
			CreatorLogin: calc.Creator.Login,
			Purity:       calc.Purity,
		}

		if calc.DateFinish.Valid {
			response[i].DateFinish = calc.DateFinish.Time.Format("02.01.2006")
		}

		if calc.Moderator.ID != 0 {
			response[i].ModeratorLogin = calc.Moderator.Login
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
		"count":  len(response),
	})
}

func (h *Handler) GetSynthesisAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	synthesis, reactions, err := h.Repository.GetSynthesisByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	synthesisFull := struct {
		ID             uint
		Status         string
		DateCreate     string
		DateUpdate     string
		DateFinish     string
		CreatorLogin   string
		ModeratorLogin string
		Purity         float32
		Reactions      []ds.Reaction
	}{
		ID:           synthesis.ID,
		Status:       synthesis.Status,
		DateCreate:   synthesis.DateCreate.Format("02.01.2006"),
		DateUpdate:   synthesis.DateUpdate.Format("02.01.2006"),
		CreatorLogin: synthesis.Creator.Login,
		Purity:       synthesis.Purity,
		Reactions:    make([]ds.Reaction, len(reactions)), // используем отдельно загруженные fuels
	}

	if synthesis.DateFinish.Valid {
		synthesisFull.DateFinish = synthesis.DateFinish.Time.Format("02.01.2006")
	}

	if synthesis.Moderator.ID != 0 {
		synthesisFull.ModeratorLogin = synthesis.Moderator.Login
	}

	for i, reaction := range reactions {
		synthesisFull.Reactions[i] = ds.Reaction{
			ID:               reaction.ID,
			Title:            reaction.Title,
			Src:              reaction.Src,
			SrcUr:            reaction.SrcUr,
			Details:          reaction.Details,
			IsDelete:         reaction.IsDelete,
			StartingMaterial: reaction.StartingMaterial,
			DensitySM:        reaction.DensitySM,
			VolumeSM:         reaction.VolumeSM,
			MolarMassSM:      reaction.MolarMassSM,
			ResultMaterial:   reaction.ResultMaterial,
			DensityRM:        reaction.DensityRM,
			VolumeRM:         reaction.VolumeRM,
			MolarMassRM:      reaction.MolarMassRM,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   synthesisFull,
	})
}

func (h *Handler) UpdateSynthesisPurityAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var input struct {
		Purity float64 `json:"purity" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.UpdateSynthesisPurity(uint(id), input.Purity)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	updatedSynthesis, _, err := h.Repository.GetSynthesisByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    updatedSynthesis,
		"message": "концентрация успешно обновлена",
	})
}

func (h *Handler) FormSynthesisAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.FormSynthesis(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	updatedSynthesis, reactions, err := h.Repository.GetSynthesisByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"data":      updatedSynthesis,
		"reactions": reactions,
		"message":   "Синтез успешно сформирован",
	})
}
