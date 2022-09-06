package znet


type Message struct {
	Id uint32
	DataLen uint32   //长度
	Data []byte
}
//new一个新的msg包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:     id,
		DataLen: uint32(len(data)),
		Data:   data,
	}
}


func (this *Message) GetMsgId() uint32{
	return this.Id
}
func (this *Message) GetMsgLen() uint32{
	return this.DataLen
}
func (this *Message) GetMsgData() []byte{
	return this.Data
}
func (this *Message) SetMsgId(id uint32) {
	this.Id = id
}
func (this *Message) SetMsgLen(len uint32) {
	this.DataLen = len
}
func (this *Message) SetMsgData(data []byte) {
	this.Data=data
}
