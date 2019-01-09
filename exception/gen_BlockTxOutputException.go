// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/eosspark/eos-go/log"
)

// template type Exception(PARENT,CODE,WHAT)

var BlockTxOutputExceptionName = reflect.TypeOf(BlockTxOutputException{}).Name()

type BlockTxOutputException struct {
	_BlockValidateException
	Elog log.Messages
}

func NewBlockTxOutputException(parent _BlockValidateException, message log.Message) *BlockTxOutputException {
	return &BlockTxOutputException{parent, log.Messages{message}}
}

func (e BlockTxOutputException) Code() int64 {
	return 3030002
}

func (e BlockTxOutputException) Name() string {
	return BlockTxOutputExceptionName
}

func (e BlockTxOutputException) What() string {
	return "Transaction outputs in block do not match transaction outputs from applying block"
}

func (e *BlockTxOutputException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e BlockTxOutputException) GetLog() log.Messages {
	return e.Elog
}

func (e BlockTxOutputException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e BlockTxOutputException) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteString(" ")
	buffer.WriteString(e.Name())
	buffer.WriteString(": ")
	buffer.WriteString(e.What())
	buffer.WriteString("\n")
	for _, l := range e.Elog {
		buffer.WriteString("[")
		buffer.WriteString(l.GetMessage())
		buffer.WriteString("] ")
		buffer.WriteString(l.GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e BlockTxOutputException) String() string {
	return e.DetailMessage()
}

func (e BlockTxOutputException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3030002,
		Name: BlockTxOutputExceptionName,
		What: "Transaction outputs in block do not match transaction outputs from applying block",
	}

	return json.Marshal(except)
}

func (e BlockTxOutputException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*BlockTxOutputException):
		callback(&e)
		return true
	case func(BlockTxOutputException):
		callback(e)
		return true
	default:
		return false
	}
}
