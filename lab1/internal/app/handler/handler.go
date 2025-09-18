//package handler
//
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/sirupsen/logrus"
//	"lab1/internal/app/repository"
//	"net/http"
//	"strconv"
//	"time"
//)
//
//type Handler struct {
//	Repository *repository.Repository
//}
//
//func NewHandler(r *repository.Repository) *Handler {
//	return &Handler{
//		Repository: r,
//	}
//}
//
//func (h *Handler) GetSteps(ctx *gin.Context) {
//	var steps []repository.Step
//	var err error
//
//	searchQuery := ctx.Query("query") // получаем значение из поля поиска
//	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
//		steps, err = h.Repository.GetSteps()
//		if err != nil {
//			logrus.Error(err)
//		}
//	} else {
//		steps, err = h.Repository.GetStepsByTitle(searchQuery) // в ином случае ищем заказ по заголовку
//		if err != nil {
//			logrus.Error(err)
//		}
//	}
//
//	cartSteps, err := h.Repository.GetStepsInCart()
//	cartCount := len(cartSteps)
//	if err != nil {
//		logrus.Error(err)
//		cartCount = 0
//	}
//
//	ctx.HTML(http.StatusOK, "index.html", gin.H{
//		"time":      time.Now().Format("15:04:05"),
//		"steps":     steps,
//		"query":     searchQuery, // передаем введенный запрос обратно на страницу
//		"cartCount": cartCount,
//		// в ином случае оно будет очищаться при нажатии на кнопку
//	})
//}
//
//func (h *Handler) GetStep(ctx *gin.Context) {
//	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
//	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
//	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
//	if err != nil {
//		logrus.Error(err)
//	}
//
//	step, err := h.Repository.GetStep(id)
//	if err != nil {
//		logrus.Error(err)
//	}
//
//	ctx.HTML(http.StatusOK, "step.html", gin.H{
//		"step": step,
//	})
//}
//
//func (h *Handler) GetStepsInCart(ctx *gin.Context) {
//	var steps []repository.Step
//	var err error
//
//	steps, err = h.Repository.GetStepsInCart()
//	if err != nil {
//		logrus.Error(err)
//	}
//
//	ctx.HTML(http.StatusOK, "cart.html", gin.H{
//		"steps": steps,
//	})
//}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab1/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/", h.GetSteps)
	router.GET("/step/:id", h.GetStep)
	router.GET("/cart/:id", h.GetCart)
	router.POST("/add-step-in-cart", h.AddStepInCart)
	router.POST("/delete/:id", h.RemoveCart)
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
