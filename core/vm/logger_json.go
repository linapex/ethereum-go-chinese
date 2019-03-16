
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450082908934144>


package vm

import (
	"encoding/json"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

type JSONLogger struct {
	encoder *json.Encoder
	cfg     *LogConfig
}

//newjsonLogger创建了一个新的EVM跟踪程序，它将执行步骤作为JSON对象打印出来。
//进入提供的流。
func NewJSONLogger(cfg *LogConfig, writer io.Writer) *JSONLogger {
	return &JSONLogger{json.NewEncoder(writer), cfg}
}

func (l *JSONLogger) CaptureStart(from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) error {
	return nil
}

//CaptureState在记录器上输出状态信息。
func (l *JSONLogger) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error {
	log := StructLog{
		Pc:            pc,
		Op:            op,
		Gas:           gas,
		GasCost:       cost,
		MemorySize:    memory.Len(),
		Storage:       nil,
		Depth:         depth,
		RefundCounter: env.StateDB.GetRefund(),
		Err:           err,
	}
	if !l.cfg.DisableMemory {
		log.Memory = memory.Data()
	}
	if !l.cfg.DisableStack {
		log.Stack = stack.Data()
	}
	return l.encoder.Encode(log)
}

//CaptureFault在记录器上输出状态信息。
func (l *JSONLogger) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error {
	return nil
}

//CaptureEnd在执行结束时触发。
func (l *JSONLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error {
	type endLog struct {
		Output  string              `json:"output"`
		GasUsed math.HexOrDecimal64 `json:"gasUsed"`
		Time    time.Duration       `json:"time"`
		Err     string              `json:"error,omitempty"`
	}
	if err != nil {
		return l.encoder.Encode(endLog{common.Bytes2Hex(output), math.HexOrDecimal64(gasUsed), t, err.Error()})
	}
	return l.encoder.Encode(endLog{common.Bytes2Hex(output), math.HexOrDecimal64(gasUsed), t, ""})
}

