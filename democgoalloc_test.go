package main

import "testing"

func TestAllocs(t *testing.T) {
	wantNoAllocs(t, getSum)
	wantNoAllocs(t, getSumDefer)

	c := NewClient()
	if got := int(testing.AllocsPerRun(1000, func() {
		if v := c.sum(); v != 42 {
			t.Fatalf("got %v; want 42", v)
		}
		if v := c.sum2(); v != 42 {
			t.Fatalf("got %v; want 42", v)
		}
		if v := goTakeFoo(c.foo); v != 42 {
			t.Fatalf("got %v; want 42", v)
		}
	})); got != 0 {
		t.Errorf("allocs = %v; want 0", got)
	}
	c.close()

}

func wantNoAllocs(t *testing.T, f func() error) {
	t.Helper()
	n := int(testing.AllocsPerRun(1000, func() {
		if err := f(); err != nil {
			t.Fatal(err)
		}
	}))
	if n != 0 {
		t.Errorf("allocs = %v; want 0", n)
	}
}

func BenchmarkSum(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := getSum(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSumDefer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := getSumDefer(); err != nil {
			b.Fatal(err)
		}
	}
}
