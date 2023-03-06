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
	"reflect"
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

const (
	AbiFirstAbi Abi = 0
	AbiSysv     Abi = 1
	AbiThiscall Abi = 3
	AbiFastcall Abi = 4
	AbiStdcall  Abi = 5
	AbiPascal   Abi = 6
	AbiRegister Abi = 7
	AbiMsCdecl  Abi = 8
	AbiDefault  Abi = AbiSysv
)

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
	cif *C.ffi_cif
	ret unsafe.Pointer
}

func NewCif(abi Abi, ret *Type, args ...*Type) (*Cif, error) {
	cif := new(Cif)

	size := mem.Sizeof(C.ffi_cif{})
	cif.cif = (*C.ffi_cif)(mem.Malloc(1, size))
	mem.Memset(unsafe.Pointer(cif.cif), 0, size)

	//malloc return val ptr
	if ret == nil {
		cif.ret = mem.Malloc(1, int(Void.size))
	} else {
		cif.ret = mem.Malloc(1, int(ret.size))
	}

	//arg typs ptr
	ats := new(unsafe.Pointer)
	if len(args) == 0 {
		_ats := mem.Malloc(1, 8)
		(*ats) = _ats
		mem.PushAt(_ats, 0, unsafe.Pointer(Void.toC()))
	} else {
		_ats := mem.Malloc(len(args), 8)
		(*ats) = _ats
		for i := range args {
			mem.PushAt(_ats, i, unsafe.Pointer(args[i].toC()))
		}
	}

	_ret := Void
	if ret != nil {
		_ret = ret
	}

	st := C.ffi_prep_cif(cif.cif, abi.toC(),
		C.uint(len(args)), _ret.toC(), (**C.ffi_type)(*ats))

	err := Status(st).Error()
	if err != nil {
		defer cif.Free()
		return nil, err
	}

	return cif, nil
}
func (cif *Cif) Free() {
	if cif.ret != nil {
		mem.Free(cif.ret)
	}
	if cif.cif.arg_types != nil {
		mem.Free(unsafe.Pointer(cif.cif.arg_types))
	}
	if cif.cif != nil {
		mem.Free(unsafe.Pointer(cif.cif))
	}
}
func (cif *Cif) Call(fn unsafe.Pointer, argAddr ...any) unsafe.Pointer {
	argc := len(argAddr)
	argv := new(unsafe.Pointer)
	if argc == 0 {
		_argv := mem.Malloc(1, 8)
		(*argv) = _argv
		defer mem.Free(_argv)
	} else {
		_argv := mem.Malloc(argc, 8)
		(*argv) = _argv
		defer mem.Free(_argv)

		for i := range argAddr {
			if argAddr[i] == nil {
				mem.PushAt(_argv, i, nil)
				continue
			}
			mem.PushAt(_argv, i, reflect.ValueOf(argAddr[i]).UnsafePointer())
		}
	}
	C.ffi_call(cif.cif, (*[0]byte)(fn), cif.ret, (*unsafe.Pointer)(*argv))
	return cif.ret
}

type Closure struct {
	cif *Cif

	closure *C.ffi_closure
	fnptr   unsafe.Pointer
	data    unsafe.Pointer
}
type ClosureConf struct {
	Abi  Abi
	Args []*Type
	Ret  *Type
}
type ClosureData struct {
	Args     []unsafe.Pointer
	Ret      unsafe.Pointer
	UserData []any
}
type closureUserData struct {
	fn       func(*ClosureData)
	argc     int
	userData *[]any
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
func NewClosure(conf ClosureConf, fn func(*ClosureData), userData ...any) *Closure {
	var err error
	cls := new(Closure)
	cls.cif, err = NewCif(conf.Abi, conf.Ret, conf.Args...)
	if err != nil {
		panic(err)
	}
	cls.fnptr = mem.Malloc(1, 8)
	cls.closure = (*C.ffi_closure)(C.ffi_closure_alloc(
		C.ulonglong(mem.Sizeof(C.ffi_closure{})), (*unsafe.Pointer)(cls.fnptr)))

	cls.data = mem.Malloc(1, mem.Sizeof(closureUserData{}))
	data := (*closureUserData)(cls.data)
	data.fn = fn
	data.argc = len(conf.Args)
	data.userData = nil
	if userData != nil {
		data.userData = &userData
	}

	C.ffi_prep_closure_loc(cls.closure, cls.cif.cif,
		(*[0]byte)(C.closure_caller), cls.data, mem.Pop(cls.fnptr))
	return cls
}
func (cls *Closure) Call(args ...any) unsafe.Pointer {
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
