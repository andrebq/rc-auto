package ui

import (
	"bytes"
	"net/http"
	"strconv"

	g "maragu.dev/gomponents"
)

func RenderNode(w http.ResponseWriter, _ *http.Request, node g.Node) {
	buf := bytes.Buffer{}
	node.Render(&buf)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("Content-Length", strconv.Itoa(buf.Len()))
	w.Write(buf.Bytes())
}
