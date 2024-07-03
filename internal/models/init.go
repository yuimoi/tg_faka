package models

type MyModel interface {
	TableName() string
}

var MyModels = []MyModel{
	&Order{},
	&ProductItem{},
	&Product{},
}
