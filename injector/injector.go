package injector

import "strings"

const reloadScriptTag = `<script src="/__live_reload.js"></script>`

func Inject(html string) string {
	if strings.Contains(html, "</body>") {
		return strings.Replace(html, "</body>", reloadScriptTag+"</body>", 1)
	}

	if strings.Contains(html, "</html") {
		return strings.Replace(html, "</html>", reloadScriptTag+"</html>", 1)
	}

	return html + reloadScriptTag
}
