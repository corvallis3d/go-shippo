package models

import (
	"encoding/json"
	"testing"
)

// These fixtures cover additions in the live changelog through October 20, 2026.
func TestCurrentShipmentExtras(t *testing.T) {
	emptyPrefix := ""
	first := 1
	no := false
	input := ShipmentInput{
		AddressFrom: "from", AddressTo: "to", Parcels: []string{"parcel"},
		Extra: &ShipmentExtra{
			Reference3: "ref-3", Reference4: "ref-4",
			HealthcareHandling: &ShipmentHealthcareHandling{Enabled: true},
			CustomerReference:  &ShipmentLabelReference{Prefix: &emptyPrefix, Value: "order", RefSort: &first},
			PONumber:           &ShipmentLabelReference{Value: "po"},
			DangerousGoods:     &ShipmentDangerousGoods{Contains: true, LithiumBatteries: &ShipmentDangerousGoodsContents{Contains: true}},
			Insurance:          &ShipmentInsurance{Amount: "100", Currency: "USD", Content: "parts"},
			AuthorityToLeave:   &no,
		},
	}
	request := jsonObject(t, input)
	var extra map[string]json.RawMessage
	if err := json.Unmarshal(request["extra"], &extra); err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{
		"reference_3": `"ref-3"`, "reference_4": `"ref-4"`,
		"healthcare_handling": `{"enabled":true}`, "authority_to_leave": "false",
		"customer_reference": `{"prefix":"","value":"order","ref_sort":1}`,
		"dangerous_goods":    `{"contains":true,"lithium_batteries":{"contains":true}}`,
	} {
		if string(extra[key]) != want {
			t.Errorf("%s = %s, want %s", key, extra[key], want)
		}
	}
	var insurance map[string]json.RawMessage
	if err := json.Unmarshal(extra["insurance"], &insurance); err != nil || insurance["provider"] != nil {
		t.Fatalf("third-party insurance must omit provider: %s (%v)", extra["insurance"], err)
	}
	parcelInsurance := jsonObject(t, ParcelInsurance{Amount: "100", Currency: "USD", Content: "parts"})
	if parcelInsurance["provider"] != nil {
		t.Fatal("parcel third-party insurance must omit provider")
	}
	var response Shipment
	if err := json.Unmarshal([]byte(`{"extra":{"healthcare_handling":{"enabled":true},"reference_3":"three","reference_4":"four"},"rates":[{"servicelevel":{"token":"synthetic_new_service"}}]}`), &response); err != nil {
		t.Fatal(err)
	}
	if !response.Extra.HealthcareHandling.Enabled || response.Extra.Reference4 != "four" || response.Rates[0].ServiceLevel.Token != "synthetic_new_service" {
		t.Fatalf("latest fields lost: %+v", response)
	}
	if ShipmentDangerousGoodsCode09 != "09" || LabelFileTypePDFA5 != "PDF_A5" {
		t.Fatal("incorrect current enum value")
	}
}
