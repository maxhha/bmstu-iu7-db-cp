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

func RunTest() <-chan error {
	done := make(chan error, 1)

	go func() {
		fmt.Println("Start test")
		// TODO: run integration test script
		cmd := exec.Command("sleep", "3")

		outfile, err := os.Create("./integration-test.log")
		if err != nil {
			done <- fmt.Errorf("create log file: %w", err)
			return
		}

		defer outfile.Close()

		cmd.Stdout = outfile
		cmd.Stderr = outfile

		if err := cmd.Start(); err != nil {
			done <- fmt.Errorf("cmd start: %w", err)
			return
		}

		if err := cmd.Wait(); err != nil {
			done <- fmt.Errorf("cmd wait: %w", err)
			return
		}

		done <- nil
	}()

	return done
}

func TestIntegration(t *testing.T) {
	port, ok := os.LookupEnv("POSTGRES_CONNECTION")
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

	if err := <-RunTest(); err != nil {
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
