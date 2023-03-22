package middleware

import (
	"context"
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadPhoto(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()

		// Upload file
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, _, err := r.FormFile("photo")

		if err != nil && r.Method == "PATCH" {
			ctx := context.WithValue(r.Context(), "dataPhoto", "false")
			c.SetRequest(r.WithContext(ctx))
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
		const MAX_UPLOAD_SIZE = 10 << 20 // 10MB
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(MAX_UPLOAD_SIZE)
		if r.ContentLength > MAX_UPLOAD_SIZE {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"code":    http.StatusBadRequest,
				"message": "Max size in 10mb",
			})
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
		// filePhoto := data[8:] // split uploads/

		// add filename to ctx
		c.Set("dataPhoto", data)

		return next(c)
	}
}
