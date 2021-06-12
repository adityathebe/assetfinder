package assetfinder

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type crtShResult struct {
	Name string `json:"common_name"`
}

func fetchCrtSh(domain string) ([]string, error) {
	var results []crtShResult
	var url = fmt.Sprintf("https://crt.sh/?q=%%.%s&output=json", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	output := make([]string, len(results))
	for _, res := range results {
		output = append(output, res.Name)
	}

	return output, nil
}
