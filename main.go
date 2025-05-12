package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	cmd "github.com/IAmRiteshKoushik/pulse/cmd"
	c "github.com/IAmRiteshKoushik/pulse/controllers"
	mw "github.com/IAmRiteshKoushik/pulse/middleware"
	pkg "github.com/IAmRiteshKoushik/pulse/pkg"
)

func StartApp() {
	failMsg := "Could not initialize app\n%w"

	// Initialize global environment variables
	env, err := cmd.NewEnvConfig()
	if err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
	cmd.EnvVars = env
	log.Println("[OK]: Environment variables configured successfully")

	// Initialize logger
	f, err := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
	defer f.Close()
	cmd.Log = cmd.NewLoggerService(cmd.EnvVars.Environment, f)
	cmd.Log.Info("[OK]: Logging service configured successfully.")

	// Initialize Regex
	if err := pkg.InitRegex(); err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
	cmd.Log.Info("[OK]: All Regular Expressions created successfully.")

	// Setup RSA (if not exists) + Initialize PASETO
	err = pkg.CheckRSAKeyPairExists()
	if err != nil {
		err = pkg.GenerateRSAKeyPair()
		if err != nil {
			panic(fmt.Errorf(failMsg, err))
		}
		cmd.Log.Info("[OK]: RSA keypair generated and saved successfully.")
	} else {
		cmd.Log.Info("[OK]: Using existing RSA keypair.")
	}
	err = pkg.InitPaseto()
	if err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
	cmd.Log.Info("[OK]: PASETO initialized successfully.")

	// Initialize database connection pool
	cmd.DBPool, err = cmd.InitDB()
	if err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
	cmd.Log.Info("[OK]: DB initialized successfully.")

	// Initialize redis (cache)

	// Initialize Server
	server := InitServer()
	port := strconv.Itoa(cmd.EnvVars.Port)
	cmd.Log.Info("[OK]: Server configured and starting on PORT " + port)
	err = server.Run(":" + port)
	if err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
}

func InitServer() *gin.Engine {

	config := cors.Config{
		AllowOrigins:              []string{cmd.EnvVars.Domain},
		AllowWildcard:             true,
		AllowMethods:              []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:              []string{"X-Csrf-Token", "Origin", "Content-Type"},
		AllowCredentials:          true,
		OptionsResponseStatusCode: 204,
		MaxAge:                    12 * time.Hour,
	}

	// Gin Configurations for Logging, Monitoring and Panic-Recovery
	ginLogs, err := os.Create("gin.log")
	if err != nil {
		cmd.Log.Fatal("Error creating log file for Gin", err)
		return nil
	}
	defer ginLogs.Close()
	multiWriter := io.MultiWriter(os.Stdout, ginLogs)
	gin.DefaultWriter = multiWriter
	gin.DefaultErrorWriter = multiWriter
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(cors.New(config)) // Setup CORS() first before other middlewares
	router.Use(mw.RecoveryMiddleware)
	router.Use(gin.Logger())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is LIVE",
		})
		return
	})

	v1 := router.Group("/api/v1")
	Routes(v1)
	return router
}

func Routes(engine *gin.RouterGroup) {
	engine.GET("/user/login", c.LoginUserCsrf)
	engine.POST("/user/login", mw.VerifyCsrf, c.LoginUser)
	engine.GET("/user/register", c.RegisterUserAccountCsrf)
	engine.POST("/user/register", mw.VerifyCsrf, c.RegisterUserAccount)
	engine.GET("/user/register/otp", c.VerifyOtpCsrf)
	engine.POST("/user/register/otp", mw.VerifyCsrf, c.VerifyOtp)
	engine.GET("/user/register/otp/resend", c.ResendUserOtpCsrf)
	engine.POST("/user/register/otp/resend", mw.VerifyCsrf, c.ResendUserOtp)
	engine.GET("/user/session", mw.Auth, c.UserSession)
	engine.GET("/user/logout", mw.Auth, c.LogoutUser)

	engine.GET("/github/login", c.GithubLoginHandler)
	engine.GET("/github/callback", c.GithubCallbackHandler)

}

func main() {
	StartApp()
}
