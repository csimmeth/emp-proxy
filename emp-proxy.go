package main 

import (
  "log"
  "bytes"
  "os"
  "regexp"
  "strings"
  "net/url"
  "net/http"
  "net/http/httputil"
  "io/ioutil"
  "crypto/tls"
  "github.com/SpotterRF/httpDigestAuth"
)

var ts_regexp *regexp.Regexp = regexp.MustCompile("(\\w+\\.ts)")
var endpoint empEndpoint = loadConfig() 


type empEndpoint struct {
  username string
  password string
  empURL   string
}


func loadConfig() empEndpoint {
  return empEndpoint{
    username: os.Getenv("EMP_USERNAME"),
    password: os.Getenv("EMP_PASSWORD"),
      empURL: os.Getenv("EMP_PATH"),
  }
}


func createFilename(url string) string {
  filename := url[1:] //Strip leading '/'
  filename = strings.Replace(filename, "/", "_", -1)
  if filename == "index.m3u8" {
    filename = "channel.m3u8"
  }
  return filename
}


func processManifest(body []byte) []byte {
  body = bytes.Replace(body, []byte("/"), []byte("_"), -1)
  body = ts_regexp.ReplaceAll(body, []byte("channel_$1"))
  return body
}


func logResponse(r *http.Response) error {
  dump, _ := httputil.DumpResponse(r, true)
  log.Printf("Response: %q", dump)
  return nil
}


// The director modifies the request path, body, and headers
func director (r *http.Request) {
  dump, _ := httputil.DumpRequest(r, false)
  log.Printf("Request Received: %q", dump)

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    dump, _ := httputil.DumpRequest(r, true)
    log.Printf("Error reading body: %v", err)
    log.Printf("Invalid Request: %q", dump)
    return
  }

  // Create a new filename and path
  filename := createFilename(r.URL.String())
  proxyPath := endpoint.empURL + filename 

  // If this is a manifest, create a new body 
  if strings.HasSuffix(filename, ".m3u8") {
    body = processManifest(body)
    r.ContentLength = -1 //Reset the content length because it might have changed
  }

  // Update the original request to reflect the changes
  r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
  r.URL, err = url.Parse(proxyPath)
  if err != nil {
    log.Printf("Error parsing URL: %v", err)
  }
  r.Host = "" 

  // Add HTTP Digest Auth headers
  DigestAuth := &httpDigestAuth.DigestHeaders{}
  DigestAuth, err = DigestAuth.Auth(endpoint.username, endpoint.password, proxyPath) 
  DigestAuth.ApplyAuth(r)

  dump, err = httputil.DumpRequestOut(r, false)
  log.Printf("Outbound Request: %q", dump)
}


func main() {
  log.Printf("Endpoint: %+v", endpoint)
  port := os.Getenv("PROXY_PORT")

  // Disable TLS verify as some EMP endpoints fail
  transCfg := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }

  proxy := &httputil.ReverseProxy {
    Director: director,
    Transport: transCfg,
    ModifyResponse: logResponse,
  }

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    proxy.ServeHTTP(w, r)
  })

  log.Printf("Server starting on localhost port %s", port)
  log.Fatal(http.ListenAndServe(":" + port, nil))
}
