package models

type Metadata struct {
	Network    string `json:"network"`
	Type       string `json:"type"`
	Host       string `json:"host"`
	SrcIP      string `json:"sourceIP"`
	SrcPortStr string `json:"sourcePort"`
	DstIP      string `json:"destinationIP"`
	DstPortStr string `json:"destinationPort"`

	SrcPort int `json:"-"`
	DstPort int `json:"-"`
}

type Connection struct {
	ID          string   `json:"id"`
	Meta        Metadata `json:"metadata"`
	Upload      int      `json:"upload"`
	Download    int      `json:"download"`
	StartStr    string   `json:"start"`
	Chains      []string `json:"chains"`
	Rule        string   `json:"rule"`
	RulePayload string   `json:"rulePayload"`
}

// WSConnection 连接相关信息提取。
type WSConnection struct {
	UpTotal   int          `json:"uploadTotal"`
	DownTotal int          `json:"downloadTotal"`
	Conns     []Connection `json:"connections"`
}
