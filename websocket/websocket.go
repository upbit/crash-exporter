// Package websocket 实现与Crash的websocket通信处理
package websocket

import (
	"crash_exporter/models"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// BaseCrash Crash的websocket通信。
type BaseCrash struct {
	Addr  string
	Token string

	logger              *logrus.Logger
	registry            *prometheus.Registry
	trafficMetrics      prometheus.CounterVec
	trafficTotalMetrics prometheus.GaugeVec
	logTargetMertics    prometheus.CounterVec
}

// NewCrash 初始化Crash的链接。
func NewCrash(addr string, token string, reg *prometheus.Registry, logger *logrus.Logger) (*BaseCrash, error) {
	crash := &BaseCrash{
		Addr:     addr,
		Token:    token,
		logger:   logger,
		registry: reg,
	}
	return crash, nil
}

func (c *BaseCrash) Registers(logLevel string) error {
	if err := c.RegisterTraffic(); err != nil {
		return err
	}
	if err := c.RegisterLogs(logLevel); err != nil {
		return err
	}
	if err := c.RegisterConnections(); err != nil {
		return err
	}
	return nil
}

// Connect 连接目标websocket，用于订阅数据。
func (c *BaseCrash) Connect(endpoint, params string) (*websocket.Conn, <-chan []byte, error) {
	url := fmt.Sprintf("ws://%s/%s?%s", c.Addr, endpoint, params)
	c.logger.Debugf("connecting to %s", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		c.logger.Errorf("Dial(%s) error: %v", url, err)
	}
	conn.SetReadLimit(models.DefaultWSReadLimit)

	ch := make(chan []byte, models.DefaultChannelSize)
	go func() {
		for {
			_, msg, errRead := conn.ReadMessage()
			if errRead != nil {
				c.logger.Errorf("Read(%s) failed: %v", url, errRead)

				conn = c.reconnect(url)
				if conn == nil {
					c.logger.Fatalf("Retry reconnect(%s) failed", url)
				}
				c.logger.Infof("reconnect() success: %s", url)
			}
			ch <- msg
		}
	}()
	return conn, ch, nil
}

func (c *BaseCrash) reconnect(url string) *websocket.Conn {
	for i := range models.DefaultReconnectNum {
		time.Sleep(models.DefaultReconnectWait)
		c.logger.Warnf("[%d] Try connect to %s", i+1, url)

		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			c.logger.Errorf("Dial(%s) error: %v", url, err)
			continue
		}
		conn.SetReadLimit(models.DefaultWSReadLimit)
		return conn
	}
	return nil
}

func (c *BaseCrash) GetToken() string {
	return "token=" + c.Token
}
