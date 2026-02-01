package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"GoMFT/internal/engine"
)

func main() {
	mftPtr := flag.String("mft", "", "Path to $MFT file")
	journalPtr := flag.String("journal", "", "Path to $UsnJrnl:$J file")
	outPtr := flag.String("out", "output.jsonl", "Path to output JSONL file")

	flag.Parse()

	if *mftPtr == "" || *journalPtr == "" {
		fmt.Println("Usage: gomft -mft <path> -journal <path> -output <path>")
		os.Exit(1)
	}

	fmt.Println("[*] Starting GoMFT Engine...")

	orch := engine.NewOrchestrator(*mftPtr, *journalPtr, *outPtr)

	if err := orch.BuildPathMap(); err != nil {
		log.Fatalf("[!] Failed to build path map: %v", err)
	}

	if err := orch.Run(); err != nil {
		log.Fatalf("[!] Engine runtime error: %v", err)
	}

	fmt.Println("[*] Analysis Complete.")
}
