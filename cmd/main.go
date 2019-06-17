package main

import (
	"flag"
	"fmt"
	metadata "github.com/monimena/audio-metadata"
	"os"
)

func main() {

	flag.Parse()

	audioFile := flag.Arg(0)

	fmt.Printf("audio file: %s\n", audioFile)

	if len(audioFile) < 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	f, err := os.Open(audioFile)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error opening file %s: %v", audioFile, err)
		os.Exit(1)
	}

	info, err := metadata.Parse(f)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error parsing file %s: %v", audioFile, err)
		os.Exit(1)
	}

	fmt.Printf("%#v\n", info)

}
