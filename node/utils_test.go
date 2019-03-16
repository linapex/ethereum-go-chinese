
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450102731214848>


//包含测试使用的一批实用程序类型声明。作为节点
//在独特的类型上操作，需要很多类型来检查各种特性。

package node

import (
	"reflect"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

//noopservice是服务接口的一个简单实现。
type NoopService struct{}

func (s *NoopService) Protocols() []p2p.Protocol { return nil }
func (s *NoopService) APIs() []rpc.API           { return nil }
func (s *NoopService) Start(*p2p.Server) error   { return nil }
func (s *NoopService) Stop() error               { return nil }

func NewNoopService(*ServiceContext) (Service, error) { return new(NoopService), nil }

//一组服务都包装基NoopService，从而产生相同的方法
//签名，但外部类型不同。
type NoopServiceA struct{ NoopService }
type NoopServiceB struct{ NoopService }
type NoopServiceC struct{ NoopService }

func NewNoopServiceA(*ServiceContext) (Service, error) { return new(NoopServiceA), nil }
func NewNoopServiceB(*ServiceContext) (Service, error) { return new(NoopServiceB), nil }
func NewNoopServiceC(*ServiceContext) (Service, error) { return new(NoopServiceC), nil }

//instructedService是一种服务的实现，所有接口
//方法既可以检测返回值，也可以检测事件挂钩。
type InstrumentedService struct {
	protocols []p2p.Protocol
	apis      []rpc.API
	start     error
	stop      error

	protocolsHook func()
	startHook     func(*p2p.Server)
	stopHook      func()
}

func NewInstrumentedService(*ServiceContext) (Service, error) { return new(InstrumentedService), nil }

func (s *InstrumentedService) Protocols() []p2p.Protocol {
	if s.protocolsHook != nil {
		s.protocolsHook()
	}
	return s.protocols
}

func (s *InstrumentedService) APIs() []rpc.API {
	return s.apis
}

func (s *InstrumentedService) Start(server *p2p.Server) error {
	if s.startHook != nil {
		s.startHook(server)
	}
	return s.start
}

func (s *InstrumentedService) Stop() error {
	if s.stopHook != nil {
		s.stopHook()
	}
	return s.stop
}

//instructingwrapper是一种专门化返回的服务构造函数的方法
//将通用的检测服务转换为返回包装特定服务的服务。
type InstrumentingWrapper func(base ServiceConstructor) ServiceConstructor

func InstrumentingWrapperMaker(base ServiceConstructor, kind reflect.Type) ServiceConstructor {
	return func(ctx *ServiceContext) (Service, error) {
		obj, err := base(ctx)
		if err != nil {
			return nil, err
		}
		wrapper := reflect.New(kind)
		wrapper.Elem().Field(0).Set(reflect.ValueOf(obj).Elem())

		return wrapper.Interface().(Service), nil
	}
}

//一组服务，所有这些服务都包装基本的instructedService，导致
//方法签名相同，但外部类型不同。
type InstrumentedServiceA struct{ InstrumentedService }
type InstrumentedServiceB struct{ InstrumentedService }
type InstrumentedServiceC struct{ InstrumentedService }

func InstrumentedServiceMakerA(base ServiceConstructor) ServiceConstructor {
	return InstrumentingWrapperMaker(base, reflect.TypeOf(InstrumentedServiceA{}))
}

func InstrumentedServiceMakerB(base ServiceConstructor) ServiceConstructor {
	return InstrumentingWrapperMaker(base, reflect.TypeOf(InstrumentedServiceB{}))
}

func InstrumentedServiceMakerC(base ServiceConstructor) ServiceConstructor {
	return InstrumentingWrapperMaker(base, reflect.TypeOf(InstrumentedServiceC{}))
}

//OneMethodAPI是测试服务返回的单个方法API处理程序。
type OneMethodAPI struct {
	fun func()
}

func (api *OneMethodAPI) TheOneMethod() {
	if api.fun != nil {
		api.fun()
	}
}

