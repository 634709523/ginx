package znet

type Message struct {
	Id uint32 // 消息的ID
	DataLen uint32 // 消息的长度
	Data []byte // 消息的内容
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}

// 创建一个Message消息包
func NewMsgPackage(id uint32,data []byte)*Message{
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
