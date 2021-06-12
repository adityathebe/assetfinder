package assetfinder

import (
	"strings"
	"sync"

	"github.com/adityathebe/assetfinder/sources"
)

type fetchFn func(string) ([]string, error)

func Scan(domain string) []string {
	sources := []fetchFn{
		sources.FetchCertSpotter,
		sources.FetchHackerTarget,
		sources.FetchThreatCrowd,
		sources.FetchCrtSh,
		sources.FetchFacebook,
		sources.FetchWayback,
		sources.FetchVirusTotal,
		sources.FetchFindSubDomains,
		sources.FetchUrlscan,
		sources.FetchBufferOverrun,
	}

	out := make(chan string)
	var wg sync.WaitGroup

	domain = strings.ToLower(domain)

	// call each of the source workers in a goroutine
	for _, source := range sources {
		wg.Add(1)
		fn := source

		go func() {
			defer wg.Done()
			// funcName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()

			names, err := fn(domain)
			if err != nil {
				// log.Println(err, funcName)
				return
			}

			// log.Printf("[%s]:: Results: %d\n", funcName, len(names))

			for _, n := range names {
				n = cleanDomain(n)
				if !strings.HasSuffix(n, domain) {
					continue
				}
				out <- n
			}
		}()
	}

	// close the output channel when all the workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	// track what we've already printed to avoid duplicates
	printed := make(map[string]bool)

	for n := range out {
		if _, ok := printed[n]; ok {
			continue
		}
		printed[n] = true
	}

	subdomains := make([]string, 0, len(printed))
	for k := range printed {
		subdomains = append(subdomains, k)
	}
	return subdomains
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
