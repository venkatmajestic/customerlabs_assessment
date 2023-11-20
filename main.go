package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

var process_chan = make(chan map[string]string)

func request_handler(w http.ResponseWriter, r *http.Request) {

	var val map[string]string
	err := json.NewDecoder(r.Body).Decode(&val)
	fmt.Printf("%+v\n", val)
	process_chan <- val
	if err != nil {
		fmt.Printf("Error Occured %v", err)
	}
	w.WriteHeader(202)

}

func main() {

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT)

	r := mux.NewRouter()
	r.HandleFunc("/", request_handler).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    ":8090",
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("Error Occured during Listen and Serve err: %v", err)
		}
	}()

	go func() {
		for {
			res, ok := <-process_chan
			if !ok {
				fmt.Println("Channel Closed", ok)
				break
			}
			fmt.Println("Channel Open", ok)
			go worker(res)
		}
	}()

	<-done
	fmt.Println("Received Interrupt Signal")
	srv.Close()
	close(process_chan)
	fmt.Println("END")

}
