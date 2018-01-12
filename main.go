package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

// License is a license to apply
type License interface {
	LicenseText() string
	FileHeader() string
}

var selectedLicense string
var listLicenses bool
var availableLicenses []License
var owner string

const filePattern = "*.go"

func init() {
	flag.StringVar(&selectedLicense, "license", "apache2.0", "The selected license")
	flag.StringVar(&owner, "owner", "", "Name of license owner")
	flag.Parse()

	availableLicenses = append(availableLicenses, NewApache20License("2018", owner))
}

func addCommentBlock(licenseText string) string {
	return "//" + strings.Replace(licenseText, "\n", "\n//", -1)
}

func addHeaderToFile(f *os.File, absPath string, license License) {
	var output []byte
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Printf("Unable to read file %s: %v\n", f.Name(), err)
			return
		}
		output = append(output, []byte(scanner.Text()+"\n")...)
		if strings.HasPrefix(scanner.Text(), "package") {
			output = append(output, []byte(addCommentBlock(license.FileHeader()))...)
		}
	}
	if err := os.Remove(absPath); err != nil {
		fmt.Printf("Unable to remove original file %s: %v\n", absPath, err)
		return
	}
	newFile, err := os.Create(absPath)
	if err != nil {
		fmt.Printf("Got error creating %s: %v\n", absPath, err)
		return
	}
	_, err = newFile.Write(output)
	if err != nil {
		fmt.Printf("Got error writing to %s: %v\n", absPath, err)
		return
	}
	newFile.Close()
}
func applyLicenseToFiles(f *os.File, dirPath string, license License) {
	files, err := f.Readdir(-1)
	if err != nil {
		fmt.Printf("Unable to read directory %s: %v\n", f.Name(), err)
		return
	}
	fmt.Println("Reading directory ", f.Name())
	for _, v := range files {
		filePath := path.Join(dirPath, v.Name())
		item, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Unable to open %s: %v\n", filePath, err)
			continue
		}
		if v.IsDir() {
			applyLicenseToFiles(item, filePath, license)
		} else {
			if path.Ext(v.Name()) == ".go" {
				addHeaderToFile(item, filePath, license)
			}
		}
		item.Close()
	}
}

func main() {
	currentLicense := availableLicenses[0]
	// Do this the idiot way and not do anything clever. Create the LICENSE
	// file at the root, find all source files, look for the line that starts
	// with "package" and insert the header right below it.

	_, err := os.Stat("LICENSE")
	if err == nil {
		fmt.Println("There is already a LICENSE file in the current directory")
		os.Exit(1)
	}
	license, err := os.Create("LICENSE")
	if err != nil {
		fmt.Println("Unable to create LICENSE file: ", err)
		os.Exit(2)
	}
	if _, err := license.Write([]byte(currentLicense.LicenseText())); err != nil {
		fmt.Println("Unable to write the license text to LICENSE: ", err)
		os.Exit(3)
	}
	license.Close()

	curDir, err := os.Open(".")
	if err != nil {
		fmt.Println("Unable to read current directory: ", err)
		os.Exit(4)
	}
	applyLicenseToFiles(curDir, ".", currentLicense)
}
