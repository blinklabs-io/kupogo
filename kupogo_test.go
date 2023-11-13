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
