package models

import (
	"encoding/json"
	"testing"
	"time"
)

func jsonObject(t *testing.T, value interface{}) map[string]json.RawMessage {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		t.Fatal(err)
	}
	return object
}

func TestTransactionPurchaseShape(t *testing.T) {
	if wire := jsonObject(t, TransactionInput{Rate: "rate-1"}); wire["label_file_type"] != nil {
		t.Fatal("unset optional label format should use Shippo's default")
	}
	ratePurchase := jsonObject(t, TransactionInput{Rate: "rate-1", Order: "order-1", LabelFileType: LabelFileTypePDF4X6})
	for key, want := range map[string]string{"rate": `"rate-1"`, "order": `"order-1"`, "label_file_type": `"PDF_4x6"`, "async": "false"} {
		if got := string(ratePurchase[key]); got != want {
			t.Errorf("%s = %s, want %s", key, got, want)
		}
	}
	if _, ok := ratePurchase["servicelevel_token"]; ok {
		t.Fatal("rate purchase contains an unset service level")
	}

	instantPurchase := jsonObject(t, TransactionInput{
		Order: "order-1", CarrierAccount: "carrier-1", ServerLevelToken: "ups_ground",
		Shipment: &ShipmentInput{
			AddressFrom: "address-1", AddressTo: "address-2", Parcels: []string{"parcel-1", "parcel-2"},
			Extra: &ShipmentExtra{QRCodeRequested: true},
		},
	})
	if got := string(instantPurchase["servicelevel_token"]); got != `"ups_ground"` {
		t.Fatalf("service level = %s", got)
	}
	if _, ok := instantPurchase["servericelevel_token"]; ok {
		t.Fatal("malformed service-level key was serialized")
	}
	if _, ok := instantPurchase["qr_code_requested"]; ok {
		t.Fatal("QR request belongs under shipment.extra")
	}
	var shipment struct {
		Parcels      []string        `json:"parcels"`
		Extra        ShipmentExtra   `json:"extra"`
		ShipmentDate json.RawMessage `json:"shipment_date"`
	}
	if err := json.Unmarshal(instantPurchase["shipment"], &shipment); err != nil {
		t.Fatal(err)
	}
	if len(shipment.Parcels) != 2 || !shipment.Extra.QRCodeRequested || shipment.ShipmentDate != nil {
		t.Fatalf("unexpected nested shipment: %+v", shipment)
	}
}

func TestTransactionResponseShapes(t *testing.T) {
	for _, test := range []struct {
		name, rate, wantID string
		expanded           bool
	}{
		{"ID", `"rate-1"`, "rate-1", false},
		{"expanded", `{"object_id":"rate-1","amount":"12.50","currency":"USD","provider":"UPS","servicelevel_token":"ups_ground"}`, "rate-1", true},
		{"null", `null`, "", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := `{"object_id":"transaction-1","status":"SUCCESS","order":"order-1","parcel":"parcel-2","rate":` + test.rate + `,"tracking_number":"tracking-2","label_url":"https://labels.example/2.pdf","qr_code_url":"https://labels.example/2.png","eta":null}`
			var transaction Transaction
			if err := json.Unmarshal([]byte(fixture), &transaction); err != nil {
				t.Fatal(err)
			}
			if transaction.Rate != test.wantID || (transaction.RateDetails != nil) != test.expanded {
				t.Fatalf("rate fields: %+v", transaction)
			}
			if transaction.Order != "order-1" || transaction.Parcel != "parcel-2" || transaction.TrackingNumber != "tracking-2" || transaction.QRCodeURL == "" || transaction.Status != "SUCCESS" {
				t.Fatalf("lost response fields: %+v", transaction)
			}
			if test.expanded && (transaction.RateDetails.Amount != "12.50" || transaction.RateDetails.ServiceLevelToken != "ups_ground") {
				t.Fatalf("lost expanded rate: %+v", transaction.RateDetails)
			}
			wire := jsonObject(t, transaction)
			if string(wire["parcel"]) != `"parcel-2"` || wire["label_url"] == nil {
				t.Fatal("marshaling transaction lost output fields")
			}
		})
	}
	var transaction Transaction
	if err := json.Unmarshal([]byte(`{"rate":123}`), &transaction); err == nil {
		t.Fatal("accepted invalid rate type")
	}
}

func TestCurrentShipmentAndRates(t *testing.T) {
	fixture := `{"object_id":"shipment-1","status":"SUCCESS","metadata":"order-1","shipment_date":"2026-10-26T12:00:00Z","carrier_accounts":["carrier-1"],"extra":{"qr_code_requested":true},"parcels":[{"object_id":"parcel-1","weight":"1","mass_unit":"lb"},{"object_id":"parcel-2","weight":"2","mass_unit":"lb"}],"rates":[{"object_id":"rate-1","shipment":"shipment-1","amount":"12.50","currency":"USD","estimated_days":3,"arrives_by":"17:00:00","included_insurance_price":null,"servicelevel":{"name":"Ground","token":"ups_ground","extended_token":"ups_ground_special","parent_servicelevel":{"token":"ups_ground"}}}]}`
	var shipment Shipment
	if err := json.Unmarshal([]byte(fixture), &shipment); err != nil {
		t.Fatal(err)
	}
	if len(shipment.Parcels) != 2 || shipment.Parcels[1].ObjectID != "parcel-2" || !shipment.Extra.QRCodeRequested || shipment.Metadata != "order-1" {
		t.Fatalf("shipment = %+v", shipment)
	}
	rate := shipment.Rates[0]
	if rate.EstimatedDays != 3 || rate.Shipment != "shipment-1" || rate.ArrivesBy != "17:00:00" || rate.ServiceLevel.ParentServiceLevel.Token != "ups_ground" {
		t.Fatalf("rate = %+v", rate)
	}
	date := time.Date(2026, time.October, 26, 12, 0, 0, 0, time.UTC)
	wire := jsonObject(t, ShipmentInput{ShipmentDate: date})
	if string(wire["shipment_date"]) != `"2026-10-26T12:00:00Z"` {
		t.Fatalf("explicit shipment date = %s", wire["shipment_date"])
	}
	parcel := jsonObject(t, ParcelInput{Template: "USPS_SmallFlatRateBox", Weight: "2", MassUnit: "lb"})
	for _, dimension := range []string{"length", "width", "height", "distance_unit"} {
		if _, ok := parcel[dimension]; ok {
			t.Errorf("template parcel contains blank %s", dimension)
		}
	}
}

func TestOrderAndAddressShapes(t *testing.T) {
	address := Address{AddressInput: AddressInput{Name: "Recipient", Country: "US", State: "OR"}, CommonOutputFields: CommonOutputFields{ObjectID: "address-1"}}
	order := jsonObject(t, OrderInput{ToAddress: address, FromAddress: &address, PlacedAt: time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC)})
	for _, key := range []string{"to_address", "from_address"} {
		var nested map[string]json.RawMessage
		if err := json.Unmarshal(order[key], &nested); err != nil {
			t.Fatal(err)
		}
		if string(nested["state"]) != `"OR"` || nested["object_created"] != nil || nested["object_id"] != nil || nested["validation_results"] != nil {
			t.Fatalf("order address = %s", order[key])
		}
	}
	if wire := jsonObject(t, AddressInput{Country: "GB"}); wire["state"] != nil {
		t.Fatalf("optional state was not omitted: %s", wire["state"])
	}
	var response OrderListResponse
	fixture := `{"count":1,"next":null,"previous":null,"results":[{"object_id":"order-1","placed_at":"2026-10-26T00:00:00Z","transactions":[{"object_id":"transaction-1","parcel":"parcel-1","rate":{"object_id":"rate-1"}}]}]}`
	if err := json.Unmarshal([]byte(fixture), &response); err != nil {
		t.Fatal(err)
	}
	if response.Count != 1 || response.Results[0].Transactions[0].Rate != "rate-1" {
		t.Fatalf("order response = %+v", response)
	}
}

func TestOrderLineItemCreateShape(t *testing.T) {
	input := OrderInput{LineItems: []LineItem{{ObjectID: "old-item-1", Title: "Part", Quantity: 1}}}
	wire := jsonObject(t, input)
	var items []map[string]json.RawMessage
	if err := json.Unmarshal(wire["line_items"], &items); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"object_id", "weight_unit", "weight", "total_price", "currency"} {
		if _, ok := items[0][key]; ok {
			t.Errorf("order creation included unset or response-only %s", key)
		}
	}
	if string(items[0]["title"]) != `"Part"` || string(items[0]["quantity"]) != "1" {
		t.Fatalf("line item = %+v", items[0])
	}
	if input.LineItems[0].ObjectID != "old-item-1" {
		t.Fatal("marshaling modified the caller's line item")
	}
	if output := jsonObject(t, input.LineItems[0]); string(output["object_id"]) != `"old-item-1"` {
		t.Fatal("standalone line item lost its response ID")
	}
}

func TestOrderValidationResponse(t *testing.T) {
	for _, payload := range []string{`{"to_address":["Invalid address"]}`, `"Invalid order"`, `null`} {
		var order OrderResponse
		fixture := `{"object_id":"order-1","validation_status":"INVALID","error":` + payload + `}`
		if err := json.Unmarshal([]byte(fixture), &order); err != nil {
			t.Fatal(err)
		}
		if order.ValidationStatus != "INVALID" || string(order.Error) != payload {
			t.Fatalf("lost validation response: %+v", order)
		}
	}
}

func TestTrackingPackageDetails(t *testing.T) {
	input := jsonObject(t, TrackingStatusInput{Carrier: "ups", TrackingNumber: "tracking-1", IncludePackageDetails: true})
	if string(input["include_package_details"]) != "true" {
		t.Fatalf("include_package_details = %s", input["include_package_details"])
	}
	legacy := jsonObject(t, TrackingStatusInput{Carrier: "ups", TrackingNumber: "tracking-1"})
	if _, ok := legacy["include_package_details"]; ok {
		t.Fatal("unset package details request was not omitted")
	}
	fixture := `{"carrier":"ups","tracking_number":"tracking-1","weight":[{"value":"2.2","unit":"LB"},{"value":"1.0","unit":"KG"}],"dimensions":[{"length":"1","width":"2","height":"3","unit":"in"},{"length":"3","unit":"cm"}]}`
	var tracking TrackingStatus
	if err := json.Unmarshal([]byte(fixture), &tracking); err != nil {
		t.Fatal(err)
	}
	if len(tracking.Weight) != 2 || tracking.Weight[0].Value != "2.2" || tracking.Weight[1].Unit != "KG" || len(tracking.Dimensions) != 2 || tracking.Dimensions[0].Height != "3" {
		t.Fatalf("package details = %+v", tracking)
	}
	partial := jsonObject(t, tracking.Dimensions[1])
	if partial["width"] != nil || partial["height"] != nil || string(partial["unit"]) != `"cm"` {
		t.Fatalf("partial dimensions = %+v", partial)
	}
}

func TestCustomsDutyAndTaxIdentification(t *testing.T) {
	input := CustomsDeclarationInput{
		B13AFilingOption: CustomsB13ANotRequired, B13ANumber: "exemption-1", IsVATCollected: true,
		DutiesPayor: &CustomsDutiesPayor{
			Type: ShipmentBillingTypeThirdParty, Account: "account-1",
			Address: &CustomsDutiesPayorAddress{Name: "Importer", Zip: "80331", Country: "DE"},
		},
		ExporterIdentification: &CustomsExporterIdentification{
			EORINumber: "DE1234567890",
			TaxID:      &CustomsTaxIdentification{Type: CustomsTaxIDTypeVAT, Number: "123456789"},
		},
	}
	wire := jsonObject(t, input)
	if string(wire["is_vat_collected"]) != "true" || string(wire["b13a_filing_option"]) != `"NOT_REQUIRED"` || string(wire["b13a_number"]) != `"exemption-1"` {
		t.Fatalf("customs fields = %+v", wire)
	}
	var duties CustomsDutiesPayor
	if err := json.Unmarshal(wire["duties_payor"], &duties); err != nil {
		t.Fatal(err)
	}
	if duties.Account != "account-1" || duties.Address.Country != "DE" || duties.Type != "THIRD_PARTY" {
		t.Fatalf("duties = %+v", duties)
	}
	var declaration CustomsDeclaration
	fixture := `{"object_id":"customs-1","b13a_filing_option":"NOT_REQUIRED","b13a_number":"exemption-1","is_vat_collected":true,"exporter_identification":{"eori_number":"DE1234567890","tax_id":{"type":"VAT","number":"123456789"}}}`
	if err := json.Unmarshal([]byte(fixture), &declaration); err != nil {
		t.Fatal(err)
	}
	if declaration.ExporterIdentification.EORINumber != input.ExporterIdentification.EORINumber || declaration.ExporterIdentification.TaxID.Type != CustomsTaxIDTypeVAT || !declaration.IsVATCollected {
		t.Fatalf("customs response = %+v", declaration)
	}
}

func TestTrackingCustomsAndRefundShapes(t *testing.T) {
	var tracking TrackingStatus
	fixture := `{"carrier":"ups","tracking_number":"tracking-1","transaction":"transaction-1","eta":null,"original_eta":"2026-10-29T12:00:00Z","messages":[],"tracking_status":{"object_id":"event-1","status":"PRE_TRANSIT","substatus":{"code":"information_received","text":"Information received","action_required":false}},"tracking_history":[]}`
	if err := json.Unmarshal([]byte(fixture), &tracking); err != nil {
		t.Fatal(err)
	}
	if tracking.Transaction != "transaction-1" || tracking.OriginalETA.IsZero() || tracking.TrackingStatus.Status != TrackingStatusStatusPreTransit || tracking.TrackingStatus.ObjectID != "event-1" || tracking.TrackingStatus.Substatus.Code != "information_received" {
		t.Fatalf("tracking = %+v", tracking)
	}
	customs := jsonObject(t, CustomsDeclarationInput{
		Certify: true, CertifySigner: "Sender", NonDeliveryOption: CustomsNonDeliveryOptionReturn,
		ContentsType: CustomsContentsTypeMerchandise, ImporterReference: "importer-1",
		Items: []*CustomsItemInput{{Description: "Part", Quantity: 1, NetWeight: "1", MassUnit: "lb", ValueAmount: "10", ValueCurrency: "USD", OriginCountry: "US", HSCode: "3926.90", SKUCode: "part-1", ManufacturerCode: "maker-1"}},
	})
	if string(customs["importer_reference"]) != `"importer-1"` || string(customs["contents_type"]) != `"MERCHANDISE"` {
		t.Fatalf("customs = %+v", customs)
	}
	var refund Refund
	if err := json.Unmarshal([]byte(`{"object_id":"refund-1","transaction":"transaction-1","status":"SUCCESS","test":true}`), &refund); err != nil {
		t.Fatal(err)
	}
	if refund.Transaction != "transaction-1" || refund.Status != StatusSuccess {
		t.Fatalf("refund = %+v", refund)
	}
}
