package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	cp "github.com/otiai10/copy"
	"image/color"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	lib "tinyappcatalogmanager/lib"
	"tinyappcatalogmanager/ui"
)

func getPageAsString(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

type Item struct {
	Category    string "category"
	Subcategory string "subcategory"
	Screenshot  string "screenshot"
	URI         string "uri"
	Name        string "name"
	Info        string "info"
	Size        string "size"
	Site        string "site"
	Downloads   string "downloads"
	Supreme     string "supreme"
	Sourcecode  string "sourcecode"
	Shareware   string "shareware"
	Noinstall   string "noinstall"
}

var allItems []Item

func convertLibItemToItem(libItem lib.Item, subcat string) Item {
	return Item{
		Category:    subcat,
		Subcategory: libItem.Category,
		Screenshot:  libItem.Screenshot,
		URI:         libItem.URI,
		Name:        libItem.Name,
		Info:        libItem.Info,
		Size:        libItem.Size,
		Site:        libItem.Site,
		Downloads:   libItem.Downloads,
		Supreme:     libItem.Supreme,
		Sourcecode:  libItem.Sourcecode,
		Shareware:   libItem.Shareware,
		Noinstall:   libItem.Noinstall,
	}
}

func parseCategories(cat string) {
	url := "https://tinyapps.org/"

	pageString, err := getPageAsString(url + "/" + cat + ".html")
	if err != nil {
		fmt.Println("Error reading page:", err)
		return
	}

	// Parse page
	libItems, err := lib.ParseWebPage(pageString, cat)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	for _, libItem := range libItems {
		item := convertLibItemToItem(libItem, cat)
		allItems = append(allItems, item)
	}

	//for _, item := range items {
	//	fmt.Println("Category:", cat)
	//	fmt.Println("Subcategory:", item.Category)
	//	fmt.Println("Screenshot:", item.Screenshot)
	//	fmt.Println("URI:", item.URI)
	//	fmt.Println("Name:", item.Name)
	//	fmt.Println("Info:", item.Info)
	//	fmt.Println()
	//}
}

func saveAllToJson() {
	file, err := os.Create("items.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(allItems)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	fmt.Println("Items saved as items.json")
}

func revertFallbackData() {
	cp.Copy("data/latest/items.json", "items.json")
	cp.Copy("data/latest/items.txt", "items.txt")
	cp.Copy("data/latest/screenshots", "../screenshots/")
}

func cleanFiles() {
	// cp.Copy("data/latest/items.txt", "items.txt")
	os.Remove("items.json")
	os.RemoveAll("../screenshots")
	os.Create("items.json")
	os.Mkdir("../screenshots", os.ModePerm)
}

func readJsonItems(file string) []lib.ItemJsonStruct {
	items, err := lib.ReadJSONFromFile(file)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return items
}

func fileExistsAndNotEmpty(filename string) (bool, error) {
	fileInfo, err := os.Stat(filename)

	if os.IsNotExist(err) {
		// File does not exist
		return false, nil
	} else if err != nil {
		// Error occurred while checking file
		return false, err
	}

	// Check if it's a regular file and not a directory
	if fileInfo.Mode().IsRegular() && fileInfo.Size() > 0 {
		return true, nil
	}

	// File exists but is empty
	return false, nil
}

func buildMenu(win fyne.Window, app fyne.App) *fyne.MainMenu {
	nt2 := fyne.NewMenuItem("Update Latest", func() {
		progressBar := widget.NewProgressBar()
		progressBar.Max = 100

		progressLabel := widget.NewLabel("")
		updatingLabel := widget.NewLabel("Please wait. Fetching latest catalog from Tinyapps.org")
		closeButton := widget.NewButton("Close App", func() {
			win.Close()
		})

		closeButton.Disabled()
		progressBarContainer := container.NewVBox(
			updatingLabel,
			progressBar,
			progressLabel,
			closeButton,
		)
		dialog.ShowCustom("Updating catalog", "", progressBarContainer, win)

		go func() {
			progressBar.SetValue(1)
			cleanFiles()
			progressBar.SetValue(2)

			// fetchApi Start
			parsedData, err := lib.ReadAndParseFile("items.txt", "\r\n")
			progressBar.SetValue(3)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			for i, line := range parsedData {
				progressBar.SetValue(float64(i) / float64(len(parsedData)) * 100)
				parseCategories(line[0])
			}
			progressBar.SetValue(99)
			saveAllToJson()
			progressBar.SetValue(100)
			// fetchApi End

			progressBar.SetValue(100)
			progressLabel.SetText("Catalog updated successfully. Please close and re-run application")
			updatingLabel.SetText("Finished")
			closeButton.Enable()
		}()
	})
	nt3 := fyne.NewMenuItem("Revert Fallback Locally", func() {
		revertFallbackData()
		dialog.ShowInformation("Updating catalog", "Catalog updated from local database and loaded successfully.", win)
	})
	manageStore := fyne.NewMenu("Manage Catalog", nt2, nt3)

	themedef := fyne.NewMenuItem("Tinyapps Default", func() {
		app.Settings().SetTheme(ui.TinyappsUi())
		//app.Preferences().SetString("theme", "default")
	})
	thememin := fyne.NewMenuItem("Tinyapps Minimal", func() {
		app.Settings().SetTheme(ui.TinyappsUiRecessed())
		//app.Preferences().SetString("theme", "min")
	})
	def := fyne.NewMenuItem("System", func() {
		app.Settings().SetTheme(theme.DefaultTheme())
	})
	themeMenus := fyne.NewMenuItem("Theme", nil)
	themeMenus.ChildMenu = fyne.NewMenu(
		"",
		themedef,
		thememin,
		def,
	)
	filterMenus := fyne.NewMenuItem("Filter Apps", nil)

	// No filter
	filterMenuNoFilterLabel := "✔  No filter (all apps) "
	filterMenuNoFilter := fyne.NewMenuItem(filterMenuNoFilterLabel, func() {
		app.Preferences().SetString("filter-supreme", "0")
		app.Preferences().SetString("filter-shareware", "0")
		app.Preferences().SetString("filter-noinstall", "0")
		app.Preferences().SetString("filter-run", "0")
		if app.Preferences().String("filter-no") == "1" {
			app.Preferences().SetString("filter-no", "0")
			filterMenuNoFilterLabel = "   No filter"
		} else {
			app.Preferences().SetString("filter-no", "1")
			filterMenuNoFilterLabel = "✔  No filter"
		}
	})

	// filter supreme
	filterMenuSupremeLabel := "  Supreme 🌱 "
	filterSupreme := fyne.NewMenuItem(filterMenuSupremeLabel, func() {
		app.Preferences().SetString("filter-no", "0")
		if app.Preferences().String("filter-supreme") == "1" {
			app.Preferences().SetString("filter-supreme", "0")
			filterMenuSupremeLabel = "   Supreme 🌱"
		} else {
			app.Preferences().SetString("filter-supreme", "1")
			filterMenuSupremeLabel = "✔  Supreme 🌱"
		}
	})

	filterMenuSharewareLabel := "  Shareware 💰 "
	filterShareware := fyne.NewMenuItem(filterMenuSharewareLabel, func() {
		app.Preferences().SetString("filter-no", "0")
		if app.Preferences().String("filter-shareware") == "1" {
			app.Preferences().SetString("filter-shareware", "0")
			filterMenuSharewareLabel = "   Shareware 💰"
		} else {
			app.Preferences().SetString("filter-shareware", "1")
			filterMenuSharewareLabel = "✔  Shareware 💰"
		}
	})

	filterMenuNoinstallLabel := "  No install ▶ "
	filterNoinstall := fyne.NewMenuItem(filterMenuNoinstallLabel, func() {
		app.Preferences().SetString("filter-no", "0")
		if app.Preferences().String("filter-noinstall") == "1" {
			app.Preferences().SetString("filter-noinstall", "0")
			filterMenuNoinstallLabel = "   No install ▶"
		} else {
			app.Preferences().SetString("filter-noinstall", "1")
			filterMenuNoinstallLabel = "✔  No install ▶"
		}
	})

	filterMenuRunLabel := "  Source code 📄 "
	filterRun := fyne.NewMenuItem(filterMenuRunLabel, func() {
		app.Preferences().SetString("filter-no", "0")
		if app.Preferences().String("filter-run") == "1" {
			app.Preferences().SetString("filter-run", "0")
			filterMenuRunLabel = "   Source code 📄"
		} else {
			app.Preferences().SetString("filter-run", "1")
			filterMenuRunLabel = "✔  Source code 📄"
		}
	})

	filterMenus.ChildMenu = fyne.NewMenu(
		"",
		filterMenuNoFilter,
		filterSupreme,
		filterShareware,
		filterNoinstall,
		filterRun,
	)
	settings := fyne.NewMenu("Settings", themeMenus, filterMenus)

	website := fyne.NewMenuItem("Open TinyApps Official Site", func() {
		u, err := url.Parse("https://tinyapps.org/")
		if err != nil {
			fmt.Println("error", err)
		}
		app.OpenURL(u)
	})
	websiteapp := fyne.NewMenuItem("Open Catalog Manager", func() {
		u, err := url.Parse("https://github.com/reactorcoder/tinyappscatalogmanager")
		if err != nil {
			fmt.Println("error", err)
		}
		app.OpenURL(u)
	})
	aboutitem := fyne.NewMenuItem("About Catalog Apps", func() {
		// imagePath := filepath.Join("../", selectedApp.Screenshot)
		//			emptyCanvas = canvas.NewImageFromFile(imagePath)
		//			emptyCanvas.SetMinSize(fyne.Size{
		//				Width:  200,
		//				Height: 200,
		//			})
		textCanvas := canvas.NewText("Developed freely to fetch lists of apps from tinyapps.org as catalog apps manager.\n\n\n\n "+
			"Lists of apps are maintained by tinyapps.org. \n\n "+
			"Apps catalog manager are available source from GitHub: https://github.com/reactorcoder/tinyappscatalogmanager/tree/main \n\n"+
			"2023 \n "+
			"Version: 1.0", color.White)

		dialog.ShowInformation("About Catalog Apps", textCanvas.Text, win)
	})
	ab := fyne.NewMenu("About", website, websiteapp, aboutitem)

	m := fyne.NewMainMenu(manageStore, settings, ab)
	return m
}

var selFilter = widget.NewLabel("")

func main() {
	a := app.NewWithID("app.codervio.tinyapps.catalogapp")

	myApp := app.New()

	myApp.Settings().SetTheme(ui.TinyappsUi())
	myApp.Preferences().SetString("filter-no", "1")
	myApp.Preferences().SetString("filter-supreme", "0")
	myApp.Preferences().SetString("filter-shareware", "0")
	myApp.Preferences().SetString("filter-noinstall", "0")
	myApp.Preferences().SetString("filter-run", "0")
	myWindow := myApp.NewWindow("TinyApps - Apps Catalog Manager")
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.SetMainMenu(buildMenu(myWindow, myApp))
	iconPath := "img/favicon-tbn.png"
	iconContents, err := os.ReadFile(iconPath)
	if err != nil {
	} else {
		iconRes := fyne.NewStaticResource("tinyapps.org app catalog", iconContents)
		myWindow.SetIcon(iconRes)
	}

	widgetPane := container.New(layout.NewMaxLayout())

	selectedAppPane := container.New(layout.NewVBoxLayout())

	var emptyCanvas *canvas.Image

	var selTab string

	leftContent := container.NewMax()
	//rightContent := container.NewMax()

	exists, err := fileExistsAndNotEmpty("items.json")
	if err != nil {
		// error
		dialog.ShowInformation("Updating catalog data",
			"There are error to update catalog data: "+fmt.Sprint(err),
			myWindow)
	} else if exists {
		// skip, exists, not empty
	} else {
		// not exists/empty
		revertFallbackData()
		dialog.ShowInformation("Updating catalog data",
			"Tinyapps catalog data has been updated from locally. You can run Manage Catalog > Update Latest to get latest updates.",
			myWindow)
	}

	items := readJsonItems("items.json")

	uniqueCategories := lib.GetUniqueCategories(items)

	tabs := container.NewAppTabs()
	//tabsleft := container.NewAppTabs()
	//tabsInside := container.NewAppTabs()
	//appslists := container.NewHBox()

	var itemCategories []string

	// Left pane cat lists
	for _, category := range uniqueCategories {
		cat := category
		itemCategories = append(itemCategories, cat)

		// Tabs
		//uniqueCategories := lib.FilterItemsAndUniqueSubcategories(items, cat)
		//tabs := container.NewAppTabs()
		//for _, subcat := range uniqueCategories {
		//tabsleft.Append(container.NewTabItem(subcat, widget.NewLabel(subcat)))
		//}
		//tabs.SetTabLocation(container.TabLocationTop)
	}

	// Category on left pane
	listCats := widget.NewList(
		func() int {
			return len(itemCategories)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(itemCategories[i])
		},
	)

	// Event on left click category
	listCats.OnSelected = func(id widget.ListItemID) {
		selectedItem := itemCategories[id]

		selTab = selectedItem

		uniqueCategories := lib.FilterItemsAndUniqueSubcategories(items, selectedItem)
		tabs.Items = nil
		tabs.Append(container.NewTabItem("", widget.NewLabel("")))
		for _, subcat := range uniqueCategories {
			// Render name inside, make tabs on right
			tabs.Append(container.NewTabItem(subcat, widget.NewLabelWithStyle(subcat, fyne.TextAlign(fyne.TextAlignCenter), fyne.TextStyle{
				Bold:      true,
				Italic:    false,
				Monospace: false,
				Symbol:    false,
				TabWidth:  1,
			})))
		}
		tabs.DisableIndex(0)
		//tabs.SetTabLocation(container.TabLocationTop)

		widgetPane.RemoveAll()
	}

	// On clicked right top tab show lists of apps
	tabs.OnSelected = func(ti *container.TabItem) {
		widgetPane.RemoveAll()

		filterNo := myApp.Preferences().String("filter-no") == "1"
		filterSupreme := myApp.Preferences().String("filter-supreme") == "1"
		filterShareware := myApp.Preferences().String("filter-shareware") == "1"
		filterNoinstall := myApp.Preferences().String("filter-noinstall") == "1"
		filterSourcecode := myApp.Preferences().String("filter-run") == "1"
		appLists := lib.FilterItemsByCategoryAndSubcategoryFiltered(
			items,
			selTab,
			ti.Text,
			filterNo,
			filterSupreme,
			filterShareware,
			filterNoinstall,
			filterSourcecode)

		if filterSupreme || filterShareware || filterNoinstall || filterSourcecode {
			activeFilter := ""
			if filterSupreme {
				activeFilter += "🌱"
			}
			if filterShareware {
				activeFilter += "💰"
			}
			if filterNoinstall {
				activeFilter += "▶"
			}
			if filterSourcecode {
				activeFilter += "📄"
			}

			selFilter.SetText("FILTER ACTIVE: " + activeFilter)
		} else {
			selFilter.SetText("")
		}

		listwidget := widget.NewList(
			func() int {
				return len(appLists)
			},
			func() fyne.CanvasObject {
				return container.NewVBox(
					widget.NewLabelWithStyle("", fyne.TextAlign(fyne.TextWrapOff), fyne.TextStyle{
						Bold:      true,
						Italic:    false,
						Monospace: false,
						Symbol:    false,
						TabWidth:  0,
					}),
					widget.NewLabelWithStyle("", fyne.TextAlign(fyne.TextAlignLeading), fyne.TextStyle{
						Bold:      false,
						Italic:    false,
						Monospace: false,
						Symbol:    false,
						TabWidth:  0,
					}),
				)
			},
			func(index int, item fyne.CanvasObject) {
				listItem := appLists[index]

				if vbox, ok := item.(*fyne.Container); ok {
					labels := vbox.Objects
					if len(labels) >= 2 {
						if nameLabel, ok := labels[0].(*widget.Label); ok {
							iconsList := ""
							if listItem.Supreme != "" {
								iconsList += "🌱"
							}
							if listItem.Shareware != "" {
								iconsList += "💰"
							}
							if listItem.Noinstall != "" {
								iconsList += "▶"
							}
							if listItem.Sourcecode != "" {
								iconsList += "📄"
							}

							nameLabel.SetText(fmt.Sprintf("%s   %s", listItem.Name, iconsList))
						}
						if valueLabel, ok := labels[1].(*widget.Label); ok {
							valueLabel.SetText(listItem.Info)
							valueLabel.Wrapping = fyne.TextWrapOff
						}
					}
				}
			})

		listwidget.OnSelected = func(id widget.ListItemID) {
			selectedApp := appLists[id]

			selectedAppPane.RemoveAll()

			imagePath := filepath.Join("../", selectedApp.Screenshot)
			emptyCanvas = canvas.NewImageFromFile(imagePath)
			emptyCanvas.SetMinSize(fyne.Size{
				Width:  200,
				Height: 200,
			})

			//selectedAppPane.Add(widget.NewLabel(selectedApp.URI))
			selectedAppPane.Add(widget.NewLabel(selectedApp.Name))
			selectedAppPane.Add(container.NewMax(emptyCanvas))
			selectedAppPane.Add(widget.NewLabel("Size: " + selectedApp.Size))
			selectedAppPane.Add(widget.NewButton("Open", func() {
				u, err := url.Parse(selectedApp.URI)
				if err != nil {
					fmt.Println("error", err)
				}
				a.OpenURL(u)
			}))

			var dwbtn = widget.NewButton("", nil)
			if selectedApp.Downloads != "" {
				dwbtn.SetText("Download")
				dwbtn.OnTapped = func() {
					u, err := url.Parse("https://tinyapps.org" + selectedApp.Downloads)
					if err != nil {
						fmt.Println("Error parsing URL:", err)
						// Optionally, you can display an error message to the user here
						return
					}

					if err := a.OpenURL(u); err != nil {
						fmt.Println("Error opening URL:", err)
						// Optionally, you can display an error message to the user here
					}
				}
			}
			selectedAppPane.Add(dwbtn)

			// Imagesearch btn
			var imgsrcbtn = widget.NewButton("", nil)
			if selectedApp.Screenshot != "" {
				imgsrcbtn.SetText("Image Search")
				imgsrcbtn.OnTapped = func() {
					u, err := url.Parse("https://lens.google.com/uploadbyurl?url=https://tinyapps.org" + selectedApp.Screenshot)
					if err != nil {
						fmt.Println("Error parsing URL:", err)
						// Optionally, you can display an error message to the user here
						return
					}

					if err := a.OpenURL(u); err != nil {
						fmt.Println("Error opening URL:", err)
						// Optionally, you can display an error message to the user here
					}
				}
			}
			selectedAppPane.Add(imgsrcbtn)

			// Site btn
			var sitebtn = widget.NewButton("", nil)
			if selectedApp.Site != "" {
				sitebtn.SetText("Site")
				sitebtn.OnTapped = func() {
					u, err := url.Parse(selectedApp.Site)
					if err != nil {
						fmt.Println("Error parsing URL:", err)
						// Optionally, you can display an error message to the user here
						return
					}

					if err := a.OpenURL(u); err != nil {
						fmt.Println("Error opening URL:", err)
						// Optionally, you can display an error message to the user here
					}
				}
			}
			selectedAppPane.Add(sitebtn)
			selectedAppPane.Add(selFilter)
		}

		widgetPane.Add(listwidget)
	}

	// Left content
	leftpanes := container.NewVSplit(listCats, selectedAppPane)
	leftContent.Add(leftpanes)
	leftpanes.Offset = 0.6

	wpane := container.NewVSplit(tabs, widgetPane)
	wpane.Offset = 0.1

	ri := container.NewBorder(
		container.NewMax(widget.NewSeparator()), nil, nil, nil, wpane,
	)

	res := container.NewHSplit(leftContent, ri)
	res.Offset = 0.1

	myWindow.SetContent(res)
	myWindow.SetMaster()
	myWindow.ShowAndRun()
}
