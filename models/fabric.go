package models

type NodeInfo struct {
	// NodeInfoId   string
	KafkaAddr    string `json:"KafkaAddr"`
	PeerNodeName string `json:"PeerNodeName"`
	LeftStorage  string `json:"LeftStorage"`
	LocationX    string `json:"LocationX"`
	LocationY    string `json:"LocationY"`
}
