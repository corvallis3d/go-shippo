package models

import "encoding/json"

// UnmarshalJSON accepts the customs declaration ID shown in Shippo's
// international guide as well as the expanded object in the current schema.
func (s *Shipment) UnmarshalJSON(data []byte) error {
	type shipment Shipment
	decoded := shipment(*s)
	wire := struct {
		*shipment
		CustomsDeclaration json.RawMessage `json:"customs_declaration"`
	}{shipment: &decoded}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if len(wire.CustomsDeclaration) != 0 {
		decoded.CustomsDeclaration = nil
		if wire.CustomsDeclaration[0] == '"' {
			var id string
			if err := json.Unmarshal(wire.CustomsDeclaration, &id); err != nil {
				return err
			}
			decoded.CustomsDeclaration = &CustomsDeclaration{CommonOutputFields: CommonOutputFields{ObjectID: id}}
		} else if err := json.Unmarshal(wire.CustomsDeclaration, &decoded.CustomsDeclaration); err != nil {
			return err
		}
	}
	*s = Shipment(decoded)
	return nil
}
