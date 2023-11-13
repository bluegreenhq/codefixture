package codefixture

type ModelRelation struct {
	TargetRef  ModelRef
	ForeignRef ModelRef
	Connector  func(target any, dependent any)
}
