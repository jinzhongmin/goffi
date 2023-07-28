package ffi

/*
#include <ffi.h>
unsigned int __GOFFI_FFI_DEFAULT_ABI = 9999;
unsigned int __GOFFI_FFI_EFI64 = 9999;
unsigned int __GOFFI_FFI_FASTCALL = 9999;
unsigned int __GOFFI_FFI_FIRST_ABI = 9999;
unsigned int __GOFFI_FFI_GNUW64 = 9999;
unsigned int __GOFFI_FFI_LAST_ABI = 9999;
unsigned int __GOFFI_FFI_MS_CDECL = 9999;
unsigned int __GOFFI_FFI_PASCAL = 9999;
unsigned int __GOFFI_FFI_REGISTER = 9999;
unsigned int __GOFFI_FFI_STDCALL = 9999;
unsigned int __GOFFI_FFI_SYSV = 9999;
unsigned int __GOFFI_FFI_THISCALL = 9999;
unsigned int __GOFFI_FFI_UNIX64 = 9999;
unsigned int __GOFFI_FFI_WIN64 = 9999;

void init_abi(){
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

func init() {
	C.init_abi()
	AbiDefault = Abi(C.__GOFFI_FFI_DEFAULT_ABI)
	AbiFirst = Abi(C.__GOFFI_FFI_FIRST_ABI)
	AbiLast = Abi(C.__GOFFI_FFI_LAST_ABI)
	AbiEfi64 = Abi(C.__GOFFI_FFI_EFI64)
	AbiFastcall = Abi(C.__GOFFI_FFI_FASTCALL)
	AbiGnuw64 = Abi(C.__GOFFI_FFI_GNUW64)
	AbiMsCdecl = Abi(C.__GOFFI_FFI_MS_CDECL)
	AbiPascal = Abi(C.__GOFFI_FFI_PASCAL)
	AbiRegister = Abi(C.__GOFFI_FFI_REGISTER)
	AbiStdcall = Abi(C.__GOFFI_FFI_STDCALL)
	AbiSysv = Abi(C.__GOFFI_FFI_SYSV)
	AbiThiscall = Abi(C.__GOFFI_FFI_THISCALL)
	AbiUnix64 = Abi(C.__GOFFI_FFI_UNIX64)
	AbiWin64 = Abi(C.__GOFFI_FFI_WIN64)
}

type Abi C.uint

func (a Abi) ffi_abi() C.ffi_abi {
	if a == 9999 {
		panic("bed abi")
	}
	return C.ffi_abi(a)
}

var (
	AbiDefault  Abi = Abi(C.__GOFFI_FFI_DEFAULT_ABI)
	AbiFirst    Abi = Abi(C.__GOFFI_FFI_FIRST_ABI)
	AbiLast     Abi = Abi(C.__GOFFI_FFI_LAST_ABI)
	AbiEfi64    Abi = Abi(C.__GOFFI_FFI_EFI64)
	AbiFastcall Abi = Abi(C.__GOFFI_FFI_FASTCALL)
	AbiGnuw64   Abi = Abi(C.__GOFFI_FFI_GNUW64)
	AbiMsCdecl  Abi = Abi(C.__GOFFI_FFI_MS_CDECL)
	AbiPascal   Abi = Abi(C.__GOFFI_FFI_PASCAL)
	AbiRegister Abi = Abi(C.__GOFFI_FFI_REGISTER)
	AbiStdcall  Abi = Abi(C.__GOFFI_FFI_STDCALL)
	AbiSysv     Abi = Abi(C.__GOFFI_FFI_SYSV)
	AbiThiscall Abi = Abi(C.__GOFFI_FFI_THISCALL)
	AbiUnix64   Abi = Abi(C.__GOFFI_FFI_UNIX64)
	AbiWin64    Abi = Abi(C.__GOFFI_FFI_WIN64)
)
