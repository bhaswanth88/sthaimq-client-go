package objects

type MQMessage struct {
	msgType   string            `json:"msgType"`
	clientId  string            `json:"clientId"`
	payload   map[string]string `json:"payload"`
	messageId string            `json:"messageId"`
}

func (M *MQMessage) MsgType() string {
	return M.msgType
}

func (M *MQMessage) SetMsgType(msgType string) {
	M.msgType = msgType
}

func (M *MQMessage) ClientId() string {
	return M.clientId
}

func (M *MQMessage) SetClientId(clientId string) {
	M.clientId = clientId
}

func (M *MQMessage) Payload() map[string]string {
	return M.payload
}

func (M *MQMessage) SetPayload(payload map[string]string) {
	M.payload = payload
}

func (M *MQMessage) MessageId() string {
	return M.messageId
}

func (M *MQMessage) SetMessageId(messageId string) {
	M.messageId = messageId
}
