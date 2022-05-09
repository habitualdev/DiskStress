//go:generate goversioninfo -icon=resources/icon.ico -manifest=resources/diskstress.exe.manifest

package main

import (
	"DiskStress/utils"
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var progressBuffer = make(chan int, 2)

var progressLoop = true

type LogProgressWriter struct{}

func (pw LogProgressWriter) Write(data []byte) (int, error) {
	progressBuffer <- len(data)
	return len(data), nil
}

func main() {

	nCpu := runtime.NumCPU() / 2
	if _, err := os.Stat("FireStarter"); err == nil {
		err := os.RemoveAll("FireStarter")
		if err != nil {
			println(err.Error())
			return
		}
	}

	quit := make(chan bool, nCpu * 2)

	err := os.Mkdir("FireStarter", 0777)
	if err != nil {
		println(err.Error())
		return
	}

	for i := 0; i < nCpu; i++ {
		if _, err := os.Stat("FireStarter/" + strconv.Itoa(i)); errors.Is(err, os.ErrNotExist) {
			os.Mkdir("FireStarter/"+strconv.Itoa(i), 0777)
		}
	}

	in, err := os.OpenFile("FireStarter/master.rand", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("DirectIO error:", err)
		os.Exit(1)
	}

	free := utils.GetDiskFreeSpace() / 10

	tee := io.TeeReader(rand.Reader, LogProgressWriter{})

	bar := pb.New64(free)
	bar.Set(pb.Bytes, true)
	bar.Set(pb.SIBytesPrefix, true)
	bar.Start()

	go func() {
		for progressLoop {
			tempProgress := <-progressBuffer
			bar.Add(tempProgress)
		}
		bar.Finish()
	}()
	color.Blue("Creating initial seed file...\n")
	color.Green("Seed file is 10% of the disk space\n\n")
	io.CopyN(in, tee, free)
	progressLoop = false
	close(progressBuffer)
	time.Sleep(100 * time.Millisecond)

	color.Green("Seed file created\n\n")
	color.Red("Begin stress test?")
	utils.Confirm("", 3)

	in.Close()

	seed, _ := os.OpenFile("FireStarter/master.rand", os.O_RDONLY, 0777)



	for i := 0; i < nCpu; i++ {
		go Churn(i, seed, int(free), nCpu, quit)
	}

	color.Red("Press enter to exit")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	color.Blue("Stopping stress test...")

	for i := 0; i < nCpu; i++ {
		quit <- true
	}

	close(quit)

	seed.Close()

	color.Blue("Cleaning up...")

	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			os.Remove(path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	os.Remove("FireStarter")

	color.Blue("Done")

}

func Churn(dir int, file *os.File, freespace int, workers int, quit chan bool) {
	for {

		select {
		case quit <- true:
			return
		default:
			dirString := "FireStarter/" + strconv.Itoa(dir)
			child, _ := os.OpenFile(dirString+"/file.rand", os.O_CREATE|os.O_RDWR, 0777)
			io.CopyN(child, file, int64(freespace/workers))
			child.Close()
			os.Remove(dirString + "/file.rand")
		}
	}

}
