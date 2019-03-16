
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450114752090112>


package simulation

import (
	"testing"
)

func TestService(t *testing.T) {
	sim := New(noopServiceFuncMap)
	defer sim.Close()

	id, err := sim.AddNode()
	if err != nil {
		t.Fatal(err)
	}

	_, ok := sim.Service("noop", id).(*noopService)
	if !ok {
		t.Fatalf("service is not of %T type", &noopService{})
	}

	_, ok = sim.RandomService("noop").(*noopService)
	if !ok {
		t.Fatalf("service is not of %T type", &noopService{})
	}

	_, ok = sim.Services("noop")[id].(*noopService)
	if !ok {
		t.Fatalf("service is not of %T type", &noopService{})
	}
}

