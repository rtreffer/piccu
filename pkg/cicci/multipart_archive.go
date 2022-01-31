package cicci

import (
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"
)

func CreateMultipartArchive(files []ExpandedFile) (string, error) {
	buffer := &strings.Builder{}
	multipartWriter := multipart.NewWriter(buffer)
	fileHeader := fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%v\"\nMIME-Version: 1.0\n\n", multipartWriter.Boundary())
	if _, err := buffer.Write([]byte(fileHeader)); err != nil {
		return "", err
	}
	for _, file := range files {
		filetype := "text/cloud-config"
		if file.IsScript {
			filetype = "text/x-shellscript"
		}
		headers := make(textproto.MIMEHeader)
		// TODO: should Content-Transfer-Ecoding be 8bit or binary? shoud Content-Type be ...; charset=UTF-8 ?
		headers.Add("MIME-Version", "1.0")
		headers.Add("Merge-Type", "list(append)+dict(recurse_array)+str()")
		headers.Add("Content-Type", filetype+"; charset=\"utf-8\"")
		headers.Add("Content-Diposition", "attachment; filename=\""+file.Filename+"\"")
		headers.Add("Content-Transfer-Encoding", "7bit")

		partWriter, err := multipartWriter.CreatePart(headers)
		if err != nil {
			return "", fmt.Errorf("can't add %s: %s", file.Filename, err)
		}
		_, err = partWriter.Write([]byte(file.Content))
		if err != nil {
			return "", fmt.Errorf("can't add %s: %s", file.Filename, err)
		}
	}
	err := multipartWriter.Close()
	return buffer.String(), err
}
