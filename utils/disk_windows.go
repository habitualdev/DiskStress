package utils
import (
	"fmt"
	"syscall"
	"unsafe"
)
func GetDiskFreeSpace() int64 {
	kernelDLL := syscall.MustLoadDLL("kernel32.dll")
	GetDiskFreeSpaceExW := kernelDLL.MustFindProc("GetDiskFreeSpaceExW")

	var free, total, avail int64

	path := "c:\\"
	GetDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&free)),
		uintptr(unsafe.Pointer(&total)),
		uintptr(unsafe.Pointer(&avail)),
	)


	fmt.Println("Free:", free/1024/1024, "GB | Total:", total/1024/1024, "GB | Available:", avail/1024/1024, "GB")
	return free
}
