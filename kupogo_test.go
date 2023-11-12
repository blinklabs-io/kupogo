package kupogo

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_GetPattern(t *testing.T) {
	t.Run("Successful request and unmarshaling of response", func(t *testing.T) {
		t.Parallel()

		// Create a mock server for the successful case
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/patterns/*" {
				patterns := []string{
					"addr_vk1x7da0l25j04my8sej5ntrgdn38wmshxhplxdfjskn07ufavsgtkqn5hljl/*",
					"*/script1cda3khwqv60360rp5m7akt50m6ttapacs8rqhn5w342z7r35m37",
					"*/dca1e44765b9f80c8b18105e17de90d4a07e4d5a83de533e53fee32e0502d17e/*",
					"*/4fc6bb0c93780ad706425d9f7dc1d3c5e3ddbf29ba8486dce904a5fc",
					"*/*",
				}
				respBody, _ := json.Marshal(patterns)
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write(respBody); err != nil {
					log.Printf("Error writing response: %v", err)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		client := &Client{KupoUrl: server.URL}

		patterns, err := client.GetPattern("*")
		expectedPatterns := []string{
			"addr_vk1x7da0l25j04my8sej5ntrgdn38wmshxhplxdfjskn07ufavsgtkqn5hljl/*",
			"*/script1cda3khwqv60360rp5m7akt50m6ttapacs8rqhn5w342z7r35m37",
			"*/dca1e44765b9f80c8b18105e17de90d4a07e4d5a83de533e53fee32e0502d17e/*",
			"*/4fc6bb0c93780ad706425d9f7dc1d3c5e3ddbf29ba8486dce904a5fc",
			"*/*",
		}
		if !reflect.DeepEqual(patterns, expectedPatterns) {
			t.Errorf("Expected patterns %v, got %v", expectedPatterns, patterns)
		}
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}
	})

	t.Run("Failed unmarshaling", func(t *testing.T) {
		t.Parallel()

		// Create a mock server for the failed unmarshaling case
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/patterns/*" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("invalid json"))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		client := &Client{KupoUrl: server.URL}

		_, err := client.GetPattern("*")
		expectedErrMsg := "failed to unmarshal pattern: invalid character 'i' looking for beginning of value"
		if err == nil {
			t.Error("Expected an error, got nil")
		} else if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
		}
	})
}

func TestClient_GetPatterns(t *testing.T) {
	t.Run("Successful request and unmarshaling of response", func(t *testing.T) {
		t.Parallel()

		// Create a mock server for the successful case
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/patterns" {
				patterns := []string{
					"*",
				}
				respBody, _ := json.Marshal(patterns)
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write(respBody); err != nil {
					log.Printf("Error writing response: %v", err)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		client := &Client{KupoUrl: server.URL}

		patterns, err := client.GetPatterns()
		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}

		log.Printf("Received patterns: %v", patterns)

		expectedPatterns := []string{"*"}
		if !reflect.DeepEqual(patterns, expectedPatterns) {
			t.Errorf("Expected patterns %v, got %v", expectedPatterns, patterns)
		}
	})

	t.Run("Failed unmarshaling", func(t *testing.T) {
		t.Parallel()

		// Create a mock server for the failed unmarshaling case
		invalidServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("invalid json"))
		}))
		defer invalidServer.Close()

		client := &Client{KupoUrl: invalidServer.URL}

		_, err := client.GetPatterns()
		expectedErrMsg := "failed to unmarshal patterns: invalid character 'i' looking for beginning of value"
		if err == nil {
			t.Error("Expected an error, got nil")
		} else if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
		}
	})
}

func TestClient_GetScriptByHash(t *testing.T) {
	t.Run("Successful request and unmarshaling of response", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/scripts/4fc6bb0c93780ad706425d9f7dc1d3c5e3ddbf29ba8486dce904a5fc" {
				response := ScriptResponse{
					Language: "plutus:v2",
					Script:   "8201838200581c3c07030e36bfffe67e2e2ec09e5293d384637cd2f004356ef320f3fe8204186482051896",
				}
				respBody, _ := json.Marshal(response)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(respBody)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		client := &Client{KupoUrl: server.URL}
		scriptResponse, err := client.GetScriptByHash("4fc6bb0c93780ad706425d9f7dc1d3c5e3ddbf29ba8486dce904a5fc")
		expectedResponse := &ScriptResponse{
			Language: "plutus:v2",
			Script:   "8201838200581c3c07030e36bfffe67e2e2ec09e5293d384637cd2f004356ef320f3fe8204186482051896",
		}
		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}
		if !reflect.DeepEqual(scriptResponse, expectedResponse) {
			t.Errorf("Expected response %v, got %v", expectedResponse, scriptResponse)
		}
	})

	t.Run("Successful request returning null", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("null"))
		}))
		defer server.Close()

		client := &Client{KupoUrl: server.URL}
		scriptResponse, err := client.GetScriptByHash("4fc6bb0c93780ad706425d9f7dc1d3c5e3ddbf29ba8486dce904a5fc")
		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}
		if scriptResponse != nil {
			t.Errorf("Expected null response, got %v", scriptResponse)
		}
	})

	t.Run("Failed unmarshaling missing key", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"language": "plutus:v2"}`))
		}))
		defer server.Close()

		client := &Client{KupoUrl: server.URL}
		_, err := client.GetScriptByHash("4fc6bb0c93780ad706425d9f7dc1d3c5e3ddbf29ba8486dce904a5fc")
		expectedErrMsg := "failed to validate script response: Key: 'ScriptResponse.Script' Error:Field validation for 'Script' failed on the 'required' tag"
		if err == nil {
			t.Fatalf("Expected no error, got %s", err)
		} else if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
		}
	})
}

func TestClient_GetDatumByHash(t *testing.T) {
	t.Run("Successful request and unmarshaling of response", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/datums/34215ad90b1ade84f5b4fe3c0a16cb3afeae468210535e0305efd93931f35059" {
				response := DatumResponse{
					Datum: "d87980",
				}
				respBody, _ := json.Marshal(response)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(respBody)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		client := &Client{KupoUrl: server.URL}
		datumResponse, err := client.GetDatumByHash("34215ad90b1ade84f5b4fe3c0a16cb3afeae468210535e0305efd93931f35059")
		expectedResponse := &DatumResponse{
			Datum: "d87980",
		}
		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}
		if !reflect.DeepEqual(datumResponse, expectedResponse) {
			t.Errorf("Expected response %v, got %v", expectedResponse, datumResponse)
		}
	})

	t.Run("Successful request returning null", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("null"))
		}))
		defer server.Close()

		client := &Client{KupoUrl: server.URL}
		datumResponse, err := client.GetDatumByHash("34215ad90b1ade84f5b4fe3c0a16cb3afeae468210535e0305efd93931f35059")
		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}
		if datumResponse != nil {
			t.Errorf("Expected null response, got %v", datumResponse)
		}
	})

	t.Run("Failed unmarshaling missing key", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		}))
		defer server.Close()

		client := &Client{KupoUrl: server.URL}
		_, err := client.GetDatumByHash("34215ad90b1ade84f5b4fe3c0a16cb3afeae468210535e0305efd93931f35059")
		expectedErrMsg := "failed to validate datum response: Key: 'DatumResponse.Datum' Error:Field validation for 'Datum' failed on the 'required' tag"
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		} else if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
		}
	})
}
