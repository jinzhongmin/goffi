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

	typs := usf.MallocN(uint64(len(elms)), 8)
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

var _zeroCCif = C.ffi_cif{}
var _zeroCCifSize uint64 = 0
var _zeroClosure = Closure{}
var _zeroClosureSize uint64 = 0
var _zeroCClosure = C.ffi_closure{}
var _zeroCClosureSize uint64 = 0
var _voidParams = []Type{}      // void params
var _voidArgs = []interface{}{} // void args
var _nilArg = usf.MallocN(1, 8) // nil arg
func init() {
	usf.Memset(_nilArg, 0, 8)
	_zeroCCifSize = usf.SizeOf(_zeroCCif)
	_zeroClosureSize = usf.SizeOf(_zeroClosure)
	_zeroCClosureSize = usf.SizeOf(_zeroCClosure)
}

func NewCif(abi Abi, output Type, inputs []Type) (*Cif, error) {
	cif := new(Cif)
	cif.cif = (*C.ffi_cif)(usf.MallocN(1, _zeroCCifSize))
	usf.Memset(unsafe.Pointer(cif.cif), 0, _zeroCCifSize)

	_inputs := inputs
	if inputs == nil {
		_inputs = _voidParams
	}

	inLen := uint64(len(_inputs))
	if inLen > 0 {
		cif.params = usf.MallocN(inLen, 8)
		usf.Memset(cif.params, 0, inLen*8)
	}
	for i := uint64(0); i < inLen; i++ {
		usf.PushAt(cif.params, i, unsafe.Pointer(_inputs[i]))
	}

	retTyp := output
	if output == nil {
		retTyp = Void
	}

	st := C.ffi_prep_cif(cif.cif, abi.abi(),
		C.uint(inLen), retTyp, (**C.ffi_type)(cif.params))

	err := Status(st).Error()
	if err != nil {
		defer cif.Free()
		return nil, err
	}

	return cif, nil
}
func (cif *Cif) Call(fn unsafe.Pointer, args []interface{}, ret unsafe.Pointer) {
	arg := args
	if args == nil {
		arg = _voidArgs
	}

	argc := uint64(len(arg))
	argv := unsafe.Pointer(nil)
	argv_free := false
	if argc > 0 {
		argv = usf.MallocN(argc, 8)
		argv_free = true

		src := *(*[]unsafe.Pointer)(usf.Slice(unsafe.Pointer(&args[0]), argc*2))
		dst := *(*[]unsafe.Pointer)(usf.Slice(argv, argc))

		for i, ii := uint64(0), uint64(1); i < argc; i++ {
			if arg[i] == nil {
				dst[i] = _nilArg
				ii += 2
				continue
			}
			dst[i] = src[ii]
			ii += 2
		}
	}

	C.ffi_call(cif.cif, (*[0]byte)(fn), ret, (*unsafe.Pointer)(argv))
	if argv_free {
		usf.Free(argv)
	}
}
func (cif *Cif) Free() {
	// if cif.ret != nil {
	// 	usf.Free(cif.ret)
	// }
	if cif == nil {
		return
	}
	if cif.params != nil {
		usf.Free(cif.params)
	}
	if cif.cif != nil {
		usf.Free(unsafe.Pointer(cif.cif))
	}
}

func NewPtr() unsafe.Pointer {
	p := (usf.MallocN(1, 8))
	usf.Memset(unsafe.Pointer(p), 0, 8)
	return p
}

type Closure struct {
	Cfunc   unsafe.Pointer
	cif     *Cif
	closure *C.ffi_closure

	Callback      func(args []unsafe.Pointer, ret unsafe.Pointer)
	callback_argc int
}

//export closure_caller
func closure_caller(cif *C.ffi_cif, ret, args, userData unsafe.Pointer) {
	cls := (*Closure)(userData)
	cls.Callback(*(*[]unsafe.Pointer)(usf.Slice(args, uint64(cls.callback_argc))), ret)
}
func NewClosure(aib Abi, outType Type, inTypes []Type, callback func(args []unsafe.Pointer, ret unsafe.Pointer)) *Closure {
	var err error
	cls := (*Closure)(usf.MallocN(1, _zeroClosureSize))
	cls.cif, err = NewCif(aib, outType, inTypes)
	if err != nil {
		panic(err)
	}

	cfn := usf.MallocN(1, 8)
	cls.closure = (*C.ffi_closure)(C.ffi_closure_alloc(
		C.uint64_t(_zeroCClosureSize), (*unsafe.Pointer)(cfn)))

	cls.Cfunc = usf.Pop(cfn)
	usf.Free(cfn)

	cls.Callback = callback
	cls.callback_argc = 0
	if inTypes != nil {
		cls.callback_argc = len(inTypes)
	}

	C.ffi_prep_closure_loc(cls.closure, cls.cif.cif,
		(*[0]byte)(C.closure_caller), unsafe.Pointer(cls), cls.Cfunc)
	return cls
}
func (cls *Closure) Call(args []interface{}, ret unsafe.Pointer) {
	cls.cif.Call(cls.Cfunc, args, ret)
}
func (cls *Closure) Free() {
	if cls == nil {
		return
	}
	if cls.closure != nil {
		C.ffi_closure_free(unsafe.Pointer(cls.closure))
	}
	if cls.Cfunc != nil {
		usf.Free(cls.Cfunc)
	}
	if cls.Callback != nil {
		cls.Callback = nil
	}
	if cls.cif != nil {
		cls.cif.Free()
	}
	usf.Free(unsafe.Pointer(cls))
}
