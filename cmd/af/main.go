package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adityathebe/assetfinder/assetfinder"
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
		subdomains := assetfinder.Scan(sc.Text())
		for _, sub := range subdomains {
			fmt.Println(sub)
		}
	}
}
