### vlog

### 集成
对logrus的简单封装,开箱即用,直接使用
```bash
// go 1.11启用go modules的版本,直接
go get github.com/golyu/vlog
```
其它版本的go,还需要下载对应的logrus依赖包

### 使用
```go
// 1 在main中初始化
_,err := Init("logs", "debug",365,24)
if err != nil {
	fmt.println("xx:", err.Error())
}
// 2 在程序任意地方,均可以调用
vlog.Debug("这是调试日志%s", "xx")
vlog.Warn("这是警告日志%d", 0)
vlog.Info("这是信息日志%s", "xx")
vlog.Error("这是错误日志%s", "xx")
```