package injector

import (
	"fmt"
	"net/http"
)

func ServeScript(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set(
		"Content-Type",
		"application/javascript",
	)

	fmt.Fprint(w, `const source = new EventSource("/__live");
		source.onmessage = () => {
			location.reload();
		};`)
}
