package main

//export GOPROXY=https://goproxy.io
//性能测试 wrk -t8 -c200 -d30s --latency  "http://127.0.0.1:8080/sd/health"
//查看函数性能参数： go tool pprof http://127.0.0.1:8080/debug/pprof/profile
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go_restful_api/config"
	"go_restful_api/model"
	v "go_restful_api/pkg/version"
	"go_restful_api/router"
	"go_restful_api/router/middleware"
	"net/http"
	"os"
	"time"
)

var cfg = pflag.StringP("config", "c", "", "apiserver config file path")
var version = pflag.BoolP("version", "v", false, "show version info.")

func main() {

	pflag.Parse()

	//显示version
	if *version {
		v := v.Get()
		marshalled, err := json.MarshalIndent(&v, "", " ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(marshalled))
		return
	}

	//读取配置项
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	//init db
	model.DB.Init()
	defer model.DB.Close()

	//创建gin引擎
	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()

	//add middlewares
	router.Load(
		g,
		//middleware.Logging(), //todo 影响性能，可以去掉
		middleware.RequestId(), //todo 影响性能，可以去掉
	)

	//启动一个协程，健康检查
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()

	//启动https服务
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")

	if cert != "" && key != "" {
		go func() {
			log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
			log.Info(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
		}()
	}

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
