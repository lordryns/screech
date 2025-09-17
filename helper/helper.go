package helper

import (
    "bytes"
    "image/jpeg"
    "image/png"
    "os"
)

// NOTE: This function is ai generated 

// ConvertJpgToPng converts JPG image bytes to a PNG file at the specified path.
func ConvertJpgToPng(jpgBytes []byte, pngPath string) error {
    // Decode the JPG bytes
    img, err := jpeg.Decode(bytes.NewReader(jpgBytes))
    if err != nil {
        return err
    }

    // Create the output PNG file
    out, err := os.Create(pngPath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Encode as PNG
    return png.Encode(out, img)
}
