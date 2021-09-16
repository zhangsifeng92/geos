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

var BlockLogExceptionName = reflect.TypeOf(BlockLogException{}).Name()

type BlockLogException struct {
	_BlockLogException
	Elog log.Messages
}

func NewBlockLogException(parent _BlockLogException, message log.Message) *BlockLogException {
	return &BlockLogException{parent, log.Messages{message}}
}

func (e BlockLogException) Code() int64 {
	return 3190000
}

func (e BlockLogException) Name() string {
	return BlockLogExceptionName
}

func (e BlockLogException) What() string {
	return "Block log exception"
}

func (e *BlockLogException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e BlockLogException) GetLog() log.Messages {
	return e.Elog
}

func (e BlockLogException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e BlockLogException) DetailMessage() string {
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

func (e BlockLogException) String() string {
	return e.DetailMessage()
}

func (e BlockLogException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3190000,
		Name: BlockLogExceptionName,
		What: "Block log exception",
	}

	return json.Marshal(except)
}

func (e BlockLogException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*BlockLogException):
		callback(&e)
		return true
	case func(BlockLogException):
		callback(e)
		return true
	default:
		return false
	}
}
