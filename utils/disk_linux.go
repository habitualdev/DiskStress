package utils

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

func GetDiskFreeSpace() int64{


	var stat unix.Statfs_t

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = unix.Statfs(wd, &stat)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return stat.Bsize
}
