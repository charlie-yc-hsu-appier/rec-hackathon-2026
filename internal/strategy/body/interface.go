package body

type Params struct {
	UserID    string
	ClickID   string
	ImgWidth  int
	ImgHeight int
	BundleID  string
	SubID     string
}

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=body

type Strategy interface {
	GenerateBody(params Params) any
}
