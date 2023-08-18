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
	opera_func := c.NewCallback(c.AbiDefault, c.I32, []c.Type{c.I32, c.I32})

	//opera_func is base goffi.Closure
	//CallbackCvt wrap from goffi.Closure.Callback
	//CallbackCvt for convert goffi.Closure.Callback to real func CallbackFunc
	opera_func.CallbackCvt = func(callback *c.Callback, args []*c.Value, ret *c.Value) {
		fn, ok := callback.CallbackFunc.(func(int32, int32) int32)
		if ok {
			a := args[0].I32()
			b := args[1].I32()
			ret.SetI32(fn(a, b))
		}
	}

	a := C.int(100)
	b := C.int(200)

	//define opera_func to add
	//call path, call goffi.Closure.Callback -> call CallbackCvt -> call CallbackFunc
	opera_func.CallbackFunc = func(a int32, b int32) int32 {
		return a + b
	}
	add := C.opera(opera_func.Cfunc, a, b)
	fmt.Println("a + b = ", add)

	//change opera_func to mul
	opera_func.CallbackFunc = func(a int32, b int32) int32 {
		return a * b
	}
	mul := C.opera(opera_func.Cfunc, a, b)
	fmt.Println("a * b = ", mul)

}

```