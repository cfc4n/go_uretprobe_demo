package main

import (
	"bytes"
	"debug/elf"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"go_uretprobe_demo/highlight"
	"golang.org/x/arch/x86/x86asm"
	"os"
	"strings"
)

var (
	// ErrorSymbolNotFound
	//ErrorGoBINNotFound  = errors.New("GO application not found")
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
	var s *bytes.Buffer
	s = bytes.NewBufferString("")
	for i := 0; i < len(instHex); {
		inst, err := x86asm.Decode(instHex[i:], 64)
		//fmt.Printf("%04X\t%s\n", i, inst.String())
		//s.WriteString(inst.String())
		s.WriteString(fmt.Sprintf("%04X\t%s", i, inst.String()))
		s.WriteString("\n")
		if err != nil {
			return nil, err
		}
		if inst.Op == x86asm.RET {
			offsets = append(offsets, i)
		}
		i += inst.Len
	}

	asmCode = s.String()
	return offsets, nil
}

func asmCodeDisplay() error {
	fmt.Println("assembly code of the hooked function:")
	//fmt.Println(asmCode)
	//return
	syntaxFile, _ := os.ReadFile("highlight/asm.yaml")
	var header *highlight.Header
	var err error
	header, err = highlight.MakeHeaderYaml(syntaxFile)
	if err != nil {
		return err
	}
	file, err := highlight.ParseFile(syntaxFile)
	if err != nil {
		return err
	}
	var syndef *highlight.Def
	syndef, err = highlight.ParseDef(file, header)

	h := highlight.NewHighlighter(syndef)
	matches := h.HighlightString(asmCode)
	lines := strings.Split(asmCode, "\n")
	for lineN, l := range lines {
		for colN, c := range l {
			// Check if the group changed at the current position
			if group, ok := matches[lineN][colN]; ok {
				// Check the group name and set the color accordingly (the colors chosen are arbitrary)
				if group == highlight.Groups["statement"] {
					color.Set(color.FgGreen)
				} else if group == highlight.Groups["preproc"] {
					color.Set(color.FgHiRed)
				} else if group == highlight.Groups["special"] {
					color.Set(color.FgBlue)
				} else if group == highlight.Groups["constant.string"] {
					color.Set(color.FgCyan)
				} else if group == highlight.Groups["constant.specialChar"] {
					color.Set(color.FgHiMagenta)
				} else if group == highlight.Groups["type"] {
					color.Set(color.FgYellow)
				} else if group == highlight.Groups["constant.number"] {
					color.Set(color.FgCyan)
				} else if group == highlight.Groups["comment"] {
					color.Set(color.FgHiGreen)
				} else {
					color.Unset()
				}
			}
			// Print the character
			fmt.Print(string(c))
		}
		// This is at a newline, but highlighting might have been turned off at the very end of the line so we should check that.
		if group, ok := matches[lineN][len(l)]; ok {
			if group == highlight.Groups["default"] || group == highlight.Groups[""] {
				color.Unset()
			}
		}
		fmt.Print("\n")
	}
	color.Unset()
	return nil
}
