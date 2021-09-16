// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/zhangsifeng92/geos/log"
)

// template type Exception(PARENT,CODE,WHAT)

var WrongProducerName = reflect.TypeOf(WrongProducer{}).Name()

type WrongProducer struct {
	_BlockValidateException
	Elog log.Messages
}

func NewWrongProducer(parent _BlockValidateException, message log.Message) *WrongProducer {
	return &WrongProducer{parent, log.Messages{message}}
}

func (e WrongProducer) Code() int64 {
	return 3030009
}

func (e WrongProducer) Name() string {
	return WrongProducerName
}

func (e WrongProducer) What() string {
	return "Block is not signed by expected producer"
}

func (e *WrongProducer) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e WrongProducer) GetLog() log.Messages {
	return e.Elog
}

func (e WrongProducer) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e WrongProducer) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteByte(' ')
	buffer.WriteString(e.Name())
	buffer.Write([]byte{':', ' '})
	buffer.WriteString(e.What())
	buffer.WriteByte('\n')
	for _, l := range e.Elog {
		buffer.WriteByte('[')
		buffer.WriteString(l.GetMessage())
		buffer.Write([]byte{']', ' '})
		buffer.WriteString(l.GetContext().String())
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func (e WrongProducer) String() string {
	return e.DetailMessage()
}

func (e WrongProducer) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3030009,
		Name: WrongProducerName,
		What: "Block is not signed by expected producer",
	}

	return json.Marshal(except)
}

func (e WrongProducer) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*WrongProducer):
		callback(&e)
		return true
	case func(WrongProducer):
		callback(e)
		return true
	default:
		return false
	}
}
