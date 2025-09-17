package main

import (
	"horaires-piscine/internal/piscine"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	piscineMap := piscine.LoadPiscineMap()

	r.GET("/api/piscine/:piscine", piscine.PiscineHandler(piscineMap))
	r.Run()
}
