package bus

import "encoding/json"

func SendData(out Transmitter, to uint64, kind string, value any) error {
	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}
	p := Payload{Value: buf}
	p.Meta.Kind = kind
	out.Transmit(to, p)
	return nil
}
