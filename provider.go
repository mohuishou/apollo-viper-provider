package apollo

import (
	"io"
	"net/http"

	"github.com/spf13/viper"
)

// ConfigProvider ConfigProvider
type ConfigProvider struct{}

// Init 初始化
func Init() {
	conf := viper.Sub("apollo")
	http.Get("")
}

// Get get config
func (ac ConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	return getConfig()
}

// Watch watch config
func (ac ConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	return getConfig()
}

// WatchChannel watch config channel
func (ac ConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	for {

	}
}
