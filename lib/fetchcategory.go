package lib

import (
	"strings"

	"golang.org/x/net/html"
)

type Item struct {
	MainCategory string
	Category     string
	Screenshot   string
	URI          string
	Name         string
	Info         string
	Size         string
	Site         string
	Downloads    string
	Supreme      string
	Sourcecode   string
	Shareware    string
	Noinstall    string
}

type ItemJsonStruct struct {
	Category    string
	Subcategory string
	Screenshot  string
	URI         string
	Name        string
	Info        string
	Size        string
	Site        string
	Downloads   string
	Supreme     string
	Sourcecode  string
	Shareware   string
	Noinstall   string
}

func ParseWebPage(htmlString string, maincategory string) ([]Item, error) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return nil, err
	}

	var items []Item
	var category string

	var parseNode func(*html.Node)
	parseNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h2" {
			category = ""
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					category += c.Data
				}
			}
			category = strings.TrimSpace(category)
		}

		if n.Type == html.ElementNode && n.Data == "p" && category != "" {
			var item Item
			item.Category = category
			item.MainCategory = maincategory

			// Parse the item URI, name, and info
			parseItemNode(n, &item)

			items = append(items, item)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c)
		}
	}

	parseNode(doc)
	return items, nil
}

func parseItemNode(n *html.Node, item *Item) {
	htmlStringInput := renderNodeToString(n)

	//htmlString := strings.ReplaceAll(htmlStringInput, "<span class=\"icon\">🌱</span> ", "")

	//regex := `<p><a href="([^"]+)">([^<]+)<\/a>\s*\[.*?\]\s*\+?\s*(.*?)\s*<a href="([^"]+)"\s*class="icon">.+?<\/a><\/p>`
	res := parseHTML(htmlStringInput)
	//re := regexp.MustCompile(regex)
	//
	//fmt.Println("-----")
	//fmt.Printf("Supreme: %v\n", res["supreme"])
	//fmt.Printf("Title: %v\n", res["title"])
	//fmt.Printf("URI: %v\n", res["uri"])
	//fmt.Printf("Size: %v\n", res["size"])
	//fmt.Printf("Sourcecode: %v\n", res["sourcecode"])
	//fmt.Printf("Shareware: %v\n", res["shareware"])
	//fmt.Printf("Noinstall: %v\n", res["noinstall"])
	//fmt.Printf("Description: %v\n", res["description"])
	//fmt.Printf("Screenshot: %v\n", res["screenshot"])
	//fmt.Printf("Downloads: %v\n", res["downloads"])
	//fmt.Printf("Site: %v\n", res["site"])
	//fmt.Println("-----")

	//match := re.FindStringSubmatch(htmlString)
	if res != nil {
		item.URI = res["uri"]
		item.Name = res["title"]
		item.Info = res["description"]
		item.Screenshot = res["screenshot"]
		item.Size = res["size"]
		item.Site = res["site"]
		item.Downloads = res["downloads"]

		item.Supreme = res["supreme"]
		item.Sourcecode = res["sourcecode"]
		item.Shareware = res["shareware"]
		item.Noinstall = res["noinstall"]

		if item.Screenshot != "" {
			Downloadimage(item.Screenshot)
		}
	}
}

func renderNodeToString(n *html.Node) string {
	var sb strings.Builder
	html.Render(&sb, n)
	return sb.String()
}
