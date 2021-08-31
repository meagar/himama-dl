package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"sync"

	"github.com/meagar/himama-dl/internal/himama"
)

func main() {
	username, password, err := fetchCredentials()
	if err != nil {
		fmt.Println("Error colleting credentials:", err)
		return
	}

	fmt.Println("himama-dl")

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

func scrape(client *himama.Client, child himama.Child) error {
	mkdir("./" + child.Name)

	work := make(chan himama.Activity, 20)

	wg := sync.WaitGroup{}
	for n := 0; n < 20; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for activity := range work {
				res, err := http.Get(activity.MediaURL)
				if err != nil {
					panic(err)
				}
				filename := activity.SuggestedLocalFilename()
				fh, err := os.OpenFile("./"+child.Name+"/"+filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
				if err != nil {
					panic(err)
				}
				if _, err := io.Copy(fh, res.Body); err != nil {
					panic(err)
				}
				res.Body.Close()
				fh.Close()
				fmt.Println(filename)
			}
		}()
	}

	for page := 1; ; page++ {
		activities, err := client.Activities(child, page)
		if err != nil {
			panic(err)
		}
		if len(activities) == 0 {
			break
		}

		for _, activity := range activities {
			work <- activity
		}
	}

	close(work)
	wg.Wait()

	return nil
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
