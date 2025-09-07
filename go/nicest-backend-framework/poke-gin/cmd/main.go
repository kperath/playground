package main

import (
	"log"
	"net/http"
	"os"

	"playground/nicest-backend-framework/pokemon"

	"github.com/gin-gonic/gin"
)

func main() {
	f, err := os.Open("data.json")
	if err != nil {
		panic(err)
	}
	dex, err := pokemon.New(f)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"pokemon": dex.Pokemon,
		})
	})
	log.Fatal(r.Run())
}
