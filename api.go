package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{"message": "Welcome to the API."})
	})

	router.GET("/buildid/:channel", func(c *gin.Context) {
		channel := c.Param("channel")
		info, err := get_build_id(channel)
		if err != nil {
			c.IndentedJSON(400, err.Error())
		} else {
			c.IndentedJSON(200, gin.H{"build_id": info})
		}
	})

	router.GET("/version/:channel", func(c *gin.Context) {
		release_channel := c.Param("channel")
		build_info, err := get_build_info(release_channel)
		if err != nil {
			c.IndentedJSON(400, err.Error())
		} else {
			c.IndentedJSON(200, build_info)
		}
	})

	err := router.Run("localhost:3000")
	if err != nil {
		fmt.Println(err)
	}
}
