package main

import (
	"log"
	"os"

	"github.com/bcc-code/bccm-utils/mediainfoserver"
	"github.com/gin-gonic/gin"
)

func mediaInfoHandler(c *gin.Context) {
	filePath := c.Query("file")
	result, err := mediainfoserver.GetInfo(filePath)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, result)
}

func simpleInfoHandler(c *gin.Context) {
	filePath := c.Query("file")
	vantageCompat := c.Query("vantageCompat")

	result, err := mediainfoserver.GetInfo(filePath)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
	}

	simpleResult := mediainfoserver.SimpleInfo(result)

	// Vantage used to return everything as strings, so we have to do the same if
	// this should be a drop-in replacement.
	if vantageCompat == "true" {
		c.JSON(200, simpleResult.AsStringly())
		return
	}

	c.JSON(200, simpleResult)
}

func main() {
	s := gin.New()
	s.GET("/mediainfo", mediaInfoHandler)
	s.GET("/mediainfoSimple", simpleInfoHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	s.Run(":" + port)
}
