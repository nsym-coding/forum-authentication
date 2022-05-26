package userimages

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

//Saves uploaded images from the forum within the userImages folder
func SaveImage(file multipart.File, filename string) {

	// Create the directory to store the images
	root := filepath.Join(".", "/templates/userImages")

	//MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
	//If path is already a directory, MkdirAll does nothing and returns nil.
	folderr := os.MkdirAll(root, os.ModePerm)
	if folderr != nil {
		fmt.Println("Cannot create requested folder")
	}

	fullPath := root + "/" + filename

	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	//copy data to file
	_, errr := io.Copy(f, file)
	if errr != nil {
		fmt.Println("Error writing image to file")
	}

}
