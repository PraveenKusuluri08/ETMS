package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Praveenkusuluri08/api/routes"
	"github.com/Praveenkusuluri08/bootstrap"
	_ "github.com/Praveenkusuluri08/docs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var wg sync.WaitGroup

type ENV struct {
	AppEnv      string `mapstructure:"APP_ENV"`
	MONGODB_URI string `mapstructure:"MONGODB_URI"`
	PORT        string `mapstructure:"PORT"`
	FromEmail   string
}

// @title			ETMS API
// @version		1.0
// @description	This is a simple ETMS API's for managing expenses by own or by others with groups
// @termsofservice	https://swagger.io/terms
// @contact.name	Praveen
// @host			localhost:8080
func main() {
	env := GetEnvConfig()

	app := bootstrap.App(env.MONGODB_URI, env.AppEnv == "production")
	fmt.Println(app)
	router := gin.Default()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := env.PORT

	if port == "" {
		port = "8080"
	}

	router.GET("/api/v1/test", testEndpoint)

	wg.Add(1)

	routes.SetUp(router)

	go func() {
		fmt.Println("Server is running in port 8080")
		if err := router.Run(":" + port); err != nil {
			log.Fatal(err)
			return
		}
	}()
	wg.Wait()
}

// Test route

// @Summary		Test endpoint to check all the connections are good
// @Description	returns just message
// @Produce		json
// @Success		200	{string}	string	"OK"
// @Router			/api/v1/test [get]
// @Tags			Test
// @Security		ApiKeyAuth
func testEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ETMS",
	})
}
func GetEnvConfig() *ENV {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	env := ENV{}
	env.MONGODB_URI = os.Getenv("MONGODB_URI")
	env.PORT = os.Getenv("PORT")
	env.AppEnv = os.Getenv("APP_ENV")
	return &env
}
