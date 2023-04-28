package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"horizon/config"
	"horizon/model"
	"horizon/router"
	"horizon/tasks"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var Static embed.FS

func main() {
	config.InitConfig()
	model.InitDb()
	var r *gin.Engine
	if config.Conf.General.Environment == "dev" {
		r = router.InitRouter()
	} else {
		// prod mode
		// gin release mode
		gin.SetMode(gin.ReleaseMode)
		// init prod router
		r = router.InitRouterPack()
		// load frontend
		sub, _ := fs.Sub(Static, "static")
		r.StaticFS("/static", http.FS(sub))
		// load index.html template
		staticTemplate := template.Must(template.New("").ParseFS(Static, "static/*.html"))
		r.SetHTMLTemplate(staticTemplate)
	}
	tasks.InitTasks()
	log.Printf(fmt.Sprintf("start horizon in :%d", config.Conf.General.Port))
	err := r.Run(fmt.Sprintf(":%d", config.Conf.General.Port))
	log.Fatal(err.Error())
}
