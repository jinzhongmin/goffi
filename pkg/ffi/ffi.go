package ffi

/*
#cgo pkg-config: libffi
#include <ffi.h>
#include <stdint.h>
extern void closure_caller(ffi_cif* cif, void* ret, void* args, void* user_data);

static unsigned int __GOFFI_FFI_DEFAULT_ABI = 9999;
static unsigned int __GOFFI_FFI_EFI64 = 9999;
static unsigned int __GOFFI_FFI_FASTCALL = 9999;
static unsigned int __GOFFI_FFI_FIRST_ABI = 9999;
static unsigned int __GOFFI_FFI_GNUW64 = 9999;
static unsigned int __GOFFI_FFI_LAST_ABI = 9999;
static unsigned int __GOFFI_FFI_MS_CDECL = 9999;
static unsigned int __GOFFI_FFI_PASCAL = 9999;
static unsigned int __GOFFI_FFI_REGISTER = 9999;
static unsigned int __GOFFI_FFI_STDCALL = 9999;
static unsigned int __GOFFI_FFI_SYSV = 9999;
static unsigned int __GOFFI_FFI_THISCALL = 9999;
static unsigned int __GOFFI_FFI_UNIX64 = 9999;
static unsigned int __GOFFI_FFI_WIN64 = 9999;

static size_t __GOFFI_const_c_cif_size;
static size_t __GOFFI_const_c_ffi_type_size;
static size_t __GOFFI_const_c_closure_size;
static void*  __GOFFI_const_arg_nil;


static size_t _FUNC__GOFFI_const_c_cif_size()     {return __GOFFI_const_c_cif_size;}
static size_t _FUNC__GOFFI_const_c_ffi_type_size(){return __GOFFI_const_c_ffi_type_size;}
static size_t _FUNC__GOFFI_const_c_closure_size() {return __GOFFI_const_c_closure_size;}
static void*  _FUNC__GOFFI_const_arg_nil()        {return __GOFFI_const_arg_nil;}

static unsigned int _FUNC__GOFFI_FFI_DEFAULT_ABI(){return __GOFFI_FFI_DEFAULT_ABI;}
static unsigned int _FUNC__GOFFI_FFI_EFI64(){return __GOFFI_FFI_EFI64;}
static unsigned int _FUNC__GOFFI_FFI_FASTCALL(){return __GOFFI_FFI_FASTCALL;}
static unsigned int _FUNC__GOFFI_FFI_FIRST_ABI(){return __GOFFI_FFI_FIRST_ABI;}
static unsigned int _FUNC__GOFFI_FFI_GNUW64(){return __GOFFI_FFI_GNUW64;}
static unsigned int _FUNC__GOFFI_FFI_LAST_ABI(){return __GOFFI_FFI_LAST_ABI;}
static unsigned int _FUNC__GOFFI_FFI_MS_CDECL(){return __GOFFI_FFI_MS_CDECL;}
static unsigned int _FUNC__GOFFI_FFI_PASCAL(){return __GOFFI_FFI_PASCAL;}
static unsigned int _FUNC__GOFFI_FFI_REGISTER(){return __GOFFI_FFI_REGISTER;}
static unsigned int _FUNC__GOFFI_FFI_STDCALL(){return __GOFFI_FFI_STDCALL;}
static unsigned int _FUNC__GOFFI_FFI_SYSV(){return __GOFFI_FFI_SYSV;}
static unsigned int _FUNC__GOFFI_FFI_THISCALL(){return __GOFFI_FFI_THISCALL;}
static unsigned int _FUNC__GOFFI_FFI_UNIX64(){return __GOFFI_FFI_UNIX64;}
static unsigned int _FUNC__GOFFI_FFI_WIN64(){return __GOFFI_FFI_WIN64;}

static void init(){

  __GOFFI_const_c_cif_size = sizeof(ffi_cif);
  __GOFFI_const_c_ffi_type_size = sizeof(ffi_type);
  __GOFFI_const_c_closure_size = sizeof(ffi_closure);
  __GOFFI_const_arg_nil = NULL;

#if defined(X86_WIN64)
  __GOFFI_FFI_FIRST_ABI = FFI_FIRST_ABI;
  __GOFFI_FFI_WIN64 = FFI_WIN64;
  __GOFFI_FFI_GNUW64 = FFI_GNUW64;
  __GOFFI_FFI_LAST_ABI = FFI_LAST_ABI;
#ifdef __GNUC__
  __GOFFI_FFI_DEFAULT_ABI = FFI_DEFAULT_ABI;
#else
  __GOFFI_FFI_DEFAULT_ABI = FFI_DEFAULT_ABI;
#endif

#elif defined(X86_64) || (defined (__x86_64__) && defined (X86_DARWIN))
  __GOFFI_FFI_FIRST_ABI = FFI_FIRST_ABI;
  __GOFFI_FFI_UNIX64 = FFI_UNIX64;
  __GOFFI_FFI_WIN64 = FFI_WIN64;
  __GOFFI_FFI_EFI64 = FFI_EFI64;
  __GOFFI_FFI_GNUW64 = FFI_GNUW64;
  __GOFFI_FFI_LAST_ABI = FFI_LAST_ABI;
  __GOFFI_FFI_DEFAULT_ABI = FFI_DEFAULT_ABI;
#elif defined(X86_WIN32)
  __GOFFI_FFI_FIRST_ABI = FFI_FIRST_ABI
  __GOFFI_FFI_SYSV = FFI_SYSV
  __GOFFI_FFI_STDCALL = FFI_STDCALL
  __GOFFI_FFI_THISCALL = FFI_THISCALL
  __GOFFI_FFI_FASTCALL = FFI_FASTCALL
  __GOFFI_FFI_MS_CDECL = FFI_MS_CDECL
  __GOFFI_FFI_PASCAL = FFI_PASCAL
  __GOFFI_FFI_REGISTER = FFI_REGISTER
  __GOFFI_FFI_LAST_ABI, = FFI_LAST_ABI,
  __GOFFI_FFI_DEFAULT_ABI = FFI_DEFAULT_ABI
#else
  __GOFFI_FFI_FIRST_ABI = FFI_FIRST_ABI
  __GOFFI_FFI_SYSV = FFI_SYSV
  __GOFFI_FFI_STDCALL = FFI_STDCALL
  __GOFFI_FFI_THISCALL = FFI_THISCALL
  __GOFFI_FFI_FASTCALL = FFI_FASTCALL
  __GOFFI_FFI_MS_CDECL = FFI_MS_CDECL
  __GOFFI_FFI_PASCAL = FFI_PASCAL
  __GOFFI_FFI_REGISTER = FFI_REGISTER
  __GOFFI_FFI_LAST_ABI, = FFI_LAST_ABI,
  __GOFFI_FFI_DEFAULT_ABI = FFI_DEFAULT_ABI
#endif
}
*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"unsafe"

	"github.com/jinzhongmin/usf"
)

type (
	Type   *C.ffi_type
	Status C.ffi_status
	Abi    C.ffi_abi
)

func (a Abi) cvt() C.ffi_abi {
	if a == 9999 {
		panic("abi not define int this system")
	}
	return C.ffi_abi(a)
}

var (
	AbiDefault  Abi
	AbiFirst    Abi
	AbiLast     Abi
	AbiEfi64    Abi
	AbiFastcall Abi
	AbiGnuw64   Abi
	AbiMsCdecl  Abi
	AbiPascal   Abi
	AbiRegister Abi
	AbiStdcall  Abi
	AbiSysv     Abi
	AbiThiscall Abi
	AbiUnix64   Abi
	AbiWin64    Abi

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
}

var (
	_voidParams = []Type{}        // void params
	_voidArgs   = []interface{}{} // void args

	_sizeOfCFfiType_ uint64
	_sizeOfCCif_     uint64
	_sizeOfCClosure_ uint64
	_cNil_           unsafe.Pointer
)

func _sizeOfCFfiType() uint64 { return _sizeOfCFfiType_ }
func _sizeOfCCif() uint64     { return _sizeOfCCif_ }
func _sizeOfCClosure() uint64 { return _sizeOfCClosure_ }
func _cNil() unsafe.Pointer   { return _cNil_ }

func init() {
	C.init()
	_sizeOfCFfiType_ = uint64(C._FUNC__GOFFI_const_c_ffi_type_size())
	_sizeOfCCif_ = uint64(C._FUNC__GOFFI_const_c_cif_size())
	_sizeOfCClosure_ = uint64(C._FUNC__GOFFI_const_c_closure_size())
	_cNil_ = C._FUNC__GOFFI_const_arg_nil()

	AbiDefault = Abi(C._FUNC__GOFFI_FFI_DEFAULT_ABI())
	AbiFirst = Abi(C._FUNC__GOFFI_FFI_FIRST_ABI())
	AbiLast = Abi(C._FUNC__GOFFI_FFI_LAST_ABI())
	AbiEfi64 = Abi(C._FUNC__GOFFI_FFI_EFI64())
	AbiFastcall = Abi(C._FUNC__GOFFI_FFI_FASTCALL())
	AbiGnuw64 = Abi(C._FUNC__GOFFI_FFI_GNUW64())
	AbiMsCdecl = Abi(C._FUNC__GOFFI_FFI_MS_CDECL())
	AbiPascal = Abi(C._FUNC__GOFFI_FFI_PASCAL())
	AbiRegister = Abi(C._FUNC__GOFFI_FFI_REGISTER())
	AbiStdcall = Abi(C._FUNC__GOFFI_FFI_STDCALL())
	AbiSysv = Abi(C._FUNC__GOFFI_FFI_SYSV())
	AbiThiscall = Abi(C._FUNC__GOFFI_FFI_THISCALL())
	AbiUnix64 = Abi(C._FUNC__GOFFI_FFI_UNIX64())
	AbiWin64 = Abi(C._FUNC__GOFFI_FFI_WIN64())
}

func NewCif(abi Abi, output Type, inputs []Type) (*Cif, error) {
	cif := new(Cif)
	cifSize := _sizeOfCCif()
	cif.cif = (*C.ffi_cif)(usf.Malloc(cifSize))
	usf.Memset(unsafe.Pointer(cif.cif), 0, cifSize)

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

	st := C.ffi_prep_cif(cif.cif, abi.cvt(),
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
				dst[i] = _cNil()
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

type Closure struct {
	ptr     unsafe.Pointer //function pointer on the c-side
	cif     *Cif
	closure *C.ffi_closure

	fn       func(args []unsafe.Pointer, ret unsafe.Pointer)
	fnHandle cgo.Handle
}

var func_store = make(map[cgo.Handle]uint64) // map[func]argc

//export closure_caller
func closure_caller(cif *C.ffi_cif, ret, args, userData unsafe.Pointer) {
	hd := *(*cgo.Handle)(userData)
	fn := hd.Value().(func(args []unsafe.Pointer, ret unsafe.Pointer))
	argc := *(*[]unsafe.Pointer)(usf.Slice(args, func_store[hd]))
	fn(argc, ret)
}

func (cif *Cif) CreateClosure(fn func(args []unsafe.Pointer, ret unsafe.Pointer)) *Closure {
	cls := (*Closure)(usf.Malloc(_sizeOfCCif()))
	cls.cif = cif

	cls.ptr = usf.MallocN(1, 8)
	cls.closure = (*C.ffi_closure)(C.ffi_closure_alloc(
		C.uint64_t(_sizeOfCClosure()), (*unsafe.Pointer)(cls.ptr)))

	cls.fn = fn
	cls.fnHandle = cgo.NewHandle(cls.fn)
	func_store[cls.fnHandle] = uint64(cif.cif.nargs)
	C.ffi_prep_closure_loc(cls.closure, cif.cif,
		(*[0]byte)(C.closure_caller), unsafe.Pointer(&cls.fnHandle), usf.Pop(cls.ptr))

	return cls
}

func (cls *Closure) CFuncPtr() unsafe.Pointer { return usf.Pop(cls.ptr) }
func (cls *Closure) Free() {
	if cls == nil {
		return
	}

	delete(func_store, cls.fnHandle)
	cls.fnHandle.Delete()
	cls.fn = nil
	if cls.closure != nil {
		C.ffi_closure_free(unsafe.Pointer(cls.closure))
	}
	if cls.ptr != nil {
		usf.Free(cls.ptr)
	}
}
