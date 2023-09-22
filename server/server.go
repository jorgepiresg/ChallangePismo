package server

import (
	"log"
	"os"

	"github.com/jorgepiresg/ChallangePismo/api"
	"github.com/jorgepiresg/ChallangePismo/app"
	"github.com/jorgepiresg/ChallangePismo/config"
	"github.com/jorgepiresg/ChallangePismo/store"
	"github.com/jorgepiresg/ChallangePismo/utils"
	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"

	_ "github.com/jorgepiresg/ChallangePismo/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server interface {
	Start()
}

type server struct {
	echo   *echo.Echo
	config config.Config
	store  store.Store
}

func New(cfg config.Config) Server {
	return &server{
		config: cfg,
	}
}

// @title Challange Pismo
// @version 1.0

// @contact.name Jorge Pires
// @contact.email jorgewpgomes@gmail.com

// @host localhost:8080/
// @BasePath api/v1
func (s *server) Start() {

	s.echo = echo.New()
	s.echo.HTTPErrorHandler = createHTTPErrorHandler()

	s.echo.Use(emiddleware.BodyLimit("2M"))
	s.echo.Use(emiddleware.Recover())
	s.echo.Use(emiddleware.RequestID())
	s.echo.Use(emiddleware.CORS())
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	s.startStore()

	app := app.New(app.Options{
		Store: s.store,
	})

	api.New(api.Options{
		Group: s.echo.Group("/api"),
		App:   app,
	})

	log.Println("Start server PID: ", os.Getpid())
	if err := s.echo.Start(s.config.ServerPort); err != nil {
		log.Println("cannot starting server ", err.Error())
	}
}

func createHTTPErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if err := c.JSON(utils.GetHTTPCode(err), err); err != nil {
			log.Println(err)
		}
	}
}

func (s *server) startStore() {
	s.store = store.New(store.Options{
		DB: s.createSqlConn(),
	})
}
