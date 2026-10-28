package models

import (
	"encoding/json"
	"testing"
)

func TestShipmentParcelRequestShapes(t *testing.T) {
	parcel := &ParcelInput{Weight: "1", MassUnit: "lb"}
	for _, test := range []struct {
		name    string
		parcels interface{}
		want    string
	}{
		{"ID", "parcel-1", `["parcel-1"]`},
		{"object", parcel, `[{"weight":"1","mass_unit":"lb"}]`},
		{"IDs", []string{"parcel-1", "parcel-2"}, `["parcel-1","parcel-2"]`},
		{"objects", []*ParcelInput{parcel}, `[{"weight":"1","mass_unit":"lb"}]`},
		{"nil", nil, `null`},
		{"nil object", (*ParcelInput)(nil), `null`},
		{"nil slice", []string(nil), `null`},
	} {
		t.Run(test.name, func(t *testing.T) {
			input := &ShipmentInput{Parcels: test.parcels}
			before, err := json.Marshal(input.Parcels)
			if err != nil {
				t.Fatal(err)
			}
			wire := jsonObject(t, input)
			if got := string(wire["parcels"]); got != test.want {
				t.Fatalf("parcels = %s, want %s", got, test.want)
			}
			after, err := json.Marshal(input.Parcels)
			if err != nil || string(after) != string(before) {
				t.Fatalf("marshaling changed caller's parcels: %s -> %s (%v)", before, after, err)
			}
		})
	}

	shipment := &ShipmentInput{Parcels: "parcel-1"}
	for _, input := range []interface{}{
		TransactionInput{Shipment: shipment},
		BatchShipmentInput{ShipmentDetails: shipment},
	} {
		wire := jsonObject(t, input)
		var nested struct {
			Parcels []string `json:"parcels"`
		}
		if err := json.Unmarshal(wire["shipment"], &nested); err != nil {
			t.Fatal(err)
		}
		if len(nested.Parcels) != 1 || nested.Parcels[0] != "parcel-1" {
			t.Fatalf("%T nested parcels = %v", input, nested.Parcels)
		}
	}
}

func TestShipmentCustomsDeclarationResponseShapes(t *testing.T) {
	var shipment Shipment
	for _, test := range []struct {
		name, customs, wantID string
	}{
		{"expanded", `{"object_id":"customs-1","contents_type":"MERCHANDISE","items":["item-1"]}`, "customs-1"},
		{"ID", `"customs-2"`, "customs-2"},
		{"null", `null`, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := `{"object_id":"shipment-1","status":"SUCCESS","parcels":[{"object_id":"parcel-1"}],"rates":[{"object_id":"rate-1"}],"customs_declaration":` + test.customs + `}`
			if err := json.Unmarshal([]byte(fixture), &shipment); err != nil {
				t.Fatal(err)
			}
			if test.wantID == "" {
				if shipment.CustomsDeclaration != nil {
					t.Fatalf("null customs declaration = %+v", shipment.CustomsDeclaration)
				}
			} else if shipment.CustomsDeclaration == nil || shipment.CustomsDeclaration.ObjectID != test.wantID {
				t.Fatalf("customs declaration = %+v", shipment.CustomsDeclaration)
			}
			if shipment.ObjectID != "shipment-1" || shipment.Status != "SUCCESS" || len(shipment.Parcels) != 1 || len(shipment.Rates) != 1 {
				t.Fatalf("lost shipment fields: %+v", shipment)
			}
			if test.name == "ID" && shipment.CustomsDeclaration.ContentsType != "" {
				t.Fatal("ID response retained a previous expanded declaration")
			}
		})
	}
	if err := json.Unmarshal([]byte(`{"customs_declaration":123}`), &shipment); err == nil {
		t.Fatal("accepted invalid customs declaration type")
	}
}
