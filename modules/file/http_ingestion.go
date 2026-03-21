// modules/file/http_ingestion.go
package file

import (
	"context"
	"io"
	"net/http"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func FileUploadHandler(bus *schema.MessageBus) http.HandlerFunc {

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
				_ = bus.Publish(r.Context(), "file.chunk", buf[:n])
			}
			if err == io.EOF {
				break
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

type IngestionModule struct {
	repo FileRepository
}

func (m *IngestionModule) Handle(ctx context.Context, payload []byte) error {
	return m.repo.StoreChunk(ctx, payload)
}

/*This supports:
• CSV ingestion
• JSON ingestion
• Binary blob storage*/
