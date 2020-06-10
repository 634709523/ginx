package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct {
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

//封包方法(压缩数据)
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte,error) {
	var (
		err error
	)
	// 创建一个存在[]byte字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	// 写dataLen
	if err = binary.Write(dataBuff,binary.LittleEndian,msg.GetDataLen());err != nil{
		return nil,err
	}
	// 写msgId
	if err = binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgId());err != nil{
		return nil,err
	}
	// 写数据
	if err = binary.Write(dataBuff,binary.LittleEndian,msg.GetData());err != nil{
		return nil,err
	}

	return dataBuff.Bytes(),err
}

func (dp *DataPack) UnPack(binaryByte []byte) (ziface.IMessage, error) {
	var (
		err error
		msg *Message
	)
	// 创建一个存在[]byte字节的缓冲
	dataBuff := bytes.NewReader(binaryByte)
	msg = &Message{}
	if err = binary.Read(dataBuff,binary.LittleEndian,&msg.DataLen);err != nil{
		return nil,err
	}
	if err = binary.Read(dataBuff,binary.LittleEndian,&msg.Id);err != nil{
		return nil,err
	}
	// 判断数据是否超过参数中的最大数据包大小
	if (utils.GlobalObject.MaxPacketSize > 0) && utils.GlobalObject.MaxPacketSize < msg.DataLen{
		return nil, errors.New("Too Large msg data received")
	}
	return msg,nil
}

func NewDataPack()*DataPack{
	return &DataPack{}
}

