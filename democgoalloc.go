package main

/*
#include <stdlib.h>

typedef struct Foo{
	int x;
	int y;
	int* ptr;
} Foo;

static int foo_sum(Foo* f) {
	return f->x + f->y;
}

static int foo_sum_other_args(Foo* f, int* x, int* y) {
	return f->x + f->y;
}

static Foo* new_foo() {
	Foo* f = (Foo*) malloc(sizeof(Foo));
	f->ptr = NULL;
	f->x = 40;
	f->y = 2;
	return f;
}

static void close_foo(Foo* f) {
	free(f);
}

*/
import "C"

import (
	"fmt"
	"sync"
)

type Client struct {
	foo *C.Foo

	closeOnce sync.Once
}

func (c *Client) sum() int  { return int(C.foo_sum(c.foo)) }
func (c *Client) sum2() int { return int(C.foo_sum_other_args(c.foo, nil, nil)) }
func (c *Client) close() {
	c.closeOnce.Do(func() {
		C.close_foo(c.foo)
	})
}

var escClient *Client

func (c *Client) leak() { escClient = c }

func NewClient() *Client {
	return &Client{foo: C.new_foo()}
}

func main() {
	f := C.new_foo()
	fmt.Println(C.foo_sum(f))
	C.close_foo(f)
}

func getSum() error {
	f := C.new_foo()
	v := C.foo_sum(f)
	C.close_foo(f)
	if v != 42 {
		return fmt.Errorf("not 42; got %v", v)
	}
	return nil
}

func getSumDefer() error {
	f := C.new_foo()
	defer C.close_foo(f)
	v := C.foo_sum(f)
	if v != 42 {
		return fmt.Errorf("not 42; got %v", v)
	}
	return nil
}

func goTakeFoo(f *C.Foo) int {
	return int(C.foo_sum(f))
}
