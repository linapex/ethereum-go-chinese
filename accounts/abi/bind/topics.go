
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:31</date>
//</624450061387960320>


package bind

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//maketopics将筛选器查询参数列表转换为筛选器主题集。
func makeTopics(query ...[]interface{}) ([][]common.Hash, error) {
	topics := make([][]common.Hash, len(query))
	for i, filter := range query {
		for _, rule := range filter {
			var topic common.Hash

//尝试基于简单类型生成主题
			switch rule := rule.(type) {
			case common.Hash:
				copy(topic[:], rule[:])
			case common.Address:
				copy(topic[common.HashLength-common.AddressLength:], rule[:])
			case *big.Int:
				blob := rule.Bytes()
				copy(topic[common.HashLength-len(blob):], blob)
			case bool:
				if rule {
					topic[common.HashLength-1] = 1
				}
			case int8:
				blob := big.NewInt(int64(rule)).Bytes()
				copy(topic[common.HashLength-len(blob):], blob)
			case int16:
				blob := big.NewInt(int64(rule)).Bytes()
				copy(topic[common.HashLength-len(blob):], blob)
			case int32:
				blob := big.NewInt(int64(rule)).Bytes()
				copy(topic[common.HashLength-len(blob):], blob)
			case int64:
				blob := big.NewInt(rule).Bytes()
				copy(topic[common.HashLength-len(blob):], blob)
			case uint8:
				blob := new(big.Int).SetUint64(uint64(rule)).Bytes()
				copy(topic[common.HashLength-len(blob):], blob)
			case uint16:
				blob := new(big.Int).SetUint64(uint64(rule)).Bytes()
				copy(topic[common.HashLength-len(blob):], blob)
			case uint32:
				blob := new(big.Int).SetUint64(uint64(rule)).Bytes()
				copy(topic[common.HashLength-len(blob):], blob)
			case uint64:
				blob := new(big.Int).SetUint64(rule).Bytes()
				copy(topic[common.HashLength-len(blob):], blob)
			case string:
				hash := crypto.Keccak256Hash([]byte(rule))
				copy(topic[:], hash[:])
			case []byte:
				hash := crypto.Keccak256Hash(rule)
				copy(topic[:], hash[:])

			default:
//尝试从时髦类型生成主题
				val := reflect.ValueOf(rule)

				switch {
				case val.Kind() == reflect.Array && reflect.TypeOf(rule).Elem().Kind() == reflect.Uint8:
					reflect.Copy(reflect.ValueOf(topic[common.HashLength-val.Len():]), val)

				default:
					return nil, fmt.Errorf("unsupported indexed type: %T", rule)
				}
			}
			topics[i] = append(topics[i], topic)
		}
	}
	return topics, nil
}

//大批量的主题重构反射类型。
var (
	reflectHash    = reflect.TypeOf(common.Hash{})
	reflectAddress = reflect.TypeOf(common.Address{})
	reflectBigInt  = reflect.TypeOf(new(big.Int))
)

//ParseTopics将索引主题字段转换为实际的日志字段值。
//
//注意，动态类型无法重建，因为它们被映射到keccak256
//哈希值作为主题值！
func parseTopics(out interface{}, fields abi.Arguments, topics []common.Hash) error {
//健全性检查字段和主题是否匹配
	if len(fields) != len(topics) {
		return errors.New("topic/field count mismatch")
	}
//遍历所有字段并从主题重新构造它们
	for _, arg := range fields {
		if !arg.Indexed {
			return errors.New("non-indexed field in topic reconstruction")
		}
		field := reflect.ValueOf(out).Elem().FieldByName(capitalise(arg.Name))

//尝试根据基元类型将主题解析回字段
		switch field.Kind() {
		case reflect.Bool:
			if topics[0][common.HashLength-1] == 1 {
				field.Set(reflect.ValueOf(true))
			}
		case reflect.Int8:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int8(num.Int64())))

		case reflect.Int16:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int16(num.Int64())))

		case reflect.Int32:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int32(num.Int64())))

		case reflect.Int64:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(num.Int64()))

		case reflect.Uint8:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint8(num.Uint64())))

		case reflect.Uint16:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint16(num.Uint64())))

		case reflect.Uint32:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint32(num.Uint64())))

		case reflect.Uint64:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(num.Uint64()))

		default:
//基本类型已用完，请尝试自定义类型
			switch field.Type() {
case reflectHash: //还包括所有动态类型
				field.Set(reflect.ValueOf(topics[0]))

			case reflectAddress:
				var addr common.Address
				copy(addr[:], topics[0][common.HashLength-common.AddressLength:])
				field.Set(reflect.ValueOf(addr))

			case reflectBigInt:
				num := new(big.Int).SetBytes(topics[0][:])
				field.Set(reflect.ValueOf(num))

			default:
//没有自定义类型，请尝试疯狂
				switch {
				case arg.Type.T == abi.FixedBytesTy:
					reflect.Copy(field, reflect.ValueOf(topics[0][common.HashLength-arg.Type.Size:]))

				default:
					return fmt.Errorf("unsupported indexed type: %v", arg.Type)
				}
			}
		}
		topics = topics[1:]
	}
	return nil
}

