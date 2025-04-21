package router

import "github.com/gin-gonic/gin"

func Initialize() {
	router := gin.Default()
	initializeRoutes(router)
	if err := router.Run(); err != nil {
		panic(err)
	}
}
