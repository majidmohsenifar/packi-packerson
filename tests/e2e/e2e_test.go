package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"testing"
	"time"
)

const (
	serverURL    = "http://127.0.0.1:8080"
	startTimeout = 15 * time.Second
)

func startServer(t *testing.T) (context.CancelFunc, *exec.Cmd) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "go", "run", "./cmd/main.go")
	cmd.Env = append(os.Environ(), "CONFIG_PATH=./config/config.toml", "GIN_MODE=release")
	cmd.Dir = "../../"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	t.Cleanup(func() {
		cancel()
		if cmd.Process != nil {
			if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
				t.Logf("failed to kill server process: %v", err)
			}
		}
	})
	return cancel, cmd
}

func waitForServer(url string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		client := http.Client{Timeout: 1 * time.Second}
		resp, err := client.Get(url + "?x=1")
		if err == nil && resp.StatusCode == 200 {
			_ = resp.Body.Close()
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

func normalizeMapKeysToInt(m map[string]interface{}) map[int]int {
	res := make(map[int]int)
	for k, v := range m {
		ki, _ := strconv.Atoi(k)
		vi := int(v.(float64))
		res[ki] = vi
	}
	return res
}

func normalizeListOfMaps(arr []interface{}) []map[int]int {
	res := make([]map[int]int, len(arr))
	for i, item := range arr {
		m := item.(map[string]interface{})
		res[i] = normalizeMapKeysToInt(m)
	}
	return res
}

func Test(t *testing.T) {
	cancel, _ := startServer(t)
	defer cancel()

	if !waitForServer(serverURL+"/correct", startTimeout) {
		t.Fatalf("server did not start in %v", startTimeout)
	}

	correctCases := []struct {
		ordered int
		want    map[int]int
	}{
		{1, map[int]int{250: 1}},
		{250, map[int]int{250: 1}},
		{500, map[int]int{500: 1}},
		{750, map[int]int{500: 1, 250: 1}},
		{1000, map[int]int{1000: 1}},
		{12001, map[int]int{5000: 2, 2000: 1, 250: 1}},
		{4999, map[int]int{2000: 2, 1000: 1}},
		{5001, map[int]int{5000: 1, 250: 1}},
	}

	for _, tc := range correctCases {
		t.Run(fmt.Sprintf("correct_%d", tc.ordered), func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/correct?x=%d", serverURL, tc.ordered))
			if err != nil {
				t.Fatalf("request error: %v", err)
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Logf("failed to close response body: %v", err)
				}
			}()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status %d", resp.StatusCode)
			}
			var body struct {
				Data struct {
					Packs map[string]interface{} `json:"packs"`
				} `json:"data"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("decode error: %v", err)
			}
			got := normalizeMapKeysToInt(body.Data.Packs)
			if fmt.Sprint(got) != fmt.Sprint(tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}

	incorrectCases := []struct {
		ordered  int
		wantList []map[int]int
	}{
		{1, []map[int]int{{5000: 1}, {2000: 1}, {1000: 1}, {500: 1}}},
		{251, []map[int]int{{5000: 1}, {2000: 1}, {1000: 1}, {250: 2}}},
		{501, []map[int]int{{5000: 1}, {2000: 1}, {1000: 1}, {500: 2}, {250: 3}}},
		{12001, []map[int]int{{5000: 3}, {2000: 7}, {1000: 13}, {500: 25}, {250: 49}}},
	}

	for _, tc := range incorrectCases {
		t.Run(fmt.Sprintf("incorrect_%d", tc.ordered), func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/incorrect?x=%d", serverURL, tc.ordered))
			if err != nil {
				t.Fatalf("request error: %v", err)
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Logf("failed to close response body: %v", err)
				}
			}()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status %d", resp.StatusCode)
			}
			var body struct {
				Data struct {
					Packs []interface{} `json:"packs"`
				} `json:"data"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("decode error: %v", err)
			}
			got := normalizeListOfMaps(body.Data.Packs)
			if fmt.Sprint(got) != fmt.Sprint(tc.wantList) {
				t.Errorf("got %v, want %v", got, tc.wantList)
			}
		})
	}

	// Validation errors
	t.Run("validation_missing_x", func(t *testing.T) {
		resp, _ := http.Get(serverURL + "/correct")
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})
	t.Run("validation_non_integer", func(t *testing.T) {
		resp, _ := http.Get(serverURL + "/correct?x=abc")
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})
	t.Run("validation_negative", func(t *testing.T) {
		resp, _ := http.Get(serverURL + "/correct?x=-1")
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})
}
