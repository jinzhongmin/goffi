//go:build darwin && !amd64
// +build darwin,!amd64

package ffi

//#include <ffi.h>
import "C"

const (
	AbiFirst    Abi = C.FFI_FIRST_ABI
	AbiSysv     Abi = C.FFI_SYSV
	AbiStdcall  Abi = C.FFI_STDCALL
	AbiThiscall Abi = C.FFI_THISCALL
	AbiFastcall Abi = C.FFI_FASTCALL
	AbiMsCdecl  Abi = C.FFI_MS_CDECL
	AbiPascal   Abi = C.FFI_PASCAL
	AbiRegister Abi = C.FFI_REGISTER
	AbiLastAbi  Abi = C.FFI_LAST_ABI
	AbiDefault  Abi = C.FFI_DEFAULT_ABI
)
