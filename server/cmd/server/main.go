package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/spf13/viper"

	"github.com/whitemonarch/gungi-server/server/internal/api"
	"github.com/whitemonarch/gungi-server/server/internal/auth"
	"github.com/whitemonarch/gungi-server/server/internal/ws"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	DB_SOURCE           string `mapstructure:"DB_SOURCE"`
	SUPABASE_JWT_SECRET string `mapstructure:"SUPABASE_JWT_SECRET"`
	PORT                string `mapstructure:"PORT"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig() (config Config, err error) {
	var cfg Config
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		cfg.DB_SOURCE = viper.Get("DB_SOURCE").(string)
		cfg.SUPABASE_JWT_SECRET = viper.Get("SUPABASE_JWT_SECRET").(string)
		cfg.PORT = viper.Get("PORT").(string)
		return cfg, nil
	}

	err = viper.Unmarshal(&cfg)
	return cfg, err
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalln("Cannot load config", err)
	}
	port := config.PORT
	if port == "" {
		port = "localhost:5080"
	} else {
		// port = ":" + port
		port = ":" + config.PORT
	}

	db := api.DBConn{}
	maxRetries := 5
	sleepDuration := 2 * time.Second
	for i := 1; i <= maxRetries; i++ {
		err = db.PostgresConnect(config.DB_SOURCE)
		if err == nil {
			break
		}
		if i < maxRetries {
			log.Println("Connection failed, retrying...")
			time.Sleep(sleepDuration)
			sleepDuration *= 2
		}
	}
	if err != nil {
		log.Fatalln("Failed to establish a database connection: ", err)
	}
	defer db.Conn.Close()

	e := echo.New()
	m := melody.New()
	m.Config.MaxMessageSize = 2048

	// e.Use(middleware.VerifySupabaseTokenMiddleware)
	verify := e.Group("", VerifySupabaseTokenMiddleware)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world")
	})

	e.GET("/game/:id", db.GetGameWithUndoRoute)
	// verify.GET("/getongoinggamelist", db.GetOngoingGameList)
	verify.GET("/overview", db.GetOverview)
	verify.GET("/getongoinggamelist", db.GetGameHistory)

	e.GET("/ws", ws.WSHandler(m, &db))

	verify.GET("/username", db.GetUsername)
	verify.PUT("/username", db.PutUsername)
	// verify.GET("/user/onboarding", db.GetOnboarding)
	// verify.PUT("/user/onboarding", db.PutOnboarding)
	// verify.PUT("/user/changename", db.ChangeUsername)

	e.Logger.Fatal(e.Start(port))
}

func VerifySupabaseTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := auth.AuthenticateSupabaseToken(c.Request().Header.Get("Authorization"))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err)
		}

		c.Set("sub", claims["sub"])

		return next(c)
	}
}
