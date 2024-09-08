package main

import (
	_ "PaintBackend/docs"
	ginHandlers "PaintBackend/internal/api"
	"PaintBackend/internal/config"
	"PaintBackend/internal/middleware"
	"PaintBackend/internal/storage"
	"PaintBackend/internal/telegram"
	"context"
	"database/sql"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	telegramHandlers "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/inlinequery"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func setupRouter(cfg *config.Config, h *ginHandlers.Handlers) *gin.Engine {
	if cfg.Env != "prod" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.Use(middleware.RequestProcessingMiddleware())

	api := r.Group("/api")
	{
		api.POST("/images/upload", h.UploadImage)
		api.DELETE("/images/:id", h.DeleteImage)
	}
	if cfg.Env != "prod" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		gin.SetMode(gin.DebugMode)

	}

	return r
}

func setupBot(token string, app *App) {

	bot, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout,
				APIURL:  gotgbot.DefaultAPIURL,
			},
		},
	})

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	dispatcher.AddHandler(telegramHandlers.NewInlineQuery(inlinequery.All, app.TelegramHandlers.Source))

	dispatcher.AddHandler(telegramHandlers.NewCommand("start", app.TelegramHandlers.Start))

	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}

	go updater.Idle()
}

type App struct {
	DB               *sql.DB
	FileStorage      *storage.PostgresFileRepository
	GinApiHandlers   *ginHandlers.Handlers
	TelegramHandlers *telegram.Handlers
}

func NewApp(db *sql.DB) *App {
	fileStorage := storage.NewPostgresFileRepository(db)
	return &App{
		GinApiHandlers:   ginHandlers.NewHandler(fileStorage),
		TelegramHandlers: telegram.NewHandler(fileStorage),
	}
}

func setupLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

func loadConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration: ", err)
	}
	return cfg
}

func setupDatabase(cfg *config.Config) *sql.DB {
	ctx := context.Background()
	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		slog.Error("Failed to connect to Postgres")
	}
	err = goose.RunContext(ctx, "up", db, "migrations", cfg.DatabaseDSN)
	if err != nil {
		slog.Error("Failed to run migrations")
		panic(err)
	}
	return db
}

// @title           Image Upload API
// @version         1.0
// @description     This is a backend API for Image Service

// @host      127.0.0.1:8080
// @BasePath  /api/
func main() {
	setupLogger()
	cfg := loadConfig()
	db := setupDatabase(cfg)
	defer db.Close()

	app := NewApp(db)
	setupBot(cfg.BotToken, app)

	r := setupRouter(cfg, app.GinApiHandlers)
	r.Run(":8080")

}
