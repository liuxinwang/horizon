package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"horizon/config"
	"horizon/model"
	"horizon/router"
	"horizon/tasks"
	"horizon/utils"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed static
var Static embed.FS

func init() {
	help := utils.HelpInit()
	config.InitConfig(help.ConfigFile)
	utils.LogInit()
	model.InitDb()
}

func main() {
	var r *gin.Engine
	if config.Conf.General.Environment == "dev" {
		r = router.InitRouter()
	} else {
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
	log.Infof("start horizon in :%d", config.Conf.General.Port)
	err := r.Run(fmt.Sprintf(":%d", config.Conf.General.Port))
	log.Fatal(err.Error())
}
