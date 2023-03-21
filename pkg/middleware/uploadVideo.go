package middleware

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadVideo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Upload file
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, _, err := c.Request().FormFile("video")

		if err != nil && c.Request().Method == "PATCH" {
			ctx := context.WithValue(c.Request().Context(), "dataVideo", "false")
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
		const MAX_UPLOAD_SIZE = 100 << 20 // 100MB
		// Parse our multipart form, 100 << 20 specifies a maximum
		// upload of 100 MB files.
		c.Request().ParseMultipartForm(MAX_UPLOAD_SIZE)
		if c.Request().ContentLength > MAX_UPLOAD_SIZE {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"code":    http.StatusBadRequest,
				"message": "Max size in 100mb",
			})
		}

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("uploads", "video-*.mp4")
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
		// fileVideo := data[8:] // split uploads/

		// add filename to ctx
		ctx := context.WithValue(c.Request().Context(), "dataVideo", data)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
