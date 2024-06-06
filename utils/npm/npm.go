package npm

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aydinnyunus/PackageSpy/utils"
)

func DownloadAllNpmPackages(username string) {
	client := &http.Client{}
	page := 0

	for {
		req, err := http.NewRequest("GET", "https://www.npmjs.com/~"+username+"?page="+strconv.Itoa(page), nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("authority", "www.npmjs.com")
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("accept-language", "en-US,en;q=0.9")
		req.Header.Set("cache-control", "no-cache")
		req.Header.Set("cookie", "__cf_bm=xpvYXMPXoYTv5C3fSp8M0e4DSfRGF6au3Vlu3t1qSmc-1706034109-1-AWUVJT1BQ0MeSSdStOGwsGPFTpJTjDDhFHe+7lH3WwJ27oIOybiR7v7hp9URhwczGlWMaWEwhfmGftvLW5Hth6g=")
		req.Header.Set("pragma", "no-cache")
		req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"macOS"`)
		req.Header.Set("sec-fetch-dest", "document")
		req.Header.Set("sec-fetch-mode", "navigate")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("sec-fetch-user", "?1")
		req.Header.Set("upgrade-insecure-requests", "1")
		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		tokenizer := html.NewTokenizer(resp.Body)

		foundPackage := false
		fmt.Println("PAGE NUMBER IS ")
		fmt.Println(page)
	InnerLoop:
		for {
			tokenType := tokenizer.Next()
			switch tokenType {
			case html.ErrorToken:
				if !foundPackage {
					// If no package is found on this page, exit the loop
					extractNpmPackages()
					utils.RunGitleaks("npm_packages")

					return
				}
				page++
				break InnerLoop
			case html.StartTagToken, html.SelfClosingTagToken:
				token := tokenizer.Token()
				if token.Data == "a" {
					for _, attr := range token.Attr {
						if attr.Key == "target" && attr.Val == "_self" {
							// Found an anchor tag with id="package-snippet"
							// Now, let's extract the href attribute
							for _, attr := range token.Attr {
								if attr.Key == "href" && strings.Contains(attr.Val, "package") {
									//todo:
									parts := strings.Split(attr.Val, "/")
									lastPart := ""
									if strings.Contains(parts[2], "@") {
										lastPart = parts[len(parts)-1] + parts[len(parts)-2]
									} else {
										lastPart = parts[len(parts)-1]
									}
									getAllDownloadUrlsNpm(lastPart)
									foundPackage = true

								}
							}
						}
					}
				}
			}
		}
	}

	// Once all pages have been processed, you can add any further processing here
}

func getAllDownloadUrlsNpm(packageName string) {
	cmd := exec.Command("npm", "pack", packageName, "--pack-destination=npm_packages")

	// Set environment variables or other necessary configurations for Docker if needed

	// Execute the command
	_, err := cmd.CombinedOutput()
	fmt.Println(packageName)
	if err != nil {
		fmt.Printf("Error executing npm pack: %v\n", err)
		return
	}

	// Print the npm pack output
	//fmt.Printf("npm pack output:\n%s\n", output)

}

func extractNpmPackages() {
	// Specify the directory containing npm packages with .tgz files.
	directory := "./npm_packages"

	// Create a list of .tgz files in the specified directory.
	tgzFiles, err := utils.FindTgzFiles(directory)
	if err != nil {
		fmt.Println("Error finding .tgz files:", err)
		return
	}

	// Extract each .tgz file.
	for _, tgzFile := range tgzFiles {
		err := utils.ExtractTgzFile(tgzFile, directory+"/"+tgzFile)
		if err != nil {
			fmt.Printf("Error extracting %s: %v\n", tgzFile, err)
		} else {
			fmt.Printf("Extracted %s successfully.\n", tgzFile)
		}
	}
}
