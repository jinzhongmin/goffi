package c

import "C"
import (
	"errors"
	"unsafe"

	"github.com/jinzhongmin/goffi/pkg/dlfcn"
	"github.com/jinzhongmin/goffi/pkg/ffi"
	"github.com/jinzhongmin/usf"
)

type LibMode int32

const (
	ModeNow    LibMode = LibMode(dlfcn.RTLDNow)
	ModeLazy   LibMode = LibMode(dlfcn.RTLDLazy)
	ModeGlobal LibMode = LibMode(dlfcn.RTLDGlobal)
	ModeLocal  LibMode = LibMode(dlfcn.RTLDLocal)
)

type Val struct {
	p unsafe.Pointer
}

func (v *Val) Free() {
	if v != nil && v.p != nil {
		usf.Free(v.p)
		v.p = nil
	}
}
func (v *Val) Ptr() unsafe.Pointer       { return (*(*[1]unsafe.Pointer)(v.p))[0] }
func (v *Val) Bool() bool                { return v.I32() != 0 }
func (v *Val) U8() uint8                 { return (*(*[1]uint8)(v.p))[0] }
func (v *Val) I8() int8                  { return (*(*[1]int8)(v.p))[0] }
func (v *Val) U16() uint16               { return (*(*[1]uint16)(v.p))[0] }
func (v *Val) I16() int16                { return (*(*[1]int16)(v.p))[0] }
func (v *Val) U32() uint32               { return (*(*[1]uint32)(v.p))[0] }
func (v *Val) I32() int32                { return (*(*[1]int32)(v.p))[0] }
func (v *Val) U64() uint64               { return (*(*[1]uint64)(v.p))[0] }
func (v *Val) I64() int64                { return (*(*[1]int64)(v.p))[0] }
func (v *Val) F32() float32              { return (*(*[1]float32)(v.p))[0] }
func (v *Val) F64() float64              { return (*(*[1]float64)(v.p))[0] }
func (v *Val) Str() string               { return GoStr(v.Ptr()) }
func (v *Val) SetPtr(val unsafe.Pointer) { (*(*[1]unsafe.Pointer)(v.p))[0] = val }
func (v *Val) SetU8(val uint8)           { (*(*[1]uint8)(v.p))[0] = val }
func (v *Val) SetI8(val int8)            { (*(*[1]int8)(v.p))[0] = val }
func (v *Val) SetU16(val uint16)         { (*(*[1]uint16)(v.p))[0] = val }
func (v *Val) SetI16(val int16)          { (*(*[1]int16)(v.p))[0] = val }
func (v *Val) SetU32(val uint32)         { (*(*[1]uint32)(v.p))[0] = val }
func (v *Val) SetI32(val int32)          { (*(*[1]int32)(v.p))[0] = val }
func (v *Val) SetU64(val uint64)         { (*(*[1]uint64)(v.p))[0] = val }
func (v *Val) SetI64(val int64)          { (*(*[1]int64)(v.p))[0] = val }
func (v *Val) SetF32(val float32)        { (*(*[1]float32)(v.p))[0] = val }
func (v *Val) SetF64(val float64)        { (*(*[1]float64)(v.p))[0] = val }

func CStr(s string) unsafe.Pointer  { return unsafe.Pointer(C.CString(s)) }
func GoStr(p unsafe.Pointer) string { return C.GoString((*C.char)(p)) }
func CBool(v bool) int32 {
	if v {
		return 1
	}
	return 0
}
func GoBool(v int32) bool { return v == 1 }

// type UChar C.uchar
// type Char C.char
// func (uc *UChar) GoStr() string        { return C.GoString((*C.char)((unsafe.Pointer)(uc))) }
// func (uc *UChar) U8S(n uint64) []uint8 { return *((*[]uint8)(usf.Slice(unsafe.Pointer(uc), n))) }
// func (c *Char) GoStr() string          { return C.GoString((*C.char)((unsafe.Pointer)(c))) }
// func (c *Char) U8S(n uint64) []uint8   { return *((*[]uint8)(usf.Slice(unsafe.Pointer(c), n))) }

type Type unsafe.Pointer

var (
	Void    Type = Type(ffi.Void)
	Pointer Type = Type(ffi.Pointer)
	U8      Type = Type(ffi.Uint8)
	I8      Type = Type(ffi.Int8)
	U16     Type = Type(ffi.Uint16)
	I16     Type = Type(ffi.Int16)
	U32     Type = Type(ffi.Uint32)
	I32     Type = Type(ffi.Int32)
	U64     Type = Type(ffi.Uint64)
	I64     Type = Type(ffi.Int64)

	F32               Type = Type(ffi.Float)
	F64               Type = Type(ffi.Double)
	F128              Type = Type(ffi.LongDouble)
	ComplexFloat      Type = Type(ffi.ComplexFloat)
	Complexdouble     Type = Type(ffi.Complexdouble)
	ComplexLongdouble Type = Type(ffi.ComplexLongdouble)
)

func TypeFree(typ Type) {
	ffi.TypeFree(ffi.Type(typ))
}
func TypeStruct(size uint64, alignment uint16, elms []Type) Type {
	return Type(ffi.Struct(size, alignment, *(*[]ffi.Type)(unsafe.Pointer(&elms))))
}

type fn struct {
	ptr unsafe.Pointer
	cif *ffi.Cif
}

type Lib struct {
	handle *dlfcn.Handle

	cacheFlag bool
	cache     map[string]*fn
	cachePtr  map[unsafe.Pointer]*fn
}

func newLib(handle *dlfcn.Handle, cacheEnable bool) *Lib {
	var fns map[string]*fn = nil
	var fnsP map[unsafe.Pointer]*fn = nil
	if cacheEnable {
		fns = make(map[string]*fn)
		fnsP = make(map[unsafe.Pointer]*fn)
	}

	return &Lib{handle: handle, cacheFlag: cacheEnable, cache: fns, cachePtr: fnsP}
}

// init lib by shared library,libpath look like ./libxxx.dll
func NewLib(libpath string, mod LibMode, cacheEnable bool) (*Lib, error) {
	l, err := dlfcn.Open(libpath, dlfcn.Mode(mod))
	if err != nil {
		return nil, errors.Join(errors.New("load lib error"), err)
	}

	return newLib(l, cacheEnable), nil
}

// init lib by other inited library
func NewLibFrom(lib *Lib, cacheEnable bool) *Lib {
	return newLib(lib.handle, cacheEnable)
}

// get lib symbol
func (lib *Lib) Symbol(fn string) (unsafe.Pointer, error) {
	return lib.handle.Symbol(fn)
}

// public return value ptr
var _retVoid = ffi.NewCifRetPtr()

func (lib *Lib) CallSymbolPtr(symbolPtr unsafe.Pointer, outType Type, inTypes []Type, args []interface{}) *Val {
	shouldStore := false

	//cache and has
	if lib.cacheFlag {
		if fn, ok := lib.cachePtr[symbolPtr]; ok {
			if outType == Void {
				fn.cif.Call(symbolPtr, args, _retVoid)
				return nil
			}

			ret := ffi.NewCifRetPtr()
			fn.cif.Call(symbolPtr, args, ret)
			return &Val{unsafe.Pointer(ret)}
		}
		shouldStore = true
	}

	//create cif
	var inTyps []ffi.Type = nil
	if inTypes != nil {
		inTyps = *(*[]ffi.Type)(unsafe.Pointer(&inTypes))
	}
	cif, err := ffi.NewCif(ffi.AbiDefault, ffi.Type(outType), inTyps)
	if err != nil {
		panic(err)
	}

	defer func() {
		if shouldStore {
			lib.cachePtr[symbolPtr] = &fn{
				ptr: symbolPtr,
				cif: cif,
			}
			return
		}
		cif.Free()
	}()

	if outType == Void {
		cif.Call(symbolPtr, args, _retVoid)
		return nil
	}

	ret := ffi.NewCifRetPtr()
	cif.Call(symbolPtr, args, ret)
	return &Val{unsafe.Pointer(ret)}
}
func (lib *Lib) Call(symbol string, outType Type, inTypes []Type, args []interface{}) *Val {
	shouldStore := false

	//cache and has, call and return
	if lib.cacheFlag {
		if fn, ok := lib.cache[symbol]; ok {
			if outType == Void {
				fn.cif.Call(fn.ptr, args, _retVoid)
				return nil
			}

			ret := ffi.NewCifRetPtr()
			fn.cif.Call(fn.ptr, args, ret)
			return &Val{unsafe.Pointer(ret)}
		}
		shouldStore = true
	}

	//load symbol
	ptr, err := lib.Symbol(symbol)
	if err != nil {
		panic(err)
	}

	//create cif
	var inTyps []ffi.Type = nil
	if inTypes != nil {
		inTyps = *(*[]ffi.Type)(unsafe.Pointer(&inTypes))
	}
	cif, err := ffi.NewCif(ffi.AbiDefault, ffi.Type(outType), inTyps)
	if err != nil {
		panic(err)
	}

	defer func() {
		if shouldStore {
			lib.cache[symbol] = &fn{
				ptr: ptr,
				cif: cif,
			}
			return
		}
		cif.Free()
	}()

	if outType == Void {
		cif.Call(ptr, args, _retVoid)
		return nil
	}

	ret := ffi.NewCifRetPtr()
	cif.Call(ptr, args, ret)
	return &Val{unsafe.Pointer(ret)}
}

func Call(ptr unsafe.Pointer, outType Type, inTypes []Type, args []interface{}) *Val {
	var inTyps []ffi.Type = nil
	if inTypes != nil {
		inTyps = *(*[]ffi.Type)(unsafe.Pointer(&inTypes))
	}

	cif, err := ffi.NewCif(ffi.AbiDefault, ffi.Type(outType), inTyps)
	if err != nil {
		panic(err)
	}

	defer cif.Free()

	if outType == Void {
		cif.Call(ptr, args, _retVoid)
		return nil
	}

	ret := ffi.NewCifRetPtr()
	cif.Call(ptr, args, ret)
	return &Val{unsafe.Pointer(ret)}
}

type Abi uint32

var (
	AbiDefault  Abi = Abi(ffi.AbiDefault)
	AbiFirst    Abi = Abi(ffi.AbiFirst)
	AbiLast     Abi = Abi(ffi.AbiLast)
	AbiEfi64    Abi = Abi(ffi.AbiEfi64)
	AbiFastcall Abi = Abi(ffi.AbiFastcall)
	AbiGnuw64   Abi = Abi(ffi.AbiGnuw64)
	AbiMsCdecl  Abi = Abi(ffi.AbiMsCdecl)
	AbiPascal   Abi = Abi(ffi.AbiPascal)
	AbiRegister Abi = Abi(ffi.AbiRegister)
	AbiStdcall  Abi = Abi(ffi.AbiStdcall)
	AbiSysv     Abi = Abi(ffi.AbiSysv)
	AbiThiscall Abi = Abi(ffi.AbiThiscall)
	AbiUnix64   Abi = Abi(ffi.AbiUnix64)
	AbiWin64    Abi = Abi(ffi.AbiWin64)
)

type Fn struct {
	cls *ffi.Closure
}

func NewFn(abi Abi, outType Type, inTypes []Type, relFn func(args []Val, ret *Val)) *Fn {
	it := make([]ffi.Type, 0)
	for _, t := range inTypes {
		it = append(it, ffi.Type(t))
	}
	cls := ffi.NewClosure(ffi.ClosureConf{
		Abi:    ffi.Abi(abi),
		Output: ffi.Type(outType),
		Inputs: it,
	}, func(cp *ffi.ClosureParams) {
		args := make([]Val, 0)
		for _, v := range cp.Args {
			args = append(args, Val{v})
		}

		relFn(args, &Val{cp.Return})
	}, []interface{}{})
	return &Fn{cls: cls}
}
func (f *Fn) Cptr() unsafe.Pointer { return f.cls.Cfn() }
func (f *Fn) Free() {
	if f != nil {
		f.cls.Free()
	}
}
