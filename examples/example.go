package main

import (
	"os"
	"fmt"
	"image/png"
	"src.techknowlogick.com/wavatars"
)

func main() {
	w := wavatars.New([]byte{0,1,2,3})
	f, err := os.OpenFile("wavatar.png", os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()
    png.Encode(f, w)
}
