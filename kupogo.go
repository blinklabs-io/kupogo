// Copyright 2023 Blink Labs, LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kupogo

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Matches []Match

type Match struct {
	TransactionIndex int       `json:"transaction_index"`
	TransactionID    string    `json:"transaction_id"`
	OutputIndex      int       `json:"output_index"`
	Address          string    `json:"address"`
	Value            Value     `json:"value"`
	DatumHash        *string   `json:"datum_hash"`
	DatumType        *string   `json:"datum_type"`
	ScriptHash       *string   `json:"script_hash"`
	CreatedAt        CreatedAt `json:"created_at"`
	SpentAt          *SpentAt  `json:"spent_at"`
}

type Assets map[string]int

type Value struct {
	Coins  int    `json:"coins"`
	Assets Assets `json:"assets"`
}

type CreatedAt struct {
	SlotNo     int    `json:"slot_no"`
	HeaderHash string `json:"header_hash"`
}

type SpentAt struct {
	SlotNo     int    `json:"slot_no"`
	HeaderHash string `json:"header_hash"`
}

type Client struct {
	KupoUrl string
}

type Metadata struct {
	Hash   string          `json:"hash"`
	Raw    []byte          `json:"-"`
	Schema json.RawMessage `json:"schema"`
}

func NewClient(url string) *Client {
	return &Client{KupoUrl: url}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed do: %s", err)
	}
	return resp, nil
}

func (c *Client) GetAllMatches() (*Matches, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/matches", c.KupoUrl),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed req: %s", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed getting all matches: %s",
				err,
			)
	}
	if resp.StatusCode != http.StatusOK {
		return nil,
			fmt.Errorf(
				"failed getting all matches: %d",
				resp.StatusCode,
			)
	}
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed getting body bytes: %s", err)
	}
	defer resp.Body.Close()
	matches := Matches{}
	err = json.Unmarshal(respBodyBytes, &matches)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal: %s", err)
	}
	return &matches, nil
}

func (c *Client) GetMatches(pattern string) (*Matches, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/matches/%s", c.KupoUrl, pattern),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed req: %s", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed getting all matches: %s",
				err,
			)
	}
	if resp.StatusCode != http.StatusOK {
		return nil,
			fmt.Errorf(
				"failed getting all matches: %d",
				resp.StatusCode,
			)
	}
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	matches := Matches{}
	err = json.Unmarshal(respBodyBytes, &matches)
	if err != nil {
		return nil, fmt.Errorf("fail unmarshal: %s", err)
	}
	return &matches, nil
}

func (c *Client) GetMetadata(slotNo int, transactionID string) (*Metadata, error) {
	url := fmt.Sprintf("%s/metadata/%d", c.KupoUrl, slotNo)
	// Add the transaction_id query parameter if provided
	if transactionID != "" {
		url += fmt.Sprintf("?transaction_id=%s", transactionID)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %s", err)
	}
	defer resp.Body.Close()

	// Check for 304 Not Modified
	if resp.StatusCode == http.StatusNotModified {
		return nil, fmt.Errorf("metadata not modified since last request")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get metadata: status code %d", resp.StatusCode)
	}

	var response struct {
		Hash   string          `json:"hash"`
		RawHex string          `json:"raw"`
		Schema json.RawMessage `json:"schema"`
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %s", err)
	}

	rawBytes, err := hex.DecodeString(response.RawHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode raw data: %s", err)
	}

	metadata := &Metadata{
		Hash:   response.Hash,
		Raw:    rawBytes,
		Schema: response.Schema,
	}

	return metadata, nil
}
