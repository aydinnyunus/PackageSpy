package pypi

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/secretScanner/utils"
)

func getPyPIDownloadUrl(packageURL string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pypi.org"+packageURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "pypi.org")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("cookie", "_ga=GA1.2.451339872.1696321050; _ga_RW7D75DF8V=GS1.1.1696321049.1.0.1696321053.0.0.0; _ga_B0F3Y2XW9M=GS1.1.1696321049.1.0.1696321053.0.0.0; session_id=Tzhu4VSTJBL4iAVBD46rLtwsXXhrcJM0WZIovcXs0-c.Za04sw.d4XY4e5OlOP0Ij8zGf6zUR8hlJWyPixZmJuX1f309aCOrG5aZFmRwbbwlYyw_BWxDA1nqFMnKUYG3p39mfaE5w")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return ""
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				// Found an anchor tag with id="package-snippet"
				// Now, let's extract the href attribute
				for _, attr := range token.Attr {
					//fmt.Println(attr)
					if attr.Key == "href" && strings.Contains(attr.Val, "files.pythonhosted.org") {
						return attr.Val
					}
				}

			}
		}
	}
}

func downloadPyPIPackage(url, username string) {

	// Define the directory where you want to save the file

	// Ensure the download directory exists, or create it if it doesn't
	if err := os.MkdirAll("pypi_packages", 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		os.Exit(1)
	}

	downloadDir := "pypi_packages/" + username

	// Ensure the download directory exists, or create it if it doesn't
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		os.Exit(1)
	}

	// Check if there are any parts after splitting
	parts := strings.Split(url, "/")

	latestPart := "unknown.tar.gz"
	if len(parts) > 0 {
		// Get the latest part (last element in the slice)
		latestPart = parts[len(parts)-1]
		fmt.Println("Latest Part:", latestPart)
	} else {
		fmt.Println("No parts found after splitting.")
	}
	// Create the file path for the downloaded file
	filename := filepath.Join(downloadDir, latestPart)
	fmt.Println(filename)
	// Create or open the file for writing
	outFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// Perform the HTTP GET request to download the file
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// Copy the response body to the output file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		fmt.Println("Error copying file:", err)
		os.Exit(1)
	}

	fmt.Printf("File downloaded to %s\n", filename)

	ext := filepath.Ext(filename)
	if ext == ".gz" {
		// It's a tar.gz file, so extract it
		if err := utils.ExtractTarGz(filename, downloadDir+"/"); err != nil {
			fmt.Println("Error extracting tar.gz:", err)
			return
		}
		fmt.Println("Extracted tar.gz successfully.")
	} else if ext == ".whl" {
		// It's a .whl file, so extract it (assuming it's a ZIP archive)
		if err := utils.ExtractWheel(filename, downloadDir+"/"); err != nil {
			fmt.Println("Error extracting .whl:", err)
			return
		}
		fmt.Println("Extracted .whl successfully.")
	} else {
		fmt.Println("Unsupported file type:", ext)
	}

}

func DownloadAllPyPIPackages(username string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pypi.org/user/"+username, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "pypi.org")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("cookie", "_ga=GA1.2.451339872.1696321050; _ga_RW7D75DF8V=GS1.1.1696321049.1.0.1696321053.0.0.0; _ga_B0F3Y2XW9M=GS1.1.1696321049.1.0.1696321053.0.0.0; session_id=Tzhu4VSTJBL4iAVBD46rLtwsXXhrcJM0WZIovcXs0-c.Za04sw.d4XY4e5OlOP0Ij8zGf6zUR8hlJWyPixZmJuX1f309aCOrG5aZFmRwbbwlYyw_BWxDA1nqFMnKUYG3p39mfaE5w")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	//todo @atlaskit/avatar vs avatar check @ and parse.
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return // End of the document
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "class" && attr.Val == "package-snippet" {
						// Found an anchor tag with id="package-snippet"
						// Now, let's extract the href attribute
						for _, attr := range token.Attr {
							if attr.Key == "href" {
								fmt.Println("Href:", "https://pypi.org"+attr.Val)
								downloadUrl := getPyPIDownloadUrl(attr.Val)
								downloadPyPIPackage(downloadUrl, username)
								utils.RunGitleaks("pypi_packages/" + username)

							}
						}
					}
				}
			}
		}
	}
}
