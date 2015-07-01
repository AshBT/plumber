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
package main // import "github.com/qadium/plumber/manager"

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"fmt"
	"strings"
	"time"
)

func forwardData(dest string, body io.ReadCloser) (io.ReadCloser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", dest, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req) // this will close the body, too
	if err != nil {
		// close the body before returning an error
		defer resp.Body.Close()
		return nil, err
	}
	if resp.StatusCode != 200 {
		msg, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%v]: Status '%v' | Error '%v'", dest, resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	return resp.Body, nil
}

func createHandler(args []string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		start := time.Now()
		log.Printf("Received connection request from '%v'.", r.RemoteAddr)
		defer log.Printf("Completed request in %v", time.Since(start))

		if r.Body == nil || r.Method != "POST" {
			http.NotFound(w, r)
		} else {
			body := r.Body
			defer body.Close()

			for _, host := range args {
				body, err = forwardData(host, body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

			final, err := ioutil.ReadAll(body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(final)
		}
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	listener, err := net.Listen("tcp", ":9800")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	go func() {
		<-c
		log.Printf("Received termination; qutting.")
		// by setting the exit status to 0, we don't cause any parent
		// processes to think this was an unexpected termination
		listener.Close()
	}()

	// sanitize args to make sure they contain valid urls
	args := []string{}
	for _, arg := range os.Args[1:] {
		parsedUrl, err := url.Parse(arg)
		if err == nil && parsedUrl.Scheme == "http" {
			args = append(args, arg)
		} else {
			log.Printf("'%v' is not a valid url, discarding", arg)
		}
	}
	log.Printf("Forwarding JSONs to '%v'.", args)

	http.HandleFunc("/", createHandler(args))
	http.Serve(listener, nil)
}
