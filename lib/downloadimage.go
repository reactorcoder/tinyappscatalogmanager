package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Downloadimage(inputFilename string) {
	filename := strings.Replace(inputFilename, "screenshots/", "", -1)

	imageURL := "https://tinyapps.org/screenshots" + filename
	outputPath := filepath.Join("../screenshots", filename)

	// Send an HTTP GET request to the image URL
	response, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	// Check if the request was successful (status code 200)
	if response.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code %d\n", response.StatusCode)
		return
	}

	if err := os.MkdirAll("../screenshots", os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Create a new file to save the image
	print(outputPath)
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Copy the image data from the HTTP response body to the local file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error copying image data:", err)
		return
	}

	fmt.Println("Image downloaded and saved as " + file.Name())
}
