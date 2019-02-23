package chrome

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const listeningURL string = "127.0.0.1"

type forwardingProxy struct {
	targetURL *url.URL
	server    *httputil.ReverseProxy
	listener  net.Listener
	port      int
}

func (proxy *forwardingProxy) start() error {
	// *Dont* verify remote certificates.
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Start the proxy and assign our custom Transport
	proxy.targetURL.Path = "/" // set the path to / as this becomes the base path
	proxy.server = httputil.NewSingleHostReverseProxy(proxy.targetURL)
	proxy.server.Transport = transport

	// Get an open port for this proxy instance to run on.
	var err error
	proxy.listener, err = net.Listen("tcp", listeningURL+":0")
	if err != nil {
		return err
	}

	// Set the port we used so that the caller of this method
	// can discover where to find this proxy instance.
	proxy.port = proxy.listener.Addr().(*net.TCPAddr).Port

	// Finally, the goroutine for the proxy service.
	go func() {
		// Create an isolated ServeMux
		//  ref: https://golang.org/pkg/net/http/#ServeMux
		httpServer := http.NewServeMux()
		httpServer.HandleFunc("/", proxy.handle)

		if err := http.Serve(proxy.listener, httpServer); err != nil {

			// Probably a better way to handle these cases. Meh.
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}

			// Looks like something is actually wrong
			log.Fatal(err)
		}

	}()

	return nil
}

// handle gets called on each request. We use this to update the host header.
func (proxy *forwardingProxy) handle(w http.ResponseWriter, r *http.Request) {
	// Replace the host so that the Host: header is correct
	r.Host = proxy.targetURL.Host

	proxy.server.ServeHTTP(w, r)
}

func (proxy *forwardingProxy) stop() {
	proxy.listener.Close()
}
