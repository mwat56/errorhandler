/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package errorhandler

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"net/http"
	"net/http/httptest"
)

type (
	// TErrorPager is an interface requiring a function to return
	// the error text of an HTTP error message page.
	TErrorPager interface {
		// GetErrorPage returns a HTML page for `aStatus`.
		// The return value is expected to be a valid HTML page.
		//
		//	`aData` is the original error text.
		//	`aStatus` is the number of the actual HTTP error status.
		GetErrorPage(aData []byte, aStatus int) []byte
	}

	// `tErrorWriter` embeds a `ResponseWriter` and includes
	// error page handling.
	tErrorWriter struct {
		http.ResponseWriter             // used to construct an HTTP response.
		errPager            TErrorPager // provider of error page contents
		status              int         // HTTP status code of current request
	}
)

// WriteHeader sends an HTTP response header with the provided
// status code.
//
//	`aStatus` the current request's status code.
func (ew *tErrorWriter) WriteHeader(aStatus int) {
	ew.status = aStatus
	if (200 != aStatus) && (nil != ew.errPager) {
		// The other error pages are send as "text/plain" by default.
		// Since we expect "text/html" from our `TErrorPager` we
		// have to make sure we send the right header line.
		ew.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	}

	defer func() {
		// Needed because `httptest.ResponseRecorder` doesn't
		// allow all response codes but panics in certain cases
		// which is very inconvenient during testing.
		if r := recover(); r != nil {
			if rr, ok := ew.ResponseWriter.(*httptest.ResponseRecorder); ok {
				rr.Code = aStatus
			}
			ew.status = aStatus
		}
	}()

	ew.ResponseWriter.WriteHeader(aStatus)
} // WriteHeader()

// Write writes the data to the connection as part of an HTTP reply.
//
//	`aData` is the data (usually text) to send to the remote client.
func (ew *tErrorWriter) Write(aData []byte) (int, error) {
	if 0 == ew.status {
		// for pages not initialised yet we assume success
		ew.status = 200
	}
	if (200 != ew.status) && (nil != ew.errPager) {
		if txt := ew.errPager.GetErrorPage(aData, ew.status); 0 < len(txt) {
			// replace the given text with our customised page
			aData = txt
		}
	}

	return ew.ResponseWriter.Write(aData)
} // Write()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// Wrap returns a handler function that includes error page handling,
// wrapping the given `aNext` and calling it internally.
//
//	`aNext` The HTTP handler responding to the actual web request.
//	`aPager` The provider of error message pages.
func Wrap(aNext http.Handler, aPager TErrorPager) http.Handler {
	return http.HandlerFunc(
		func(aWriter http.ResponseWriter, aRequest *http.Request) {
			ew := &tErrorWriter{
				aWriter,
				aPager,
				0,
			}
			aNext.ServeHTTP(ew, aRequest)
		})
} // Wrap()

/* _EoF_ */
