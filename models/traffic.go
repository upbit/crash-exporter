package models

// WSTraffic WebSocket返回的流量结构。
type WSTraffic struct {
	Up   int `json:"up"`
	Down int `json:"down"`
}
