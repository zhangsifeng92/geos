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

var BlockResourceExhaustedName = reflect.TypeOf(BlockResourceExhausted{}).Name()

type BlockResourceExhausted struct {
	_BlockValidateException
	Elog log.Messages
}

func NewBlockResourceExhausted(parent _BlockValidateException, message log.Message) *BlockResourceExhausted {
	return &BlockResourceExhausted{parent, log.Messages{message}}
}

func (e BlockResourceExhausted) Code() int64 {
	return 3030005
}

func (e BlockResourceExhausted) Name() string {
	return BlockResourceExhaustedName
}

func (e BlockResourceExhausted) What() string {
	return "Block exhausted allowed resources"
}

func (e *BlockResourceExhausted) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e BlockResourceExhausted) GetLog() log.Messages {
	return e.Elog
}

func (e BlockResourceExhausted) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e BlockResourceExhausted) DetailMessage() string {
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

func (e BlockResourceExhausted) String() string {
	return e.DetailMessage()
}

func (e BlockResourceExhausted) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3030005,
		Name: BlockResourceExhaustedName,
		What: "Block exhausted allowed resources",
	}

	return json.Marshal(except)
}

func (e BlockResourceExhausted) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*BlockResourceExhausted):
		callback(&e)
		return true
	case func(BlockResourceExhausted):
		callback(e)
		return true
	default:
		return false
	}
}
