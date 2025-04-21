package ui

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func MainLayout(title string, content ...g.Node) g.Node {
	return h.Doctype(h.HTML(
		h.Data("theme", "dark"),
		head(title),
		body(content...)))
}

func TextEditorLayout(title string, filepath string) g.Node {
	return h.Doctype(h.HTML(
		h.Data("theme", "dark"),
		h.Class("flood-fill"),
		textEditorHead(title),
		body(h.Main(h.Data("filepath", filepath)))))
}

func textEditorHead(title string) g.Node {
	return h.Head(
		h.Title(title),
		h.Link(h.Href("/assets/css/layout.css"), h.Rel("stylesheet"), h.Type("text/css")),
		h.Script(h.Src("/assets/js/socket.js")),
		h.Script(h.Src("/assets/js/check-version.js"), g.Attr("async")),
		h.Script(h.Src("/assets/js/editor.bundle.js")),
		h.Script(h.Src("/assets/js/editor.main.js"), g.Attr("async")),
	)
}

func head(title string) g.Node {
	return h.Head(
		h.Title(title),
		h.Link(h.Href("/assets/css/classless.css"), h.Rel("stylesheet"), h.Type("text/css")),
		h.Link(h.Href("/assets/css/themes.css"), h.Rel("stylesheet"), h.Type("text/css")),
		h.Link(h.Href("/assets/css/layout.css"), h.Rel("stylesheet"), h.Type("text/css")),
		h.Link(h.Href("/assets/css/main.css"), h.Rel("stylesheet"), h.Type("text/css")),
		h.Script(h.Src("/assets/js/socket.js")),
		h.Script(h.Src("/assets/js/check-version.js"), g.Attr("async")),
		h.Script(h.Src("/assets/js/main.js"), g.Attr("async")),
	)
}

func body(content ...g.Node) g.Node {
	content = append(content,
		h.Script(h.Src("/assets/js/alpine-morph.js")),
		h.Script(h.Src("/assets/js/alpine.js")))
	return h.Body(content...)
}
