package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab1/internal/app/handler"
	"lab1/internal/app/repository"
	"log"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")
	// слева название папки, в которую выгрузится наша статика
	// справа путь к папке, в которой лежит статика

	r.GET("/reaction", handler.GetReactions)
	r.GET("/reaction/:id", handler.GetReaction)
	r.GET("/synthesis/:id", handler.GetReactionsInSynthesis)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}
