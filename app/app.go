package app

import (
	"fmt"
	"hrserver/api"
	"hrserver/util/config"
	"hrserver/util/log"
	"net/http"

	"github.com/astaxie/beego/orm"
	gin "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const MasterName = "hrserver"

type ConfigData struct {
	Port    string
	LogDir  string
	IsHttps bool
	Cert    string
	Key     string
	SqlConn string
}

type App struct {
	handlers *api.Handlers
}

var g_ConfigData *ConfigData

func OnInitFlag(c *config.Config) (err error) {

	g_ConfigData = new(ConfigData)
	g_ConfigData.Port = c.GetString("port")
	g_ConfigData.LogDir = c.GetString("logdir")
	g_ConfigData.IsHttps = c.GetBool("https")
	g_ConfigData.Cert = c.GetString("https_cert")
	g_ConfigData.Key = c.GetString("https_key")
	g_ConfigData.SqlConn = c.GetString("db_url")

	if "" == g_ConfigData.Port || "" == g_ConfigData.LogDir {
		return fmt.Errorf("config not right")
	}
	return

}

func (app *App) OnStart(c *config.Config) error {

	if err := OnInitFlag(c); err != nil {
		return err
	}

	if _, err := log.NewLog(g_ConfigData.LogDir, MasterName, 0); err != nil {
		return err
	}

	if len(g_ConfigData.SqlConn) > 0 {
		if err := orm.RegisterDataBase("default", "mysql", g_ConfigData.SqlConn); err != nil {
			return err
		}
	}

	app.handlers = api.NewHandlers()

	router := gin.Default()

	router.Handle("GET", "health", func(c *gin.Context) { c.String(200, "success") })

	v1 := router.Group("/v1").Use(app.handlers.UserAuthrize)
	{
		v1.POST("/post", app.handlers.HandlerPost)
		v1.GET("/get", app.handlers.HandlerGet)
	}

	fmt.Println("Listen:", g_ConfigData.Port)

	if g_ConfigData.IsHttps {
		http.ListenAndServeTLS(":"+g_ConfigData.Port, g_ConfigData.Cert, g_ConfigData.Key, router)
	} else {
		http.ListenAndServe(":"+g_ConfigData.Port, router)
	}

	return nil
}

func (app *App) Shutdown() {
	app.handlers.OnClose()
	fmt.Println("server shutdown")
}
