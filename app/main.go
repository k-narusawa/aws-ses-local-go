package main

import (
	"aws-ses-local-go/internal/repository"
	"aws-ses-local-go/internal/rest"
	"aws-ses-local-go/usecase/aws"
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
		log.Fatal("Error loading .env file")
	}
}

func main() {
	address := os.Getenv("SERVER_ADDRESS")

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.Renderer = t

	mailRepo := repository.NewMailRepository()
	awsSvc := aws.NewService(mailRepo)

	rest.NewAwsHandler(e, awsSvc)

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
