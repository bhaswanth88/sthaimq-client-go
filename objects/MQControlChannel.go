package objects


type MQControlChannel struct {
	messageType int // 1 - onOpen, 2-onError, 3-OnClose, 4, OnAuthenticate
	messagePayload *string
}

func (M *MQControlChannel) MessagePayload() *string {
	return M.messagePayload
}

func (M *MQControlChannel) SetMessagePayload(messagePayload *string) {
	M.messagePayload = messagePayload
}

func (M *MQControlChannel) MessageType() int {
	return M.messageType
}

func (M *MQControlChannel) SetMessageType(messageType int) {
	M.messageType = messageType
}
