package objects

type MQMessage struct {
	MsgType   string            `json:"msgType"`
	ClientId  string            `json:"clientId"`
	Payload   map[string]string `json:"payload"`
	MessageId string            `json:"messageId"`
}
