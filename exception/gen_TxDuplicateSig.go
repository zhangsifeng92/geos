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

var TxDuplicateSigName = reflect.TypeOf(TxDuplicateSig{}).Name()

type TxDuplicateSig struct {
	_AuthorizationException
	Elog log.Messages
}

func NewTxDuplicateSig(parent _AuthorizationException, message log.Message) *TxDuplicateSig {
	return &TxDuplicateSig{parent, log.Messages{message}}
}

func (e TxDuplicateSig) Code() int64 {
	return 3090001
}

func (e TxDuplicateSig) Name() string {
	return TxDuplicateSigName
}

func (e TxDuplicateSig) What() string {
	return "Duplicate signature included"
}

func (e *TxDuplicateSig) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e TxDuplicateSig) GetLog() log.Messages {
	return e.Elog
}

func (e TxDuplicateSig) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e TxDuplicateSig) DetailMessage() string {
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

func (e TxDuplicateSig) String() string {
	return e.DetailMessage()
}

func (e TxDuplicateSig) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3090001,
		Name: TxDuplicateSigName,
		What: "Duplicate signature included",
	}

	return json.Marshal(except)
}

func (e TxDuplicateSig) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*TxDuplicateSig):
		callback(&e)
		return true
	case func(TxDuplicateSig):
		callback(e)
		return true
	default:
		return false
	}
}
