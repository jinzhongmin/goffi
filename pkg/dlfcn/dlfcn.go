package dlfcn

/*
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include <dlfcn.h>
void *rtdl_default(){ return (void *)0; }
void *rtdl_next(){ return (void *)-1; }
*/
import "C"
import (
	"errors"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var charsetTransformer transform.Transformer

func SetDefaultCharset(t transform.Transformer) { charsetTransformer = t }

func init() {
	DlsymDefault = C.rtdl_default()
	DlsymNext = C.rtdl_next()
	charsetTransformer = simplifiedchinese.GBK.NewDecoder()
}

type Mode C.int

const (
	/* Relocations are performed when the object is loaded. */
	RTLDNow Mode = C.RTLD_NOW
	/* Relocations are performed at an implementation-defined time.
	 * Windows API does not support lazy symbol resolving (when first reference
	 * to a given symbol occurs). So RTLD_LAZY implementation is same as RTLD_NOW.
	 */
	RTLDLazy   Mode = C.RTLD_LAZY
	RTLDGlobal Mode = C.RTLD_GLOBAL /* All symbols are available for relocation processing of other modules. */
	RTLDLocal  Mode = C.RTLD_LOCAL  /* All symbols are not made available for relocation processing by other modules. */

)

// /* Get diagnostic information. */
// DLFCN_EXPORT char *dlerror(void);
func Error() error {
	e := C.dlerror()
	if e == nil {
		return nil
	}
	r, _, err := transform.String(charsetTransformer, C.GoString(e))
	if err != nil {
		return errors.New(C.GoString(e))
	}
	return errors.New(r)
}

type Handle struct {
	c unsafe.Pointer
}

// /* Open a symbol table handle. */
// DLFCN_EXPORT void *dlopen(const char *file, int mode);
func Open(file string, mod Mode) (*Handle, error) {
	f := C.CString(file)
	defer C.free(unsafe.Pointer(f))

	h := C.dlopen(f, C.int(mod))
	if h == nil {
		return nil, Error()
	}

	return &Handle{c: h}, nil
}

func GetDefaultHandle() *Handle { return &Handle{c: DlsymDefault} }
func GetNextHandle() *Handle    { return &Handle{c: DlsymNext} }

func (hd *Handle) Ptr() unsafe.Pointer { return hd.c }

// /* Close a symbol table handle. */
// DLFCN_EXPORT int dlclose(void *handle);
func (hd *Handle) Close() {
	if hd != nil && hd.c != DlsymDefault && hd.c != DlsymNext {
		C.dlclose(hd.c)
	}
}

// /* Get the address of a symbol from a symbol table handle. */
// DLFCN_EXPORT void *dlsym(void *handle, const char *name);
func (hd Handle) Symbol(name string) (unsafe.Pointer, error) {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	p := C.dlsym(hd.c, n)
	if p == nil {
		return nil, Error()
	}

	return p, nil
}

// /* Structure filled in by dladdr() */
// type Info struct {
// 	Fname string
// 	Fbase unsafe.Pointer
// 	Sname string
// 	Saddr unsafe.Pointer
// }

// /* Translate address to symbolic information (no POSIX standard) */
// DLFCN_EXPORT int dladdr(const void *addr, Dl_info *info);
// func Address(addr unsafe.Pointer) (*Info, error) {
// 	inf := (*C.Dl_info)(C.malloc(32))
// 	defer C.free(unsafe.Pointer(inf))
// 	i := C.dladdr(addr, inf)
// 	if i == 0 {
// 		return nil, errors.New("not found")
// 	}
// 	return &Info{
// 		Fname: C.GoString(inf.dli_fname),
// 		Fbase: inf.dli_fbase,
// 		Sname: C.GoString(inf.dli_sname),
// 		Saddr: inf.dli_saddr,
// 	}, nil
// }

// The symbol lookup happens in the normal global scope.
// #define RTLD_DEFAULT    ((void *)0)
var DlsymDefault unsafe.Pointer

// Specifies the next object after this one that defines name.
// #define RTLD_NEXT       ((void *)-1)
var DlsymNext unsafe.Pointer

// Dlsym
func Dlsym(lookup unsafe.Pointer, name string) (unsafe.Pointer, error) {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	r := C.dlsym(unsafe.Pointer(lookup), n)
	if r == nil {
		return nil, Error()
	}
	return r, nil
}
