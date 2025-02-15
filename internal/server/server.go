package server

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/pidanou/c1-core/internal/ui"
)

func getFileSystem(useOS bool, embededFiles embed.FS) http.FileSystem {
	if useOS {
		return http.FS(os.DirFS("internal/ui/static"))
	}

	fsys, err := fs.Sub(embededFiles, "internal/ui/static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func Start() {
	e := echo.New()
	useOS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useOS, ui.StaticFiles))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	e.GET("/", GetHome)
	e.Logger.Fatal(e.Start(":1323"))
}
