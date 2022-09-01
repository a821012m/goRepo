package main

import (
	"fmt"
	"line/controllers"

	"line/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	configInit()

	router := gin.Default()

	// router.SetTrustedProxies([]string{"192.168.50.192", "192.168.50.95"})
	useSawaggerUI(router)
	useCors(router)
	controllers.NewLineController(router)
	port := viper.GetString("System.Port")
	router.Run(":" + port)
}

/*
設定檔讀取初始化
*/
func configInit() {
	viper.SetConfigName("appSettings") // name of config file (without extension)
	viper.SetConfigType("json")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")           // optionally look for config in the working directory
	err := viper.ReadInConfig()        // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

/*
使用swagger ui
*/
func useSawaggerUI(router *gin.Engine) {
	docs.SwaggerInfo.Title = "Line Bot API"
	docs.SwaggerInfo.Description = "Line Bot Api Server "
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "jameschang821012.asuscomm.com/"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"https"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// 使用cors
func useCors(router *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))
}
