# apollo-viper-provider

使用[apollo](https://github.com/ctripcorp/apollo)作为[viper](https://github.com/spf13/viper)的远端配置

## Usage
`InitViper` 传入一个只有包含`server`以及`app_id`的viper实例，
返回一个配置好远端的viper实例
```go
viper.SetDefault("apollo.server", "xxx")
viper.SetDefault("apollo.app_id", "xx")
v, err := InitViper(viper.Sub("apollo"))
// 如果需要监听配置状态
v.WatchRemoteConfigOnChannel()
// 获取配置信息
v.GetBool("test") // true
```

其余方法参考[viper](https://github.com/spf13/viper)

## TODO
- [ ] 支持多Namespace