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

var CfaInsideGeneratedTxName = reflect.TypeOf(CfaInsideGeneratedTx{}).Name()

type CfaInsideGeneratedTx struct {
	_TransactionException
	Elog log.Messages
}

func NewCfaInsideGeneratedTx(parent _TransactionException, message log.Message) *CfaInsideGeneratedTx {
	return &CfaInsideGeneratedTx{parent, log.Messages{message}}
}

func (e CfaInsideGeneratedTx) Code() int64 {
	return 3040010
}

func (e CfaInsideGeneratedTx) Name() string {
	return CfaInsideGeneratedTxName
}

func (e CfaInsideGeneratedTx) What() string {
	return "Context free action is not allowed inside generated transaction"
}

func (e *CfaInsideGeneratedTx) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e CfaInsideGeneratedTx) GetLog() log.Messages {
	return e.Elog
}

func (e CfaInsideGeneratedTx) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e CfaInsideGeneratedTx) DetailMessage() string {
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

func (e CfaInsideGeneratedTx) String() string {
	return e.DetailMessage()
}

func (e CfaInsideGeneratedTx) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3040010,
		Name: CfaInsideGeneratedTxName,
		What: "Context free action is not allowed inside generated transaction",
	}

	return json.Marshal(except)
}

func (e CfaInsideGeneratedTx) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*CfaInsideGeneratedTx):
		callback(&e)
		return true
	case func(CfaInsideGeneratedTx):
		callback(e)
		return true
	default:
		return false
	}
}
