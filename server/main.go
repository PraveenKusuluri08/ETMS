package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Praveenkusuluri08/api/routes"
	"github.com/Praveenkusuluri08/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var wg sync.WaitGroup

type ENV struct {
	AppEnv      string `mapstructure:"APP_ENV"`
	MONGODB_URI string `mapstructure:"MONGODB_URI"`
	PORT        string `mapstructure:"PORT"`
	FromEmail   string
}

func main() {
	env := GetEnvConfig()

	app := bootstrap.App(env.MONGODB_URI)
	fmt.Println(app)
	router := gin.Default()

	port := env.PORT

	if port == "" {
		port = "8080"
	}
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ETMS",
		})
	})
	wg.Add(1)

	routes.GroupRouter(router.Group("/api/v1/groups"))
	routes.UserRoutes(router.Group("/api/v1/users"))

	go func() {
		fmt.Println("Server is running in port 8080")
		if err := router.Run(":" + port); err != nil {
			log.Fatal(err)
			return
		}
	}()
	wg.Wait()
}
func GetEnvConfig() *ENV {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	env := ENV{}
	env.MONGODB_URI = os.Getenv("MONGODB_URI")
	env.PORT = os.Getenv("PORT")
	env.AppEnv = os.Getenv("APP_ENV")
	env.AppEnv = os.Getenv("APP_ENV")
	return &env
}
