package model

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/RedHatInsights/sources-api-go/util"
)

type Event interface {
	ToEvent() interface{}
}

type updatingAttributes struct {
	ApplicationsCollection               interface{}            `json:"applications"`
	ApplicationAuthenticationsCollection interface{}            `json:"application_authentications"`
	AuthenticationsCollection            interface{}            `json:"authentications"`
	EndpointsCollection                  interface{}            `json:"endpoints"`
	SourceRelation                       interface{}            `json:"source"`
	Updated                              map[string]interface{} `json:"updated"`
}

type EventModelDao interface {
	BulkMessage(resource util.Resource) (map[string]interface{}, error)
	FetchAndUpdateBy(resource util.Resource, updateAttributes map[string]interface{}) (interface{}, error)
	ToEventJSON(resource util.Resource) ([]byte, error)
}

func UpdateMessage(eventObject EventModelDao, resource util.Resource, attributes []string) ([]byte, error) {
	updatedAttributes := updatingAttributes{}

	resourceID := ""
	if resource.ResourceUID != "" && resource.ResourceID == 0 {
		resourceID = resource.ResourceUID
	} else {
		resourceID = strconv.Itoa(int(resource.ResourceID))
	}

	updatedAttributes.Updated = map[string]interface{}{resource.ResourceType: map[string]interface{}{resourceID: attributes}}

	bulkMessage, err := eventObject.BulkMessage(resource)
	if err != nil {
		return nil, fmt.Errorf("error in BulkMessage: %v", err.Error())
	}

	updatedAttributes.ApplicationAuthenticationsCollection = bulkMessage["application_authentications"]
	updatedAttributes.ApplicationsCollection = bulkMessage["applications"]
	updatedAttributes.AuthenticationsCollection = bulkMessage["authentications"]
	updatedAttributes.EndpointsCollection = bulkMessage["endpoints"]
	updatedAttributes.SourceRelation = bulkMessage["source"]

	data, err := json.Marshal(updatedAttributes)
	if err != nil {
		return nil, fmt.Errorf("error in parsing update attributes: %v", err)
	}

	return data, nil
}
