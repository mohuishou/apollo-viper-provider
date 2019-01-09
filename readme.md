# apollo-viper-provider

使用[apollo](https://github.com/ctripcorp/apollo)作为[viper](https://github.com/spf13/viper)的远端配置

## Usage
`InitViper` 传入一个包含配置的viper实例，
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

## Config

参考: [其它语言客户端接入指南
](https://github.com/ctripcorp/apollo/wiki/%E5%85%B6%E5%AE%83%E8%AF%AD%E8%A8%80%E5%AE%A2%E6%88%B7%E7%AB%AF%E6%8E%A5%E5%85%A5%E6%8C%87%E5%8D%97)

| Key | 必须 | 默认 | 说明 | 
|----|----|----|----|
| server | 是 | | 服务器地址|
| app_id | 是 | | app_id |
| cluster_name | 否  | default | 集群名|
| namespace_name | 否  | application | 命名空间 |
| ip | 否 | | ip |