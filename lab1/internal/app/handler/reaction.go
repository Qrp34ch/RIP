package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"lab1/internal/app/ds"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	userId, err := h.GetUserID(ctx)
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"synthesisCount": h.Repository.GetReactionsInSynthesis(userId),
		"reactions":      reactions,
		"id":             h.Repository.FindUserSynthesis(userId),
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

	userId, err := h.GetUserID(ctx)

	err = h.Repository.AddReactionInSynthesis(uint(id), userId)
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
	userId, err := h.GetUserID(ctx)

	ctx.HTML(http.StatusOK, "synthesis.html", gin.H{
		"synthesisReactions": synthesisReactions,
		"id":                 id,
		"user":               h.Repository.GetUserNameByID(userId),
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
		//VolumeSM:         reactionInput.VolumeSM,
		MolarMassSM:    reactionInput.MolarMassSM,
		ResultMaterial: reactionInput.ResultMaterial,
		DensityRM:      reactionInput.DensityRM,
		//VolumeRM:         reactionInput.VolumeRM,
		MolarMassRM: reactionInput.MolarMassRM,
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
		//VolumeSM:         reactionInput.VolumeSM,
		MolarMassSM:    reactionInput.MolarMassSM,
		ResultMaterial: reactionInput.ResultMaterial,
		DensityRM:      reactionInput.DensityRM,
		//VolumeRM:         reactionInput.VolumeRM,
		MolarMassRM: reactionInput.MolarMassRM,
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
	userId, err := h.GetUserID(ctx)
	err = h.Repository.AddReactionInSynthesis(uint(id), userId)
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
	userID, err := h.GetUserID(ctx)
	if err != nil {
		fmt.Printf("nuiladno")
	}

	synthesisID := h.Repository.GetSynthesisID(userID)
	synthesisCount := h.Repository.GetSynthesisCount(userID)

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
	userID, err := h.GetUserID(ctx)

	syntheses, err := h.Repository.GetSyntheses(filter.Status, filter.StartDate, filter.EndDate, userID)
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
			//VolumeSM:         reaction.VolumeSM,
			MolarMassSM:    reaction.MolarMassSM,
			ResultMaterial: reaction.ResultMaterial,
			DensityRM:      reaction.DensityRM,
			//VolumeRM:         reaction.VolumeRM,
			MolarMassRM: reaction.MolarMassRM,
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

func (h *Handler) CompleteOrRejectSynthesisAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var input struct {
		NewStatus bool `json:"new_status" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	moderatorID := uint(2)

	err = h.Repository.CompleteOrRejectSynthesis(uint(id), moderatorID, input.NewStatus)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	updatedSynthesis, reactions, err := h.Repository.GetSynthesisByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	message := "Заявка отклонена"
	if input.NewStatus {
		message = "Заявка завершена"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"data":      updatedSynthesis,
		"reactions": reactions,
		"message":   message,
	})
}

func (h *Handler) DeleteSynthesisAPI(ctx *gin.Context) {
	userId, err := h.GetUserID(ctx)
	id := h.Repository.GetSynthesisID(userId)

	err = h.Repository.DeleteSynthesis(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Синтез успешно удален",
	})
}

func (h *Handler) RemoveReactionFromSynthesisAPI(ctx *gin.Context) {
	userId, err := h.GetUserID(ctx)
	synthesisID := h.Repository.GetSynthesisID(userId)
	reactionIDStr := ctx.Query("reaction_id")
	reactionID, err := strconv.Atoi(reactionIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.RemoveReactionFromSynthesis(uint(synthesisID), uint(reactionID))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Реакция удалена из синтеза",
	})
}

func (h *Handler) UpdateReactionInSynthesisAPI(ctx *gin.Context) {
	userId, err := h.GetUserID(ctx)
	synthesisID := h.Repository.GetSynthesisID(userId)

	var input struct {
		ReactionID uint    `json:"reaction_id" binding:"required"`
		VolumeSM   float64 `json:"volume_sm" binding:"required"`
	}
	//var err error
	if err = ctx.ShouldBindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.UpdateReactionInSynthesis(uint(synthesisID), input.ReactionID, input.VolumeSM)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Данные реакции обновлены в синтезе",
	})
}

type RegisterReq struct {
	Login string `json:"login"` // лучше назвать то же самое что login
	Pass  string `json:"pass"`
	FIO   string `json:"fio"`
}

type RegisterResp struct {
	Ok bool `json:"ok"`
}

func (h *Handler) RegisterUserAPI(gCtx *gin.Context) {
	req := &RegisterReq{}

	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if req.Pass == "" {
		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("pass is empty"))
		return
	}

	if req.Login == "" {
		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("name is empty"))
		return
	}

	err = h.Repository.RegisterUser(req.Login, generateHashString(req.Pass), req.FIO)
	if err != nil {
		gCtx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	gCtx.JSON(http.StatusOK, &RegisterResp{
		Ok: true,
	})
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (h *Handler) GetUserProfileAPI(ctx *gin.Context) {
	userID, err := h.GetUserID(ctx)

	user, err := h.Repository.GetUserProfile(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResp struct {
	ExpiresIn   int       `json:"expires_in"`
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	User        *UserInfo `json:"user"`
}

type UserInfo struct {
	ID          uint   `json:"id" example:"1"`
	Login       string `json:"login" example:"admin"`
	FIO         string `json:"fio" example:"Иванов Иван Иванович"`
	IsModerator bool   `json:"is_moderator" example:"true"`
}

// LoginUserAPI godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param input body LoginReq true "Login credentials"
// @Success 200 {object} LoginResp
// @Failure 400 {object} object{status=string,description=string} "Bad Request"
// @Failure 403 {object} object{status=string,description=string} "Forbidden"
// @Failure 500 {object} object{status=string,description=string} "Internal Server Error"
// @Router /API/users/login [post]
func (h *Handler) LoginUserAPI(ctx *gin.Context) {
	req := &LoginReq{}

	err := json.NewDecoder(ctx.Request.Body).Decode(req)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Аутентифицируем пользователя
	user, err := h.Repository.AuthUser(req.Login, generateHashString(req.Password))
	fmt.Printf(generateHashString(req.Password))
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	cfg := h.Config

	// Генерируем JWT токен
	token := jwt.NewWithClaims(cfg.JWT.SigningMethod, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cfg.JWT.ExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "bitop-admin",
			Subject:   strconv.FormatUint(uint64(user.ID), 10), // добавляем ID пользователя
		},
		UserUUID: uuid.New(),
		IsAdmin:  user.IsModerator,
		UserID:   user.ID,
		Scopes:   []string{},
	})

	strToken, err := token.SignedString([]byte(cfg.JWT.Token))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
		return
	}

	ctx.JSON(http.StatusOK, LoginResp{
		ExpiresIn:   int(cfg.JWT.ExpiresIn.Seconds()), // конвертируем в секунды
		AccessToken: strToken,
		TokenType:   "Bearer",
		User: &UserInfo{
			ID:          user.ID,
			Login:       user.Login,
			FIO:         user.FIO,
			IsModerator: user.IsModerator,
		},
	})
}

func (h *Handler) LogoutUserAPI(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Выход выполнен",
		})
		return
	}

	if !strings.HasPrefix(tokenString, jwtPrefix) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Выход выполнен",
		})
		return
	}

	// Отрезаем префикс
	tokenString = tokenString[len(jwtPrefix):]

	// Парсим токен чтобы получить expiration time
	claims := &ds.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.Token), nil
	})

	if err != nil || !token.Valid {
		// Если токен невалиден, все равно считаем выход успешным
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Выход выполнен",
		})
		return
	}

	// Добавляем токен в черный список
	if h.Repository.RedisClient != nil {
		// Время жизни в черном списке = оставшееся время жизни токена
		remainingTTL := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
		if remainingTTL > 0 {
			err = h.Repository.RedisClient.WriteJWTToBlacklist(ctx.Request.Context(), tokenString, remainingTTL)
			if err != nil {
				// Логируем ошибку, но все равно возвращаем успех
				logrus.Errorf("Ошибка добавления токена в черный список: %v", err)
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Выход выполнен успешно",
	})
}

func (h *Handler) UpdateUserAPI(ctx *gin.Context) {
	userID, err := h.GetUserID(ctx)

	var input struct {
		Login    *string `json:"login,omitempty"`
		Name     *string `json:"name,omitempty"`
		Password *string `json:"password,omitempty"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	updates := make(map[string]interface{})
	if input.Login != nil {
		updates["login"] = *input.Login
	}
	if input.Name != nil {
		updates["name"] = *input.Name
	}

	//if input.Password != nil {
	//	input.Password = generateHashString(input.Password)
	//	updates["password"] = *input.Password
	//}
	if len(updates) == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("нет полей для обновления"))
		return
	}
	user, err := h.Repository.UpdateUser(userID, updates)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    user,
		"message": "Данные обновлены",
	})
}

// GetUserID получает userID из контекста (установленного middleware)
func (h *Handler) GetUserID(ctx *gin.Context) (uint, error) {
	userID, exists := ctx.Get("userID")
	if !exists {
		return 0, fmt.Errorf("требуется авторизация")
	}

	return userID.(uint), nil
}
