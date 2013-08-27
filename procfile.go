package main

import (
	"bufio"
	"fmt"
	_ "github.com/kr/pretty"
	"io"
	"os"
	"regexp"
)

var procfileEntryRegexp = regexp.MustCompile("^([A-Za-z0-9_]+):\\s*(.+)$")

type ProcfileEntry struct {
	Name    string
	Command string
}

type Procfile struct {
	Entries []ProcfileEntry
}

func ReadProcfile(filename string) (*Procfile, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return parseProcfile(fd)
}

func (pf *Procfile) HasProcess(name string) (exists bool) {
	for _, entry := range pf.Entries {
		if name == entry.Name {
			return true
		}
	}
	return false
}

func (pf *Procfile) LongestProcessName() (longest int) {
	longest = 6 // length of forego
	for _, entry := range pf.Entries {
		if len(entry.Name) > longest {
			longest = len(entry.Name)
		}
	}
	return
}

func parseProcfile(r io.Reader) (*Procfile, error) {
	pf := new(Procfile)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := procfileEntryRegexp.FindStringSubmatch(scanner.Text())
		pf.Entries = append(pf.Entries, ProcfileEntry{parts[1], parts[2]})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Reading Procfile: %s", err)
	}
	return pf, nil
}
