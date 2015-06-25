/**
 * Copyright 2015 Qadium, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"
)

// Test that the handler with no args just returns the data sent
func TestHandlerNoArgs(t *testing.T) {
	handler := createHandler(nil)

	req, err := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString("{'foo': 3}"))
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != 200 || w.Body.String() != "{'foo': 3}" {
		t.Error("Did not get an idempotent response")
	}
}

// Test that the handler with no args returns 404 if we try a GET or
// POST with no data
func TestHandlerInvalidRequest(t *testing.T) {
	handler := createHandler(nil)

	req, err := http.NewRequest("GET", "http://foobar.com", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != 404 {
		t.Error("Expected an error code of 404 from a GET request.")
	}

	req, err = http.NewRequest("POST", "http://foobar.com", nil)
	if err != nil {
		t.Error(err)
	}
	handler(w, req)

	if w.Code != 404 {
		t.Error("Expected an error code of 404 from an empty POST request.")
	}
}

func makeTestHandler(response string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		w.Write([]byte(response + buf.String()))
	}
}

// Test that the handler forwards data to servers *in order*
func TestHandlerForwardsData(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(makeTestHandler("foo")))
	defer ts1.Close()

	ts2 := httptest.NewServer(http.HandlerFunc(makeTestHandler("first")))
	defer ts2.Close()

	handler := createHandler([]string{ts1.URL, ts2.URL})
	req, err := http.NewRequest("POST", "http://foobar.com", bytes.NewBufferString("{'foo': 3}"))
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != 200 || w.Body.String() != "firstfoo{'foo': 3}" {
		t.Errorf("Got '%s'; did not get expected response", w.Body.String())
	}
}

func TestMainRunnerExitsGracefully(t *testing.T) {
	// set the interrupt handler to go off after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	go func() {
		time.Sleep(50 * time.Millisecond)
		resp, err := http.Post("http://localhost:9800", "application/json", bytes.NewBufferString("{'foo': 3}"))
		if err != nil {
			t.Error(err)
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		result := buf.String()
		if result != "{'foo': 3}" {
			t.Errorf("Got '%s'; did not get expected response", result)
		}
	}()
	main()
}

// Benchmark to see how many requests the "empty" handler can handle
// Obviously, this is not accurate in production since we forward to
// multiple handlers.
func BenchmarkHandler(b *testing.B) {
	handler := createHandler(nil)

	req, err := http.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString("{'foo': 3}"))
	if err != nil {
		b.Error(err)
	}

	w := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		handler(w, req)
	}
}
