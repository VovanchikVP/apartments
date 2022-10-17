package entities

type OperationGroups struct {
	ID   int
	Name string
}

type OperationGroupsTest struct {
	Name string `json:"Name"`
}

type Operation struct {
	ID           int
	Date         string
	Type         string
	Group        OperationGroups
	Value        float64
	Proof        bool
	Descriptions string
}
