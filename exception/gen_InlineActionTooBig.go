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

var InlineActionTooBigName = reflect.TypeOf(InlineActionTooBig{}).Name()

type InlineActionTooBig struct {
	_ActionValidateException
	Elog log.Messages
}

func NewInlineActionTooBig(parent _ActionValidateException, message log.Message) *InlineActionTooBig {
	return &InlineActionTooBig{parent, log.Messages{message}}
}

func (e InlineActionTooBig) Code() int64 {
	return 3050009
}

func (e InlineActionTooBig) Name() string {
	return InlineActionTooBigName
}

func (e InlineActionTooBig) What() string {
	return "Inline Action exceeds maximum size limit"
}

func (e *InlineActionTooBig) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e InlineActionTooBig) GetLog() log.Messages {
	return e.Elog
}

func (e InlineActionTooBig) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e InlineActionTooBig) DetailMessage() string {
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

func (e InlineActionTooBig) String() string {
	return e.DetailMessage()
}

func (e InlineActionTooBig) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3050009,
		Name: InlineActionTooBigName,
		What: "Inline Action exceeds maximum size limit",
	}

	return json.Marshal(except)
}

func (e InlineActionTooBig) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*InlineActionTooBig):
		callback(&e)
		return true
	case func(InlineActionTooBig):
		callback(e)
		return true
	default:
		return false
	}
}
