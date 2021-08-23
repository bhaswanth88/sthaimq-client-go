package objects


type MQControlChannel struct {
	MessageType int // 1 - onOpen, 2-onError, 3-OnClose, 4, OnAuthenticate, 5-OnData
	Body *MQMessage

}
