package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func parseHTML(html string) map[string]string {
	result := make(map[string]string)

	// Define regular expressions to match patterns
	// <span class="icon">🌱</span>
	patternSupreme := `<span class="icon">(.*?)</span>`
	//patternURI := `<p>(?:<a href="(.*?)">(.*?)<\/a>)|<p>(.*)\[`
	patternTitle := `<p>(?:<a href="(.*?)">(.*?)<\/a>)|<p>(.*)\[`
	patternSize := `\[(.*?)\]`
	patternSourcecode := `\ (?:([+{S}$]{1,}))\ `
	patternShareware := `\ (?:([+{S}$]{1,}))\ `
	patternNoinstall := `\ (?:([+{S}$]{1,}))\ `
	patternDescription := `(?m)\] (.*?)\ (.*?)(<a href="\/|<a class='icon' href='\/)`
	patternScreenshot := `(/screenshots/(.*?))"`
	patternDownloads := `(/downloads/(.*?))"`
	patternSite := `<a href="([^"]*)" class="icon">🌎</a>`

	// Compile regular expressions
	regexSupreme := regexp.MustCompile(patternSupreme)
	//regexURI := regexp.MustCompile(patternURI)
	regexTitle := regexp.MustCompile(patternTitle)
	regexSize := regexp.MustCompile(patternSize)
	regexSourcecode := regexp.MustCompile(patternSourcecode)
	regexShareware := regexp.MustCompile(patternShareware)
	regexNoinstall := regexp.MustCompile(patternNoinstall)
	regexDescription := regexp.MustCompile(patternDescription)
	regexScreenshot := regexp.MustCompile(patternScreenshot)
	regexDownloads := regexp.MustCompile(patternDownloads)
	regexSite := regexp.MustCompile(patternSite)

	// Extract information using regular expressions
	result["supreme"] = ""
	if matches := regexSupreme.FindStringSubmatch(html); len(matches) > 1 {
		result["supreme"] = matches[1]
	}

	result["uri"] = ""
	result["title"] = ""
	htmltitle := strings.ReplaceAll(html, "<span class=\"icon\">🌱</span> ", "")
	if matches := regexTitle.FindStringSubmatch(htmltitle); len(matches) > 0 {
		result["uri"] = matches[1]
		if len(matches) >= 3 && matches[2] != "" {
			result["title"] = matches[2]
		} else if len(matches) >= 4 && matches[3] != "" {
			result["title"] = matches[3]
		}
	}

	result["size"] = ""
	if matches := regexSize.FindStringSubmatch(html); len(matches) > 1 {
		result["size"] = matches[1]
	}

	result["sourcecode"] = "false"
	if matches := regexSourcecode.FindStringSubmatch(html); len(matches) > 1 {
		if strings.Contains(matches[1], "{S}") {
			result["sourcecode"] = "true"
		}
	}

	result["shareware"] = "false"
	if matches := regexShareware.FindStringSubmatch(html); len(matches) > 1 {
		if strings.Contains(matches[1], "$") {
			result["shareware"] = "true"
		}
	}

	result["noinstall"] = "false"
	if matches := regexNoinstall.FindStringSubmatch(html); len(matches) == 1 {
		if strings.Contains(matches[1], "+") {
			result["noinstall"] = "true"
		}
	}

	result["description"] = ""
	if matches := regexDescription.FindStringSubmatch(html); len(matches) > 1 {
		result["description"] = matches[2]
	}

	result["screenshot"] = ""
	if matches := regexScreenshot.FindStringSubmatch(html); len(matches) > 1 {
		result["screenshot"] = matches[1]
	}

	result["downloads"] = ""
	if matches := regexDownloads.FindStringSubmatch(html); len(matches) > 1 {
		result["downloads"] = matches[1]
	}

	result["site"] = ""
	if matches := regexSite.FindStringSubmatch(html); len(matches) > 1 {
		result["site"] = matches[1]
	}

	return result
}

// ParseAndStoreItems fetches the HTML content from the given URL, extracts the list items, and stores them in a file.
func ParseAndStoreItems(url string, filePath string) error {
	// Send an HTTP GET request to fetch the HTML content
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch the HTML content: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read the response body: %v", err)
	}

	// Parse the HTML content
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Extract the list items
	var items []string
	var extractItems func(*html.Node)
	extractItems = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && strings.HasPrefix(attr.Val, "/") {
					items = append(items, n.FirstChild.Data)
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractItems(c)
		}
	}
	extractItems(doc)

	// Store the items to a file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	for _, item := range items {
		_, err := file.WriteString(item + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	fmt.Println("Items have been stored to", filePath)

	return nil
}
