package models

// WSLog 日志相关信息提取。
type WSLog struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type LogTarget struct {
	Src   string `json:"src"`
	Dst   string `json:"dst"`
	Match string `json:"match"`
	Type  string `json:"type"`
}
