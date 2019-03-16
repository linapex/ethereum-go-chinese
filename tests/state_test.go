
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450122457026560>


package tests

import (
	"bytes"
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/core/vm"
)

func TestState(t *testing.T) {
	t.Parallel()

	st := new(testMatcher)
//长期测试：
	st.slow(`^stAttackTest/ContractCreationSpam`)
	st.slow(`^stBadOpcode/badOpcodes`)
	st.slow(`^stPreCompiledContracts/modexp`)
	st.slow(`^stQuadraticComplexityTest/`)
	st.slow(`^stStaticCall/static_Call50000`)
	st.slow(`^stStaticCall/static_Return50000`)
	st.slow(`^stStaticCall/static_Call1MB`)
	st.slow(`^stSystemOperationsTest/CallRecursiveBomb`)
	st.slow(`^stTransactionTest/Opcodes_TransactionInit`)
//断裂试验：
st.skipLoad(`^stTransactionTest/OverflowGasRequire\.json`) //气体限制>256位
st.skipLoad(`^stTransactionTest/zeroSigTransa[^/]*\.json`) //还不支持EIP-86
//预期故障：
	st.fails(`^stRevertTest/RevertPrecompiledTouch\.json/EIP158`, "bug in test")
	st.fails(`^stRevertTest/RevertPrecompiledTouch\.json/Byzantium`, "bug in test")
	st.fails(`^stRevertTest/RevertPrecompiledTouch.json/Constantinople`, "bug in test")

	st.walk(t, stateTestDir, func(t *testing.T, name string, test *StateTest) {
		for _, subtest := range test.Subtests() {
			subtest := subtest
			key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
			name := name + "/" + key
			t.Run(key, func(t *testing.T) {
				withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
					_, err := test.Run(subtest, vmconfig)
					return st.checkFailure(t, name, err)
				})
			})
		}
	})
}

//GasLimit高于此值的事务在失败时不会获得VM跟踪。
const traceErrorLimit = 400000

//接受--vm.*命令行参数的状态测试的vm配置。
var testVMConfig = func() vm.Config {
	vmconfig := vm.Config{}
	flag.StringVar(&vmconfig.EVMInterpreter, utils.EVMInterpreterFlag.Name, utils.EVMInterpreterFlag.Value, utils.EVMInterpreterFlag.Usage)
	flag.StringVar(&vmconfig.EWASMInterpreter, utils.EWASMInterpreterFlag.Name, utils.EWASMInterpreterFlag.Value, utils.EWASMInterpreterFlag.Usage)
	flag.Parse()
	return vmconfig
}()

func withTrace(t *testing.T, gasLimit uint64, test func(vm.Config) error) {
	err := test(testVMConfig)
	if err == nil {
		return
	}
	t.Error(err)
	if gasLimit > traceErrorLimit {
		t.Log("gas limit too high for EVM trace")
		return
	}
	tracer := vm.NewStructLogger(nil)
	err2 := test(vm.Config{Debug: true, Tracer: tracer})
	if !reflect.DeepEqual(err, err2) {
		t.Errorf("different error for second run: %v", err2)
	}
	buf := new(bytes.Buffer)
	vm.WriteTrace(buf, tracer.StructLogs())
	if buf.Len() == 0 {
		t.Log("no EVM operation logs generated")
	} else {
		t.Log("EVM operation log:\n" + buf.String())
	}
	t.Logf("EVM output: 0x%x", tracer.Output())
	t.Logf("EVM error: %v", tracer.Error())
}

