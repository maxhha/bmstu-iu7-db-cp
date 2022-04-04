//go:build integration

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func RunTest() error {
	fmt.Println("Start test")
	cmd := exec.Command("./integration_test/run.sh")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd start: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("cmd wait: %w", err)
	}

	return nil
}

func TestIntegration(t *testing.T) {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: Init(),
	}
	done := false

	go func() {
		err := server.ListenAndServe()
		if !done {
			t.Error(err)
		}
	}()

	if err := RunTest(); err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(5)*time.Second,
	)
	defer cancel()

	done = true
	if err := server.Shutdown(ctx); err != nil {
		t.Error(err)
	}
}
