package models

import (
	"encoding/json"
	"testing"
)

func TestCurrentBatchShapes(t *testing.T) {
	input := BatchShipmentInput{
		Shipment: "legacy-shipment", CarrierAccount: "carrier-1", ServiceLevelToken: "ups_ground",
		ShipmentDetails: &ShipmentInput{AddressFrom: "from-1", AddressTo: "to-1", Parcels: []string{"parcel-1"}},
	}
	wire := jsonObject(t, input)
	var shipment ShipmentInput
	if err := json.Unmarshal(wire["shipment"], &shipment); err != nil {
		t.Fatal(err)
	}
	if shipment.AddressFrom != "from-1" || string(wire["servicelevel_token"]) != `"ups_ground"` {
		t.Fatalf("batch shipment request = %s", wire)
	}
	input.ShipmentDetails = nil
	if legacy := jsonObject(t, input); string(legacy["shipment"]) != `"legacy-shipment"` {
		t.Fatalf("legacy shipment ID lost: %s", legacy["shipment"])
	}

	var batch Batch
	fixture := `{"object_id":"batch-1","status":"VALID","test":true,"batch_shipments":{"next":null,"results":[{"object_id":"batch-shipment-1","shipment":"shipment-1","status":"TRANSACTION_FAILED","transaction":"transaction-1","messages":[{"text":"declined"}]}]}}`
	if err := json.Unmarshal([]byte(fixture), &batch); err != nil {
		t.Fatal(err)
	}
	item := batch.BatchShipments.Results[0]
	if !batch.Test || item.Status != "TRANSACTION_FAILED" || item.Shipment != "shipment-1" {
		t.Fatalf("batch response = %+v, item = %+v", batch, item)
	}
	// The request marshaler must not hide fields on the embedded response type.
	response := jsonObject(t, item)
	for key, want := range map[string]string{"object_id": `"batch-shipment-1"`, "shipment": `"shipment-1"`, "status": `"TRANSACTION_FAILED"`, "transaction": `"transaction-1"`} {
		if got := string(response[key]); got != want {
			t.Errorf("%s = %s, want %s", key, got, want)
		}
	}
	if response["messages"] == nil {
		t.Fatal("response lost messages")
	}
}

func TestCurrentManifestAndCarrierAccountShapes(t *testing.T) {
	var manifest Manifest
	if err := json.Unmarshal([]byte(`{"object_id":"manifest-1","status":"ERROR","errors":["invalid_address"],"documents":[]}`), &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Status != "ERROR" || len(manifest.Errors) != 1 || manifest.Errors[0] != "invalid_address" {
		t.Fatalf("manifest = %+v", manifest)
	}
	var account CarrierAccount
	fixture := `{"object_id":"carrier-1","carrier":"fedex","carrier_name":"FedEx","account_id":"****","active":true,"is_shippo_account":false,"metadata":"warehouse","service_levels":[{"name":"Ground","token":"fedex_ground","supports_return_labels":true}],"object_info":{"authentication":{"type":"default","status":"connected"},"verification":{"policy":"required"}}}`
	if err := json.Unmarshal([]byte(fixture), &account); err != nil {
		t.Fatal(err)
	}
	if account.CarrierName != "FedEx" || account.Metadata != "warehouse" || account.ObjectInfo.Authentication.Status != "connected" || account.ObjectInfo.Verification.Policy != "required" || !account.ServiceLevels[0].SupportsReturnLabels {
		t.Fatalf("carrier account = %+v", account)
	}
}
