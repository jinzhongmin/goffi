package ffi

/*
#cgo pkg-config: libffi
#include <ffi.h>
#include <stdint.h>
extern void closure_caller(ffi_cif* cif, void* ret, void* args, void* user_data);
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/jinzhongmin/usf"
)

type Type *C.ffi_type

var (
	Void    Type = &C.ffi_type_void
	Pointer Type = &C.ffi_type_pointer
	Uint8   Type = &C.ffi_type_uint8
	Int8    Type = &C.ffi_type_sint8
	Uint16  Type = &C.ffi_type_uint16
	Int16   Type = &C.ffi_type_sint16
	Uint32  Type = &C.ffi_type_uint32
	Int32   Type = &C.ffi_type_sint32
	Uint64  Type = &C.ffi_type_uint64
	Int64   Type = &C.ffi_type_sint64

	Float             Type = &C.ffi_type_float
	Double            Type = &C.ffi_type_double
	LongDouble        Type = &C.ffi_type_longdouble
	ComplexFloat      Type = &C.ffi_type_complex_float
	Complexdouble     Type = &C.ffi_type_complex_double
	ComplexLongdouble Type = &C.ffi_type_complex_longdouble
)

func TypeFree(typ Type) {
	if typ != nil && ((*C.ffi_type)(typ))._type == C.FFI_TYPE_STRUCT && typ.elements != nil {
		usf.Free(unsafe.Pointer(typ.elements))
	}
}

// create custom type for Struct, need free with TypeFree()
func Struct(size uint64, alignment uint16, elms []Type) Type {
	t := (*C.ffi_type)(usf.MallocOf(1, C.ffi_type{}))
	t.size = C.uint64_t(size)
	t.alignment = C.ushort(alignment)
	t._type = C.FFI_TYPE_STRUCT

	typs := usf.Malloc(uint64(len(elms)), 8)
	copy(*(*[]Type)(usf.Slice(typs, uint64(len(elms)))), elms) //copy elms to typs

	t.elements = (**C.ffi_type)(typs)
	return Type(t)
}

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
	params unsafe.Pointer
	// ret    unsafe.Pointer
}

// void params
var _voidParams = []Type{}

func NewCif(abi Abi, output Type, inputs []Type) (*Cif, error) {
	// if inputs == nil || output == nil {
	// 	panic("Type cannot be nil")
	// }

	cif := new(Cif)
	cifSize := usf.Sizeof(C.ffi_cif{})
	cif.cif = (*C.ffi_cif)(usf.Malloc(1, cifSize))
	usf.Memset(unsafe.Pointer(cif.cif), 0, cifSize)

	_inputs := inputs
	if inputs == nil {
		_inputs = _voidParams
	}
	inLen := len(_inputs)
	if inLen > 0 {
		cif.params = usf.Malloc(uint64(inLen), 8)
		usf.Memset(cif.params, 0, uint64(inLen)*8)
	}
	for i := range _inputs {
		usf.PushAt(cif.params, uint64(i), unsafe.Pointer(_inputs[i]))
	}

	retTyp := output
	if output == nil {
		retTyp = Void
	}
	// cif.ret = usf.Malloc(1, uint64(retTyp.size))
	// usf.Memset(cif.ret, 0, uint64(retTyp.size))

	st := C.ffi_prep_cif(cif.cif, abi.ffi_abi(),
		C.uint(inLen), retTyp, (**C.ffi_type)(cif.params))

	err := Status(st).Error()
	if err != nil {
		defer cif.Free()
		return nil, err
	}

	return cif, nil
}
func (cif *Cif) Free() {
	// if cif.ret != nil {
	// 	usf.Free(cif.ret)
	// }
	if cif.params != nil {
		usf.Free(cif.params)
	}
	if cif.cif != nil {
		usf.Free(unsafe.Pointer(cif.cif))
	}
}

func NewCifRetPtr() unsafe.Pointer {
	p := (usf.Malloc(1, 8))
	usf.Memset(unsafe.Pointer(p), 0, 8)
	return p
}

// void args
var _voidArgs = []interface{}{}

// nil arg
var _nilArg = usf.Malloc(1, 8)

func init() {
	usf.Memset(_nilArg, 0, 8)
}

func (cif *Cif) Call(fn unsafe.Pointer, args []interface{}, returnPtr unsafe.Pointer) {
	_args := args
	if args == nil {
		_args = _voidArgs
	}

	argc := len(_args)
	var argv unsafe.Pointer
	if argc > 0 {
		argv = usf.Malloc(uint64(argc), 8)
		// usf.Memset(argv, 0, 8*uint64(argc))
		defer usf.Free(argv)

		ptrs := *(*[]unsafe.Pointer)(usf.Slice(unsafe.Pointer(&args[0]), uint64(len(args)*2)))
		for i := range _args {
			if _args[i] == nil {
				usf.PushAt(argv, uint64(i), _nilArg)
				continue
			}
			// lv := (*(*[2]unsafe.Pointer)(unsafe.Pointer(&_args[i])))
			usf.PushAt(argv, uint64(i), ptrs[i*2+1])
		}
	}

	// for i := range _args {
	// 	if _args[i] == nil {
	// 		usf.PushAt(argv, uint64(i), _nilArg)
	// 		continue
	// 	}
	// 	lv := (*(*[2]unsafe.Pointer)(unsafe.Pointer(&_args[i])))
	// 	usf.PushAt(argv, uint64(i), lv[1])
	// }
	C.ffi_call(cif.cif, (*[0]byte)(fn), returnPtr, (*unsafe.Pointer)(argv))
}

type Closure struct {
	cif                *Cif
	cfn                unsafe.Pointer
	closure            *C.ffi_closure
	callback           func(*ClosureParams)
	callback_user_data []interface{}
	callback_data      unsafe.Pointer
}

// for NewClosure configure
type ClosureConf struct {
	Abi    Abi
	Inputs []Type
	Output Type
}

// for callback call
type ClosureParams struct {
	Args     []unsafe.Pointer
	Return   unsafe.Pointer
	UserData []interface{}
}
type closure_Data struct {
	callback func(*ClosureParams)
	argc     int
	userData *[]interface{}
}

//export closure_caller
func closure_caller(cif *C.ffi_cif, ret, args, userData unsafe.Pointer) {
	data := (*closure_Data)(userData)

	input := new(ClosureParams)
	input.Args = *(*[]unsafe.Pointer)(usf.Slice(args, uint64(data.argc)))
	input.Return = ret
	if data.userData != nil {
		input.UserData = *data.userData
	}

	data.callback(input)
}
func NewClosure(conf ClosureConf, callback func(*ClosureParams), userData []interface{}) *Closure {
	var err error
	cls := new(Closure)
	cls.cif, err = NewCif(conf.Abi, conf.Output, conf.Inputs)
	if err != nil {
		panic(err)
	}

	cls.cfn = usf.Malloc(1, 8)
	cls.closure = (*C.ffi_closure)(C.ffi_closure_alloc(
		C.uint64_t(usf.Sizeof(C.ffi_closure{})), (*unsafe.Pointer)(cls.cfn)))

	cls.callback = callback
	cls.callback_user_data = userData
	cls.callback_data = (usf.MallocOf(1, closure_Data{}))
	(*closure_Data)(cls.callback_data).callback = callback
	(*closure_Data)(cls.callback_data).argc = len(conf.Inputs)
	(*closure_Data)(cls.callback_data).userData = &cls.callback_user_data

	C.ffi_prep_closure_loc(cls.closure, cls.cif.cif,
		(*[0]byte)(C.closure_caller), unsafe.Pointer(cls.callback_data), usf.Pop(cls.cfn))
	return cls
}
func (cls *Closure) Call(args []interface{}, ret unsafe.Pointer) {
	cls.cif.Call(usf.Pop(cls.cfn), args, ret)
}
func (cls *Closure) Cfn() unsafe.Pointer {
	return (*(*[1]unsafe.Pointer)(cls.cfn))[0]
}
func (cls *Closure) Free() {
	if cls.closure != nil {
		C.ffi_closure_free(unsafe.Pointer(cls.closure))
	}
	if cls.callback_data != nil {
		(*closure_Data)(cls.callback_data).callback = nil
		(*closure_Data)(cls.callback_data).userData = nil
		usf.Free(cls.callback_data)
	}
	if cls.cfn != nil {
		usf.Free(cls.cfn)
	}
	if cls.callback != nil {
		cls.callback = nil
	}
	if cls.callback_user_data != nil {
		cls.callback_user_data = nil
	}
	if cls.cif != nil {
		cls.cif.Free()
	}
}
