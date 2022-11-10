package hubspot

// Association client
type Association struct {
	Client
}

type AssociationLabels struct {
	Results []struct {
		Category string      `json:"category"`
		TypeId   int         `json:"typeId"`
		Label    interface{} `json:"label"`
	} `json:"results"`
}

type AssociationProps struct {
	AssociationCategory string `json:"associationCategory,omitempty"`
	AssociationTypeId   int    `json:"associationTypeId,omitempty"`
}

// Association constructor (from Client)
func (c Client) Association() Association {
	return Association{
		Client: c,
	}
}

func (a Association) GetAssociations(fromObjectType string, toObjectType string) (AssociationLabels, error) {
	// /crm/v4/associations/{fromObjectType}/{toObjectType}/labels

	resp := AssociationLabels{}
	err := a.Request("GET", "/crm/v4/associations/"+fromObjectType+"/"+toObjectType+"/labels", nil, &resp)
	if err != nil {
		return AssociationLabels{}, err
	}

	return resp, nil

}
