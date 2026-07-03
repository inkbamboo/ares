// Package sls 提供阿里云 SLS（Simple Log Service）日志服务的 logrus Hook 实现。
package sls

import (
	"fmt"

	slssdk "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/sirupsen/logrus"
)

// 确保 SLSHook 实现了 logrus.Hook 接口
var _ logrus.Hook = (*SLSHook)(nil)

// OptionFunc 是 SLSHook 配置选项的函数类型。
type OptionFunc func(*Option)

// Option 存储 SLSHook 的可配置参数。
type Option struct {
	project  string
	logstore string
	topic    string
	source   string
}

// SetProject 设置 SLS 项目名称。
func SetProject(name string) OptionFunc {
	return func(o *Option) {
		o.project = name
	}
}

// SetLogstore 设置 SLS 日志库名称。
func SetLogstore(name string) OptionFunc {
	return func(o *Option) {
		o.logstore = name
	}
}

// SetTopic 设置日志主题。
func SetTopic(name string) OptionFunc {
	return func(o *Option) {
		o.topic = name
	}
}

// SetSource 设置日志来源标识。
func SetSource(name string) OptionFunc {
	return func(o *Option) {
		o.source = name
	}
}

// NewSLSHook 创建阿里云 SLS 日志 Hook 实例。
// 使用 Producer 客户端将日志发送到 SLS 服务。
func NewSLSHook(accessKeyId, accessKeySecret, endpoint, allowLogLevel string, opts ...OptionFunc) *SLSHook {
	opt := &Option{}
	if len(opts) > 0 {
		for _, fun := range opts {
			if fun != nil {
				fun(opt)
			}
		}
	}
	pc := producer.GetDefaultProducerConfig()
	pc.IsJsonType = true
	pc.Endpoint = endpoint
	pc.CredentialsProvider = slssdk.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")
	pc.AllowLogLevel = allowLogLevel
	p, err := producer.NewProducer(pc)
	if err != nil {
		fmt.Println("producer.NewProducer err: ", err)
		return &SLSHook{opt, nil}
	}
	p.Start()
	return &SLSHook{opt, p}
}

// SLSHook 实现 logrus.Hook 接口，将日志发送到阿里云 SLS。
type SLSHook struct {
	opt      *Option
	producer *producer.Producer
}

// Fire 在日志触发时调用，将日志条目发送到 SLS。
func (hook *SLSHook) Fire(entry *logrus.Entry) error {
	if hook.producer == nil {
		return fmt.Errorf("SLS producer is nil, hook not initialized")
	}
	ts := uint32(entry.Time.Unix())
	var contents []*slssdk.LogContent
	for key, value := range entry.Data {
		keyCopy := key
		valueCopy := value
		k := fmt.Sprint(keyCopy)
		v := fmt.Sprint(valueCopy)
		contents = append(contents, &slssdk.LogContent{
			Key:   &k,
			Value: &v,
		})
	}

	timeKey := "time"
	timeContent := entry.Time.Format("2006-01-02 15:04:05")
	contents = append(contents, &slssdk.LogContent{
		Key:   &timeKey,
		Value: &timeContent,
	})

	levelKey := "level"
	levelContent := entry.Level.String()
	contents = append(contents, &slssdk.LogContent{
		Key:   &levelKey,
		Value: &levelContent,
	})

	msgKey := "message"
	msgContent := entry.Message
	contents = append(contents, &slssdk.LogContent{
		Key:   &msgKey,
		Value: &msgContent,
	})

	log := &slssdk.Log{
		Time:     &ts,
		Contents: contents,
	}

	err := hook.producer.SendLog(hook.opt.project, hook.opt.logstore, hook.opt.topic, hook.opt.source, log)
	if err != nil {
		return err
	}
	return nil
}

// Levels 返回 Hook 支持的日志级别，此处为所有级别。
func (hook *SLSHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// SafeClose 安全关闭 SLS Producer 客户端。
func (hook *SLSHook) SafeClose() {
	if hook.producer != nil {
		hook.producer.SafeClose()
	}
}

// Close 关闭 SLS Producer 客户端，timeoutMs 为超时时间（毫秒）。
func (hook *SLSHook) Close(timeoutMs int64) error {
	if hook.producer != nil {
		return hook.producer.Close(timeoutMs)
	}
	return nil
}
