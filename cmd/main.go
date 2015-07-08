package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	lib "github.com/russmack/hoist/lib"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	rootPath = "../www/"
)

func main() {
	router := httprouter.New()
	router.HandlerFunc("GET", "/index.html", indexHandler)
	router.HandlerFunc("GET", "/images.html", imagesHandler)
	router.HandlerFunc("GET", "/containers.html", containersHandler)
	router.GET("/images/:endpoint", imagesEndpointsHandler)
	router.HandlerFunc("GET", "/", indexHandler)
	router.ServeFiles("/static/*filepath", http.Dir(rootPath))

	log.Fatal(http.ListenAndServe(":8100", router))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(rootPath, "index.html"))
}
func imagesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(rootPath, "images.html"))
}
func containersHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(rootPath, "containers.html"))
}
func imagesEndpointsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "This is the image list response data: %s\n", ps.ByName("endpoint"))
}
func elseHandler(w http.ResponseWriter, r *http.Request) {
	p := path.Join(rootPath, r.URL.Path)
	fmt.Println(p)
	http.ServeFile(w, r, p)
}

func query() {
	certPath := os.Getenv("DOCKER_CERT_PATH")

	//caCert, _ := getCaCert(certPath + "/ca.pem")
	sslCert, _ := lib.GetSslCert(certPath + "/cert.pem")
	sslKey, _ := lib.GetSslKey(certPath + "/key.pem")
	tlsConfig, err := lib.GetTLSConfig(nil, sslCert, sslKey)
	if err != nil {
		log.Fatal("Error getting TLS config.", err)
	}
	//fd, err := net.Dial("unix", "/var/run/docker.sock")
	//fd, err := net.Dial("tcp", "192.168.59.103:2375")
	fmt.Println("Dialing...")

	// START REQUEST

	tlsConfig.InsecureSkipVerify = true

	transport := http.Transport{
		Dial:            lib.DialTimeout,
		TLSClientConfig: tlsConfig,
	}
	status := 0
	client := http.Client{
		Transport: &transport,
	}
	addr := "https://192.168.59.103:2376"
	//uri := fmt.Sprintf("%s/_ping", addr)
	//uri := fmt.Sprintf("%s/containers/json", addr)
	uri := fmt.Sprintf("%s/images/json", addr)
	resp, err := client.Get(uri)
	if err != nil {
		//return 0, err
		log.Fatal("Error getting http resource.", err)
	} else {
		defer resp.Body.Close()
		status = resp.StatusCode
	}
	fmt.Println("Status:", status)
	fmt.Println("Resp:")
	body := ""
	if status == 200 {
		bodyBuf, _ := lib.ReadHttpResponseBody(resp)
		body = string(bodyBuf)
	} else {
		body = "No body in http response."
	}
	fmt.Println(body)
	// END REQUEST

	fmt.Println("Starting read routine...")
	//go readConn(fd)

	//done := make(chan bool)
	//<-done
}
