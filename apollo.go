package apollo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("apollo.cluster_name", "default")
	viper.SetDefault("apollo.namespace_name", "application")
}

// Notification apollo的通知结构
type Notification struct {
	NamespaceName  string `json:"namespaceName"`
	NotificationID int    `json:"notificationId"`
}

// Apollo 使用apollo作为viper的Remote Provider
type Apollo struct {
	ClusterName   string         `mapstructure:"cluster_name"`
	NamespaceName string         `mapstructure:"namespace_name"`
	Server        string         `mapstructure:"server"`
	AppID         string         `mapstructure:"app_id"`
	ReleaseKey    string         `mapstructure:"release_key"`
	IP            string         `mapstructure:"ip"`
	Notifications []Notification `mapstructure:"notifications"`
}

// configResponse 返回的配置结构
type configResponse struct {
	Configurations json.RawMessage `json:"configurations"`
	ReleaseKey     string          `json:"releaseKey"`
	Cluster        string          `json:"cluster"`
	NamespaceName  string          `json:"namespaceName"`
	AppID          string          `json:"appId"`
}

// InitViper 初始化Viper
func InitViper(v *viper.Viper) (*viper.Viper, error) {
	// 初始化
	a, err := New(v)
	if err != nil {
		return nil, err
	}
	viper.RemoteConfig = a

	v = viper.New()
	err = v.AddRemoteProvider("consul", a.Server, a.AppID)
	if err != nil {
		return nil, err
	}
	v.SetConfigType("json")
	if err := v.ReadRemoteConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

// New 新建一个远程配置
func New(v *viper.Viper) (*Apollo, error) {
	a := Apollo{}
	err := v.Unmarshal(&a)
	a.Notifications = []Notification{
		{
			NamespaceName:  a.NamespaceName,
			NotificationID: -1,
		},
	}
	return &a, err
}

func (a Apollo) Get(rp viper.RemoteProvider) (io.Reader, error) {
	b, err := a.getNoCache()
	r := bytes.NewReader(b)
	return r, err
}

func (a Apollo) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	b, err := a.getWithCache()
	r := bytes.NewReader(b)
	return r, err
}

func (a Apollo) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	quit := make(chan bool)
	viperResponseCh := make(chan *viper.RemoteResponse)
	go func(vc chan<- *viper.RemoteResponse, quit <-chan bool) {
		for {
			select {
			case <-quit:
				return
			default:
				// check 配置信息是否有变化
				isChange, err := a.getNotifications()
				if err != nil {
					vc <- &viper.RemoteResponse{Error: err}
					return
				}

				// 配置有变化，获取最新的配置
				if isChange {
					value, err := a.getNoCache()
					vc <- &viper.RemoteResponse{Value: value, Error: err}
				}
			}
		}
	}(viperResponseCh, quit)
	return viperResponseCh, quit
}

func (a Apollo) notificationsStr() string {
	b, err := json.Marshal(a.Notifications)
	if err != nil {
		return ""
	}
	return string(b)
}

func (a *Apollo) getWithCache() ([]byte, error) {
	uri := fmt.Sprintf(
		"%s/configfiles/json/%s/%s/%s",
		a.Server,
		a.AppID,
		a.ClusterName,
		a.NamespaceName,
	)

	params := url.Values{}
	if a.IP != "" {
		params.Add("ip", a.IP)
		uri = uri + "?" + params.Encode()
	}
	return a.get(uri)
}

func (a *Apollo) getNoCache() ([]byte, error) {
	uri := fmt.Sprintf(
		"%s/configs/%s/%s/%s",
		a.Server,
		a.AppID,
		a.ClusterName,
		a.NamespaceName,
	)

	params := url.Values{}
	if a.IP != "" {
		params.Add("ip", a.IP)
		uri = uri + "?" + params.Encode()
	}

	return a.get(uri)
}

// get 通过接口获取配置信息
func (a *Apollo) get(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var configResp configResponse
	if err := json.Unmarshal(b, &configResp); err != nil {
		return nil, err
	}

	a.ReleaseKey = configResp.ReleaseKey
	return configResp.Configurations, nil
}

// getNotifications 获取通知信息
func (a *Apollo) getNotifications() (bool, error) {
	params := url.Values{}
	params.Add("appId", a.AppID)
	params.Add("cluster", a.ClusterName)
	params.Add("notifications", a.notificationsStr())
	resp, err := http.Get(fmt.Sprintf(
		"%s/notifications/v2?%s",
		a.Server,
		params.Encode(),
	))
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		return false, nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, &a.Notifications)
	return true, err
}
