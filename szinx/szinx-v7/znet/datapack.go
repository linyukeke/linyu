package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"szinx/utils"
	"szinx/ziface"
)

//封包拆包的模块
type DataPack struct {}

func NewDataPack() ziface.IDataPack{
	return &DataPack{}
}

func (this *DataPack)GetHeadLen() uint32 {
	return 8   //uint32是4字节，id和datalen也是uint32
}
//封包方法
func (this *DataPack) Pack(msg ziface.IMessage)([]byte,error)  {
	//创建一个buf，存放字节的缓冲
	databuff := bytes.NewBuffer([]byte{})   //声明一个byte类型的切片并初始化
	//databuff.Bytes()是一个二进制的序列化
	//将datalen写入databuff
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetMsgLen());err != nil{
		//有大段和小段两种模式，读的时候必须一致，此处选择小段
		fmt.Println("write len secces")
		return nil, err
	}
	//将msgid写入databuff
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetMsgId());err != nil{
		fmt.Println("write id fail")
		return nil, err
	}
	//将data数据写入databuff
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetMsgData());err != nil{
		fmt.Println("write data success")
		return nil, err
	}
	return databuff.Bytes(),nil
}
//拆包方法，只需要读包的head信息，再根据head信息里的data长度，在进行一次读取
func (this *DataPack) Unpack(binaryData []byte) (ziface.IMessage,error)  {
	databuff := bytes.NewReader(binaryData)

	msg := &Message{}

	//读datalen
	if err := binary.Read(databuff, binary.LittleEndian,&msg.DataLen); err != nil {
		return nil, err
	}
	//读msgid
	if err := binary.Read(databuff,binary.LittleEndian,&msg.Id);err != nil{
		return nil, err
	}
	//判断datalen是否超出了最大的包长度
	if (utils.G.MaxPackageSize >0 && msg.DataLen >utils.G.MaxPackageSize){
		return nil,errors.New("too large msg data recv!")
	}
	return msg,nil    //此时msg中只有长度和id

}



