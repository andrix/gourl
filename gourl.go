package main

import (
    "os"
    "flag"
    "fmt"
    "io"
    "net/url"
    "net/http"
    "crypto/tls"
    "strings"
    "log"
    "encoding/base64"
)

func main() {
    var headers = flag.Bool("I", false, "Print headers information")
    var insecure = flag.Bool("k", false, "Allow to perform insecure SSL connections")
    var user_agent = flag.String("A", "gourl/1.0", "Set the User-Agent")
    var output = flag.String("o", "", "Write output to a file instead of Stdout")
    var proxy = flag.String("x", "", "HTTP Proxy in the form HOST:PORT")
    var proxy_user = flag.String("U", "", "Proxy user and password in the form USER:PASS")

    flag.Parse()
    log.SetFlags(0)
    if len(flag.Args()) <= 0 {
        fmt.Printf("Usage: gourl [options] url\n")
    } else {
        aurl := flag.Arg(0)
        req, err := http.NewRequest("GET", aurl, nil)
        if err != nil {
            log.Fatalf("Error creating request: %s\n", err)
        }
        req.Header.Add("User-Agent", *user_agent)
        if *proxy_user != "" {
            authtoken := "Basic " + base64.StdEncoding.EncodeToString([]byte(*proxy_user))
            req.Header.Add("Proxy-Authorization", authtoken)
        }

        // Create TLS config
        tlsConfig := tls.Config{RootCAs: nil}

        // If insecure, skip CA verfication
        if *insecure {
            tlsConfig.InsecureSkipVerify = true
        }

        tr := &http.Transport{
            TLSClientConfig:    &tlsConfig,
            DisableCompression: true,
        }

        // If -x/Proxy - set up HTTP proxy
        if *proxy != "" {
            proxyUrl, err := url.Parse("http://" + *proxy)
            if err != nil {
                panic(err)
            }
            tr.Proxy = http.ProxyURL(proxyUrl)
        }

        client := &http.Client{Transport: tr}
        resp, err := client.Do(req)

        if err != nil {
            log.Fatalf("An error ocurred: %s\n", err)
        }
        if *headers {
            fmt.Printf("%s %s\n", resp.Proto, resp.Status)
            for k, vals := range resp.Header {
                fmt.Printf("%s: %s\n", k, strings.Join(vals, ";"))
            }
            fmt.Println()
        } else {
            defer resp.Body.Close()
            var out *os.File
            if *output == "" {
                out = os.Stdout
            } else {
                out, err = os.Create(*output)
                if err != nil {
                    log.Fatalf("Error writing to file: %s\n", err)
                }
            }
            defer func() {
                if err := out.Close(); err != nil {
                    panic(err)
                }
            }()

            // Create buffer
            buf := make([]byte, 1024)

            for {
                n, err := resp.Body.Read(buf)
                if err != nil && err != io.EOF { panic(err) }
                if n == 0 { break }

                if _, err := out.Write(buf[:n]); err != nil {
                    panic(err)
                }
            }
        }
    }
}
