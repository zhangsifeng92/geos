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

var TooManyTxAtOnceName = reflect.TypeOf(TooManyTxAtOnce{}).Name()

type TooManyTxAtOnce struct {
	_TransactionException
	Elog log.Messages
}

func NewTooManyTxAtOnce(parent _TransactionException, message log.Message) *TooManyTxAtOnce {
	return &TooManyTxAtOnce{parent, log.Messages{message}}
}

func (e TooManyTxAtOnce) Code() int64 {
	return 3040012
}

func (e TooManyTxAtOnce) Name() string {
	return TooManyTxAtOnceName
}

func (e TooManyTxAtOnce) What() string {
	return "Pushing too many transactions at once"
}

func (e *TooManyTxAtOnce) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e TooManyTxAtOnce) GetLog() log.Messages {
	return e.Elog
}

func (e TooManyTxAtOnce) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e TooManyTxAtOnce) DetailMessage() string {
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
		buffer.WriteString("]")
		buffer.WriteString("\n")
		buffer.WriteString(l.GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e TooManyTxAtOnce) String() string {
	return e.DetailMessage()
}

func (e TooManyTxAtOnce) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3040012,
		Name: TooManyTxAtOnceName,
		What: "Pushing too many transactions at once",
	}

	return json.Marshal(except)
}

func (e TooManyTxAtOnce) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*TooManyTxAtOnce):
		callback(&e)
		return true
	case func(TooManyTxAtOnce):
		callback(e)
		return true
	default:
		return false
	}
}