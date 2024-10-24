package websocket

import (
	"crash_exporter/models"
	"encoding/json"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

var (
	reNormalLog = regexp.MustCompile(`(?P<src>[\d.]+):\d+ --> (?P<dst>[^\s]+) match (?P<match>[^\s]+) using (?P<type>[^\[]+)`)   //nolint:lll
	reErrorLog  = regexp.MustCompile(`dial (?P<type>[^\[\s]+)[^\)]+\) to (?P<dst>[^\s]+) error: [^:]+:[^:]+: (?P<match>[^\n]+)`) //nolint:lll
	reDNSLog    = regexp.MustCompile(`\[(?P<type>DNS)\] (resolve ){0,1}(?P<dst>[^\s]+) (error:|-->) (?P<match>.*$)`)
)

// RegisterLogs 订阅日志推送。
func (c *BaseCrash) RegisterLogs(logLevel string) error {
	// 注册指标
	c.logMertics = *promauto.NewCounterVec(prometheus.CounterOpts{
		Name: models.MetricLogName,
		Help: "The number of log(connection) from log system",
	}, []string{models.LogLabelLevel})
	c.registry.MustRegister(c.logMertics)

	// 初始化连接
	conn, ch, err := c.Connect("logs", c.GetToken()+"&level="+logLevel)
	if err != nil {
		c.logger.Errorf("Connect error: %v", err)
		return err
	}

	go func() {
		defer conn.Close()

		for {
			data := <-ch
			var obj models.WSLog
			err = json.Unmarshal(data, &obj)
			if err != nil {
				c.logger.Errorf("Parse Traffic error: %s\n%s", err, data)
				continue
			}
			c.logMertics.With(prometheus.Labels{models.LogLabelLevel: obj.Type}).Inc()

			// 解析日志内容并上报
			target, logType := c.MatchLogTarget(obj.Payload)
			if target == nil {
				c.logger.Warnf("Unknown log: %+v", obj)
				continue
			}
			finalLogger := c.logger.WithFields(logrus.Fields{
				"log_type": logType,
				"src":      target.Src,
				"dst":      target.Dst,
				"match":    target.Match,
				"type":     target.Type,
			})
			switch obj.Type {
			case "debug":
				finalLogger.Debugln(obj.Payload)
			case "info":
				finalLogger.Infoln(obj.Payload)
			case "warning":
				finalLogger.Warnln(obj.Payload)
			case "error":
				finalLogger.Errorln(obj.Payload)
			default:
				finalLogger.Warnf("Unknown level[%s]: %s", obj.Type, obj.Payload)
			}
		}
	}()
	return nil
}

func (c *BaseCrash) MatchLogTarget(message string) (*models.LogTarget, string) {
	var matches []string
	var names []string
	var logType string
	for {
		// DNS
		matches = reDNSLog.FindStringSubmatch(message)
		if matches != nil {
			names = reDNSLog.SubexpNames()
			logType = models.LogTypeDNS
			break
		}

		// Normal
		matches = reNormalLog.FindStringSubmatch(message)
		if matches != nil {
			names = reNormalLog.SubexpNames()
			logType = models.LogTypeNormal
			break
		}

		// Error
		matches = reErrorLog.FindStringSubmatch(message)
		if matches != nil {
			names = reErrorLog.SubexpNames()
			logType = models.LogTypeError
			break
		}

		return nil, "" //nolint:staticcheck // no match
	}

	result := &models.LogTarget{
		Type: "DIRECT",
	}
	for i, value := range matches {
		switch names[i] {
		case "src":
			result.Src = value
		case "dst":
			result.Dst = value
		case "match":
			result.Match = value
		case "type":
			result.Type = value
		}
	}
	return result, logType
}
