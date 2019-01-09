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

var ProducerDoubleConfirmName = reflect.TypeOf(ProducerDoubleConfirm{}).Name()

type ProducerDoubleConfirm struct {
	_ProducerException
	Elog log.Messages
}

func NewProducerDoubleConfirm(parent _ProducerException, message log.Message) *ProducerDoubleConfirm {
	return &ProducerDoubleConfirm{parent, log.Messages{message}}
}

func (e ProducerDoubleConfirm) Code() int64 {
	return 3170003
}

func (e ProducerDoubleConfirm) Name() string {
	return ProducerDoubleConfirmName
}

func (e ProducerDoubleConfirm) What() string {
	return "Producer is double confirming known range"
}

func (e *ProducerDoubleConfirm) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ProducerDoubleConfirm) GetLog() log.Messages {
	return e.Elog
}

func (e ProducerDoubleConfirm) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e ProducerDoubleConfirm) DetailMessage() string {
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

func (e ProducerDoubleConfirm) String() string {
	return e.DetailMessage()
}

func (e ProducerDoubleConfirm) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3170003,
		Name: ProducerDoubleConfirmName,
		What: "Producer is double confirming known range",
	}

	return json.Marshal(except)
}

func (e ProducerDoubleConfirm) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ProducerDoubleConfirm):
		callback(&e)
		return true
	case func(ProducerDoubleConfirm):
		callback(e)
		return true
	default:
		return false
	}
}
