package main
//export GOPROXY=https://goproxy.io
import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go_restful_api/config"
	"go_restful_api/router"
	"net/http"
	"time"
)

var cfg = pflag.StringP("config", "c", "", "apiserver config file path")

func main()  {

	pflag.Parse()

	//读取配置项
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	//创建gin引擎
	gin.SetMode(viper.GetString("runmode"))
	g := gin.New()
	middlewares := []gin.HandlerFunc{}

	router.Load(g, middlewares...)

	//启动一个协程，健康检查
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()

	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		resp, err := http.Get(viper.GetString("url") + viper.GetString("addr") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}

	return errors.New("Cannot connect to the router.")
}