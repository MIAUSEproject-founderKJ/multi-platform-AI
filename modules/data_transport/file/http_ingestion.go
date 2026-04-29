// modules/data_transport/file/http_ingestion.go

package transport_file

import (
	"context"
	"io"
	"net/http"

	runtime_bus "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/bus"
)

func FileUploadHandler(bus *runtime_bus.MessageBus) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer file.Close()

		buf := make([]byte, 64*1024)

		for {
			n, err := file.Read(buf)

			if n > 0 {
				bus.Publish(runtime_bus.Message{
					Topic: "file.chunk",
					Data:  buf[:n],
				})
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

type FileRepository interface {
	StoreChunk(ctx context.Context, data []byte) error
}
