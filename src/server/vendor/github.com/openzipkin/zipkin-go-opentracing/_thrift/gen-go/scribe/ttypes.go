// Autogenerated by Thrift Compiler (0.9.3)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package scribe

import (
	"bytes"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var GoUnusedProtection__ int

type ResultCode int64

const (
	ResultCode_OK        ResultCode = 0
	ResultCode_TRY_LATER ResultCode = 1
)

func (p ResultCode) String() string {
	switch p {
	case ResultCode_OK:
		return "OK"
	case ResultCode_TRY_LATER:
		return "TRY_LATER"
	}
	return "<UNSET>"
}

func ResultCodeFromString(s string) (ResultCode, error) {
	switch s {
	case "OK":
		return ResultCode_OK, nil
	case "TRY_LATER":
		return ResultCode_TRY_LATER, nil
	}
	return ResultCode(0), fmt.Errorf("not a valid ResultCode string")
}

func ResultCodePtr(v ResultCode) *ResultCode { return &v }

func (p ResultCode) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *ResultCode) UnmarshalText(text []byte) error {
	q, err := ResultCodeFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

// Attributes:
//  - Category
//  - Message
type LogEntry struct {
	Category string `thrift:"category,1" json:"category"`
	Message  string `thrift:"message,2" json:"message"`
}

func NewLogEntry() *LogEntry {
	return &LogEntry{}
}

func (p *LogEntry) GetCategory() string {
	return p.Category
}

func (p *LogEntry) GetMessage() string {
	return p.Message
}
func (p *LogEntry) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *LogEntry) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Category = v
	}
	return nil
}

func (p *LogEntry) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Message = v
	}
	return nil
}

func (p *LogEntry) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("LogEntry"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *LogEntry) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("category", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:category: ", p), err)
	}
	if err := oprot.WriteString(string(p.Category)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.category (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:category: ", p), err)
	}
	return err
}

func (p *LogEntry) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("message", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:message: ", p), err)
	}
	if err := oprot.WriteString(string(p.Message)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.message (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:message: ", p), err)
	}
	return err
}

func (p *LogEntry) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("LogEntry(%+v)", *p)
}
