// Package models 内部结构与常理定义
package models

import "time"

const (
	DefaultReconnectNum  = 5 * 4            // 最大重连次数
	DefaultReconnectWait = 15 * time.Second // 重连等待时间
	DefaultChannelSize   = 10               // DefaultChannelSize 返回消息的队列长度
	DefaultWSReadLimit   = 10 * 1024 * 1024 // DefaultWSReadLimit 默认WebSocket的最大读取长度
)

const (
	MerticTrafficName     = "crash_traffic"
	TrafficLabelDirection = "direction"
	TrafficDirectionUp    = "up"
	TrafficDirectionDown  = "down"

	MetricLogTargetName = "crash_log_request_target"
	LogTargetLabelSrc   = "src"
	LogTargetLabelDst   = "dst"
	LogTargetLabelMatch = "match"
	LogTargetLabelType  = "type"

	MerticTrafficTotalName     = "crash_traffic_total"
	TrafficTotalLabelDirection = "direction"
	TrafficTotalDirectionUp    = "up"
	TrafficTotalDirectionDown  = "down"
)
