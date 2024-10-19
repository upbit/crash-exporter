package websocket

import (
	"crash_exporter/models"
	"encoding/json"
)

// RegisterLogs 订阅日志推送。
func (c *BaseCrash) RegisterLogs(logLevel string) error {
	// 注册指标

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
			}

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
	}()
	return nil
}
