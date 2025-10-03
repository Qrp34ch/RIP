package handler

import (
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
