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

import (
	"fmt"
	"time"

	"github.com/jinzhongmin/goffi/pkg/c"
	"github.com/jinzhongmin/goffi/pkg/dlfcn"
	"github.com/jinzhongmin/usf"
)

func main() {

	//###################
	// hello world.
	//###################
	printf, _ := dlfcn.Dlsym(dlfcn.DlsymDefault, "printf")

	str := c.CStr("hello %s .\n")
	defer usf.Free(str)

	world := c.CStr("world")
	defer usf.Free(world)

	c.Call(printf, c.Void, []c.Type{c.Pointer, c.Pointer}, // function prototype
		[]interface{}{&str, &world}) //args

	//###################
	// callback
	//###################

	signal, _ := dlfcn.Dlsym(dlfcn.DlsymDefault, "signal") //#signal from <signal.h>

	callback := c.NewFn(c.AbiDefault, c.Void, []c.Type{c.I32}, //definition function prototype

		func(args []c.Val, ret *c.Val) { //real function call
			sigIntCode := args[0].I32()
			fmt.Println("I've been called. code is: ", sigIntCode)
		})

	defer callback.Free() //need free

	SIGINT := int32(2) //from <signal.h>
	callback_cptr := callback.Cptr() //real c function pointer

	//Call signal and register the callback function
	c.Call(signal, c.Void, []c.Type{c.I32, c.Pointer},
		[]interface{}{&SIGINT, &callback_cptr})

	fmt.Println("input ctrl+c, will callback")

	for {
		fmt.Println("wait ctrl+c")
		time.Sleep(time.Second)
	}
}

```