package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"golang.org/x/text/encoding/charmap"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func usage() {
	fmt.Println("Usage: agnog cvf_file [output]")
	flag.PrintDefaults()
}

const chunkSize = 403

func main() {
	// Set up command line arguments
	verbosePtr := flag.Bool("v", false, "Verbose")
	flag.Usage = usage

	flag.Parse()

	if flag.NArg() != 1 && flag.NArg() != 2 {
		usage()
		os.Exit(2)
	}

	// Open input file
	in, err := os.Open(flag.Args()[0])
	check(err)
	defer in.Close()

	// Set up output
	out := os.Stdout

	if flag.NArg() > 1 {
		out, err = os.Create(flag.Args()[1])
		check(err)
		defer out.Close()
	}

	// Process input
	if *verbosePtr {
		fmt.Fprintln(os.Stderr, "Reading file...")
	}

	var buf [chunkSize]byte

	for i := 1; ; i++ {
		if *verbosePtr {
			fmt.Fprintf(os.Stderr, "\rProcessing chunk %d...", i)
		}

		// Read chunk
		bytesRead, err := in.Read(buf[:])
		check(err)
		if bytesRead < chunkSize {
			fmt.Fprintln(os.Stderr)
			if bytesRead != 0 && *verbosePtr {
				fmt.Fprintf(os.Stderr, "Expected %d bytes, instead got %d\n", chunkSize, bytesRead)
			}
			break
		}

		// Strip tail null bytes and decode hebrew charset
		strEnd := bytes.IndexByte(buf[:bytesRead], 0)
		decoder := charmap.ISO8859_8.NewDecoder()
		str, err := decoder.String(string(buf[:strEnd]))
		check(err)

		// Write line to output
		fmt.Fprintln(out, str)

		if flag.NArg() > 1 {
			out.Sync()
		}
	}

	if *verbosePtr {
		fmt.Fprintln(os.Stderr, "Done")
	}
}
