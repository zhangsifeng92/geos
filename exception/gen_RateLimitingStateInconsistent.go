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

var RateLimitingStateInconsistentName = reflect.TypeOf(RateLimitingStateInconsistent{}).Name()

type RateLimitingStateInconsistent struct {
	_MiscException
	Elog log.Messages
}

func NewRateLimitingStateInconsistent(parent _MiscException, message log.Message) *RateLimitingStateInconsistent {
	return &RateLimitingStateInconsistent{parent, log.Messages{message}}
}

func (e RateLimitingStateInconsistent) Code() int64 {
	return 3100001
}

func (e RateLimitingStateInconsistent) Name() string {
	return RateLimitingStateInconsistentName
}

func (e RateLimitingStateInconsistent) What() string {
	return "Internal state is no longer consistent"
}

func (e *RateLimitingStateInconsistent) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e RateLimitingStateInconsistent) GetLog() log.Messages {
	return e.Elog
}

func (e RateLimitingStateInconsistent) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e RateLimitingStateInconsistent) DetailMessage() string {
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

func (e RateLimitingStateInconsistent) String() string {
	return e.DetailMessage()
}

func (e RateLimitingStateInconsistent) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3100001,
		Name: RateLimitingStateInconsistentName,
		What: "Internal state is no longer consistent",
	}

	return json.Marshal(except)
}

func (e RateLimitingStateInconsistent) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*RateLimitingStateInconsistent):
		callback(&e)
		return true
	case func(RateLimitingStateInconsistent):
		callback(e)
		return true
	default:
		return false
	}
}
