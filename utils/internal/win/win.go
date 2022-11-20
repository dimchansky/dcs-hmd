//go:build windows

package win

import (
	"syscall"
	"unsafe"

	// #include <wtypes.h>
	// #include <winable.h>
	"C"
)

type (
	ATOM            uint16
	BOOL            int32
	COLORREF        uint32
	DWM_FRAME_COUNT uint64
	DWORD           uint32
	HACCEL          HANDLE
	HANDLE          uintptr
	HBITMAP         HANDLE
	HBRUSH          HANDLE
	HCURSOR         HANDLE
	HDC             HANDLE
	HDROP           HANDLE
	HDWP            HANDLE
	HENHMETAFILE    HANDLE
	HFONT           HANDLE
	HGDIOBJ         HANDLE
	HGLOBAL         HANDLE
	HGLRC           HANDLE
	HHOOK           HANDLE
	HICON           HANDLE
	HIMAGELIST      HANDLE
	HINSTANCE       HANDLE
	HKEY            HANDLE
	HKL             HANDLE
	HMENU           HANDLE
	HMODULE         HANDLE
	HMONITOR        HANDLE
	HPEN            HANDLE
	HRESULT         int32
	HRGN            HANDLE
	HRSRC           HANDLE
	HTHUMBNAIL      HANDLE
	HWND            HANDLE
	LPARAM          uintptr
	LPCVOID         unsafe.Pointer
	LRESULT         uintptr
	PVOID           unsafe.Pointer
	QPC_TIME        uint64
	SIZE_T          uintptr
	TRACEHANDLE     uintptr
	ULONG_PTR       uintptr
	WPARAM          uintptr
	WNDENUMPROC     uintptr
)

var (
	modUser32                    = syscall.NewLazyDLL("user32.dll")
	procGetWindowThreadProcessId = modUser32.NewProc("GetWindowThreadProcessId")
	procEnumWindows              = modUser32.NewProc("EnumWindows")
	procGetWindowLong            = modUser32.NewProc("GetWindowLongW")
	procSetWindowLong            = modUser32.NewProc("SetWindowLongW")

	modKernel32             = syscall.NewLazyDLL("kernel32.dll")
	procGetCurrentProcessId = modKernel32.NewProc("GetCurrentProcessId")
)

// GetWindowLong and GetWindowLongPtr constants
const (
	GWL_EXSTYLE     = -20
	GWL_STYLE       = -16
	GWL_WNDPROC     = -4
	GWLP_WNDPROC    = -4
	GWL_HINSTANCE   = -6
	GWLP_HINSTANCE  = -6
	GWL_HWNDPARENT  = -8
	GWLP_HWNDPARENT = -8
	GWL_ID          = -12
	GWLP_ID         = -12
	GWL_USERDATA    = -21
	GWLP_USERDATA   = -21
)

// Extended window style constants
const (
	WS_EX_DLGMODALFRAME    = 0x00000001
	WS_EX_NOPARENTNOTIFY   = 0x00000004
	WS_EX_TOPMOST          = 0x00000008
	WS_EX_ACCEPTFILES      = 0x00000010
	WS_EX_TRANSPARENT      = 0x00000020
	WS_EX_MDICHILD         = 0x00000040
	WS_EX_TOOLWINDOW       = 0x00000080
	WS_EX_WINDOWEDGE       = 0x00000100
	WS_EX_CLIENTEDGE       = 0x00000200
	WS_EX_CONTEXTHELP      = 0x00000400
	WS_EX_RIGHT            = 0x00001000
	WS_EX_LEFT             = 0x00000000
	WS_EX_RTLREADING       = 0x00002000
	WS_EX_LTRREADING       = 0x00000000
	WS_EX_LEFTSCROLLBAR    = 0x00004000
	WS_EX_RIGHTSCROLLBAR   = 0x00000000
	WS_EX_CONTROLPARENT    = 0x00010000
	WS_EX_STATICEDGE       = 0x00020000
	WS_EX_APPWINDOW        = 0x00040000
	WS_EX_OVERLAPPEDWINDOW = 0x00000100 | 0x00000200
	WS_EX_PALETTEWINDOW    = 0x00000100 | 0x00000080 | 0x00000008
	WS_EX_LAYERED          = 0x00080000
	WS_EX_NOINHERITLAYOUT  = 0x00100000
	WS_EX_LAYOUTRTL        = 0x00400000
	WS_EX_NOACTIVATE       = 0x08000000
)

func GetWindowThreadProcessId(hwnd HWND) (HANDLE, uint32) {
	var processId uint32
	ret, _, _ := procGetWindowThreadProcessId.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&processId)))

	return HANDLE(ret), processId
}

func IsErrSuccess(err error) bool {
	if e, ok := err.(syscall.Errno); ok {
		if e == 0 {
			return true
		}
	}

	return false
}

func EnumWindows(lpEnumFunc WNDENUMPROC, lParam LPARAM) bool {
	if _, _, err := procEnumWindows.Call(uintptr(lpEnumFunc), uintptr(lParam)); IsErrSuccess(err) {
		return true
	}

	return false
}

func GetWindowLong(hwnd HWND, index int) uint32 {
	ret, _, _ := procGetWindowLong.Call(
		uintptr(hwnd),
		uintptr(index))

	return uint32(ret)
}

func SetWindowLong(hwnd HWND, index int, value uint32) uint32 {
	ret, _, _ := procSetWindowLong.Call(
		uintptr(hwnd),
		uintptr(index),
		uintptr(value))

	return uint32(ret)
}

func GetCurrentProcessId() uint32 {
	r, _, err := procGetCurrentProcessId.Call()
	if !IsErrSuccess(err) {
		return 0
	}

	return uint32(r)
}

func GetCurrentProcessWindows() []HWND {
	var r []HWND

	if id := GetCurrentProcessId(); id != 0 {
		EnumWindows(WNDENUMPROC(syscall.NewCallback(func(hwnd HWND, lParam LPARAM) uintptr {
			if _, i := GetWindowThreadProcessId(hwnd); i == id {
				r = append(r, hwnd)
			}

			return 1
		})), 0)
	}

	return r
}

func EnableCurrentProcessWindowClickThrough() {
	for _, hWnd := range GetCurrentProcessWindows() {
		EnableWindowClickThrough(hWnd)
	}
}

func EnableWindowClickThrough(hwnd HWND) {
	exStyle := GetWindowLong(hwnd, GWL_EXSTYLE)
	SetWindowLong(hwnd, GWL_EXSTYLE, exStyle|WS_EX_LAYERED|WS_EX_TRANSPARENT)
}

func DisableWindowClickThrough(hwnd HWND) {
	exStyle := GetWindowLong(hwnd, GWL_EXSTYLE)
	SetWindowLong(hwnd, GWL_EXSTYLE, exStyle&^WS_EX_LAYERED&^WS_EX_TRANSPARENT)
}
