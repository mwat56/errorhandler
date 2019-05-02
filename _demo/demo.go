package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/mwat56/go-errorhandler"
)

type tErrPage bool

func (ep tErrPage) GetErrorPage(aData []byte, aStatus int) []byte {
	//
	// Here you can prepare the error page you want to return
	//

	return aData
} // GetErrorPage()

// `myHandler()` is a dummy for demonstration purposes.
func myHandler(aWriter http.ResponseWriter, aRequest *http.Request) {
	io.WriteString(aWriter, "Hello world!")
} // myHandler()

func main() {
	var ep tErrPage

	pageHandler := http.NewServeMux()
	pageHandler.HandleFunc("/", myHandler)

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: errorhandler.Wrap(pageHandler, ep),
	}

	if err := server.ListenAndServe(); nil != err {
		log.Fatalf("%s: %v", os.Args[0], err)
	}
} // main()
