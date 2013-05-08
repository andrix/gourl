package main

import (
    "os"
    "flag"
    "fmt"
    "io"
    "net/http"
    "crypto/tls"
    "strings"
)

func main() {
    var headers = flag.Bool("I", false, "Print headers information")
    var insecure = flag.Bool("k", false, "Allow to perform insecure SSL connections")
    var user_agent = flag.String("A", "gourl/1.0", "Set the User-Agent")
    var output = flag.String("o", "", "Write output to a file instead of Stdout")

    flag.Parse()
    if len(flag.Args()) <= 0 {
        fmt.Printf("Usage: gourl [options] url\n")
    } else {
        url := flag.Arg(0)
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            fmt.Printf("Error creating request: %s\n", err)
            return
        }
        req.Header.Add("User-Agent", *user_agent)

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
        client := &http.Client{Transport: tr}
        resp, err := client.Do(req)

        if err != nil {
            fmt.Printf("An error ocurred: %s\n", err)
            return
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
                    fmt.Printf("Error writing to file: %s\n", err)
                    return
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
