package websocket

import (
	"crash_exporter/models"
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ReportTraffic 订阅流量推送。
func (c *BaseCrash) RegisterTraffic() error {
	// 注册指标
	c.trafficMetrics = *promauto.NewCounterVec(prometheus.CounterOpts{
		Name: models.MerticTrafficName,
		Help: "The number of current up/down traffic",
	}, []string{models.TrafficLabelDirection})
	c.registry.MustRegister(c.trafficMetrics)

	// 初始化连接
	conn, ch, err := c.Connect("traffic", c.GetToken())
	if err != nil {
		c.logger.Errorf("Connect error: %v", err)
		return err
	}

	go func() {
		defer conn.Close()

		for {
			data := <-ch
			var obj models.WSTraffic
			err = json.Unmarshal(data, &obj)
			if err != nil {
				c.logger.Errorf("Parse Traffic error: %s\n%s", err, data)
			}
			// c.logger.Debugf("Traffic: %+v", obj)
			c.trafficMetrics.With(prometheus.Labels{models.TrafficLabelDirection: models.TrafficDirectionUp}).
				Add(float64(obj.Up))
			c.trafficMetrics.With(prometheus.Labels{models.TrafficLabelDirection: models.TrafficDirectionDown}).
				Add(float64(obj.Down))
		}
	}()
	return nil
}
