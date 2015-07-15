package main

import (
	//"bufio"
	//"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	lib "github.com/russmack/hoist/lib"
	"io"
	//"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
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
	case "search":
		fmt.Fprintf(w, imageSearch(cfg, ps.ByName("id")))
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
	case "delete":
		fmt.Fprintf(w, containerDelete(cfg, ps.ByName("id")))
	}
}

type Response map[string]interface{}

func monitorEndpointsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	switch ps.ByName("endpoint") {
	case "info":
		fmt.Fprintf(w, monitorInfo(cfg))
	case "version":
		fmt.Fprintf(w, monitorVersion(cfg))
	case "ping":
		fmt.Fprintf(w, monitorPing(cfg))
	case "events":
		// if r.Header.Get("Accept") == "application/json" {
		// 	w.Header().Set("Content-Type", "application/json")
		// 	fmt.Fprint(w, Response{"success": true, "message": "OK"})
		// }
		monitorEvents(cfg, w)
	}
}
func elseHandler(w http.ResponseWriter, r *http.Request) {
	p := path.Join(rootPath, r.URL.Path)
	fmt.Println(p)
	http.ServeFile(w, r, p)
}

func imageList(cfg Config) string {
	uri := fmt.Sprintf("%s/images/json", cfg.Addr)
	return getHttpString(uri)
}

func imageInspect(cfg Config, imageId string) string {
	uri := fmt.Sprintf("%s/images/%s/json", cfg.Addr, imageId)
	return getHttpString(uri)
}

func imageHistory(cfg Config, imageId string) string {
	uri := fmt.Sprintf("%s/images/%s/history", cfg.Addr, imageId)
	return getHttpString(uri)
}

func imageSearch(cfg Config, term string) string {
	uri := fmt.Sprintf("%s/images/search?term=%s", cfg.Addr, term)
	return getHttpString(uri)
}

func containerList(cfg Config) string {
	uri := fmt.Sprintf("%s/containers/json?all=true", cfg.Addr)
	return getHttpString(uri)
}

func containerInspect(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/json", cfg.Addr, containerId)
	return getHttpString(uri)
}

func containerLog(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/log", cfg.Addr, containerId)
	return getHttpString(uri)
}

func containerTop(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/top", cfg.Addr, containerId)
	return getHttpString(uri)
}

func containerStats(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/stats", cfg.Addr, containerId)
	fmt.Println("Req stats:", uri)
	s := getHttpString(uri)
	fmt.Println("Got stats:", s)
	return s
}

func containerChanges(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/changes", cfg.Addr, containerId)
	return getHttpString(uri)
}

func containerStart(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/start", cfg.Addr, containerId)
	return getHttpString(uri)
}

func containerStop(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/stop", cfg.Addr, containerId)
	return getHttpString(uri)
}

func containerDelete(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s", cfg.Addr, containerId)
	return deleteHttp(uri)
}

func monitorInfo(cfg Config) string {
	uri := fmt.Sprintf("%s/info", cfg.Addr)
	return getHttpString(uri)
}
func monitorVersion(cfg Config) string {
	uri := fmt.Sprintf("%s/version", cfg.Addr)
	return getHttpString(uri)
}
func monitorPing(cfg Config) string {
	uri := fmt.Sprintf("%s/_ping", cfg.Addr)
	body := getHttpString(uri)
	bodyJson := ""
	b, err := json.Marshal(body)
	if err != nil {
		bodyJson = "{ success: false, error: 'unknown' }"
	} else {
		bodyJson = string(b)
	}
	return bodyJson
}

func monitorEvents(cfg Config, w http.ResponseWriter) {
	done := make(chan bool)
	//uri := fmt.Sprintf("%s/events", cfg.Addr)
	//eChan := make(chan Event)
	eChan := make(chan string)

	//ch := msgBroker.Subscribe()
	//defer msgBroker.Unsubscribe(ch)

	//go func(w http.ResponseWriter, eChan chan Event) {
	go func(w http.ResponseWriter, eChan chan string) {
		f, ok := w.(http.Flusher)
		if !ok {
			//http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("TestHeader", "IWasHere")

		cn, ok := w.(http.CloseNotifier)
		if !ok {
			//http.Error(rw, "cannot stream", http.StatusInternalServerError)
			return
		}

	loop:
		for {
			//msg := <-ch
			//fmt.Fprintf(w, "data: Message: %s\n\n", msg)
			//f.Flush()
			select {
			case <-cn.CloseNotify():
				fmt.Println("done: closed connection")
				return

			case ev, more := <-eChan:
				if !more {
					fmt.Println("Finished rx from ev chan")
					break loop
				}
				fmt.Println("%+v", ev)
				//fmt.Println("rx")
				//fmt.Fprintf(w, "%+v", ev)
				fmt.Fprintf(w, "TEST")
				f.Flush()
				break loop
			}
		}
		fmt.Println("sending done")
		done <- true
	}(w, eChan)
	eChan <- "hi there"
	//getHttpStream(uri, eChan)
	<-done
	fmt.Println("Finished stream")
}

//func monitorEvents(cfg Config, eChan chan string) {
func xmonitorEvents(cfg Config, w io.Writer) {
	done := make(chan bool)
	uri := fmt.Sprintf("%s/events", cfg.Addr)
	eChan := make(chan Event)
	go func(w io.Writer, eChan chan Event) {
	loop:
		for {
			select {
			case ev, more := <-eChan:
				if !more {
					fmt.Println("Finished rx from ev chan")
					break loop
				}
				//fmt.Println("%+v", ev)
				fmt.Println("rx")
				fmt.Fprintf(w, "%+v", ev)
			}
		}
		done <- true
	}(w, eChan)
	getHttpStream(uri, eChan)
	<-done
	fmt.Println("RETURNING")
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

func getHttpString(uri string) string {
	fmt.Println("Dialing...")

	tlsConfig, err := lib.GetTLSConfig(nil, cfg.SslCert, cfg.SslKey)
	if err != nil {
		log.Fatal("Error getting TLS config.", err)
	}
	tlsConfig.InsecureSkipVerify = true

	transport := http.Transport{
		Dial:                  lib.DialTimeout,
		TLSClientConfig:       tlsConfig,
		ResponseHeaderTimeout: time.Second * 45,
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
		bodyBuf, err := lib.ReadHttpResponseBody(resp)
		if err != nil {
			fmt.Println("err reading body:", err)
			bodyStr := "{ \"success\": false, \"error\": \"" + err.Error() + "\" }"
			bodyBuf = []byte(bodyStr)
		}
		body = string(bodyBuf)
	} else {
		b, err := json.Marshal(resp)
		if err != nil {
			body = "{ success: false, error: '" + err.Error() + "' }"
		} else {
			body = string(b)
		}
	}
	fmt.Println("Body:", body)
	return body
}

func deleteHttp(uri string) string {

	fmt.Println("Dialing...   for delete")

	tlsConfig, err := lib.GetTLSConfig(nil, cfg.SslCert, cfg.SslKey)
	if err != nil {
		log.Fatal("Error getting TLS config.", err)
	}
	tlsConfig.InsecureSkipVerify = true

	transport := http.Transport{
		Dial:                  lib.DialTimeout,
		TLSClientConfig:       tlsConfig,
		ResponseHeaderTimeout: time.Second * 45,
	}
	status := 0
	client := http.Client{
		Transport: &transport,
	}
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		log.Fatal("Error creating new DELETE request.")
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error getting http resource.", err)
	} else {
		defer resp.Body.Close()
		status = resp.StatusCode
	}

	body := ""
	if status == 200 {
		bodyBuf, err := lib.ReadHttpResponseBody(resp)
		if err != nil {
			fmt.Println("err reading body:", err)
			bodyStr := "{ \"success\": false, \"error\": \"" + err.Error() + "\" }"
			bodyBuf = []byte(bodyStr)
		}
		body = string(bodyBuf)
	} else {
		b, err := json.Marshal(resp)
		if err != nil {
			body = "{ success: false, error: '" + err.Error() + "' }"
		} else {
			body = string(b)
		}
	}
	fmt.Println("Body:", body)
	return body
}

type Event struct {
	Status string `json:"Status"`
	Id     string `json:"Id"`
	From   string `json:"From"`
	Time   int64  `json:"Time"`
}

func getHttpStream(uri string, eChan chan Event) {
	fmt.Println("Dialing...")

	tlsConfig, err := lib.GetTLSConfig(nil, cfg.SslCert, cfg.SslKey)
	if err != nil {
		log.Fatal("Error getting TLS config.", err)
	}
	tlsConfig.InsecureSkipVerify = true

	transport := http.Transport{
		//Dial:            lib.DialTimeout,
		TLSClientConfig: tlsConfig,
		//ResponseHeaderTimeout: time.Second * 15,
	}
	//status := 0
	client := http.Client{
		Transport: &transport,
	}
	//req, err := http.NewRequest("GET", c.path.String(), nil)
	fmt.Println("Requesting stream...")
	req, err := http.NewRequest("GET", uri, nil)
	res, err := client.Do(req)
	fmt.Println("Reading stream...")
	go func(res *http.Response, client *http.Client) {
		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)
		for {
			var event Event
			err = decoder.Decode(&event)
			if err != nil {
				//if err == io.EOF || err == io.ErrUnexpectedEOF {
				// if c.eventMonitor.isEnabled() {
				// 	// Signal that we're exiting.
				// 	eventChan <- EOFEvent
				// }
				//fmt.Println("...broken...")
				//break
				//}
				//errChan <- err
				fmt.Println("decoder err", err)
				close(eChan)
				break
			}
			// if event.Time == 0 {
			// 	fmt.Println(".")
			// 	continue
			// }
			//if !c.eventMonitor.isEnabled() {
			//	return
			//}
			//eventChan <- &event
			//fmt.Printf("event: %+v\n", event)
			eChan <- event
		}
	}(res, &client)
}
