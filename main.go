package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/meagar/himama-dl/internal/himama"
	"golang.org/x/sync/semaphore"
)

func main() {
	fmt.Println("himama-dl v0.0.3")

	username, password, err := fetchCredentials()
	if err != nil {
		fmt.Println("Error colleting credentials:", err)
		return
	}

	client, err := himama.NewClient(username, password)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	children, err := client.FetchChildren()
	if err != nil {
		fmt.Println("Error initializing HiMama client:", err)
		return
	}

	chosenChildren, err := selectChildren(children)
	if err != nil {
		fmt.Println("Error selecting children for download:", err)
		return
	}

	for _, c := range chosenChildren {
		if err := scrape(client, c); err != nil {
			fmt.Println("Error downloaded data for", c.Name, ":", err)
			return
		}
	}
	fmt.Printf("Total: %d\nDownloaded %d\nAlready Downloaded: %d\n", total, completed, skipped)
}

func fetchCredentials() (username, password string, err error) {
	flag.StringVar(&username, "username", "", "HiMama username (ie, your email)")
	flag.StringVar(&password, "password", "", "HiMama password")
	flag.Parse()

	if username == "" {
		fmt.Print("Username: ")
		fmt.Scanf("%s", &username)
	}

	if password == "" {
		fmt.Print("Password: ")
		fmt.Scanf("%s", &password)
	}

	return
}

var total, completed, skipped int32

func scrape(client *himama.Client, child himama.Child) error {
	mkdir("./" + child.Name)

	work := scrapeActivities(client, child)

	// blocks until all downloads are finished
	spawnDownloadWorkers(child, work)

	return nil
}

func spawnDownloadWorkers(child himama.Child, work <-chan himama.Activity) {
	wg := sync.WaitGroup{}
	// These workers hit S3, so we can parallelize pretty heavily
	sem := semaphore.NewWeighted(10)

	for activity := range work {
		sem.Acquire(context.Background(), 1)
		wg.Add(1)
		go func(activity himama.Activity) {
			defer wg.Done()
			defer sem.Release(1)
			filename := activity.SuggestedLocalFilename()

			dest := "./" + child.Name + "/" + filename
			if !fileExists(dest) {
				download(activity.MediaURL, dest)
				atomic.AddInt32(&completed, 1)
			} else {
				atomic.AddInt32(&skipped, 1)
			}

			fmt.Printf("%d/%d: %s\n", completed, total, filename)
		}(activity)
	}

	wg.Wait()
}

func download(srcURL, destPath string) {
	res, err := http.Get(srcURL)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fh, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	if _, err := io.Copy(fh, res.Body); err != nil {
		panic(err)
	}
}

// Scraps the activity index pages
// Returns a channel over which activity media links (ie links to S3 objects) are sent
func scrapeActivities(client *himama.Client, child himama.Child) <-chan himama.Activity {
	work := make(chan himama.Activity, 2000)

	wg := sync.WaitGroup{}
	sem := semaphore.NewWeighted(5)

	go func() {
		done := false
		for page := 1; !done; page++ {
			wg.Add(1)
			sem.Acquire(context.Background(), 1)

			page := page
			go func() {
				defer sem.Release(1)
				defer wg.Done()

				if activities, err := client.Activities(child, page); err != nil {
					panic(err)
				} else {
					if len(activities) == 0 {
						done = true
					}

					atomic.AddInt32(&total, int32(len(activities)))
					for _, activity := range activities {
						work <- activity
					}
				}
			}()
		}

		wg.Wait()
		close(work)
	}()

	return work
}

func selectChildren(children []himama.Child) ([]himama.Child, error) {
	if len(children) == 0 {
		fmt.Println("Unable to find children")
		return nil, fmt.Errorf("no children found")
	}

	if len(children) == 1 {
		// TODO: Test this codepath
		fmt.Println("Found 1 child:")
		fmt.Printf("%s (%s)\n", children[0].Name, children[0].ID)
		fmt.Printf("Press return to continue")
		fmt.Scan()
		return nil, fmt.Errorf("TODO: Impmement single child download")
	}

	var choice int
	for {
		fmt.Println("Found multiple children. Which account to scrape?")
		for idx, child := range children {
			fmt.Printf("%d. %s (%s)\n", idx+1, child.Name, child.ID)
		}
		fmt.Printf("%d. All\n", len(children)+1)
		fmt.Scanf("%d", &choice)
		if choice >= 1 && choice <= len(children)+1 {
			break
		}
	}

	if choice == len(children)+1 {
		return children, nil
	}

	return []himama.Child{children[choice-1]}, nil
}

func mkdir(path string) {
	if err := os.Mkdir(path, 0700); err != nil {
		if !os.IsExist(err) {
			panic(fmt.Errorf("unable to create directory ./%s: %w", path, err))
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	} else if errors.Is(err, fs.ErrNotExist) {
		return false
	}

	panic(err)
}
