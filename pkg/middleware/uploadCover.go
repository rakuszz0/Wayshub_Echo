package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func UploadCover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Upload file
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, err := c.FormFile("cover")

		if err != nil && c.Request().Method == "PATCH" {
			c.Set("dataCover", "false")
			return next(c)
		}

		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, "Error Retrieving the File")
		}
		defer file.Close()

		const MAX_UPLOAD_SIZE = 10 << 20 // 10MB
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		c.Request().ParseMultipartForm(MAX_UPLOAD_SIZE)
		if c.Request().ContentLength > MAX_UPLOAD_SIZE {
			response := Result{Code: http.StatusBadRequest, Message: "Max size in 10mb"}
			return c.JSON(http.StatusBadRequest, response)
		}

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("uploads", "image-*.png")
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			return c.JSON(http.StatusInternalServerError, err)
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
		// fileCover := data[8:] // split uploads/

		// add filename to ctx
		c.Set("dataCover", data)
		return next(c)
	}
}
