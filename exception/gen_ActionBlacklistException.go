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

var ActionBlacklistExceptionName = reflect.TypeOf(ActionBlacklistException{}).Name()

type ActionBlacklistException struct {
	_WhitelistBlacklistException
	Elog log.Messages
}

func NewActionBlacklistException(parent _WhitelistBlacklistException, message log.Message) *ActionBlacklistException {
	return &ActionBlacklistException{parent, log.Messages{message}}
}

func (e ActionBlacklistException) Code() int64 {
	return 3130005
}

func (e ActionBlacklistException) Name() string {
	return ActionBlacklistExceptionName
}

func (e ActionBlacklistException) What() string {
	return "Action to execute is on the blacklist"
}

func (e *ActionBlacklistException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ActionBlacklistException) GetLog() log.Messages {
	return e.Elog
}

func (e ActionBlacklistException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e ActionBlacklistException) DetailMessage() string {
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

func (e ActionBlacklistException) String() string {
	return e.DetailMessage()
}

func (e ActionBlacklistException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3130005,
		Name: ActionBlacklistExceptionName,
		What: "Action to execute is on the blacklist",
	}

	return json.Marshal(except)
}

func (e ActionBlacklistException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ActionBlacklistException):
		callback(&e)
		return true
	case func(ActionBlacklistException):
		callback(e)
		return true
	default:
		return false
	}
}
