package middleware

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadThumbnail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Upload file
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, _, err := c.Request().FormFile("thumbnail")

		if err != nil && c.Request().Method == "PATCH" {
			ctx := context.WithValue(c.Request().Context(), "dataThumbnail", "false")
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}

		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, "Error Retrieving the File")
		}
		defer file.Close()
		// fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		// fmt.Printf("File Size: %+v\n", handler.Size)
		// fmt.Printf("MIME Header: %+v\n", handler.Header)
		const MAX_UPLOAD_SIZE = 100 << 20 // 10MB
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		c.Request().ParseMultipartForm(MAX_UPLOAD_SIZE)
		if c.Request().ContentLength > MAX_UPLOAD_SIZE {
			return c.JSON(http.StatusBadRequest, Result{Code: http.StatusBadRequest, Message: "Max size in 10mb"})
		}

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("uploads", "image-*.png")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		// write this byte array to our temporary file
		tempFile.Write(fileBytes)

		data := tempFile.Name()
		// fileThumbnail := data[8:] // split uploads/

		// add filename to ctx

		c.Set("dataThumbnail", data)

		return next(c)
	}
}
