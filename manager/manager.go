package main // import "github.com/qadium/plumber/manager"

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
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
		return nil, err
	}
	return resp.Body, nil
}

func createHandler(args []string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		if r.Body == nil || r.Method != "POST" {
			http.NotFound(w, r)
		} else {
			body := r.Body
			defer body.Close()

			for _, host := range args {
				body, err = forwardData(host, body)
				if err != nil {
					panic(err)
				}
			}

			final, err := ioutil.ReadAll(io.LimitReader(body, 1048576))
			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(final)
		}
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		log.Printf("Received termination; qutting.")
		// by setting the exit status to 0, we don't cause any parent
		// processes to think this was an unexpected termination
		os.Exit(0)
	}()

	http.HandleFunc("/", createHandler(os.Args[1:]))
	log.Fatal(http.ListenAndServe(":9800", nil))
}
