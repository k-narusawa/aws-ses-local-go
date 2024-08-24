package main

import (
	"aws-ses-local-go/config"
	"aws-ses-local-go/internal/controllers/rest"
	"aws-ses-local-go/internal/gateways/dao"
	"aws-ses-local-go/internal/gateways/query"
	"aws-ses-local-go/internal/gateways/repository"
	v1 "aws-ses-local-go/usecase/aws/v1"
	v2 "aws-ses-local-go/usecase/aws/v2"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

const (
	defaultAddress = "8080"
)

func init() {
	log.Printf("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}
}

func main() {
	address := os.Getenv("SERVER_ADDRESS")

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("template/dist/*.html")),
	}
	e.Static("/", "template/dist")

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${uri}  ${latency_human}\n",
		Output: e.Logger.Output(),
	}))

	e.Renderer = t

	db := config.DBConnect()

	mailDao := dao.NewMailDao(db)
	mailRepo := repository.NewMailRepository(*mailDao)

	v1Svc := v1.NewService(mailRepo)
	v2Svc := v2.NewService(mailRepo)
	mailQSvc := query.NewMailDtoQueryService(*mailDao)
	cntQSvc := query.NewCountQueryService(*mailDao)

	rest.NewAwsHandler(e, v1Svc, v2Svc)
	rest.NewMailHandler(e, mailQSvc, cntQSvc)

	e.GET("/health", healthCheck)

	if address == "" {
		address = defaultAddress
	}
	log.Fatal(e.Start(address))
}

func healthCheck(c echo.Context) error {
	return c.String(200, "OK")
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
