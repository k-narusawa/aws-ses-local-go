package main

import (
	"aws-ses-local-go/config"
	"aws-ses-local-go/internal/controllers/rest"
	"aws-ses-local-go/internal/gateways/dao"
	"aws-ses-local-go/internal/gateways/query"
	"aws-ses-local-go/internal/gateways/repository"
	v1 "aws-ses-local-go/usecase/aws/v1"
	v2 "aws-ses-local-go/usecase/aws/v2"
	"aws-ses-local-go/usecase/mail"
	"aws-ses-local-go/views"
	"html/template"
	"io"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

const (
	defaultAddress = ":8080"
)

func main() {
	e := echo.New()

	registerViews(e)

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${uri}  ${latency_human}\n",
		Output: e.Logger.Output(),
	}))

	db := config.DBConnect()

	mailDao := dao.NewMailDao(db)
	mailRepo := repository.NewMailRepository(*mailDao)

	v1Svc := v1.NewService(mailRepo)
	v2Svc := v2.NewService(mailRepo)
	mailSvc := mail.NewService(mailRepo)
	mailQSvc := query.NewMailDtoQueryService(*mailDao)
	cntQSvc := query.NewCountQueryService(*mailDao)

	rest.NewAwsHandler(e, v1Svc, v2Svc)
	rest.NewMailHandler(e, mailQSvc, cntQSvc, mailSvc)

	e.GET("/health", healthCheck)

	log.Fatal(e.Start(defaultAddress))
}

// フロントエンドの静的ファイルはバイナリに埋め込んだものを配信する
func registerViews(e *echo.Echo) {
	e.Renderer = &Template{
		templates: template.Must(template.ParseFS(views.FS, "*.html")),
	}
	e.StaticFS("/", views.FS)
}

func healthCheck(c echo.Context) error {
	return c.String(200, "OK")
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
