package lib

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	httpTimeout = time.Duration(1 * time.Second)
)

func DialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, httpTimeout)
}

func ReadHttpResponseBody(resp *http.Response) ([]byte, error) {
	var buf bytes.Buffer
	n, err := io.Copy(&buf, resp.Body)
	if err != nil {
		log.Println("Error while copying body, got", n, "bytes:", err)
	}
	bufbytes := buf.Bytes()
	buflen := len(bufbytes)
	log.Println(resp.Request.URL.String(), "Body len:", buflen, "; cap:", cap(bufbytes))
	minbuf := make([]byte, buflen)
	copy(minbuf, bufbytes[0:buflen])
	return minbuf, err
}

func readConn(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Error reading socket.", err)
			continue
		}
		data := buf[:n]
		fmt.Println("Received response:", string(data))
	}
}
