package main

import (
	"debug/elf"
	"errors"
	"golang.org/x/arch/x86/x86asm"
)

var (
	ErrorGoBINNotFound  = errors.New("GO application not found")
	ErrorSymbolNotFound = errors.New("symbol not found")
	ErrorNoRetFound     = errors.New("no RET instructions found")
)

func findRetOffsets(elfPath, symbolName string) (offsets []int, err error) {
	//var err error
	var goSymbs []elf.Symbol
	var goElf *elf.File
	goElf, err = elf.Open(elfPath)
	if err != nil {
		return
	}
	goSymbs, err = goElf.Symbols()
	if err != nil {
		return nil, err
	}

	var found bool
	var symbol elf.Symbol
	for _, s := range goSymbs {
		if s.Name == symbolName {
			symbol = s
			found = true
			break
		}
	}

	if !found {
		return nil, ErrorSymbolNotFound
	}

	section := goElf.Sections[symbol.Section]

	var elfText []byte
	elfText, err = section.Data()
	if err != nil {
		return nil, err
	}

	start := symbol.Value - section.Addr
	end := start + symbol.Size

	var instHex []byte
	instHex = elfText[start:end]
	offsets, err = decodeInstruction(instHex)
	if len(offsets) == 0 {
		return offsets, ErrorNoRetFound
	}
	return offsets, nil
}

// decodeInstruction Decode into assembly instructions and identify the RET instruction to return the offset.
func decodeInstruction(instHex []byte) ([]int, error) {
	var offsets []int
	for i := 0; i < len(instHex); {
		inst, err := x86asm.Decode(instHex[i:], 64)
		//fmt.Println(inst.String())
		if err != nil {
			return nil, err
		}
		if inst.Op == x86asm.RET {
			offsets = append(offsets, i)
		}
		i += inst.Len
	}
	return offsets, nil
}
