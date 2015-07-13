package main

import (
	"encoding/json"
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
	router.GET("/containers/:endpoint", containersEndpointsHandler)
	router.GET("/containers/:endpoint/:id", containersEndpointsHandler)
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
	switch ps.ByName("endpoint") {
	case "list":
		fmt.Fprintf(w, listImages(cfg))
	case "inspect":
		fmt.Fprintf(w, inspectImage(cfg, ps.ByName("id")))
	case "history":
		fmt.Fprintf(w, historyImage(cfg, ps.ByName("id")))
	}
}
func containersEndpointsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	switch ps.ByName("endpoint") {
	case "list":
		fmt.Fprintf(w, listContainers(cfg))
	case "inspect":
		fmt.Fprintf(w, inspectContainer(cfg, ps.ByName("id")))
	case "log":
		fmt.Fprintf(w, logContainer(cfg, ps.ByName("id")))
	case "top":
		fmt.Fprintf(w, topContainer(cfg, ps.ByName("id")))
	case "stats":
		fmt.Fprintf(w, statsContainer(cfg, ps.ByName("id")))
	case "changes":
		fmt.Fprintf(w, changesContainer(cfg, ps.ByName("id")))
	case "start":
		fmt.Fprintf(w, startContainer(cfg, ps.ByName("id")))
	case "stop":
		fmt.Fprintf(w, stopContainer(cfg, ps.ByName("id")))
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

func listContainers(cfg Config) string {
	uri := fmt.Sprintf("%s/containers/json?all=true", cfg.Addr)
	return sendRequest(uri)
}

func inspectContainer(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/json", cfg.Addr, containerId)
	return sendRequest(uri)
}

func logContainer(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/log", cfg.Addr, containerId)
	return sendRequest(uri)
}

func topContainer(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/top", cfg.Addr, containerId)
	return sendRequest(uri)
}

func statsContainer(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/stats", cfg.Addr, containerId)
	fmt.Println("Req stats:", uri)
	s := sendRequest(uri)
	fmt.Println("Got stats:", s)
	return s
}

func changesContainer(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/changes", cfg.Addr, containerId)
	return sendRequest(uri)
}

func startContainer(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/start", cfg.Addr, containerId)
	return sendRequest(uri)
}

func stopContainer(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/stop", cfg.Addr, containerId)
	return sendRequest(uri)
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
		log.Fatal("Error getting http resource.", err)
	} else {
		defer resp.Body.Close()
		status = resp.StatusCode
	}
	body := ""
	if status == 200 {
		bodyBuf, _ := lib.ReadHttpResponseBody(resp)
		body = string(bodyBuf)
	} else {
		b, err := json.Marshal(resp)
		if err != nil {
			body = "{ success: false, error: 'unknown' }"
		} else {
			body = string(b)
		}
	}
	fmt.Println("Body:", body)
	return body
}
