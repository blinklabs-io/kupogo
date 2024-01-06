# kupogo
Kupo client in Go

Note: you should use https://github.com/SundaeSwap-finance/kugo instead

## Usage
```golang
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blinklabs-io/kupogo"
)

func main() {
	client := kupogo.NewClient("https://kupo.example.com")
	matches, err := client.GetMatches("addr1...")
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	jsonData, err := json.MarshalIndent(matches, "", "    ")
	if err != nil {
		fmt.Println("ERROR: marshal:", err)
		os.Exit(1)
	}
	fmt.Printf("matches: %s\n", jsonData)
}
```
