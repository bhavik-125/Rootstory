package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Herb stores botanical & traceability details
type Herb struct {
	HerbID        string `json:"herbID"`
	Name          string `json:"name"`
	Scientific    string `json:"scientific"`
	Farmer        string `json:"farmer"`
	Quantity      string `json:"quantity"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	Region        string `json:"region"`     // State or region name
	PlaceName     string `json:"placeName"`  // City / Village name
	GrowthStage   string `json:"growthStage"`
	PlantingDate  string `json:"plantingDate"`
	HarvestDate   string `json:"harvestDate"`
	LabReportHash string `json:"labReportHash"`
	Status        string `json:"status"`
	Timestamp     string `json:"timestamp"`
}

// SmartContract for Herb Passport
type SmartContract struct {
	contractapi.Contract
}

// AddHerb creates a new herb record
func (s *SmartContract) AddHerb(ctx contractapi.TransactionContextInterface,
	herbID, name, scientific, farmer, quantity,
	latitude, longitude, region, placeName,
	growthStage, plantingDate string) error {

	exists, err := s.HerbExists(ctx, herbID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Herb %s already exists", herbID)
	}

	h := Herb{
		HerbID:       herbID,
		Name:         name,
		Scientific:   scientific,
		Farmer:       farmer,
		Quantity:     quantity,
		Latitude:     latitude,
		Longitude:    longitude,
		Region:       region,
		PlaceName:    placeName,
		GrowthStage:  growthStage,
		PlantingDate: plantingDate,
		Status:       "Submitted",
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}

	b, err := json.Marshal(h)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(herbID, b)
}



// UpdateGrowthStage updates growth stage & timestamp
func (s *SmartContract) UpdateGrowthStage(ctx contractapi.TransactionContextInterface, herbID, newStage string) error {
	h, err := s.GetHerb(ctx, herbID)
	if err != nil {
		return err
	}
	h.GrowthStage = newStage
	h.Timestamp = time.Now().UTC().Format(time.RFC3339)

	b, err := json.Marshal(h)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(herbID, b)
}

// UpdateLabReport stores lab report hash and status
func (s *SmartContract) UpdateLabReport(ctx contractapi.TransactionContextInterface, herbID, labHash, status string) error {
	h, err := s.GetHerb(ctx, herbID)
	if err != nil {
		return err
	}
	h.LabReportHash = labHash
	h.Status = status
	h.Timestamp = time.Now().UTC().Format(time.RFC3339)

	b, err := json.Marshal(h)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(herbID, b)
}

// GetHerb fetches details for a given herb
func (s *SmartContract) GetHerb(ctx contractapi.TransactionContextInterface, herbID string) (*Herb, error) {
	b, err := ctx.GetStub().GetState(herbID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, fmt.Errorf("Herb %s does not exist", herbID)
	}

	var h Herb
	if err := json.Unmarshal(b, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

// GetAllHerbs returns all herb entries
func (s *SmartContract) GetAllHerbs(ctx contractapi.TransactionContextInterface) ([]*Herb, error) {
	iter, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var results []*Herb
	for iter.HasNext() {
		kvs, err := iter.Next()
		if err != nil {
			return nil, err
		}
		var h Herb
		if err := json.Unmarshal(kvs.Value, &h); err != nil {
			return nil, err
		}
		results = append(results, &h)
	}
	return results, nil
}

// QueryHerbsByRegion fetches all herbs from a given region
func (s *SmartContract) QueryHerbsByRegion(ctx contractapi.TransactionContextInterface, region string) ([]*Herb, error) {
	all, err := s.GetAllHerbs(ctx)
	if err != nil {
		return nil, err
	}
	var out []*Herb
	for _, h := range all {
		if strings.EqualFold(h.Region, region) {
			out = append(out, h)
		}
	}
	return out, nil
}

// QueryHerbsByName searches by name or scientific name
func (s *SmartContract) QueryHerbsByName(ctx contractapi.TransactionContextInterface, nameSubstr string) ([]*Herb, error) {
	all, err := s.GetAllHerbs(ctx)
	if err != nil {
		return nil, err
	}
	var out []*Herb
	for _, h := range all {
		if strings.Contains(strings.ToLower(h.Name), strings.ToLower(nameSubstr)) ||
			strings.Contains(strings.ToLower(h.Scientific), strings.ToLower(nameSubstr)) {
			out = append(out, h)
		}
	}
	return out, nil
}

// HerbExists checks if herb exists
func (s *SmartContract) HerbExists(ctx contractapi.TransactionContextInterface, herbID string) (bool, error) {
	b, err := ctx.GetStub().GetState(herbID)
	if err != nil {
		return false, err
	}
	return b != nil, nil
}

// GetHistory returns update history for a herb
func (s *SmartContract) GetHistory(ctx contractapi.TransactionContextInterface, herbID string) ([]map[string]interface{}, error) {
	history, err := ctx.GetStub().GetHistoryForKey(herbID)
	if err != nil {
		return nil, err
	}
	defer history.Close()

	var arr []map[string]interface{}
	for history.HasNext() {
		mod, err := history.Next()
		if err != nil {
			return nil, err
		}
		var value map[string]interface{}
		_ = json.Unmarshal(mod.Value, &value)
		entry := map[string]interface{}{
			"TxId":      mod.TxId,
			"Timestamp": mod.Timestamp,
			"IsDelete":  mod.IsDelete,
			"Value":     value,
		}
		arr = append(arr, entry)
	}
	return arr, nil
}
