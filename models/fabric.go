package models

type NodeInfo struct {
	// NodeInfoId   string
	PeerNodeName string `json:"PeerNodeName"`
	LeftStorage  string `json:"LeftStorage"`
	LocationX    string `json:"LocationX"`
	LocationY    string `json:"LocationY"`
}
