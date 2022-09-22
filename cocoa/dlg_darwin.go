package cocoa

// #cgo darwin LDFLAGS: -framework Cocoa
// #include <stdlib.h>
// #include <sys/syslimits.h>
// #include "dlg.h"
import "C"

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

type AlertParams struct {
	p C.AlertDlgParams
}

func mkAlertParams(msg, title string, style C.AlertStyle) *AlertParams {
	a := AlertParams{C.AlertDlgParams{msg: C.CString(msg), style: style}}
	if title != "" {
		a.p.title = C.CString(title)
	}
	return &a
}

func (a *AlertParams) run() C.DlgResult {
	return C.alertDlg(&a.p)
}

func (a *AlertParams) free() {
	C.free(unsafe.Pointer(a.p.msg))
	if a.p.title != nil {
		C.free(unsafe.Pointer(a.p.title))
	}
}

func nsStr(s string) unsafe.Pointer {
	return C.NSStr(unsafe.Pointer(&[]byte(s)[0]), C.int(len(s)))
}

func YesNoDlg(msg, title string) bool {
	a := mkAlertParams(msg, title, C.MSG_YESNO)
	defer a.free()
	return a.run() == C.DLG_OK
}

func InfoDlg(msg, title string) {
	a := mkAlertParams(msg, title, C.MSG_INFO)
	defer a.free()
	a.run()
}

func ErrorDlg(msg, title string) {
	a := mkAlertParams(msg, title, C.MSG_ERROR)
	defer a.free()
	a.run()
}

const MAX_FILES = 2048
const SINGLE_BUFSIZE = C.PATH_MAX
const ALL_BUFSIZE = MAX_FILES * SINGLE_BUFSIZE

func FileDlg(save bool, title string, exts []string, relaxExt, multiple bool, startDir string, filename string) ([]string, error) {
	mode := C.LOADDLG
	if save {
		mode = C.SAVEDLG
	}
	return fileDlg(mode, title, exts, relaxExt, multiple, startDir, filename)
}

func DirDlg(title string, startDir string) (string, error) {
	ss, err := fileDlg(C.DIRDLG, title, nil, false, false, startDir, "")
	if len(ss) == 0 {
		return "", err
	}
	return ss[0], err
}

func fileDlg(mode int, title string, exts []string, relaxExt, multiple bool, startDir, filename string) ([]string, error) {
	p := C.FileDlgParams{
		mode:       C.int(mode),
		max_files:  MAX_FILES,
		single_buf: SINGLE_BUFSIZE,
	}

	if multiple {
		p.multiple = C.int(1)
	}

	written := (*C.size_t)(C.malloc(C.size_t(C.sizeof_size_t)))
	defer C.free(unsafe.Pointer(written))
	p.written = written

	p.buf = (*C.char)(C.malloc(ALL_BUFSIZE))
	defer C.free(unsafe.Pointer(p.buf))
	buf := (*(*[ALL_BUFSIZE]byte)(unsafe.Pointer(p.buf)))[:]
	if title != "" {
		p.title = C.CString(title)
		defer C.free(unsafe.Pointer(p.title))
	}
	if startDir != "" {
		p.startDir = C.CString(startDir)
		defer C.free(unsafe.Pointer(p.startDir))
	}
	if filename != "" {
		p.filename = C.CString(filename)
		defer C.free(unsafe.Pointer(p.filename))
	}
	if multiple {
		p.multiple = 1
	}
	if len(exts) > 0 {
		if len(exts) > 999 {
			panic("more than 999 extensions not supported")
		}
		ptrSize := int(unsafe.Sizeof(&title))
		p.exts = (*unsafe.Pointer)(C.malloc(C.size_t(ptrSize * len(exts))))
		defer C.free(unsafe.Pointer(p.exts))
		cext := (*(*[999]unsafe.Pointer)(unsafe.Pointer(p.exts)))[:]
		for i, ext := range exts {
			i := i
			cext[i] = nsStr(ext)
			defer C.NSRelease(cext[i])
		}
		p.numext = C.int(len(exts))
		if relaxExt {
			p.relaxext = 1
		}
	}
	switch C.fileDlg(&p) {
	case C.DLG_OK:
		fmt.Println("OK")
		s := string(buf[:(int)(*written)])
		ss := strings.Split(s, string([]byte{0}))
		return ss, nil
	case C.DLG_TOOMANY:
		return nil, errors.New("too many files selected")
	case C.DLG_CANCEL:
		fmt.Println("cancelled")

		return nil, nil
	case C.DLG_URLFAIL:
		return nil, errors.New("failed to get file-system representation for selected URL")
	}
	panic("unhandled case")
}
