package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	/* Setting up command line arguments. */
	num := flag.Int("p", 3000, "Port to run the API on.")
	flag.Parse()

	/* Reading .env file. */
	err_dot := godotenv.Load()
	if err_dot != nil {
		fmt.Println("Error loading .env file")
	}

	/* Setting gin to release mode. */
	gin.SetMode(gin.ReleaseMode)

	/* Setting up the API. */
	router := gin.Default()
	router.Use(validate_API_key())

	router.SetTrustedProxies(nil)

	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Welcome to the API."})
	})

	router.GET("/buildid/:channel", func(c *gin.Context) {
		channel := c.Param("channel")
		info, err := get_build_id(channel)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"build_id": info})
		}
	})

	router.GET("/version/:channel", func(c *gin.Context) {
		release_channel := c.Param("channel")
		build_info, err := get_build_info(release_channel)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		} else {
			c.IndentedJSON(http.StatusOK, build_info)
		}
	})

	err := router.Run(fmt.Sprintf("localhost:%d", *num))
	if err != nil {
		fmt.Println(err)
	}
}

func validate_API_key() gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key := c.Request.Header.Get("X-API-Key")
		if api_key != os.Getenv("API_KEY") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid API_KEY") /* This is really cringe, please switch to redis soon. */
		}
	}
}
