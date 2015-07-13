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
	router.HandlerFunc("GET", "/monitor.html", monitorHandler)
	router.GET("/images/:endpoint", imagesEndpointsHandler)
	router.GET("/images/:endpoint/:id", imagesEndpointsHandler)
	router.GET("/containers/:endpoint", containersEndpointsHandler)
	router.GET("/containers/:endpoint/:id", containersEndpointsHandler)
	router.GET("/monitor/:endpoint", monitorEndpointsHandler)
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
func monitorHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(rootPath, "monitor.html"))
}
func imagesEndpointsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	switch ps.ByName("endpoint") {
	case "list":
		fmt.Fprintf(w, imageList(cfg))
	case "inspect":
		fmt.Fprintf(w, imageInspect(cfg, ps.ByName("id")))
	case "history":
		fmt.Fprintf(w, imageHistory(cfg, ps.ByName("id")))
	}
}
func containersEndpointsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	switch ps.ByName("endpoint") {
	case "list":
		fmt.Fprintf(w, containerList(cfg))
	case "inspect":
		fmt.Fprintf(w, containerInspect(cfg, ps.ByName("id")))
	case "log":
		fmt.Fprintf(w, containerLog(cfg, ps.ByName("id")))
	case "top":
		fmt.Fprintf(w, containerTop(cfg, ps.ByName("id")))
	case "stats":
		fmt.Fprintf(w, containerStats(cfg, ps.ByName("id")))
	case "changes":
		fmt.Fprintf(w, containerChanges(cfg, ps.ByName("id")))
	case "start":
		fmt.Fprintf(w, containerStart(cfg, ps.ByName("id")))
	case "stop":
		fmt.Fprintf(w, containerStop(cfg, ps.ByName("id")))
	}
}
func monitorEndpointsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	switch ps.ByName("endpoint") {
	case "info":
		fmt.Fprintf(w, monitorInfo(cfg))
	case "version":
		fmt.Fprintf(w, monitorVersion(cfg))
	case "ping":
		fmt.Fprintf(w, monitorPing(cfg))
	}
}
func elseHandler(w http.ResponseWriter, r *http.Request) {
	p := path.Join(rootPath, r.URL.Path)
	fmt.Println(p)
	http.ServeFile(w, r, p)
}

func imageList(cfg Config) string {
	uri := fmt.Sprintf("%s/images/json", cfg.Addr)
	return sendRequest(uri)
}

func imageInspect(cfg Config, imageId string) string {
	uri := fmt.Sprintf("%s/images/%s/json", cfg.Addr, imageId)
	return sendRequest(uri)
}

func imageHistory(cfg Config, imageId string) string {
	uri := fmt.Sprintf("%s/images/%s/history", cfg.Addr, imageId)
	return sendRequest(uri)
}

func containerList(cfg Config) string {
	uri := fmt.Sprintf("%s/containers/json?all=true", cfg.Addr)
	return sendRequest(uri)
}

func containerInspect(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/json", cfg.Addr, containerId)
	return sendRequest(uri)
}

func containerLog(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/log", cfg.Addr, containerId)
	return sendRequest(uri)
}

func containerTop(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/top", cfg.Addr, containerId)
	return sendRequest(uri)
}

func containerStats(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/stats", cfg.Addr, containerId)
	fmt.Println("Req stats:", uri)
	s := sendRequest(uri)
	fmt.Println("Got stats:", s)
	return s
}

func containerChanges(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/changes", cfg.Addr, containerId)
	return sendRequest(uri)
}

func containerStart(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/start", cfg.Addr, containerId)
	return sendRequest(uri)
}

func containerStop(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/stop", cfg.Addr, containerId)
	return sendRequest(uri)
}

func monitorInfo(cfg Config) string {
	uri := fmt.Sprintf("%s/info", cfg.Addr)
	return sendRequest(uri)
}
func monitorVersion(cfg Config) string {
	uri := fmt.Sprintf("%s/version", cfg.Addr)
	return sendRequest(uri)
}
func monitorPing(cfg Config) string {
	uri := fmt.Sprintf("%s/_ping", cfg.Addr)
	body := sendRequest(uri)
	bodyJson := ""
	b, err := json.Marshal(body)
	if err != nil {
		bodyJson = "{ success: false, error: 'unknown' }"
	} else {
		bodyJson = string(b)
	}
	return bodyJson
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
