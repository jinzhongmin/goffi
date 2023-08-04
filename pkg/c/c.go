package c

import "C"
import (
	"errors"
	"unsafe"

	"github.com/jinzhongmin/goffi/pkg/dlfcn"
	"github.com/jinzhongmin/goffi/pkg/ffi"
	"github.com/jinzhongmin/usf"
)

type Abi uint32
type LibMode int32
type Type unsafe.Pointer

const (
	ModeNow    LibMode = LibMode(dlfcn.RTLDNow)
	ModeLazy   LibMode = LibMode(dlfcn.RTLDLazy)
	ModeGlobal LibMode = LibMode(dlfcn.RTLDGlobal)
	ModeLocal  LibMode = LibMode(dlfcn.RTLDLocal)
)

var (
	_retVoid = ffi.NewZeroPtr() // public return value ptr

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

type Value struct {
	p unsafe.Pointer
}

func NewVal(p unsafe.Pointer) *Value { return &Value{p} }
func (v *Value) Free() {
	if v != nil && v.p != nil {
		usf.Free(v.p)
		v.p = nil
	}
}
func (v *Value) Ptr() unsafe.Pointer { return (*(*[1]unsafe.Pointer)(v.p))[0] }
func (v *Value) Bool() bool          { return v.I32() != 0 }
func (v *Value) U8() uint8           { return (*(*[1]uint8)(v.p))[0] }
func (v *Value) I8() int8            { return (*(*[1]int8)(v.p))[0] }
func (v *Value) U16() uint16         { return (*(*[1]uint16)(v.p))[0] }
func (v *Value) I16() int16          { return (*(*[1]int16)(v.p))[0] }
func (v *Value) U32() uint32         { return (*(*[1]uint32)(v.p))[0] }
func (v *Value) I32() int32          { return (*(*[1]int32)(v.p))[0] }
func (v *Value) U64() uint64         { return (*(*[1]uint64)(v.p))[0] }
func (v *Value) I64() int64          { return (*(*[1]int64)(v.p))[0] }
func (v *Value) F32() float32        { return (*(*[1]float32)(v.p))[0] }
func (v *Value) F64() float64        { return (*(*[1]float64)(v.p))[0] }
func (v *Value) Str() string         { return GoStr(v.Ptr()) }
func (v *Value) PtrFree() unsafe.Pointer {
	defer v.Free()
	return (*(*[1]unsafe.Pointer)(v.p))[0]
}
func (v *Value) BoolFree() bool {
	defer v.Free()
	return v.I32Free() != 0
}
func (v *Value) U8Free() uint8 {
	defer v.Free()
	return (*(*[1]uint8)(v.p))[0]
}
func (v *Value) I8Free() int8 {
	defer v.Free()
	return (*(*[1]int8)(v.p))[0]
}
func (v *Value) U16Free() uint16 {
	defer v.Free()
	return (*(*[1]uint16)(v.p))[0]
}
func (v *Value) I16Free() int16 {
	defer v.Free()
	return (*(*[1]int16)(v.p))[0]
}
func (v *Value) U32Free() uint32 {
	defer v.Free()
	return (*(*[1]uint32)(v.p))[0]
}
func (v *Value) I32Free() int32 {
	defer v.Free()
	return (*(*[1]int32)(v.p))[0]
}
func (v *Value) U64Free() uint64 {
	defer v.Free()
	return (*(*[1]uint64)(v.p))[0]
}
func (v *Value) I64Free() int64 {
	defer v.Free()
	return (*(*[1]int64)(v.p))[0]
}
func (v *Value) F32Free() float32 {
	defer v.Free()
	return (*(*[1]float32)(v.p))[0]
}
func (v *Value) F64Free() float64 {
	defer v.Free()
	return (*(*[1]float64)(v.p))[0]
}
func (v *Value) StrFree() string {
	defer v.Free()
	return GoStr(v.PtrFree())
}
func (v *Value) SetPtr(val unsafe.Pointer) { (*(*[1]unsafe.Pointer)(v.p))[0] = val }
func (v *Value) SetU8(val uint8)           { (*(*[1]uint8)(v.p))[0] = val }
func (v *Value) SetI8(val int8)            { (*(*[1]int8)(v.p))[0] = val }
func (v *Value) SetU16(val uint16)         { (*(*[1]uint16)(v.p))[0] = val }
func (v *Value) SetI16(val int16)          { (*(*[1]int16)(v.p))[0] = val }
func (v *Value) SetU32(val uint32)         { (*(*[1]uint32)(v.p))[0] = val }
func (v *Value) SetI32(val int32)          { (*(*[1]int32)(v.p))[0] = val }
func (v *Value) SetU64(val uint64)         { (*(*[1]uint64)(v.p))[0] = val }
func (v *Value) SetI64(val int64)          { (*(*[1]int64)(v.p))[0] = val }
func (v *Value) SetF32(val float32)        { (*(*[1]float32)(v.p))[0] = val }
func (v *Value) SetF64(val float64)        { (*(*[1]float64)(v.p))[0] = val }

func CStr(s string) unsafe.Pointer  { return unsafe.Pointer(C.CString(s)) }
func GoStr(p unsafe.Pointer) string { return C.GoString((*C.char)(p)) }
func CBool(v bool) int32 {
	if v {
		return 1
	}
	return 0
}
func GoBool(v int32) bool { return v == 1 }

type fn struct {
	ptr unsafe.Pointer
	cif *ffi.Cif
}

type prototype struct {
	cif     *ffi.Cif
	outType Type
	inTypes []Type
}

func (prop prototype) same(outType Type, inTypes []Type) bool {
	if prop.outType != outType || len(prop.inTypes) != len(inTypes) {
		return false
	}
	for i := range prop.inTypes {
		if prop.inTypes[i] != inTypes[i] {
			return false
		}
	}
	return true
}

type Lib struct {
	handle *dlfcn.Handle

	props     []*prototype
	cacheFlag bool
	cache     map[string]*fn
	cachePtr  map[unsafe.Pointer]*fn
}

func newLib(handle *dlfcn.Handle, cacheEnable bool) *Lib {
	var props []*prototype = nil
	var cache map[string]*fn = nil
	var cachePtr map[unsafe.Pointer]*fn = nil
	if cacheEnable {
		props = make([]*prototype, 0)
		cache = make(map[string]*fn)
		cachePtr = make(map[unsafe.Pointer]*fn)
	}

	return &Lib{handle: handle, cacheFlag: cacheEnable, props: props, cache: cache, cachePtr: cachePtr}
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

func (lib *Lib) checkPrototype(outType Type, inTypes []Type) *prototype {
	for i := range lib.props {
		if lib.props[i].same(outType, inTypes) {
			return lib.props[i]
		}
	}
	return nil
}

func (lib *Lib) Symbol(fn string) (unsafe.Pointer, error) {
	return lib.handle.Symbol(fn)
}

func (lib *Lib) CallSymbolPtr(symbolPtr unsafe.Pointer, outType Type, inTypes []Type, args []interface{}) *Value {
	if lib.cacheFlag {
		if fn, ok := lib.cachePtr[symbolPtr]; ok {
			if outType == Void {
				fn.cif.Call(symbolPtr, args, nil)
				return nil
			}

			ret := ffi.NewZeroPtr()
			fn.cif.Call(symbolPtr, args, ret)
			return &Value{ret}
		}
	}

	if lib.cacheFlag {
		if prop := lib.checkPrototype(outType, inTypes); prop != nil {
			fn := &fn{
				ptr: symbolPtr,
				cif: prop.cif,
			}
			lib.cachePtr[symbolPtr] = fn

			if outType == Void {
				fn.cif.Call(fn.ptr, args, nil)
				return nil
			}

			ret := ffi.NewZeroPtr()
			fn.cif.Call(fn.ptr, args, ret)
			return &Value{ret}
		}
	}

	//create cif
	cif, err := ffi.NewCif(ffi.AbiDefault, ffi.Type(outType), *(*[]ffi.Type)(unsafe.Pointer(&inTypes)))
	if err != nil {
		panic(err)
	}

	if lib.cacheFlag {
		fn := &fn{
			ptr: symbolPtr,
			cif: cif,
		}
		lib.cachePtr[symbolPtr] = fn
		lib.props = append(lib.props,
			&prototype{outType: outType, inTypes: inTypes, cif: cif})
	} else {
		cif.Free()
	}

	if outType == Void {
		cif.Call(symbolPtr, args, nil)
		return nil
	}

	ret := ffi.NewZeroPtr()
	cif.Call(symbolPtr, args, ret)
	return &Value{unsafe.Pointer(ret)}
}

func (lib *Lib) Call(symbol string, outType Type, inTypes []Type, args []interface{}) *Value {

	//cache and has, call and return
	if lib.cacheFlag {
		if fn, ok := lib.cache[symbol]; ok {
			if outType == Void {
				fn.cif.Call(fn.ptr, args, nil)
				return nil
			}

			ret := ffi.NewZeroPtr()
			fn.cif.Call(fn.ptr, args, ret)
			return &Value{ret}
		}
	}

	//load symbol
	ptr, err := lib.Symbol(symbol)
	if err != nil {
		panic(err)
	}

	if lib.cacheFlag {
		if prop := lib.checkPrototype(outType, inTypes); prop != nil {
			fn := &fn{
				ptr: ptr,
				cif: prop.cif,
			}
			lib.cache[symbol] = fn

			if outType == Void {
				fn.cif.Call(fn.ptr, args, nil)
				return nil
			}

			ret := ffi.NewZeroPtr()
			fn.cif.Call(fn.ptr, args, ret)
			return &Value{ret}
		}
	}

	cif, err := ffi.NewCif(ffi.AbiDefault, ffi.Type(outType), *(*[]ffi.Type)(unsafe.Pointer(&inTypes)))
	if err != nil {
		panic(err)
	}

	if lib.cacheFlag {
		fn := &fn{
			ptr: ptr,
			cif: cif,
		}
		lib.cache[symbol] = fn
		lib.props = append(lib.props,
			&prototype{outType: outType, inTypes: inTypes, cif: cif})
	} else {
		cif.Free()
	}

	if outType == Void {
		cif.Call(ptr, args, nil)
		return nil
	}

	ret := ffi.NewZeroPtr()
	cif.Call(ptr, args, ret)
	return &Value{ret}
}

func Call(ptr unsafe.Pointer, outType Type, inTypes []Type, args []interface{}) *Value {
	cif, err := ffi.NewCif(ffi.AbiDefault, ffi.Type(outType), *(*[]ffi.Type)(unsafe.Pointer(&inTypes)))
	if err != nil {
		panic(err)
	}

	if outType == Void {
		cif.Call(ptr, args, nil)
		cif.Free()
		return nil
	}

	ret := ffi.NewZeroPtr()
	cif.Call(ptr, args, ret)
	cif.Free()
	return &Value{ret}
}

type Fn struct {
	cls *ffi.Closure
	fn  func(args []Value, ret Value)
}

func NewFn(abi Abi, outType Type, inTypes []Type, relFn func(args []Value, ret Value)) *Fn {
	cls := ffi.NewClosure(ffi.Abi(abi), ffi.Type(outType), *(*[]ffi.Type)(unsafe.Pointer(&inTypes)),
		func(args []unsafe.Pointer, ret unsafe.Pointer) {
			_args := *(*[]Value)(usf.Slice(unsafe.Pointer(&args[0]), uint64(len(args))))
			relFn(_args, Value{ret})
		})
	return &Fn{cls: cls, fn: relFn}
}
func (f *Fn) Cptr() unsafe.Pointer { return f.cls.Cfn() }
func (f *Fn) Free() {
	if f != nil {
		f.cls.Free()
	}
}
