package codefixture

type ModelRelation struct {
	TargetRef  ModelRef
	ForeignRef ModelRef
	Connector  Connector
}

type Connector func(target any, foreign any)
