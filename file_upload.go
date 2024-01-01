package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type UploadResponse struct {
	Status string `json:"status"`
	Data   struct {
		URL string `json:"url"`
	} `json:"data"`
}

func handleFileUpload(b *gotgbot.Bot, ctx *ext.Context) error {
	filesize := ctx.EffectiveMessage.Document.FileSize
	filename := ctx.EffectiveMessage.Document.FileName
	if filesize > 100000000 {
		ctx.EffectiveMessage.Reply(b, "<b>⚠️ Request Entity Too Large. 🚀 Max Upload File Size Limit: 100MB</b>", &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return nil
	} else {
		msg, _ := ctx.EffectiveMessage.Reply(b, "⬇️ Downloading to host. Please wait... ... ", &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		file, err := b.GetFile(ctx.EffectiveMessage.Document.FileId, &gotgbot.GetFileOpts{})
		if err != nil {
			return err
		}
		url := file.URL(b, &gotgbot.RequestOpts{})

		folderName := strconv.Itoa(int(ctx.EffectiveMessage.MessageId))
		outputPath := folderName + "/" + filename
		os.Mkdir(folderName, 0755)
		savedPath, err := downloadFile(url, outputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil
		}
		msg.EditText(b, "✅ Download completed. Uploading to tmpfiles.org", &gotgbot.EditMessageTextOpts{})

		response, err := uploadToTmpFiles(savedPath)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		cleanDirectory(folderName)
		directLink, _ := getDirectLink(response.Data.URL)
		shareURL := fmt.Sprintf("https://telegram.me/share/url?url=%v&text=%v", directLink, filename)
		inlineButtons := gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					gotgbot.InlineKeyboardButton{
						Text: "⚡ __Direct Download Link__",
						Url:  directLink,
					},
					gotgbot.InlineKeyboardButton{
						Text: "🌎 __Download Page__",
						Url:  response.Data.URL,
					},
				},

				{
					gotgbot.InlineKeyboardButton{
						Text: "↪️ __Share__",
						Url:  shareURL,
					},
				},
			},
		}

		if response.Status == "success" {
			msg.EditText(b, "<b>✅ Upload Succesfull\n⚡ Direct Download Link: \n</b>"+directLink, &gotgbot.EditMessageTextOpts{
				ReplyMarkup:           inlineButtons,
				ParseMode:             "html",
				DisableWebPagePreview: true,
			})
		}
	}

	return nil

}

func downloadFile(url, outputPath string) (string, error) {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file. status code: %d", response.StatusCode)
	}

	buffer := make([]byte, 1<<20)

	for {

		n, err := response.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return "", err
		}

		if _, err := outputFile.Write(buffer[:n]); err != nil {
			return "", err
		}

		if err == io.EOF {
			break
		}
	}

	return outputPath, nil
}

func uploadToTmpFiles(filePath string) (*UploadResponse, error) {
	uploadURL := "https://tmpfiles.org/api/v1/upload"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024*1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		_, err = part.Write(buffer[:n])
		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uploadURL, &requestBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to upload file. status: %v", response.Status)
	}

	var uploadResponse UploadResponse
	err = json.NewDecoder(response.Body).Decode(&uploadResponse)
	if err != nil {
		return nil, err
	}
	return &uploadResponse, nil
}

func getDirectLink(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}
	directLink := "https://" + parsedURL.Host + "/dl" + parsedURL.Path
	return directLink, nil
}
