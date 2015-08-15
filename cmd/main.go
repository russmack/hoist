package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
	lib "github.com/russmack/hoist/lib"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
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
	rootPath   = "../www/"
	dbFilename = "hoist.db"
)

var (
	templates = template.Must(template.ParseFiles(
		path.Join(rootPath, "index.html"),
		path.Join(rootPath, "images.html"),
		path.Join(rootPath, "containers.html"),
		path.Join(rootPath, "nodes.html"),
		path.Join(rootPath, "monitor.html"),
		path.Join(rootPath, "header.html"),
		path.Join(rootPath, "footer.html"),
		path.Join(rootPath, "menubar.html"),
	))
)

func init() {

	db := NewDatabase(dbFilename)
	db.Init()
}

func main() {
	initConfig()
	// httprouter is too strict with routes - consider another, or wait for v2.
	router := httprouter.New()
	router.HandlerFunc("GET", "/offline.appcache", appcacheHandler)
	router.HandlerFunc("GET", "/favicon.ico", faviconHandler)
	router.HandlerFunc("GET", "/", indexHandler)
	router.HandlerFunc("GET", "/index.html", indexHandler)
	router.HandlerFunc("GET", "/images.html", imagesHandler)
	router.HandlerFunc("GET", "/containers.html", containersHandler)
	router.HandlerFunc("GET", "/nodes.html", nodesHandler)
	router.HandlerFunc("GET", "/monitor.html", monitorHandler)
	router.GET("/images/:endpoint", imagesGetHandler)
	router.GET("/images/:endpoint/:id", imagesGetHandler)
	router.GET("/containers/:endpoint", containersGetHandler)
	router.GET("/containers/:endpoint/:id", containersGetHandler)
	router.GET("/nodes/get/:nodeid/images/list", nodeImagesGetHandler)
	router.GET("/nodes/list", nodesListHandler)
	router.GET("/monitor/:endpoint/:nodeid", monitorGetHandler)
	router.POST("/nodes", nodesPostHandler)
	router.ServeFiles("/static/*filepath", http.Dir(rootPath))

	fmt.Println("Starting server on port 8100...")
	log.Fatal(http.ListenAndServe(":8100", router))
}

func appcacheHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/cache-manifest")
	http.ServeFile(w, r, path.Join(rootPath, "offline.appcache"))
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	body, err := base64.StdEncoding.DecodeString(faviconBase64)
	if err != nil {
		fmt.Println("favicon handler decoding error:", err)
		return
	}
	w.Header().Set("content-type", "image/x-icon")
	w.Write(body)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Mainscript string
	}{
		"index",
	}
	templates.ExecuteTemplate(w, "index.html", data)
}
func imagesHandler(w http.ResponseWriter, r *http.Request) {
	nid := r.URL.Query().Get("nodeid")
	data := struct {
		Mainscript string
		NodeId     string
	}{
		"images",
		nid,
	}
	templates.ExecuteTemplate(w, "images.html", data)
}
func containersHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Mainscript string
	}{
		"containers",
	}
	templates.ExecuteTemplate(w, "containers.html", data)
}
func nodesHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Mainscript string
	}{
		"nodes",
	}
	templates.ExecuteTemplate(w, "nodes.html", data)
}

func monitorHandler(w http.ResponseWriter, r *http.Request) {
	nid := r.URL.Query().Get("nodeid")
	data := struct {
		Mainscript string
		NodeId     string
	}{
		"monitor",
		nid,
	}
	templates.ExecuteTemplate(w, "monitor.html", data)
}
func imagesGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	switch ps.ByName("endpoint") {
	case "inspect":
		fmt.Fprintf(w, imageInspect(cfg, ps.ByName("id")))
	case "history":
		fmt.Fprintf(w, imageHistory(cfg, ps.ByName("id")))
	case "search":
		fmt.Fprintf(w, imageSearch(cfg, ps.ByName("id")))
	case "delete":
		fmt.Fprintf(w, imageDelete(cfg, ps.ByName("id")))
	}
}
func containersGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	case "restart":
		fmt.Fprintf(w, containerRestart(cfg, ps.ByName("id")))
	case "delete":
		fmt.Fprintf(w, containerDelete(cfg, ps.ByName("id")))
	}
}
func nodesListHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, nodeList(cfg))
}
func nodesPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var node Node
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&node)
	if err != nil {
		fmt.Println("Unable to decode json node post.", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, nodeAdd(cfg, &node))
}

type Response map[string]interface{}

func nodeImagesGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, imageList(cfg, ps.ByName("nodeid")))
}

func monitorGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	switch ps.ByName("endpoint") {
	case "info":
		fmt.Fprintf(w, monitorInfo(cfg, ps.ByName("nodeid")))
	case "version":
		fmt.Fprintf(w, monitorVersion(cfg, ps.ByName("nodeid")))
	case "ping":
		fmt.Fprintf(w, monitorPing(cfg, ps.ByName("nodeid")))
	case "events":
		monitorEvents(cfg, ps.ByName("nodeid"), w)
	}
}
func elseHandler(w http.ResponseWriter, r *http.Request) {
	p := path.Join(rootPath, r.URL.Path)
	fmt.Println(p)
	http.ServeFile(w, r, p)
}

func imageList(cfg Config, nodeId string) string {
	fmt.Println("Getting node for id:", nodeId)
	// Get ipaddress for nodeId from db
	db := NewDatabase(dbFilename)
	nodesDb := NewNodesDataStore(db)
	n, err := strconv.ParseInt(nodeId, 10, 64)
	node, err := nodesDb.GetNodeById(n)
	if err != nil {
		fmt.Println("Unable to get node for monitor info.", err)
		return ""
	}

	fmt.Printf("Got node for monitor info: %+v\n", node)
	// Replace ip address in cfg.Addr with node.Address
	port := 2376
	fmt.Println("PORT:", node.Port)
	if node.Port != 0 {
		port = node.Port
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, port)
	uri := fmt.Sprintf("%s/images/json", addr)
	fmt.Println(" for addr:", uri)
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

func imageDelete(cfg Config, imageId string) string {
	uri := fmt.Sprintf("%s/images/%s", cfg.Addr, imageId)
	return deleteHttp(uri)
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
	return postHttp(uri, "")
}

func containerStop(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/stop", cfg.Addr, containerId)
	return postHttp(uri, "")
}

func containerRestart(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/restart", cfg.Addr, containerId)
	return postHttp(uri, "")
}

func containerDelete(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s", cfg.Addr, containerId)
	return deleteHttp(uri)
}

func nodeList(cfg Config) string {
	db := NewDatabase(dbFilename)
	nodesDb := NewNodesDataStore(db)
	nodes := nodesDb.GetNodes()
	b, err := json.Marshal(nodes)
	if err != nil {
		fmt.Println(err)
		return "err occurred"
	}
	return string(b)
}
func nodeAdd(cfg Config, h *Node) string {
	h.Created = time.Now().String()
	db := NewDatabase(dbFilename)
	nodesDb := NewNodesDataStore(db)
	node, err := nodesDb.AddNode(h)
	if err != nil {
		return fmt.Sprintf("Unable to add node.", err)
	}
	json, err := json.Marshal(node)
	if err != nil {
		return fmt.Sprintf("Unable to marshal new node json.", err)
	}
	return string(json)
}

func monitorInfo(cfg Config, nodeId string) string {
	fmt.Println("Getting node for id:", nodeId)
	// Get ipaddress for nodeId from db
	db := NewDatabase(dbFilename)
	nodesDb := NewNodesDataStore(db)
	n, err := strconv.ParseInt(nodeId, 10, 64)
	node, err := nodesDb.GetNodeById(n)
	if err != nil {
		fmt.Println("Unable to get node for monitor info.", err)
		return ""
	}
	fmt.Printf("Got node for monitor info: %+v\n", node)
	// Replace ip address in cfg.Addr with node.Address
	port := 2376
	fmt.Println("PORT:", node.Port)
	if node.Port != 0 {
		port = node.Port
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, port)

	//uri := fmt.Sprintf("%s/info", cfg.Addr)
	uri := fmt.Sprintf("%s/info", addr)
	fmt.Println("Monitoring info for addr:", uri)
	return getHttpString(uri)
}
func monitorVersion(cfg Config, nodeId string) string {
	db := NewDatabase(dbFilename)
	nodesDb := NewNodesDataStore(db)
	n, err := strconv.ParseInt(nodeId, 10, 64)
	node, err := nodesDb.GetNodeById(n)
	if err != nil {
		fmt.Println("Unable to get node for monitor info.", err)
		return ""
	}
	fmt.Printf("Got node for monitor ping: %+v\n", node)
	// Replace ip address in cfg.Addr with node.Address
	port := 2376
	if node.Port != 0 {
		port = node.Port
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, port)

	uri := fmt.Sprintf("%s/version", addr)
	return getHttpString(uri)
}
func monitorPing(cfg Config, nodeId string) string {
	db := NewDatabase(dbFilename)
	nodesDb := NewNodesDataStore(db)
	n, err := strconv.ParseInt(nodeId, 10, 64)
	node, err := nodesDb.GetNodeById(n)
	if err != nil {
		fmt.Println("Unable to get node for monitor info.", err)
		return ""
	}
	fmt.Printf("Got node for monitor ping: %+v\n", node)
	// Replace ip address in cfg.Addr with node.Address
	port := 2376
	if node.Port != 0 {
		port = node.Port
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, port)

	uri := fmt.Sprintf("%s/_ping", addr)
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

func monitorEvents(cfg Config, nodeId string, w http.ResponseWriter) {
	fmt.Println("monitoring events ...")
	db := NewDatabase(dbFilename)
	nodesDb := NewNodesDataStore(db)
	n, err := strconv.ParseInt(nodeId, 10, 64)
	node, err := nodesDb.GetNodeById(n)
	if err != nil {
		fmt.Println("Unable to get node for monitor info.", err)
		return
	}
	fmt.Printf("Got node for monitor ping: %+v\n", node)
	// Replace ip address in cfg.Addr with node.Address
	port := 2376
	if node.Port != 0 {
		port = node.Port
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, port)

	done := make(chan bool)
	uri := fmt.Sprintf("%s/events", addr)
	eChan := make(chan Event)

	f, ok := w.(http.Flusher)
	if !ok {
		//http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	cn, ok := w.(http.CloseNotifier)
	if !ok {
		//http.Error(rw, "cannot stream", http.StatusInternalServerError)
		return
	}

	go func(w http.ResponseWriter, eChan chan Event) {
	loop:
		for {
			select {
			case <-cn.CloseNotify():
				fmt.Println("done: closed connection")
				return
			case ev, more := <-eChan:
				if !more {
					fmt.Println("Finished rx from ev chan")
					break loop
				}
				fmt.Println("event: %+v", ev)
				fmt.Fprintf(w, "data: %+v\n\n", ev)
				f.Flush()
				////break loop
			}
		}
		fmt.Println("sending done")
		done <- true
	}(w, eChan)
	getHttpStream(uri, eChan)
	<-done
	fmt.Println("Finished stream")
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
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting http resource. %s\" }", err)
		log.Println(body)
		return body
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

func postHttp(uri string, data string) string { // TODO: change 'data' type.

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
	postBody := bytes.NewBuffer([]byte(data))
	req, err := http.NewRequest("POST", uri, postBody)
	if err != nil {
		log.Fatal("Error creating new POST request.")
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
	fmt.Println("Requesting stream...")
	req, err := http.NewRequest("GET", uri, nil)
	res, err := client.Do(req)
	fmt.Println("Reading stream...")
	go func(res *http.Response, client *http.Client) {
		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)
		for {
			var event Event
			fmt.Println("loop start: %+v", event)
			err = decoder.Decode(&event)
			fmt.Println("loop decoding: %+v", event)
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
			fmt.Printf("event fired: %+v\n", event)
			eChan <- event
			fmt.Println("event enqueued")
		}
	}(res, &client)
}

func (d *NodesDataStore) GetNodes() []Node {
	// TODO: maxrows should not be hardcoded.
	return selectRows(d.Db.DbName, "Nodes", "50")
}

func selectRows(dbName string, tableName string, maxRows string) []Node {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println("Error: unable to open database: " + err.Error())
		os.Exit(1)
	}
	defer db.Close()

	stmt, err := db.Prepare(
		"select rowid, name, scheme, address, port, description, created " +
			" from " + tableName + " limit " + maxRows)
	//stmt, err := db.Prepare("select * from ? limit ?")
	if err != nil {
		fmt.Println("Error: unable to prepare query: " + err.Error())
		os.Exit(1)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("Error: unable to execute query: " + err.Error())
		os.Exit(1)
	}
	defer rows.Close()
	nodes := []Node{}
	for rows.Next() {
		var rowid int
		var name string
		var scheme string
		var address string
		var port int
		var description string
		var created string
		err := rows.Scan(&rowid, &name, &scheme, &address, &port, &description, &created)
		if err != nil {
			fmt.Println("ERR: ", err)
		}
		node := &Node{
			Id:          rowid,
			Name:        name,
			Scheme:      scheme,
			Address:     address,
			Port:        port,
			Description: description,
			Created:     created,
		}
		nodes = append(nodes, *node)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("Err from rows: ", err)
	}
	return nodes
}

type Node struct {
	Id          int
	Name        string
	Scheme      string
	Address     string
	Port        int
	Description string
	Created     string
}

type Database struct {
	DbName string
}
type NodesDataStore struct {
	Db *Database
}

func NewDatabase(dbName string) *Database {
	return &Database{DbName: dbName}
}

func NewNodesDataStore(db *Database) *NodesDataStore {
	return &NodesDataStore{Db: db}
}

func (d *NodesDataStore) CreateTable() {
	stmt := ` 
			create table if not exists Nodes ( 
		        Name text, 
				Scheme text,
		        Address text not null, 
				Port integer,
		        Description text, 
		        Created text 
		    );
			`
	d.Db.CreateTable(stmt)
}

func (d *Database) Init() {
	// Ensure tables exist.
	nodesDb := NewNodesDataStore(d)
	nodesDb.CreateTable()
}

func (d *Database) CreateTable(stmt string) {
	db, err := sql.Open("sqlite3", d.DbName)
	if err != nil {
		fmt.Println("Error: unable to open database: " + err.Error())
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec(stmt)
	if err != nil {
		fmt.Println("Error: unable to create database table: " + err.Error())
		os.Exit(1)
	}
}

func (d *NodesDataStore) AddNode(n *Node) (Node, error) {
	db, err := sql.Open("sqlite3", d.Db.DbName)
	if err != nil {
		fmt.Println("Error: unable to open database: " + err.Error())
		return Node{}, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error: unable to being transaction: " + err.Error())
		return Node{}, err
	}

	stmt, err := tx.Prepare("insert into Nodes(Name, Scheme, Address, Port, Description, Created) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error: unable to prepare transaction statement: " + err.Error())
		return Node{}, err
	}
	defer stmt.Close()

	r, err := stmt.Exec(n.Name, n.Scheme, n.Address, n.Port, n.Description, n.Created)
	if err != nil {
		fmt.Println("Error: unable to insert database record: " + err.Error())
		return Node{}, err
	}
	tx.Commit()
	lastInsertedId, err := r.LastInsertId()
	if err != nil {
		return Node{}, err
	}
	node, err := d.GetNodeById(lastInsertedId)
	if err != nil {
		fmt.Println("Unable to GetNode.", err)
		return Node{}, err
	}

	/*
		for rows.Next() {
			var id int
			var name string
			var address string
			var description string
			var created string
			rows.Scan(&id, &name, &address, &description, &created)
			node := &Node{
				Id:          id,
				Name:        name,
				Address:     address,
				Description: description,
				Created:     created,
			}
			nodes = append(nodes, *node)
		}
	*/
	return node, nil
}

func (d *NodesDataStore) GetNodeById(id int64) (Node, error) {
	db, err := sql.Open("sqlite3", d.Db.DbName)
	if err != nil {
		fmt.Println("Error: unable to open database: " + err.Error())
		return Node{}, err
	}
	defer db.Close()

	stmt := "select rowid, Name, Scheme, Address, Port, Description, Created from Nodes where rowid = ?"
	row := db.QueryRow(stmt, id)
	var node Node
	row.Scan(&node.Id, &node.Name, &node.Scheme, &node.Address, &node.Port, &node.Description, &node.Created)
	switch {
	case err == sql.ErrNoRows:
		log.Println("No node with specified id.")
	case err != nil:
		log.Println("Unable to Get Node.", err)
	default:
		//
	}
	return node, err
}
