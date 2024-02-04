# gofii

## Depend on

dlfcn libffi

## Windows

install msys2

``` shell
pacman -S mingw-w64-x86_64-dlfcn
pacman -S mingw-w64-x86_64-libffi
```

## example

``` go
package main

/*
typedef int (*opera_fn)(int ,int);
int opera(void* fn, int a, int b){
	opera_fn f = (opera_fn)(fn);
	return f(a, b);
}
*/
import "C"
import (
	"fmt"

	"github.com/jinzhongmin/goffi/pkg/c"
	"github.com/jinzhongmin/goffi/pkg/dlfcn"
	"github.com/jinzhongmin/usf"
)

func main() {

	//###################
	// hello world.
	//###################
	str := c.CStr("hello %s .\n")
	defer usf.Free(str)

	world := c.CStr("world")
	defer usf.Free(world)

	fnPrototype := &c.FuncPrototype{
		Name:    "printf",
		InTypes: []c.Type{c.Pointer, c.Pointer},
		OutType: c.Void,
	}
	fnPrototype.Create(dlfcn.DlsymDefault)
	fnPrototype.Call([]interface{}{&str, &world})

	//###################
	// callback
	//###################

	//define callback prototype
	cbprototype := c.DefineCallbackPrototype(c.AbiDefault, c.I32, []c.Type{c.I32, c.I32})
	defer cbprototype.Free()

	a := C.int(100)
	b := C.int(200)

	add_fun := cbprototype.CreateCallback(func(args []*c.Value, ret *c.Value) {
		a := args[0].I32()
		b := args[1].I32()

		ret.SetI32(a + b)
	})
	add := C.opera(add_fun.CFuncPtr(), a, b)
	fmt.Println("a + b = ", add)
	add_fun.Free()

	//change opera_func to mul
	mul_fun := cbprototype.CreateCallback(func(args []*c.Value, ret *c.Value) {
		a := args[0].I32()
		b := args[1].I32()

		ret.SetI32(a * b)
	})

	mul := C.opera(mul_fun.CFuncPtr(), a, b)
	fmt.Println("a * b = ", mul)
	mul_fun.Free()

}

```