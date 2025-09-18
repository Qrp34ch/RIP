package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab1/internal/app/ds"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) GetSteps(ctx *gin.Context) {
	var steps []ds.Step
	var err error

	searchQuery := ctx.Query("query") // получаем значение из нашего поля
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		steps, err = h.Repository.GetSteps()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		steps, err = h.Repository.GetStepsByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"cartCount": h.Repository.GetStepsInCart(),
		"steps":     steps,
		"id":        h.Repository.FindUserCart(1),
		"query":     searchQuery, // передаем введенный запрос обратно на страницу
		// в ином случае оно будет очищаться при нажатии на кнопку
	})
}

func (h *Handler) GetStep(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	step, err := h.Repository.GetStep(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "step.html", gin.H{
		"step": step,
	})
}

func (h *Handler) AddStepInCart(ctx *gin.Context) {
	// считываем значение из формы, которую мы добавим в наш шаблон
	strId := ctx.PostForm("step_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	// Вызов функции добавления чата в заявку
	err = h.Repository.AddStepInCart(uint(id))
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return
	}

	// после вызова сразу произойдет обновление страницы
	ctx.Redirect(http.StatusFound, "/")
}

func (h *Handler) GetCart(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	var steps []ds.Step
	steps, err = h.Repository.GetCart(uint(id))
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"steps": steps,
		"id":    id,
	})
}

func (h *Handler) RemoveCart(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	err = h.Repository.RemoveCart(uint(id))
	ctx.Redirect(http.StatusFound, "/")
}
