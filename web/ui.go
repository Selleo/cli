package web

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UI(ctx context.Context, f embed.FS) error {
	e := echo.New()
	ff, err := fs.Sub(f, "packages/secrets-ui/dist")
	if err != nil {
		return err
	}
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(ff))))
	e.GET("/api/aws/secrets", func(c echo.Context) error {
		type response struct {
			Path    string
			Secrets []secret
		}
		data := response{
			Path: "selleo/dev/til",
			Secrets: []secret{
				{"API_KEY", "002ddaed-a985-11ed-99d9-4cedfb79ac39", true},
				{"AUTH_REDIRECT", "http://localhost:3000/auth/google/callback", false},
				{"DB_POOL", "32", false},
				{"GOOGLE_CLIENT_ID", "1234567890-abc123def456.apps.googleusercontent.com", false},
				{"GOOGLE_CLIENT_SECRET", "1234567890-abc123def456.apps.googleusercontent.com", false},
				{"SECRET_KEY_BASE", "PkFbZGp/z+1B/VZXA8I0H4RaExqlq675VLJ0OxjbRU2WkkNK50yQKorpbySeJFChz9YlDi1JmBSUed+X8idHzQ==", true},
				{"TZ", "Warsaw", false},
			},
		}

		return c.JSON(http.StatusOK, data)
	})
	return e.Start(":9999")
}

type secret struct {
	Key      string
	Value    string
	ReadOnly bool
}
