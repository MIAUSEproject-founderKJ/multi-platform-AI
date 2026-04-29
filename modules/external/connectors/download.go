// modules/external/connectors/download.go
package external_connectors

import (
	"io"
	"net/http"
	"os"
)

func DownloadModule(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, _ := os.Create(dest)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
