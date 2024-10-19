package websocket

import (
	"crash_exporter/models"
	"encoding/json"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	reNormalLog = regexp.MustCompile(`(?P<src>[\d.]+:\d+) --> (?P<dst>[^\s]+) match (?P<match>[^\s]+) using (?P<type>[^\[]+)`)   //nolint:lll
	reErrorLog  = regexp.MustCompile(`dial (?P<type>[^\[\s]+)[^\)]+\) to (?P<dst>[^\s]+) error: [^:]+:[^:]+: (?P<match>[^\n]+)`) //nolint:lll
)

// RegisterLogs 订阅日志推送。
func (c *BaseCrash) RegisterLogs(logLevel string) error {
	// 注册指标
	c.logTargetMertics = *promauto.NewCounterVec(prometheus.CounterOpts{
		Name: models.MetricLogTargetName,
		Help: "The number of log(connection) from log system",
	}, []string{
		"LT", models.LogTargetLabelSrc, models.LogTargetLabelDst, models.LogTargetLabelMatch, models.LogTargetLabelType,
	})
	c.registry.MustRegister(c.logTargetMertics)

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
			c.PrintLog(&obj) // 输出到远程日志

			// 解析日志内容并上报访问计数
			logType := "normal"
			target := c.MatchNormalLogTarget(obj.Payload)
			if target == nil {
				logType = "error"
				target = c.MatchErrorLogTarget(obj.Payload)
				if target == nil {
					c.logger.Warnf("Unknown log: %+v", obj)
					continue
				}
			}
			c.logTargetMertics.With(prometheus.Labels{
				"LT":                       logType,
				models.LogTargetLabelSrc:   target.Src,
				models.LogTargetLabelDst:   target.Dst,
				models.LogTargetLabelMatch: target.Match,
				models.LogTargetLabelType:  target.Type,
			}).Inc()
		}
	}()
	return nil
}

func (c *BaseCrash) PrintLog(obj *models.WSLog) {
	switch obj.Type {
	case "debug":
		c.logger.Debugln(obj.Payload)
	case "info":
		c.logger.Infoln(obj.Payload)
	case "warning":
		c.logger.Warnln(obj.Payload)
	case "error":
		c.logger.Errorln(obj.Payload)
	default:
		c.logger.Warnf("Unknown level[%s]: %s", obj.Type, obj.Payload)
	}
}

func (c *BaseCrash) MatchNormalLogTarget(message string) *models.LogTarget {
	matches := reNormalLog.FindStringSubmatch(message)
	if matches == nil {
		return nil
	}

	result := &models.LogTarget{}
	names := reNormalLog.SubexpNames()
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
	return result
}

func (c *BaseCrash) MatchErrorLogTarget(message string) *models.LogTarget {
	matches := reErrorLog.FindStringSubmatch(message)
	if matches == nil {
		return nil
	}

	result := &models.LogTarget{}
	names := reErrorLog.SubexpNames()
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
	return result
}
