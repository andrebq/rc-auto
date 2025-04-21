package ui

import (
	"encoding/json"
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func ButtonBar(nodes ...g.Node) g.Node {
	return h.Div(
		h.Class("button-bar"),
		g.Group(nodes),
	)
}

func EvenTrigger(kind string, value any) g.Node {
	buf, err := json.Marshal(value)
	if err != nil {
		buf = []byte("{}")
	}
	return g.Attr(fmt.Sprintf("et-%v", kind), string(buf))
}

func Controls() g.Node {
	return h.Section(
		ButtonBar(
			h.Button(g.Text("Frente"),
				h.ID("move-forward"),
				EvenTrigger("click", struct{ Move string }{Move: "forward"})),
			h.Button(g.Text("Parar"),
				h.ID("move-stop"),
				EvenTrigger("click", struct{ Move string }{Move: "stop"})),
			h.Button(g.Text("Ré"),
				h.ID("move-forward"),
				EvenTrigger("click", struct{ Move string }{Move: "forward"})),
		),
		h.Div(
			h.Input(g.Attr("type", "text"), g.Attr("placeholder", "Velocidade")),
		),
	)
}
