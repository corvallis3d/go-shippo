package models

import (
	"encoding/json"
	"testing"
)

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
