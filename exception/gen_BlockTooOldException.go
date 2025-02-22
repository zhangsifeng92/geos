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

var BlockTooOldExceptionName = reflect.TypeOf(BlockTooOldException{}).Name()

type BlockTooOldException struct {
	_BlockValidateException
	Elog log.Messages
}

func NewBlockTooOldException(parent _BlockValidateException, message log.Message) *BlockTooOldException {
	return &BlockTooOldException{parent, log.Messages{message}}
}

func (e BlockTooOldException) Code() int64 {
	return 3030006
}

func (e BlockTooOldException) Name() string {
	return BlockTooOldExceptionName
}

func (e BlockTooOldException) What() string {
	return "Block is too old to push"
}

func (e *BlockTooOldException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e BlockTooOldException) GetLog() log.Messages {
	return e.Elog
}

func (e BlockTooOldException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e BlockTooOldException) DetailMessage() string {
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

func (e BlockTooOldException) String() string {
	return e.DetailMessage()
}

func (e BlockTooOldException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3030006,
		Name: BlockTooOldExceptionName,
		What: "Block is too old to push",
	}

	return json.Marshal(except)
}

func (e BlockTooOldException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*BlockTooOldException):
		callback(&e)
		return true
	case func(BlockTooOldException):
		callback(e)
		return true
	default:
		return false
	}
}
