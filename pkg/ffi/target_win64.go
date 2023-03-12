//go:build windows && !386
// +build windows,!386

package ffi

//#include <ffi.h>
import "C"

const (
	AbiFirst Abi = C.FFI_FIRST_ABI
	//AbiUinix64
	AbiWin64 Abi = C.FFI_WIN64
	//AbiEfi64
	AbiGnuW64  Abi = C.FFI_GNUW64
	AbiLastAbi Abi = C.FFI_LAST_ABI
	AbiDefault Abi = C.FFI_DEFAULT_ABI
)
