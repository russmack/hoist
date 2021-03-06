package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
	lib "github.com/russmack/hoist/lib"
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

type Container struct {
	Command string
	Created int
	Id      string
	Image   string
	Labels  interface{}
	Names   []string
	Ports   []interface{}
	Status  string
}

const (
	version    = "0.1"
	rootPath   = "../www/"
	dbFilename = "hoist.db"
)

var (
	templates = template.Must(template.ParseFiles(
		path.Join(rootPath, "index.html"),
		path.Join(rootPath, "clusters.html"),
		path.Join(rootPath, "nodes.html"),
		path.Join(rootPath, "images.html"),
		path.Join(rootPath, "containers.html"),
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
	router.HandlerFunc("GET", "/clusters.html", clustersHandler)
	router.HandlerFunc("GET", "/nodes.html", nodesHandler)
	router.HandlerFunc("GET", "/images.html", imagesHandler)
	router.HandlerFunc("GET", "/containers.html", containersHandler)
	router.HandlerFunc("GET", "/monitor.html", monitorHandler)
	router.GET("/"+version+"/images/search/:term", nodeImageSearchGetHandler)
	router.GET("/"+version+"/clusters/:clusterid/nodes/:nodeid/images/list", nodeImagesGetHandler)
	router.GET("/"+version+"/nodes/:nodeid/images/inspect/:imageid", nodeImageInspectGetHandler)
	router.GET("/"+version+"/nodes/:nodeid/images/history/:imageid", nodeImageHistoryGetHandler)
	router.GET("/"+version+"/nodes/:nodeid/images/delete/:imageid", nodeImageDeleteHandler)
	router.GET("/"+version+"/nodes/:nodeid/containers/list", nodeContainersGetHandler)
	router.GET("/"+version+"/nodes/:nodeid/containers/inspect/:containerid", nodeContainerInspectGetHandler)
	router.GET("/"+version+"/nodes/:nodeid/containers/top/:containerid", nodeContainerTopGetHandler)
	router.GET("/"+version+"/nodes/:nodeid/containers/start/:containerid", nodeContainerStartGetHandler)
	router.GET("/"+version+"/nodes/:nodeid/containers/stop/:containerid", nodeContainerStopGetHandler)
	router.GET("/"+version+"/nodes/:nodeid/containers/restart/:containerid", nodeContainerRestartGetHandler)
	router.GET("/"+version+"/nodes/:nodeid/containers/changes/:containerid", nodeContainerChangesGetHandler)
	router.GET("/"+version+"/nodes/:nodeid/containers/delete/:containerid", nodeContainerDeleteGetHandler)
	router.GET("/"+version+"/nodes/:nodeid/containers/scaleout/:containerid", nodeContainerScaleOutGetHandler)
	router.GET("/"+version+"/clusters", clustersListHandler)
	router.GET("/"+version+"/clusters/:clusterid/nodes", nodesListHandler)
	router.GET("/"+version+"/monitor/:endpoint/:nodeid", monitorGetHandler)
	router.POST("/"+version+"/clusters", clustersPostHandler)
	router.POST("/"+version+"/nodes", nodesPostHandler)
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
	cid := r.URL.Query().Get("clusterid")
	nid := r.URL.Query().Get("nodeid")
	data := struct {
		Mainscript string
		ClusterId  string
		NodeId     string
	}{
		"images",
		cid,
		nid,
	}
	templates.ExecuteTemplate(w, "images.html", data)
}
func containersHandler(w http.ResponseWriter, r *http.Request) {
	nid := r.URL.Query().Get("nodeid")
	data := struct {
		Mainscript string
		NodeId     string
	}{
		"containers",
		nid,
	}
	templates.ExecuteTemplate(w, "containers.html", data)
}
func clustersHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Mainscript string
	}{
		"clusters",
	}
	templates.ExecuteTemplate(w, "clusters.html", data)
}
func nodesHandler(w http.ResponseWriter, r *http.Request) {
	cid := r.URL.Query().Get("clusterid")
	data := struct {
		Mainscript string
		ClusterId  string
	}{
		"nodes",
		cid,
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

func containersGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	switch ps.ByName("endpoint") {
	case "log":
		fmt.Fprintf(w, containerLog(cfg, ps.ByName("id")))
	case "stats":
		fmt.Fprintf(w, containerStats(cfg, ps.ByName("id")))
	}
}
func nodeContainersGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, containerList(cfg, ps.ByName("nodeid")))
}
func nodeContainerInspectGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, containerInspect(cfg, ps.ByName("nodeid"), ps.ByName("containerid")))
}
func nodeContainerTopGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, containerTop(cfg, ps.ByName("nodeid"), ps.ByName("containerid")))
}
func nodeContainerStartGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, containerStart(cfg, ps.ByName("nodeid"), ps.ByName("containerid")))
}
func nodeContainerStopGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, containerStop(cfg, ps.ByName("nodeid"), ps.ByName("containerid")))
}
func nodeContainerRestartGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, containerRestart(cfg, ps.ByName("nodeid"), ps.ByName("containerid")))
}
func nodeContainerChangesGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, containerChanges(cfg, ps.ByName("nodeid"), ps.ByName("containerid")))
}
func nodeContainerDeleteGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, containerDelete(cfg, ps.ByName("nodeid"), ps.ByName("containerid")))
}
func nodeContainerScaleOutGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, containerScaleOut(cfg, ps.ByName("nodeid"), ps.ByName("containerid")))
}
func nodesListHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, nodeList(cfg, ps.ByName("clusterid")))
}
func clustersListHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, clusterList(cfg))
}

func clustersPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var cluster Cluster
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cluster)
	if err != nil {
		fmt.Println("Unable to decode json cluster post.", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, clusterAdd(cfg, &cluster))
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

func nodeImageSearchGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, imageSearch(cfg, ps.ByName("term")))
}

func nodeImageInspectGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, imageInspect(cfg, ps.ByName("nodeid"), ps.ByName("imageid")))
}

func nodeImageHistoryGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, imageHistory(cfg, ps.ByName("nodeid"), ps.ByName("imageid")))
}

func nodeImageDeleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, imageDelete(cfg, ps.ByName("nodeid"), ps.ByName("imageid")))
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

func getNodeById(nodeId string) (Node, error) {
	fmt.Println("Getting node for id:", nodeId)
	// Get ipaddress for nodeId from db
	db := NewDatabase(dbFilename)
	nodesDb := NewClustersDataStore(db)
	n, err := strconv.ParseInt(nodeId, 10, 64)
	node, err := nodesDb.GetNodeById(n)
	if err != nil {
		fmt.Println("Unable to get node for images list.", err)
		return Node{}, err
	}

	fmt.Printf("Got node for images list: %+v\n", node)
	// Replace ip address in cfg.Addr with node.Address
	port := 2376
	fmt.Println("PORT:", node.Port)
	if node.Port == 0 {
		node.Port = port
	}
	return node, nil
}

func imageList(cfg Config, nodeId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/images/json", addr)
	fmt.Println(" for addr:", uri)
	return getHttpString(uri)
}

func imageInspect(cfg Config, nodeId string, imageId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/images/%s/json", addr, imageId)
	return getHttpString(uri)
}

func imageHistory(cfg Config, nodeId string, imageId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/images/%s/history", addr, imageId)
	fmt.Println(" for addr:", uri)
	return getHttpString(uri)
}

func imageSearch(cfg Config, term string) string {
	uri := fmt.Sprintf("%s/images/search?term=%s", cfg.Addr, term)
	return getHttpString(uri)
}

func imageDelete(cfg Config, nodeId string, imageId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	// TODO: seems need to use image name, not imageId
	uri := fmt.Sprintf("%s/images/%s", addr, imageId)
	fmt.Println("Delete image with uri:", uri)
	return deleteHttp(uri)
}

func containerList(cfg Config, nodeId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/containers/json?all=true", addr)
	fmt.Println(" for addr:", uri)
	return getHttpString(uri)
}

func containerInspect(cfg Config, nodeId string, containerId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/containers/%s/json", addr, containerId)
	return getHttpString(uri)
}

func containerLog(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/log", cfg.Addr, containerId)
	return getHttpString(uri)
}

func containerTop(cfg Config, nodeId string, containerId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/containers/%s/top", addr, containerId)
	return getHttpString(uri)
}

func containerStats(cfg Config, containerId string) string {
	uri := fmt.Sprintf("%s/containers/%s/stats", cfg.Addr, containerId)
	fmt.Println("Req stats:", uri)
	s := getHttpString(uri)
	fmt.Println("Got stats:", s)
	return s
}

func containerChanges(cfg Config, nodeId string, containerId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/containers/%s/changes", addr, containerId)
	return getHttpString(uri)
}

func containerStart(cfg Config, nodeId string, containerId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/containers/%s/start", addr, containerId)
	return postHttp(uri, "", nil)
}

func containerStop(cfg Config, nodeId string, containerId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/containers/%s/stop", addr, containerId)
	return postHttp(uri, "", nil)
}

func containerRestart(cfg Config, nodeId string, containerId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/containers/%s/restart", addr, containerId)
	return postHttp(uri, "", nil)
}

func containerDelete(cfg Config, nodeId string, containerId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/containers/%s", addr, containerId)
	return deleteHttp(uri)
}

type ContainerInspection struct {
	Config struct {
		Image string   `json:"Image"`
		Cmd   []string `json:"Cmd"`
	} `json:"Config"`
	Name string   `json:"Name"`
	Args []string `json:"Args"`
}

func makeContainerName(cfg Config, nodeId, containerName string, cList []Container) string {
	// Give new container a name.
	// UUID is bulky, so use container name with incremented numeric suffix.
	// Get a list of all containers and ensure the new name is unique.
	const suffixMaxLen = 3

	cloneeHasNumSuffix := true
	nameSuffix := containerName[len(containerName)-suffixMaxLen:]
	if num, err := strconv.Atoi(nameSuffix); err != nil || num%1 != 0 {
		cloneeHasNumSuffix = false
	}

	// [{"Command":"/hello","Created":1445715734,"Id":"69cb7ebdbe63b896e5703b49981eddc3bf5f49335bc4a90dcab8090ec57dbe9e",
	//   "Image":"hello-world:latest","Labels":{},"Names":["/pensive_pare001"],"Ports":[],"Status":""}
	highestCloneSuffix := 1
	for _, j := range cList { // List of existing containers
		for _, n := range j.Names { // Each container can have multiple names
			if strings.HasPrefix(n, containerName) {
				// This container name is relevant - figure out the numeric suffix.
				nameSuffix := n[len(n)-suffixMaxLen:]
				if num, err := strconv.Atoi(nameSuffix); err != nil && num%1 == 0 {
					continue
				} else {
					if num >= highestCloneSuffix {
						highestCloneSuffix = num + 1
					}
				}
			}
		}
	}

	containerName = strings.Replace(containerName, "/", "", -1)
	suffix := padString(highestCloneSuffix, suffixMaxLen)
	newName := ""
	if !cloneeHasNumSuffix {
		newName = containerName + "." + suffix
	} else {
		newName = containerName[:len(containerName)-suffixMaxLen-1] + "." + suffix
	}

	return newName
}

func padString(n int, nLen int) string {
	s := strconv.Itoa(n)
	sLen := len(s)
	for i := 0; i < nLen-sLen; i++ {
		s = "0" + s
	}
	return s
}

func containerScaleOut(cfg Config, nodeId string, containerId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	// Inspect container to get image.
	// Inspect container to get container name.
	var cInfo ContainerInspection
	info := containerInspect(cfg, nodeId, containerId)
	err = json.Unmarshal([]byte(info), &cInfo)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error parsing container inspection. %s\" }", err)
		log.Println(body)
		return body
	}

	log.Printf("CINFO: %+v\n\n", cInfo)
	cList, err := currentContainers(cfg, nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error parsing container inspection. %s\" }", err)
		log.Println(body)
		return body
	}
	newName := makeContainerName(cfg, nodeId, cInfo.Name, cList)

	// Create new container with incremented name.
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/containers/create?name=%s", addr, newName)
	body := `{"Image": "` + cInfo.Config.Image + `"}`
	log.Println("Sending: ", uri)
	log.Println("with body: ", body)

	//w.Header().Set("Content-Type", "application/json")
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	return postHttp(uri, body, headers)
}

func currentContainers(cfg Config, nodeId string) ([]Container, error) {
	cListJson := containerList(cfg, nodeId)
	var cList []Container
	err := json.Unmarshal([]byte(cListJson), &cList)
	return cList, err
}

func clusterList(cfg Config) string {
	db := NewDatabase(dbFilename)
	clustersDb := NewClustersDataStore(db)
	clusters := clustersDb.GetClusters()
	b, err := json.Marshal(clusters)
	if err != nil {
		fmt.Println(err)
		return "err occurred"
	}
	return string(b)
}
func nodeList(cfg Config, clusterId string) string {
	db := NewDatabase(dbFilename)
	nodesDb := NewClustersDataStore(db)
	nodes := nodesDb.GetNodes(clusterId)
	b, err := json.Marshal(nodes)
	if err != nil {
		fmt.Println(err)
		return "err occurred"
	}
	return string(b)
}
func clusterAdd(cfg Config, h *Cluster) string {
	h.Created = time.Now().String()
	db := NewDatabase(dbFilename)
	clustersDb := NewClustersDataStore(db)
	cluster, err := clustersDb.AddCluster(h)
	if err != nil {
		return fmt.Sprintf("Unable to add cluster.", err)
	}
	json, err := json.Marshal(cluster)
	if err != nil {
		return fmt.Sprintf("Unable to marshal new cluster json.", err)
	}
	return string(json)
}
func nodeAdd(cfg Config, h *Node) string {
	h.Created = time.Now().String()
	db := NewDatabase(dbFilename)
	nodesDb := NewClustersDataStore(db)
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
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/info", addr)
	fmt.Println("Monitoring info for addr:", uri)
	return getHttpString(uri)
}
func monitorVersion(cfg Config, nodeId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
	uri := fmt.Sprintf("%s/version", addr)
	return getHttpString(uri)
}
func monitorPing(cfg Config, nodeId string) string {
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return body
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
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
	node, err := getNodeById(nodeId)
	if err != nil {
		body := fmt.Sprintf("{ \"success\": false, \"error\": \"Error getting node. %s\" }", err)
		log.Println(body)
		return
	}
	addr := fmt.Sprintf("%s://%s:%d", node.Scheme, node.Address, node.Port)
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
	if status >= 200 && status < 300 {
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
	return body
}

func postHttp(uri string, data string, headers map[string]string) string { // TODO: change 'data' type.

	fmt.Println("Dialing...   for post")

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
	// Add any headers.
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	fmt.Printf("REQUEST : %+v\n", req)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error getting http resource.", err)
	} else {
		defer resp.Body.Close()
		status = resp.StatusCode
	}

	body := ""
	if status >= 200 && status < 300 {
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
			fmt.Println("Error marshalling to json.")
			fmt.Printf("Object to marshal : %+v\n", resp)
			body = "{ success: false, error: '" + err.Error() + "' }"
		} else {
			body = string(b)
		}
	}
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
	if status >= 200 && status < 300 {
		bodyBuf, err := lib.ReadHttpResponseBody(resp)
		if err != nil {
			fmt.Println("err reading body:", err)
			bodyStr := "{ \"success\": false, \"error\": \"" + err.Error() + "\" }"
			bodyBuf = []byte(bodyStr)
		}
		body = string(bodyBuf)
		if strings.TrimSpace(body) == "" {
			body = "{\"StatusCode\": \"OK\"}"
		}
	} else {
		b, err := json.Marshal(resp)
		if err != nil {
			body = "{ success: false, error: '" + err.Error() + "' }"
		} else {
			body = string(b)
		}
	}
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

func (d *ClustersDataStore) GetClusters() []Cluster {
	// TODO: maxrows should not be hardcoded.
	return selectClusterRows(d.Db.DbName, "Clusters", "50")
}
func (d *ClustersDataStore) GetNodes(clusterId string) []Node {
	// TODO: maxrows should not be hardcoded.
	return selectNodeRows(d.Db.DbName, "Nodes", "50", clusterId)
}

func selectClusterRows(dbName string, tableName string, maxRows string) []Cluster {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println("Error: unable to open database: " + err.Error())
		os.Exit(1)
	}
	defer db.Close()

	stmt, err := db.Prepare(
		"select rowid, name, description, created " +
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
	clusters := []Cluster{}
	for rows.Next() {
		var rowid int
		var name string
		var description string
		var created string
		err := rows.Scan(&rowid, &name, &description, &created)
		if err != nil {
			fmt.Println("ERR: ", err)
		}
		cluster := &Cluster{
			Id:          rowid,
			Name:        name,
			Description: description,
			Created:     created,
		}
		clusters = append(clusters, *cluster)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("Err from rows: ", err)
	}
	return clusters
}
func selectNodeRows(dbName string, tableName string, maxRows string, clusterId string) []Node {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println("Error: unable to open database: " + err.Error())
		os.Exit(1)
	}
	defer db.Close()

	stmt, err := db.Prepare(
		"select rowid, name, scheme, address, port, description, created " +
			" from " + tableName + " where ClusterId = " + clusterId + " limit " + maxRows)
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

type Cluster struct {
	Id          int
	Name        string
	Description string
	Created     string
}

type Node struct {
	Id          int
	Name        string
	Scheme      string
	Address     string
	Port        int
	Description string
	Created     string
	ClusterId   int
}

type Database struct {
	DbName string
}
type ClustersDataStore struct {
	Db *Database
}

func NewDatabase(dbName string) *Database {
	return &Database{DbName: dbName}
}

func NewClustersDataStore(db *Database) *ClustersDataStore {
	return &ClustersDataStore{Db: db}
}

func (d *ClustersDataStore) CreateClustersTable() {
	stmt := `
			create table if not exists Clusters (
				Name text, 
				Description text, 
				Created text
			);
			`
	d.Db.CreateTable(stmt)
}

func (d *ClustersDataStore) CreateNodesTable() {
	stmt := ` 
			create table if not exists Nodes ( 
		        Name text, 
				Scheme text,
		        Address text not null, 
				Port integer,
		        Description text, 
		        Created text, 
				ClusterId integer
		    );
			`
	d.Db.CreateTable(stmt)
}

func (d *Database) Init() {
	// Ensure tables exist.
	nodesDb := NewClustersDataStore(d)
	nodesDb.CreateClustersTable()
	nodesDb.CreateNodesTable()
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

func (d *ClustersDataStore) AddCluster(n *Cluster) (Cluster, error) {
	db, err := sql.Open("sqlite3", d.Db.DbName)
	if err != nil {
		fmt.Println("Error: unable to open database: " + err.Error())
		return Cluster{}, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error: unable to being transaction: " + err.Error())
		return Cluster{}, err
	}

	stmt, err := tx.Prepare("insert into Clusters(Name, Description, Created) values (?, ?, ?)")
	if err != nil {
		fmt.Println("Error: unable to prepare transaction statement: " + err.Error())
		return Cluster{}, err
	}
	defer stmt.Close()

	r, err := stmt.Exec(n.Name, n.Description, n.Created)
	if err != nil {
		fmt.Println("Error: unable to insert database record: " + err.Error())
		return Cluster{}, err
	}
	tx.Commit()
	lastInsertedId, err := r.LastInsertId()
	if err != nil {
		return Cluster{}, err
	}
	cluster, err := d.GetClusterById(lastInsertedId)
	if err != nil {
		fmt.Println("Unable to GetCluster.", err)
		return Cluster{}, err
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
	return cluster, nil
}
func (d *ClustersDataStore) AddNode(n *Node) (Node, error) {
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

func (d *ClustersDataStore) GetClusterById(id int64) (Cluster, error) {
	db, err := sql.Open("sqlite3", d.Db.DbName)
	if err != nil {
		fmt.Println("Error: unable to open database: " + err.Error())
		return Cluster{}, err
	}
	defer db.Close()

	stmt := "select rowid, Name, Description, Created from Clusters where rowid = ?"
	row := db.QueryRow(stmt, id)
	var cluster Cluster
	row.Scan(&cluster.Id, &cluster.Name, &cluster.Description, &cluster.Created)
	switch {
	case err == sql.ErrNoRows:
		log.Println("No cluster with specified id.")
	case err != nil:
		log.Println("Unable to Get Cluster.", err)
	default:
		//
	}
	return cluster, err
}
func (d *ClustersDataStore) GetNodeById(id int64) (Node, error) {
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
