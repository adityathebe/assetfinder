package assetfinder

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/adityathebe/assetfinder/sources"
)

type fetchFn func(string) ([]string, error)

type Source struct {
	name string
	fn   fetchFn
}

type Result struct {
	DNS string `json:"dns"`
	Src string `json:"src"`
}

func (t *Result) String() string {
	x, _ := json.Marshal(t)
	return string(x)
}

func ScanMain(domain string) <-chan Result {
	var outChan = make(chan Result)
	go scan(domain, outChan)
	return outChan
}

func scan(domain string, outChan chan<- Result) {
	defer close(outChan)

	sources := []Source{
		{"CertSpotter", sources.FetchCertSpotter},
		{"HackerTarget", sources.FetchHackerTarget},
		{"ThreatCrowd", sources.FetchThreatCrowd},
		{"CrtSh", sources.FetchCrtSh},
		{"Facebook", sources.FetchFacebook},
		{"Wayback", sources.FetchWayback},
		{"VirusTotal", sources.FetchVirusTotal},
		{"FindSubDomains", sources.FetchFindSubDomains},
		{"Urlscan", sources.FetchUrlscan},
		{"BufferOverrun", sources.FetchBufferOverrun},
	}

	out := make(chan Result)
	var wg sync.WaitGroup

	domain = strings.ToLower(domain)

	// call each of the source workers in a goroutine
	for _, source := range sources {
		wg.Add(1)

		go func(source Source) {
			defer wg.Done()

			names, err := source.fn(domain)
			if err != nil {
				return
			}

			for _, n := range names {
				n = cleanDomain(n)
				if !strings.HasSuffix(n, domain) {
					continue
				}
				out <- Result{
					DNS: n,
					Src: source.name,
				}
			}
		}(source)
	}

	// close the output channel when all the workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	var printed = make(map[string]struct{})
	for n := range out {
		if _, ok := printed[n.DNS]; ok {
			continue
		}

		printed[n.DNS] = struct{}{}
		outChan <- n
	}
}

func cleanDomain(d string) string {
	d = strings.ToLower(d)

	// no idea what this is, but we can't clean it ¯\_(ツ)_/¯
	if len(d) < 2 {
		return d
	}

	if d[0] == '*' || d[0] == '%' {
		d = d[1:]
	}

	if d[0] == '.' {
		d = d[1:]
	}

	return d
}
