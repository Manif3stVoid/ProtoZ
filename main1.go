package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	
	"github.com/chromedp/chromedp"
	"github.com/tomnomnom/unfurl"
)

func main() {
	var userJS string
	flag.StringVar(&userJS, "j", "", "js to run on all pages")
	flag.Parse()

	// Open the file containing URLs
	file, err := os.Open("urls.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Read URLs from the file into a slice
	var urls []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		urls = append(urls, sc.Text())
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	var wg sync.WaitGroup
	results := make(chan string)
  
	numWorkers := 5
	wg.Add(numWorkers)

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for url := range urls {
				formattedURLs := unfurl.Format("%s://%d%p", url)
        
				queryURL := fmt.Sprintf("%s?proto[hacker]=1337", formattedURLs)

				// Step 3: Run ProtoZ with awk in shell command
				cmd := exec.Command("bash", "-c", fmt.Sprintf(`echo "%s" | awk '{print $1 "?proto[hacker]=1337"}'`, queryURL))
				output, err := cmd.Output()
				if err != nil {
					log.Printf("error running command: %v", err)
					continue
				}
				processedURL := strings.TrimSpace(string(output))

				ctx, cancel := chromedp.NewContext(context.Background())
				defer cancel()

				var res string
				err = chromedp.Run(ctx,
					chromedp.Navigate(processedURL),
					chromedp.Evaluate(userJS, &res),
				)
				if err != nil {
					log.Printf("Error on %s : %s", processedURL, err)
					continue
				}

				// Send result to results channel
				results <- fmt.Sprintf("%s %v", processedURL, res)
			}
		}()
	}

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results
	for result := range results {
		log.Println(result)
	}
}
