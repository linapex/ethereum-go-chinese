
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450098037788672>

package metrics

import "testing"

func BenchmarkCounter(b *testing.B) {
	c := NewCounter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Inc(1)
	}
}

func TestCounterClear(t *testing.T) {
	c := NewCounter()
	c.Inc(1)
	c.Clear()
	if count := c.Count(); 0 != count {
		t.Errorf("c.Count(): 0 != %v\n", count)
	}
}

func TestCounterDec1(t *testing.T) {
	c := NewCounter()
	c.Dec(1)
	if count := c.Count(); -1 != count {
		t.Errorf("c.Count(): -1 != %v\n", count)
	}
}

func TestCounterDec2(t *testing.T) {
	c := NewCounter()
	c.Dec(2)
	if count := c.Count(); -2 != count {
		t.Errorf("c.Count(): -2 != %v\n", count)
	}
}

func TestCounterInc1(t *testing.T) {
	c := NewCounter()
	c.Inc(1)
	if count := c.Count(); 1 != count {
		t.Errorf("c.Count(): 1 != %v\n", count)
	}
}

func TestCounterInc2(t *testing.T) {
	c := NewCounter()
	c.Inc(2)
	if count := c.Count(); 2 != count {
		t.Errorf("c.Count(): 2 != %v\n", count)
	}
}

func TestCounterSnapshot(t *testing.T) {
	c := NewCounter()
	c.Inc(1)
	snapshot := c.Snapshot()
	c.Inc(1)
	if count := snapshot.Count(); 1 != count {
		t.Errorf("c.Count(): 1 != %v\n", count)
	}
}

func TestCounterZero(t *testing.T) {
	c := NewCounter()
	if count := c.Count(); 0 != count {
		t.Errorf("c.Count(): 0 != %v\n", count)
	}
}

func TestGetOrRegisterCounter(t *testing.T) {
	r := NewRegistry()
	NewRegisteredCounter("foo", r).Inc(47)
	if c := GetOrRegisterCounter("foo", r); 47 != c.Count() {
		t.Fatal(c)
	}
}

