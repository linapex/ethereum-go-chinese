
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450084133670912>

package bn256

import (
	"crypto/rand"

	"testing"
)

func TestLatticeReduceCurve(t *testing.T) {
	k, _ := rand.Int(rand.Reader, Order)
	ks := curveLattice.decompose(k)

	if ks[0].BitLen() > 130 || ks[1].BitLen() > 130 {
		t.Fatal("reduction too large")
	} else if ks[0].Sign() < 0 || ks[1].Sign() < 0 {
		t.Fatal("reduction must be positive")
	}
}

func TestLatticeReduceTarget(t *testing.T) {
	k, _ := rand.Int(rand.Reader, Order)
	ks := targetLattice.decompose(k)

	if ks[0].BitLen() > 66 || ks[1].BitLen() > 66 || ks[2].BitLen() > 66 || ks[3].BitLen() > 66 {
		t.Fatal("reduction too large")
	} else if ks[0].Sign() < 0 || ks[1].Sign() < 0 || ks[2].Sign() < 0 || ks[3].Sign() < 0 {
		t.Fatal("reduction must be positive")
	}
}

