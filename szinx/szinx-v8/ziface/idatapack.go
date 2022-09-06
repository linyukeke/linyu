package ziface

//拆包，封包，模块，直接面向tcp连接中的数据流，用于处理tcp粘包问题
type IDataPack interface {
	GetHeadLen() uint32  //获取头部长度
	Pack(msg IMessage)([]byte,error)    //封包方法
	Unpack([]byte) (IMessage,error)       //拆包方法
}



