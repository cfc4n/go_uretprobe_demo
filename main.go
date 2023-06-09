package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/gojue/ebpfmanager"
	_ "github.com/shuLhan/go-bindata"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	GO_APPLICATION_ELF_PATH = "/home/cfc4n/project/go_uretprobe_demo/bin/demo"
	COUNT_CC_SYMBOL         = "main.CountCC"
)

var asmCode string

func main() {
	var enable bool
	var goAppPath string
	flag.BoolVar(&enable, "e", false, "Use uprobe+offset address instead of uretprobe, default:disabled ")
	flag.StringVar(&goAppPath, "g", "/home/cfc4n/project/go_uretprobe_demo/bin/demo", "ELF file compiled by the Go programming language.")
	flag.Parse()

	fmt.Println("Github repo : https://github.com/cfc4n/go_uretprobe_demo")
	fmt.Printf("Use uprobe+offset address instead of uretprobe:%v\n", enable)
	fmt.Printf("traced ELF file:%s\n", goAppPath)
	fmt.Println("attach function:", COUNT_CC_SYMBOL)
	var sec = "uretprobe/countcc"
	var ebpfFunc = "uretprobe_countcc"
	var m = &manager.Manager{
		Probes: []*manager.Probe{
			{
				Section:          sec,
				EbpfFuncName:     ebpfFunc,
				AttachToFuncName: COUNT_CC_SYMBOL,
				BinaryPath:       goAppPath,
			},
		},
	}

	if enable {
		// 查找ELF文件中被HOOk函数的符号表中，RET指令的偏移量
		offsets, err := findRetOffsets(goAppPath, COUNT_CC_SYMBOL)
		if err != nil {
			log.Fatal(err)
		}
		// dwarf
		// for test
		//dwarfList(goAppPath, COUNT_CC_SYMBOL)
		//return
		//
		err = asmCodeDisplay()
		if err != nil {
			log.Fatal(err)
		}

		sec = "uprobe/countcc"
		ebpfFunc = "uprobe_countcc"
		m.Probes = m.Probes[:0] // 清空slice
		for _, offset := range offsets {
			m.Probes = append(m.Probes,
				&manager.Probe{
					Section:          sec,
					UprobeOffset:     uint64(offset),
					EbpfFuncName:     ebpfFunc,
					AttachToFuncName: COUNT_CC_SYMBOL,
					BinaryPath:       goAppPath,
					UID:              fmt.Sprintf("%s_%d", ebpfFunc, offset),
				})
			log.Printf("Golang uretprobe hook %s [RET] at 0x%X\n", COUNT_CC_SYMBOL, offset)
		}
	}

	// Initialize the manager
	buf, err := Asset("/probe.o")
	if err != nil {
		log.Fatal(errors.New(fmt.Sprintf("error:%v , couldn't find asset", err)))
	}

	if err = m.Init(bytes.NewReader(buf)); err != nil {
		log.Fatal(err)
	}

	// Start the manager
	if err = m.Start(); err != nil {
		log.Fatal(err)
	}

	log.Println("successfully started, head over to /sys/kernel/debug/tracing/trace_pipe")

	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)

	<-stopper

	// Close the manager
	if err = m.Stop(manager.CleanAll); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
