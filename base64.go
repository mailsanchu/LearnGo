package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	var input string
	input = "/home/svarkey/Pictures/Garden_Alteration.png"
	//outFile = os.Args[2]

	imgFile, err := os.Open(input) // a QR code image
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	// if you create a new image instead of loading from file, encode the image to buffer instead with png.Encode()

	// png.Encode(&buf, image)

	// convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)


	file, err := os.Create("/tmp/1.html")
	if err != nil {
		panic("Cannot create file")
	}
	defer file.Close()
	fmt.Fprintf(file, "<html><body><img src=\"data:image/png;base64," + imgBase64Str + "\" /></body></html>")


}