# RootStory: Herb Passport Smart Contract

This repository contains the Go-based smart contract for the RootStory Herb Passport System. The contract is designed for deployment on a Hyperledger Fabric network and provides functionality for recording, updating, and retrieving herb traceability information.

---

## Purpose

RootStory enables end-to-end traceability of medicinal herbs by storing cultivation details, farmer metadata, geographic coordinates, growth stages, and laboratory verification results on a tamper-resistant blockchain ledger.

---

## Files in This Repository

### `herb_contract.go`
Contains the complete smart contract implementation, including:
- Herb data structure  
- Add, update, and query functions  
- History retrieval  
- Region and name-based searches  
- State validation and timestamps  

### `main.go`
Initializes and starts the chaincode using the Fabric contract API.

---

## Smart Contract Capabilities

### Core Operations
- Create new herb entries  
- Update growth stage  
- Attach laboratory report hash and status  
- Fetch a single herb record  
- Fetch all herb records  
- Query herbs by region  
- Query herbs by name or scientific name  
- Check if a herb exists  
- Retrieve the full modification history of a herb  

---

## Deployment

To deploy this contract, place the folder containing `herb_contract.go` and `main.go` inside your Fabric chaincode directory and use the standard Fabric deployment command:
```
peer lifecycle chaincode package rootstory.tar.gz --path ./ --lang golang --label rootstory_1
```

Then follow the standard Fabric install, approve, and commit lifecycle steps.

---

## Example Invocation

### Add a herb
```
peer chaincode invoke -n rootstory -C mychannel -c
'{"Args":["AddHerb","H001","Tulsi","Ocimum sanctum","FarmerA","50kg","23.45","72.56","Gujarat","Ahmedabad","Seedling","2025-01-10"]}'
```

### Update growth stage
```
peer chaincode invoke -n rootstory -C mychannel -c
'{"Args":["UpdateGrowthStage","H001","Vegetative"]}'
```

### Update lab report
```
peer chaincode invoke -n rootstory -C mychannel -c
'{"Args":["UpdateLabReport","H001","QmExampleHash","Approved"]}'
```

### Get all herbs
```
peer chaincode query -n rootstory -C mychannel -c '{"Args":["GetAllHerbs"]}'
```

### Query by region
```
peer chaincode query -n rootstory -C mychannel -c '{"Args":["QueryHerbsByRegion","Gujarat"]}'
```

### Retrieve update history
```
peer chaincode query -n rootstory -C mychannel -c '{"Args":["GetHistory","H001"]}'
```
