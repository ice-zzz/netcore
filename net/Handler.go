package network

type Handler struct {
	funList map[uint16]RecvHandler
}

type RecvHandler func(message *MessageData) *MessageData

func (h *Handler) Init() {
	h.funList = make(map[uint16]RecvHandler)
}

func (h *Handler) AddHandler(messageType uint16, fun RecvHandler) {
	h.funList[messageType] = fun
}

type MessageData struct {
	MessageType uint16
	Message     []byte
}

func (h *Handler) Execute(data *MessageData) *MessageData {
	if v, ok := h.funList[data.MessageType]; ok {
		return v(data)
	}
	return &MessageData{
		MessageType: 500,
		Message:     []byte("非法消息类型"),
	}
}
