package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	keys := make([]string, 0, len(pages))
	for k := range pages {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var sorted []PageData

	for _, k := range keys {
		sorted = append(sorted, pages[k])
	}

	data, err := json.MarshalIndent(sorted, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0o644)
	if err != nil {
		return err
	}

	fmt.Println("json report created succesfully")
	return nil
}
