package main
import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// Check and create necessary folders
	if err := ensureFoldersExist("uploads", "outputs"); err != nil {
		log.Fatalf("Error creating folders: %s", err)
	}

	// Log the current state of the folders
	log.Println("uploads folder exists:", folderExists("uploads"))
	log.Println("outputs folder exists:", folderExists("outputs"))

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/convert", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			log.Printf("Error getting file: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the file"})
			return
		}

		// Save the file on the server
		filename := filepath.Join("uploads", file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			log.Printf("Error saving file: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
			return
		}

		// Convert the file to SRT format
		srtFilename := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)) + ".srt"
		srtPath := filepath.Join("outputs", srtFilename)
		err = convertToSRT(filename, srtPath)
		if err != nil {
			log.Printf("Error converting file: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert the file"})
			return
		}

		// Download the SRT file
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", srtFilename))
		c.Header("Content-Type", "application/octet-stream")
		c.File(srtPath)
	})

	r.Run(":8080")
}

func convertToSRT(inputFile string, outputFile string) error {
	input, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	i := 1

	for scanner.Scan() {
		line := scanner.Text()

		// Generate automatic time intervals (4 seconds per line)
		start := (i - 1) * 4000
		end := i * 4000

		// Convert milliseconds to time format (HH:MM:SS,mmm)
		startTime := fmt.Sprintf("%02d:%02d:%02d,%03d", start/3600000, (start%3600000)/60000, ((start%60000)/1000), start%1000)
		endTime := fmt.Sprintf("%02d:%02d:%02d,%03d", end/3600000, (end%3600000)/60000, ((end%60000)/1000), end%1000)

		// SRT time format (00:00:00,000 --> 00:00:00,000)
		srtTime := fmt.Sprintf("%s --> %s", startTime, endTime)

		// Create an SRT format subtitle entry
		srtEntry := fmt.Sprintf("%d\n%s\n%s\n\n", i, srtTime, line)

		// Write the subtitle entry to the output file
		output.WriteString(srtEntry)

		i++
	}

	log.Println("Conversion completed successfully") // Add a log message when the conversion is finished

	return nil
}

func ensureFoldersExist(folders ...string) error {
	for _, folder := range folders {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			if err := os.MkdirAll(folder, 0755); err != nil {
				return fmt.Errorf("error creating folder %s: %s", folder, err)
			}
		}
	}
	return nil
}

func folderExists(folder string) bool {
	_, err := os.Stat(folder)
	return !os.IsNotExist(err)
}
