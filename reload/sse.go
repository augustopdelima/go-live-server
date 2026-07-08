package reload

import (
	"fmt"
	"net/http"
)

func (r *Reloader) HandleSSE(
	writer http.ResponseWriter,
	request *http.Request,
) {
	flusher, ok := writer.(http.Flusher);

	if !ok {
		http.Error(
			writer,
			"SSE not supported",
			http.StatusInternalServerError,
		)

		return
	}

	writer.Header().Set(
		"Content-Type",
		"text/event-stream",
	)

	writer.Header().Set(
		"Cache-Control",
		"no-cache",
	)

	writer.Header().Set(
		"Connection",
		"keep-alive",
	)

	client := r.Subscribe()

	defer r.Unsubscribe(client)

	for {
		select {
			case <-request.Context().Done():
				return

			case <-client:
				fmt.Fprintf(
					writer,
					"data: reload\n\n",
				)

				flusher.Flush()
		}

	}
}
