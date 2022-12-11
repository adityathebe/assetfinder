package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adityathebe/assetfinder"
)

func main() {
	var domains io.Reader
	domains = os.Stdin

	// If a single domain is supplied
	if len(os.Args) == 2 {
		domain := os.Args[1]
		if domain != "" {
			domains = strings.NewReader(domain)
		}
	}

	sc := bufio.NewScanner(domains)

	for sc.Scan() {
		outChan := assetfinder.ScanMain(sc.Text())
		for res := range outChan {
			fmt.Println(res.String())
		}
	}
}
