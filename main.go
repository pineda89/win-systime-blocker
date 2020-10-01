package main

import (
	"fmt"
	"github.com/gonutz/w32"
	"golang.org/x/sys/windows"
	"os"
	"time"
	"unsafe"
)

const PROCESS_ALL_ACCESS = 0x1F0FFF

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procGetSystemTimeAsFileTime            = modkernel32.NewProc("GetSystemTimeAsFileTime")
	procOpenProcess = modkernel32.NewProc("OpenProcess")
	procWriteProcessMemory         = modkernel32.NewProc("WriteProcessMemory")
)

func main() {
	res := w32.GetSystemTimeAsFileTime()
	fmt.Println("current time:", res.DwLowDateTime, res.DwHighDateTime, res.ToUint64())

	timeObj := windateToTime(int64(res.ToUint64()))
	fmt.Println("current time parsed", timeObj)

	pid := os.Getpid()
	handler := OpenProcess(pid)
	ptr := []byte{0x81 , 0x01 , 0x1F , 0x40 , 0xEF , 0x45 , 0x81 , 0x41 , 0x04 , 0xF7 , 0x97 , 0xD6 , 0x02 , 0xC3 , 0xCC, 0xCC}
	WRITE(handler, procGetSystemTimeAsFileTime.Addr(), uintptr(unsafe.Pointer(&ptr[0])), uintptr(len(ptr)))

	res = w32.GetSystemTimeAsFileTime()
	fmt.Println("modified time:", res.DwLowDateTime, res.DwHighDateTime, res.ToUint64())

	timeObj = windateToTime(int64(res.ToUint64()))
	fmt.Println("modified time parsed", timeObj)

	// procGetSystemTimeAsFileTime.Addr() -> address of function GetSystemTimeAsFileTime in our process.
	// Is the same address in target process. kernel32.dll is shared :)
	// in this example, we are forcing DwLowDateTime=0x45EF401F and DwHighDateTime=0x02D697F7 . Is 2249-02-03 13:22:17.7959967 +0000 UTC
	// for modify date, simply modify the byte array
}

func windateToTime(input int64) time.Time {
	t := time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)
	d := time.Duration(input)
	for i := 0; i < 100; i++ {
		t = t.Add(d)
	}
	return t
}

func OpenProcess(pid int) uintptr {
	handle, _, _ := procOpenProcess.Call(uintptr(PROCESS_ALL_ACCESS), uintptr(1), uintptr(pid))
	return handle
}

func WRITE(hProcess uintptr, lpBaseAddress, lpBuffer, nSize uintptr) (int, bool) {
	var nBytesWritten int
	ret, _, _ := procWriteProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		lpBuffer,
		nSize,
		uintptr(unsafe.Pointer(&nBytesWritten)),
	)

	return nBytesWritten, ret != 0
}
