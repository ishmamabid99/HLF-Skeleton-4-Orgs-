package main

import (
	"fmt"
	"encoding/json"
	//"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	// "github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)


type SmartContract struct {
	contractapi.Contract
}


type PersonInfo struct {
	Identifier string `json:"identifier"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Status string `json:"status"`
}


//For Testing Purpose
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	personInfo := []PersonInfo{
		{
			Identifier          : "01",
			FirstName 			: "Ashraf",
			LastName			: "Tasin",
			Status				: "Approved",
		},
		{
			Identifier          : "02",
			FirstName 			: "Joydip",
			LastName			: "Das",
			Status				: "Approved",
		},
	}

	for _ , person := range personInfo{
		personInfoAsBytes, err := json.Marshal(person)
		
		if err != nil {
			return err
		}

		err2 := ctx.GetStub().PutState(person.Identifier, personInfoAsBytes)

		if err2 != nil {
			return fmt.Errorf("failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) AddPresonInfo(ctx contractapi.TransactionContextInterface, identifier string, firstName string, lastName string, status string) error {
	
	exists,err := s.PersonExists(ctx,identifier)

    if err != nil {
      return err
    }

    if exists {
      return fmt.Errorf("the person %s already exists", identifier)
    }
	
	personInfo := PersonInfo{
			Identifier          : identifier,
			FirstName 			: firstName,
			LastName			: lastName,
			Status				: status,
	}

	personInfoAsBytes, err2 := json.Marshal(personInfo)

	if err2 != nil {
		return err2
	}

	return ctx.GetStub().PutState(identifier, personInfoAsBytes)
}

func (s *SmartContract) QueryPersonInfo(ctx contractapi.TransactionContextInterface, identifier string) (*PersonInfo, error) {
	infoBlockAsBytes, err := ctx.GetStub().GetState(identifier)
	
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if infoBlockAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", identifier)
	}

	infoBlock := new(PersonInfo)
	err2 := json.Unmarshal(infoBlockAsBytes, infoBlock)

	if err2 != nil {
		return nil, err2
	}

	return infoBlock, nil
}

func (s *SmartContract) UpdatePresonInfo(ctx contractapi.TransactionContextInterface, identifier string, firstName string, lastName string, status string) error {
	
	exists,err := s.PersonExists(ctx,identifier)

    if err != nil {
      return err
    }

    if exists {
      return fmt.Errorf("the person %s already exists", identifier)
    }

	if !exists {
	  return fmt.Errorf("the person %s doesn't exist", identifier)
	}
	
	personInfo := PersonInfo{
			Identifier          : identifier,
			FirstName 			: firstName,
			LastName			: lastName,
			Status				: status,
	}

	personInfoAsBytes, err2 := json.Marshal(personInfo)

	if err2 != nil {
		return err2
	}

	return ctx.GetStub().PutState(identifier, personInfoAsBytes)
}

func (s *SmartContract) DeletePresonInfo(ctx contractapi.TransactionContextInterface, identifier string, firstName string, lastName string, status string) error {
	
	exists,err := s.PersonExists(ctx,identifier)

    if err != nil {
      return err
    }

	if !exists {
	  return fmt.Errorf("the person %s doesn't exist", identifier)
	}

	return ctx.GetStub().DelState(identifier)
}

func (s *SmartContract) UpdatePresonStatus(ctx contractapi.TransactionContextInterface, identifier string, status string) error {
	
	person,err := s.QueryPersonInfo(ctx,identifier)

	if err != nil {
		return err
	}
	
	person.Status = status
	personInfoAsBytes, err2 := json.Marshal(person)

	if err2 != nil {
		return err2
	}

	return ctx.GetStub().PutState(identifier, personInfoAsBytes)
}

func (s *SmartContract) PersonExists(ctx contractapi.TransactionContextInterface, identifier string) (bool,error){
	personInfoAsBytes,err := ctx.GetStub().GetState(identifier)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return personInfoAsBytes !=nil , nil
}

func (s *SmartContract) GetAllInfo(ctx contractapi.TransactionContextInterface, identifier string) ([]*PersonInfo,error) {
	start := ""
	end := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(start,end)

	if err != nil {
		return nil,err
	}
	
	defer resultsIterator.Close()

	var persons []*PersonInfo

	for resultsIterator.HasNext(){
		queryResponse, err2 := resultsIterator.Next()
		if err2 != nil {
			return nil, err2
		}

		var person PersonInfo
		err3 := json.Unmarshal(queryResponse.Value, &person)

		if err3 != nil {
			return nil, err
		}
		persons = append(persons, &person)
	}

	return persons,nil
}

func (s *SmartContract) RandomFunc(ctx contractapi.TransactionContextInterface, identifier string) ([]byte,error) {
	personInfoAsBytes, err := ctx.GetStub().GetCreator()

	if err != nil {
		return nil, err
	}
	
	// infoBlock := new(InfoBlock)

	// err2 = json.Unmarshal(personInfoAsBytes, infoBlock)

	return personInfoAsBytes, nil
}

func main() {
	smartContract := new(SmartContract)

	cc, err := contractapi.NewChaincode(smartContract)

	if err != nil {
		fmt.Printf("Error create smartcert chaincode: %s", err.Error())
		return
	}

	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting smartcert chaincode: %s", err.Error())
	}
}