package apollo

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

var conf = viper.Sub("apollo")

func init() {
	viper.SetDefault("apollo", map[string]string{
		"cluster_name":   "default",
		"namespace_name": "application",
	})
	if err := validate(); err != nil {
		panic(err)
	}
}

func validate() error {
	if conf.GetString("server") == "" {
		return nil
	}
	if conf.GetString("appid") == "" {
		return nil
	}
	return nil
}

func getConfig() (io.Reader, error) {
	url := fmt.Sprintf(
		"%s/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		conf.GetString("server"),
		conf.GetString("appid"),
		conf.GetString("cluster_name"),
		conf.GetString("namespace_name"),
		conf.GetString("release_key"),
		conf.GetString("ip"),
	)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	return resp.Body, err
}
