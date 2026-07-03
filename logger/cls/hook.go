// Package cls 提供腾讯云 CLS 日志服务的 logrus Hook 实现（对外暴露版本）。
package cls

import (
	"fmt"
	"github.com/sirupsen/logrus"
	clssdk "github.com/tencentcloud/tencentcloud-cls-sdk-go"
)

// 确保 CLSHook 实现了 logrus.Hook 接口
var _ logrus.Hook = (*CLSHook)(nil)

// OptionFunc 是 CLSHook 配置选项的函数类型。
type OptionFunc func(*Option)

// Option 存储 CLSHook 的可配置参数。
type Option struct {
	topic string
}

// SetTopic 设置 CLS 日志主题名称。
func SetTopic(name string) OptionFunc {
	return func(o *Option) {
		o.topic = name
	}
}

// NewCLSHook 创建腾讯云 CLS 日志 Hook 实例。
func NewCLSHook(accessKeyId, accessKeySecret, endpoint, allowLogLevel string, opts ...OptionFunc) *CLSHook {
	opt := &Option{}
	if len(opts) > 0 {
		for _, fun := range opts {
			if fun != nil {
				fun(opt)
			}
		}
	}
	pc := clssdk.GetDefaultAsyncProducerClientConfig()
	pc.AccessKeyID = accessKeyId
	pc.AccessKeySecret = accessKeySecret
	pc.Endpoint = endpoint

	p, err := clssdk.NewAsyncProducerClient(pc)
	if err != nil {
		fmt.Println("clssdk.NewAsyncProducerClient err: ", err)
		return &CLSHook{}
	}
	p.Start()
	return &CLSHook{opt, p}
}

// CLSHook 实现 logrus.Hook 接口，将日志发送到腾讯云 CLS。
type CLSHook struct {
	opt      *Option
	producer *clssdk.AsyncProducerClient
}

// Fire 在日志触发时调用，将日志条目发送到 CLS。
func (hook *CLSHook) Fire(entry *logrus.Entry) error {
	if hook.producer == nil {
		return fmt.Errorf("CLS producer is nil, hook not initialized")
	}
	var out = map[string]string{
		"time":    entry.Time.Format("2006-01-02 15:04:05"),
		"level":   entry.Level.String(),
		"message": entry.Message,
	}

	for key, value := range entry.Data {
		k := fmt.Sprint(key)
		v := fmt.Sprint(value)
		out[k] = v
	}

	return hook.producer.SendLog(
		hook.opt.topic,
		clssdk.NewCLSLog(
			entry.Time.Unix(),
			out,
		), nil,
	)
}

// Levels 返回 Hook 支持的日志级别。
func (hook *CLSHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Close 关闭 CLS Producer 客户端。
func (hook *CLSHook) Close(timeoutMs int64) error {
	if hook.producer != nil {
		return hook.producer.Close(timeoutMs)
	}
	return nil
}
