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

type Image struct {
	Id      string `json:"Id"`
	Created int    `json:"Created"`
	//Labels   `json:"Labels"`
	ParentId string `json:"ParentId"`
	//RepoDigests []string `json:"RepoDigests"`
	RepoTags    []string `json:"RepoTags"`
	Size        int      `json:"Size"`
	VirtualSize int      `json:"VirtualSize"`
}

const (
	rootPath = "../www/"
)

func main() {
	initConfig()
	router := httprouter.New()
	router.HandlerFunc("GET", "/index.html", indexHandler)
	router.HandlerFunc("GET", "/images.html", imagesHandler)
	router.HandlerFunc("GET", "/containers.html", containersHandler)
	router.GET("/images/:endpoint", imagesEndpointsHandler)
	router.GET("/images/:endpoint/:id", imagesEndpointsHandler)
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
	//fmt.Fprintf(w, "This is the image list response data: %s\n", ps.ByName("endpoint"))
	switch ps.ByName("endpoint") {
	case "list":
		fmt.Fprintf(w, listImages(cfg))
	case "inspect":
		//fmt.Printf("Received request for inspect: %s\n", ps.ByName("id"))
		//fmt.Fprintf(w, "This is the inspect endpoint for image: %s", ps.ByName("id"))
		fmt.Fprintf(w, inspectImage(cfg, ps.ByName("id")))
	case "history":
		fmt.Fprintf(w, historyImage(cfg, ps.ByName("id")))
	}
}
func elseHandler(w http.ResponseWriter, r *http.Request) {
	p := path.Join(rootPath, r.URL.Path)
	fmt.Println(p)
	http.ServeFile(w, r, p)
}

func listImages(cfg Config) string {
	uri := fmt.Sprintf("%s/images/json", cfg.Addr)
	return sendRequest(uri)
}

func inspectImage(cfg Config, imageId string) string {
	uri := fmt.Sprintf("%s/images/%s/json", cfg.Addr, imageId)
	return sendRequest(uri)
}

func historyImage(cfg Config, imageId string) string {
	uri := fmt.Sprintf("%s/images/%s/history", cfg.Addr, imageId)
	return sendRequest(uri)
}

func listContainers() string {
	//uri := fmt.Sprintf("%s/containers/json", addr)
	return ""
}

func ping() string {
	//uri := fmt.Sprintf("%s/_ping", addr)
	return ""
}

type Config struct {
	CertPath string
	CaCert   string
	SslCert  []byte
	SslKey   []byte
	Addr     string
}

var (
	cfg Config
)

func initConfig() {
	//cfg = new(Config)
	cfg.CertPath = os.Getenv("DOCKER_CERT_PATH")

	//caCert, _ := getCaCert(certPath + "/ca.pem")
	cfg.SslCert, _ = lib.GetSslCert(cfg.CertPath + "/cert.pem")
	cfg.SslKey, _ = lib.GetSslKey(cfg.CertPath + "/key.pem")
	//fd, err := net.Dial("unix", "/var/run/docker.sock")
	//fd, err := net.Dial("tcp", "192.168.59.103:2375")
	cfg.Addr = "https://192.168.59.103:2376"
}

func sendRequest(uri string) string {
	fmt.Println("Dialing...")

	// START REQUEST

	tlsConfig, err := lib.GetTLSConfig(nil, cfg.SslCert, cfg.SslKey)
	if err != nil {
		log.Fatal("Error getting TLS config.", err)
	}
	tlsConfig.InsecureSkipVerify = true

	transport := http.Transport{
		Dial:            lib.DialTimeout,
		TLSClientConfig: tlsConfig,
	}
	status := 0
	client := http.Client{
		Transport: &transport,
	}
	resp, err := client.Get(uri)
	if err != nil {
		//return 0, err
		log.Fatal("Error getting http resource.", err)
	} else {
		defer resp.Body.Close()
		status = resp.StatusCode
	}
	//fmt.Println("Status:", status)
	//fmt.Println("Resp:")
	body := ""
	if status == 200 {
		bodyBuf, _ := lib.ReadHttpResponseBody(resp)
		body = string(bodyBuf)
	} else {
		body = "No body in http response."
	}
	//fmt.Println(body)
	return body
	// END REQUEST

	//fmt.Println("Starting read routine...")
	//go readConn(fd)

	//done := make(chan bool)
	//<-done
}
