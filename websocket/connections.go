package websocket

import (
	"crash_exporter/models"
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// RegisterConnections 订阅连接推送。
func (c *BaseCrash) RegisterConnections() error {
	// 注册指标
	c.trafficTotalMetrics = *promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: models.MerticTrafficTotalName,
		Help: "The total number of up/down traffic from crash start.",
	}, []string{models.TrafficTotalLabelDirection})
	c.registry.MustRegister(c.trafficTotalMetrics)

	// 初始化连接
	conn, ch, err := c.Connect("connections", c.GetToken())
	if err != nil {
		c.logger.Errorf("Connect error: %v", err)
		return err
	}

	go func() {
		defer conn.Close()

		for {
			data := <-ch
			var obj models.WSConnection
			err = json.Unmarshal(data, &obj)
			if err != nil {
				c.logger.Errorf("Parse Traffic error: %s\n%s", err, data)
			}
			// c.logger.Infof("Connections: %+v", obj)
			c.trafficTotalMetrics.
				With(prometheus.Labels{models.TrafficTotalLabelDirection: models.TrafficTotalDirectionUp}).
				Set(float64(obj.UpTotal))
			c.trafficTotalMetrics.
				With(prometheus.Labels{models.TrafficTotalLabelDirection: models.TrafficTotalDirectionDown}).
				Set(float64(obj.DownTotal))
		}
	}()
	return nil
}
