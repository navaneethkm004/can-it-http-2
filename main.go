package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func main() {
	index := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var sendError string
		var success bool
		var proto string

		err := r.ParseForm()
		if err != nil {
			sendError = err.Error()
		}

		url := r.Form.Get("url")

		if url == "" {
			url = "https://repl.it"
		}

		fmt.Println("checking http/2 for", url)

		client := &http.Client{
			Timeout: time.Second * 3,
		}

		res, err := client.Head(url)
		if err != nil {
			sendError = err.Error()
		} else {
			res.Body.Close()

			proto = res.Proto
			success = res.ProtoMajor >= 2
		}

		err = index.Execute(
			w,
			struct {
				Err         string
				Placeholder string
				Success     bool
				Proto       string
			}{
				sendError,
				url,
				success,
				proto,
			},
		)
		if err != nil {
			fmt.Println("template error", err)
		}
	})

	fmt.Println("listening at :8080")
	http.ListenAndServe(":8080", nil)
}
