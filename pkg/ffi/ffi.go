package ffi

/*
#cgo windows CFLAGS:  -I../../3rdparty/windows/ffi/include
#cgo windows LDFLAGS: -L../../3rdparty/windows/ffi/lib -lffi
#cgo !windows pkg-config: libffi
#include <ffi.h>
extern void closure_caller(ffi_cif* cif, void* ret, void* args, void* user_data);
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/jinzhongmin/mem"
)

type Type C.ffi_type

func (typ *Type) toC() *C.ffi_type { return (*C.ffi_type)(typ) }

var (
	Void    *Type = (*Type)(&C.ffi_type_void)
	Pointer *Type = (*Type)(&C.ffi_type_pointer)
	Uint8   *Type = (*Type)(&C.ffi_type_uint8)
	Int8    *Type = (*Type)(&C.ffi_type_sint8)
	Uint16  *Type = (*Type)(&C.ffi_type_uint16)
	Int16   *Type = (*Type)(&C.ffi_type_sint16)
	Uint32  *Type = (*Type)(&C.ffi_type_uint32)
	Int32   *Type = (*Type)(&C.ffi_type_sint32)
	Uint64  *Type = (*Type)(&C.ffi_type_uint64)
	Int64   *Type = (*Type)(&C.ffi_type_sint64)

	Float             *Type = (*Type)(&C.ffi_type_float)
	Double            *Type = (*Type)(&C.ffi_type_double)
	LongDouble        *Type = (*Type)(&C.ffi_type_longdouble)
	ComplexFloat      *Type = (*Type)(&C.ffi_type_complex_float)
	Complexdouble     *Type = (*Type)(&C.ffi_type_complex_double)
	ComplexLongdouble *Type = (*Type)(&C.ffi_type_complex_longdouble)
)

type Abi C.ffi_abi

func (a Abi) toC() C.ffi_abi { return C.ffi_abi(a) }

type Status C.ffi_status

const (
	StatusOk         Status = C.FFI_OK
	StatusBadTypedef Status = C.FFI_BAD_TYPEDEF
	StatusBadAbi     Status = C.FFI_BAD_ABI
	StatusBadArgType Status = C.FFI_BAD_ARGTYPE
)

func (st Status) Error() error {
	switch C.ffi_status(st) {
	case C.FFI_OK:
		return nil
	case C.FFI_BAD_TYPEDEF:
		return errors.New("bad typedef")
	case C.FFI_BAD_ABI:
		return errors.New("bad abi")
	case C.FFI_BAD_ARGTYPE:
		return errors.New("bad argtype")
	default:
		return errors.New("unknow error")
	}
}

type Cif struct {
	cif    *C.ffi_cif
	ret    unsafe.Pointer
	params []*Type
}

func NewCif(abi Abi, output *Type, inputs ...*Type) (*Cif, error) {
	cif := new(Cif)

	size := mem.Sizeof(C.ffi_cif{})
	cif.cif = (*C.ffi_cif)(mem.Malloc(1, size))
	mem.Memset(unsafe.Pointer(cif.cif), 0, size)

	cif.params = make([]*Type, 0)
	if inputs == nil || (len(inputs) == 1 && inputs[0] == nil) {
		cif.params = append(cif.params, Void)
	} else {
		cif.params = append(cif.params, inputs...)
	}

	inputLen := len(cif.params)
	args := mem.Malloc(inputLen, 8)
	mem.Memset(args, 0, inputLen*8)

	for i := range cif.params {
		mem.PushAt(args, i, unsafe.Pointer(cif.params[i].toC()))
	}

	ret := Void
	if output != nil {
		ret = output
	}
	cif.ret = mem.Malloc(1, int(ret.size))

	st := C.ffi_prep_cif(cif.cif, abi.toC(),
		C.uint(uint32(inputLen)), ret.toC(), (**C.ffi_type)(unsafe.Pointer(&args)))

	err := Status(st).Error()
	if err != nil {
		defer cif.Free()
		return nil, err
	}

	return cif, nil
}
func (cif *Cif) Free() {
	cif.params = nil
	if cif.ret != nil {
		mem.Free(cif.ret)
	}
	if cif.cif.arg_types != nil {
		p := *(*unsafe.Pointer)((unsafe.Pointer)(cif.cif.arg_types))
		mem.Free(unsafe.Pointer(p))
	}
	if cif.cif != nil {
		mem.Free(unsafe.Pointer(cif.cif))
	}
}
func (cif *Cif) Call(fn unsafe.Pointer, args ...interface{}) unsafe.Pointer {
	argc := len(args)
	if cif.params[0] == Void &&
		((args == nil) || (len(args) == 1 && args[0] == nil)) {

		argc = 0
		args = nil

	} else {
		if argc != int(cif.cif.nargs) {
			panic("param len not equal arg")
		}
	}

	argcN := 1
	if argc > 1 {
		argcN = argc
	}
	argv := mem.Malloc(argcN, 8)
	mem.Memset(argv, 0, 8*argcN)
	defer mem.Free(argv)

	if argc != 0 {
		for i := range args {
			if args[i] == nil {
				continue
			}
			vs := (*(*[2]unsafe.Pointer)(unsafe.Pointer(&args[i])))
			mem.PushAt(argv, i, *(*unsafe.Pointer)(vs[1]))
		}
	}

	C.ffi_call(cif.cif, (*[0]byte)(fn), cif.ret, (*unsafe.Pointer)(argv))
	return cif.ret
}

type Closure struct {
	cif *Cif

	closure *C.ffi_closure
	fnptr   unsafe.Pointer
	data    unsafe.Pointer
}
type ClosureConf struct {
	Abi    Abi
	Inputs []*Type
	Output *Type
}
type ClosureData struct {
	Args     []unsafe.Pointer
	Ret      unsafe.Pointer
	UserData []interface{}
}
type closureUserData struct {
	fn       func(*ClosureData)
	argc     int
	userData *[]interface{}
}

//export closure_caller
func closure_caller(cif *C.ffi_cif, ret, args, userData unsafe.Pointer) {
	data := (*closureUserData)(userData)
	input := new(ClosureData)
	input.Args = *(*[]unsafe.Pointer)(mem.Slice(args, data.argc))
	input.Ret = ret
	if data.userData != nil {
		input.UserData = *data.userData
	}
	data.fn(input)
}
func NewClosure(conf ClosureConf, fn func(*ClosureData), userData ...interface{}) *Closure {
	var err error
	cls := new(Closure)
	cls.cif, err = NewCif(conf.Abi, conf.Output, conf.Inputs...)
	if err != nil {
		panic(err)
	}
	cls.fnptr = mem.Malloc(1, 8)
	cls.closure = (*C.ffi_closure)(C.ffi_closure_alloc(
		C.ulonglong(mem.Sizeof(C.ffi_closure{})), (*unsafe.Pointer)(cls.fnptr)))

	cls.data = mem.Malloc(1, mem.Sizeof(closureUserData{}))
	data := (*closureUserData)(cls.data)
	data.fn = fn
	data.argc = len(conf.Inputs)
	data.userData = nil
	if userData != nil {
		data.userData = &userData
	}

	C.ffi_prep_closure_loc(cls.closure, cls.cif.cif,
		(*[0]byte)(C.closure_caller), cls.data, mem.Pop(cls.fnptr))
	return cls
}
func (cls *Closure) Call(args ...interface{}) unsafe.Pointer {
	return cls.cif.Call(mem.Pop(cls.fnptr), args...)
}
func (cls *Closure) Cfn() unsafe.Pointer {
	return mem.Pop(cls.fnptr)
}
func (cls *Closure) Free() {
	if cls.cif != nil {
		cls.cif.Free()
	}
	if cls.fnptr != nil {
		mem.Free(cls.fnptr)
	}
	if cls.data != nil {
		data := (*closureUserData)(cls.data)
		data.fn = nil
		data.userData = nil
		mem.Free(cls.data)
	}
	if cls.closure != nil {
		C.ffi_closure_free(unsafe.Pointer(cls.closure))
	}
}
