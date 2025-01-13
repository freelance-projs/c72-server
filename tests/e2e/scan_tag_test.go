package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
)

func getURL(path string) string {
	return "http://localhost:8080/api" + path
}

func bodyBytes(req any) io.Reader {
	data, _ := json.Marshal(req)
	return bytes.NewReader(data)
}

func TestScanTag(t *testing.T) {
	httpReq, _ := http.NewRequest(http.MethodPost, getURL("/tags"), bodyBytes(dto.ScanTagRequest{
		TagIDs: []string{},
	}))

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		t.Errorf("failed to create tags: %v", err)
	}
	if httpResp.StatusCode != http.StatusCreated {
		t.Errorf("failed to create tags, got status: %d", httpResp.StatusCode)
	}

	var resp dto.Response
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Errorf("failed to create tags, got response: %v", resp)
	}
}
