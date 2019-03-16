
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450082678247424>


package vm

import (
	"testing"
)

func TestIntPoolPoolGet(t *testing.T) {
	poolOfIntPools.pools = make([]*intPool, 0, poolDefaultCap)

	nip := poolOfIntPools.get()
	if nip == nil {
		t.Fatalf("Invalid pool allocation")
	}
}

func TestIntPoolPoolPut(t *testing.T) {
	poolOfIntPools.pools = make([]*intPool, 0, poolDefaultCap)

	nip := poolOfIntPools.get()
	if len(poolOfIntPools.pools) != 0 {
		t.Fatalf("Pool got added to list when none should have been")
	}

	poolOfIntPools.put(nip)
	if len(poolOfIntPools.pools) == 0 {
		t.Fatalf("Pool did not get added to list when one should have been")
	}
}

func TestIntPoolPoolReUse(t *testing.T) {
	poolOfIntPools.pools = make([]*intPool, 0, poolDefaultCap)
	nip := poolOfIntPools.get()
	poolOfIntPools.put(nip)
	poolOfIntPools.get()

	if len(poolOfIntPools.pools) != 0 {
		t.Fatalf("Invalid number of pools. Got %d, expected %d", len(poolOfIntPools.pools), 0)
	}
}

