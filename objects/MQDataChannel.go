package objects

type MQDataChannel struct {
	body *MQMessage
}

func (M *MQDataChannel) Body() *MQMessage {
	return M.body
}

func (M *MQDataChannel) SetBody(body *MQMessage) {
	M.body = body
}


