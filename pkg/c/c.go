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

	CNilStr string = "\r\n\n\r\000\n\r\n\r\000\t\a"
)

func init() {
	_retVoid = usf.Malloc(8)
	usf.Push(_retVoid, nil)
}

var (
	_retVoid unsafe.Pointer // public return value ptr

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

func Free(p unsafe.Pointer) { usf.Free(p) }
func CStr(s string) unsafe.Pointer {
	if s == CNilStr {
		return nil
	}
	return unsafe.Pointer(C.CString(s))
}
func CBool(v bool) int32 {
	if v {
		return 1
	}
	return 0
}
func GoStr(p unsafe.Pointer) string { return C.GoString((*C.char)(p)) }
func GoBool(v int32) bool           { return v == 1 }

type (
	Bool    int32
	BoolRef *Bool
)

const (
	True  Bool = 1
	False Bool = 0
)

func BoolFrom(b bool) Bool {
	if b {
		return True
	}
	return False
}
func BoolPtrFrom(b *bool) *Bool {
	r := False
	if *b {
		r = True
	}
	return &r
}
func (b Bool) Is() bool      { return b != 0 }
func (b Bool) Ok() bool      { return b != 0 }
func (b Bool) Ref(ref *bool) { (*ref) = b.Is() }
func (b Bool) Not() Bool {
	if b == True {
		return False
	}
	return True
}

type (
	gostr struct {
		data unsafe.Pointer
		len  int
	}
	Str struct {
		data unsafe.Pointer
		len  uint64 //without \0
		cap  uint64
	}
)

func NewStr(str string, cap ...uint64) *Str {
	gs := *(*gostr)((unsafe.Pointer)(&str))

	var _cap uint64
	if cap != nil && cap[0] > uint64(gs.len)+1 {
		_cap = cap[0]
	} else {
		_cap = uint64(gs.len) + 1
	}

	_len := uint64(gs.len)
	_data := usf.Malloc(_cap)
	usf.Memset(_data, 0, _cap)
	usf.Memcpy(_data, gs.data, uint64(gs.len))
	return &Str{len: _len, cap: _cap, data: _data}
}
func NewStrNil() *Str { return &Str{} }
func (str *Str) Free() {
	if str == nil || str.data == nil {
		return
	}
	usf.Free(str.data)
	str.data = nil
}
func (str *Str) CheckLen() {
	bts := str._capBytes()
	i := uint64(0)
	for ; i < str.cap; i++ {
		if bts[i] != 0 {
			continue
		}
		break
	}
	str.len = i
	for ; i < str.cap; i++ {
		bts[i] = 0
	}
}
func (str *Str) Len() uint64         { return str.len } // without \0
func (str *Str) Cap() uint64         { return str.cap } // without \0
func (str *Str) Ptr() unsafe.Pointer { return str.data }
func (str *Str) Clear() {
	bs := str._capBytes()
	for i := 0; i < int(str.len); i++ {
		bs[i] = 0
	}
	str.len = 0
}
func (str *Str) AddrPtr() unsafe.Pointer {
	if str == nil {
		ptr := (*C.char)(nil)
		return unsafe.Pointer(&ptr)
	}
	return unsafe.Pointer(&str.data)
}
func (str *Str) LenBytes() []byte   { return *(*[]byte)(unsafe.Pointer(str)) } // without \0
func (str *Str) CapBytes() []byte   { return str._capBytes() }
func (str *Str) AddCap(size uint64) { str._capAadd(size) }
func (str *Str) Str() string        { return C.GoString((*C.char)(str.data)) }
func (str *Str) Get() string        { return C.GoString((*C.char)(str.data)) }
func (str *Str) Set(s string) *Str {
	gs := *(*gostr)((unsafe.Pointer)(&s))
	if uint64(gs.len) < str.len {
		for i := gs.len; i < int(str.len); i++ {
			str.CapBytes()[i] = 0
		}
	} else if str.cap < uint64(gs.len+1) {
		str._capAadd(uint64(gs.len) + 32 - str.cap)
	}
	str.len = uint64(gs.len)
	usf.Memcpy(str.data, gs.data, str.len)
	return str
}
func (str *Str) SetNil() {
	usf.Free(str.data)
	str.data = nil
	str.len = 0
	str.cap = 0
}
func (str *Str) Append(s string) {
	gs := *(*gostr)((unsafe.Pointer)(&s))
	all_len := str.len + uint64(gs.len)

	if str.cap < all_len+1 {
		str._capAadd(all_len + 32 - str.cap)
	}
	start := unsafe.Pointer(&str._capBytes()[str.len])
	usf.Memcpy(start, gs.data, uint64(gs.len))
	str.len = all_len
}
func (str *Str) CopyTo(dst *Str) {
	usf.Memcpy(dst.data, str.data, str.len)
}
func (str *Str) _capBytes() []byte { return *(*[]byte)(usf.Slice(str.data, str.cap)) }
func (str *Str) _capAadd(size uint64) {
	i := str.cap
	str.cap += size
	str.data = usf.Realloc(str.data, str.cap)
	bts := str.CapBytes()
	for ; i < str.cap; i++ {
		bts[i] = 0
	}
}

type (
	ListItem interface {
		Ptr() unsafe.Pointer
	}
	List struct {
		len  uint64
		cap  uint64
		data unsafe.Pointer
	}
)

func NewList(cap uint64) *List {
	data := usf.MallocN(cap, 8)
	usf.Memset(data, 0, cap*8)
	return &List{len: 0, cap: cap, data: data}
}
func NewListFrom(items []ListItem, cap ...uint64) *List {
	items_len := uint64(len(items))
	var _cap uint64
	if cap != nil && cap[0] > items_len {
		_cap = cap[0]
	} else {
		_cap = items_len
	}
	_data := usf.MallocN(_cap, 8)
	usf.Memset(_data, 0, _cap*8)
	for i := range items {
		usf.PushAt(_data, uint64(i), items[i].Ptr())
	}
	return &List{len: items_len, cap: _cap, data: _data}
}
func (li *List) Free() {
	li.len = 0
	li.cap = 0
	usf.Free(li.data)
}
func (li *List) FreeAll() {
	ps := li.CapPtrs()
	for i := 0; i < int(li.len); i++ {
		usf.Free(ps[i])
	}
	li.len = 0
	li.cap = 0
	usf.Free(li.data)
}
func (li *List) Len() uint64               { return li.len }
func (li *List) Ptr() unsafe.Pointer       { return li.data }
func (li *List) AddrPtr() unsafe.Pointer   { return unsafe.Pointer(&li.data) }
func (li *List) CapPtrs() []unsafe.Pointer { return li._capPtrs() }
func (li *List) Set(idx uint64, ptr unsafe.Pointer) {
	if !(idx < li.len) {
		panic("idx out of rang")
	}
	usf.PushAt(li.data, idx, ptr)
}
func (li *List) Get(idx uint64) unsafe.Pointer {
	if !(idx < li.cap) {
		panic("idx out of rang")
	}
	return usf.PopAt(li.data, idx)
}
func (li *List) Append(ptr unsafe.Pointer) {
	if li.len == li.cap {
		li._capAdd(8)
	}
	li._capPtrs()[li.len] = ptr
	li.len += 1
}
func (li *List) _capPtrs() []unsafe.Pointer { return *(*[]unsafe.Pointer)(usf.Slice(li.data, li.cap)) }
func (li *List) _capAdd(size uint64) {
	i := li.cap
	li.cap += size
	li.data = usf.Realloc(li.data, li.cap*8)
	ps := li._capPtrs()
	for ; i < li.cap; i++ {
		ps[i] = nil
	}
}

type (
	prototypes []*prototype
	prototype  struct {
		cif     *ffi.Cif
		outType Type
		inTypes []Type
	}
)

func (protp prototype) same(outType Type, inTypes []Type) bool {
	if protp.outType != outType || len(protp.inTypes) != len(inTypes) {
		return false
	}
	for i := range protp.inTypes {
		if protp.inTypes[i] != inTypes[i] {
			return false
		}
	}
	return true
}

func (protps prototypes) lookup(outType Type, inTypes []Type) *prototype {
	for i := range protps {
		if protps[i].same(outType, inTypes) {
			return protps[i]
		}
	}
	return nil
}

type Lib struct {
	prototypes prototypes
	handle     *dlfcn.Handle
}

func LoadLib(libpath string, mod LibMode) (*Lib, error) {
	l, err := dlfcn.Open(libpath, dlfcn.Mode(mod))
	if err != nil {
		return nil, errors.Join(errors.New("load lib error"), err)
	}
	return &Lib{handle: l, prototypes: make([]*prototype, 0)}, nil
}
func NewLibDefault() *Lib {
	return &Lib{handle: dlfcn.GetDefaultHandle(), prototypes: make([]*prototype, 0)}
}
func NewLibNext() *Lib {
	return &Lib{handle: dlfcn.GetNextHandle(), prototypes: make([]*prototype, 0)}
}

func NewLib(libpath string, mod LibMode) (*Lib, error) {
	l, err := dlfcn.Open(libpath, dlfcn.Mode(mod))
	if err != nil {
		return nil, errors.Join(errors.New("load lib error"), err)
	}
	return &Lib{handle: l, prototypes: make([]*prototype, 0)}, nil
}

func NewLibFrom(lib *Lib) *Lib {
	return &Lib{handle: lib.handle, prototypes: make([]*prototype, 0)}
}

func (lib *Lib) UnLoad() { lib.handle.Close() }
func (lib *Lib) lookup(outType Type, inTypes []Type) *prototype {
	protp := lib.prototypes.lookup(outType, inTypes)
	if protp != nil {
		return protp
	}

	cif, err := ffi.NewCif(ffi.AbiDefault, ffi.Type(outType), *(*[]ffi.Type)(unsafe.Pointer(&inTypes)))
	if err != nil {
		panic(err)
	}

	protp = &prototype{
		cif:     cif,
		outType: outType,
		inTypes: inTypes,
	}
	lib.prototypes = append(lib.prototypes, protp)

	return protp
}
func (lib *Lib) Symbol(fn string) unsafe.Pointer {
	p, err := lib.handle.Symbol(fn)
	if err != nil {
		panic(err)
	}
	return p
}
func (lib *Lib) Call(fp *FuncPrototype, args []interface{}) *Value {
	if fp.complete {
		if fp.OutType == Void {
			fp.Cif.Call(fp.Ptr, args, nil)
			return nil
		}

		ret := usf.Malloc(8)
		retV := (*Value)(ret)
		retV.v = nil
		fp.Cif.Call(fp.Ptr, args, ret)
		return retV
	}

	if fp.Ptr == nil {
		fp.Ptr = lib.Symbol(fp.Name)
	}
	if fp.Cif == nil {
		protp := lib.lookup(fp.OutType, fp.InTypes)
		fp.Cif = protp.cif
	}
	fp.complete = true

	if fp.OutType == Void {
		fp.Cif.Call(fp.Ptr, args, nil)
		return nil
	}

	ret := usf.Malloc(8)
	retV := (*Value)(ret)
	retV.v = nil
	fp.Cif.Call(fp.Ptr, args, ret)
	return retV
}

type FuncPrototype struct {
	Name    string //func name in C
	OutType Type   //return type int C
	InTypes []Type //params type int C

	Ptr unsafe.Pointer //dlfcn func pointer
	Cif *ffi.Cif       //cif

	complete bool //mean this strucr is complete
}

// dlsym:dlfcn.DlsymDefault|dlfcn.DlsymNext|dlfcn.Handle.Ptr()
func (fp *FuncPrototype) Create(dlsym unsafe.Pointer) (err error) {
	if fp.Ptr == nil {
		fp.Ptr, err = dlfcn.Dlsym(dlsym, fp.Name)
		if err != nil {
			return
		}
	}
	if fp.Cif == nil {
		fp.Cif, err = ffi.NewCif(ffi.AbiDefault, ffi.Type(fp.OutType),
			*(*[]ffi.Type)(unsafe.Pointer(&fp.InTypes)))
		if err != nil {
			return
		}
	}
	fp.complete = true

	return
}
func (fp *FuncPrototype) Free() {
	if fp != nil && fp.Cif != nil {
		fp.Cif.Free()
	}
}
func (fp *FuncPrototype) Call(args []interface{}) *Value {
	if fp.OutType == Void {
		fp.Cif.Call(fp.Ptr, args, nil)
		return nil
	}

	ret := usf.Malloc(8)
	retV := (*Value)(ret)
	retV.v = nil
	fp.Cif.Call(fp.Ptr, args, ret)
	return retV
}

type Value struct{ v unsafe.Pointer }

func (v *Value) Free()                   { usf.Free(unsafe.Pointer(v)) }
func (v *Value) U8() uint8               { return *(*uint8)(unsafe.Pointer(v)) }
func (v *Value) I8() int8                { return *(*int8)(unsafe.Pointer(v)) }
func (v *Value) U16() uint16             { return *(*uint16)(unsafe.Pointer(v)) }
func (v *Value) I16() int16              { return *(*int16)(unsafe.Pointer(v)) }
func (v *Value) U32() uint32             { return *(*uint32)(unsafe.Pointer(v)) }
func (v *Value) I32() int32              { return *(*int32)(unsafe.Pointer(v)) }
func (v *Value) U64() uint64             { return *(*uint64)(unsafe.Pointer(v)) }
func (v *Value) I64() int64              { return *(*int64)(unsafe.Pointer(v)) }
func (v *Value) F32() float32            { return *(*float32)(unsafe.Pointer(v)) }
func (v *Value) F64() float64            { return *(*float64)(unsafe.Pointer(v)) }
func (v *Value) Ptr() unsafe.Pointer     { return v.v }
func (v *Value) Str() string             { return GoStr(v.Ptr()) }
func (v *Value) Bool() bool              { return v.I32() != 0 }
func (v *Value) U8Free() uint8           { defer v.Free(); return *(*uint8)(unsafe.Pointer(v)) }
func (v *Value) I8Free() int8            { defer v.Free(); return *(*int8)(unsafe.Pointer(v)) }
func (v *Value) U16Free() uint16         { defer v.Free(); return *(*uint16)(unsafe.Pointer(v)) }
func (v *Value) I16Free() int16          { defer v.Free(); return *(*int16)(unsafe.Pointer(v)) }
func (v *Value) U32Free() uint32         { defer v.Free(); return *(*uint32)(unsafe.Pointer(v)) }
func (v *Value) I32Free() int32          { defer v.Free(); return *(*int32)(unsafe.Pointer(v)) }
func (v *Value) U64Free() uint64         { defer v.Free(); return *(*uint64)(unsafe.Pointer(v)) }
func (v *Value) I64Free() int64          { defer v.Free(); return *(*int64)(unsafe.Pointer(v)) }
func (v *Value) F32Free() float32        { defer v.Free(); return *(*float32)(unsafe.Pointer(v)) }
func (v *Value) F64Free() float64        { defer v.Free(); return *(*float64)(unsafe.Pointer(v)) }
func (v *Value) PtrFree() unsafe.Pointer { defer v.Free(); return v.v }
func (v *Value) StrFree() string         { defer v.Free(); return GoStr(v.Ptr()) }
func (v *Value) BoolFree() bool          { return v.I32Free() != 0 }
func (v *Value) SetU8(i uint8)           { *(*uint8)(unsafe.Pointer(v)) = i }
func (v *Value) SetI8(i int8)            { *(*int8)(unsafe.Pointer(v)) = i }
func (v *Value) SetU16(i uint16)         { *(*uint16)(unsafe.Pointer(v)) = i }
func (v *Value) SetI16(i int16)          { *(*int16)(unsafe.Pointer(v)) = i }
func (v *Value) SetU32(i uint32)         { *(*uint32)(unsafe.Pointer(v)) = i }
func (v *Value) SetI32(i int32)          { *(*int32)(unsafe.Pointer(v)) = i }
func (v *Value) SetU64(i uint64)         { *(*uint64)(unsafe.Pointer(v)) = i }
func (v *Value) SetI64(i int64)          { *(*int64)(unsafe.Pointer(v)) = i }
func (v *Value) SetF32(i float32)        { *(*float32)(unsafe.Pointer(v)) = i }
func (v *Value) SetF64(i float64)        { *(*float64)(unsafe.Pointer(v)) = i }
func (v *Value) SetPtr(i unsafe.Pointer) { *(*unsafe.Pointer)(unsafe.Pointer(v)) = i }

type Callback struct {
	*ffi.Closure
	//Converts the input and output variables to their real types and calls CallbackFunc
	CallbackCvt  func(callback *Callback, args []*Value, ret *Value)
	CallbackFunc interface{}
}

func (cb *Callback) FuncPtr() unsafe.Pointer { return cb.Closure.FuncPtr() }
func (cb *Callback) Free() {
	if cb == nil {
		return
	}
	cb.Closure.Free()
	cb.CallbackFunc = nil
	cb.CallbackCvt = nil

	usf.Free(unsafe.Pointer(cb))
}
func createClosureCallback(cb *Callback) func(args []unsafe.Pointer, ret unsafe.Pointer) {
	fn := func(args []unsafe.Pointer, ret unsafe.Pointer) {
		if cb == nil || cb.CallbackCvt == nil {
			return
		}
		_args := *(*[]*Value)(usf.Slice(unsafe.Pointer(&args[0]), uint64(len(args))))
		cb.CallbackCvt(cb, _args, (*Value)(ret))
	}
	return fn
}
func NewCallback(abi Abi, outType Type, inTypes []Type) *Callback {
	cb := (*Callback)(usf.MallocOf(1, Callback{}))
	cb.Closure = ffi.NewClosure(ffi.Abi(abi), ffi.Type(outType),
		*(*[]ffi.Type)(unsafe.Pointer(&inTypes)), createClosureCallback(cb))
	cb.CallbackCvt = nil
	cb.CallbackFunc = nil
	return cb
}
