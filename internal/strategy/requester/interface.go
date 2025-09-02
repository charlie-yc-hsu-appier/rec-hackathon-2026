package requester

type Params struct {
	RequestURL string
	UserID     string
	ClickID    string
	ImgWidth   int
	ImgHeight  int
}

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=requester
type Strategy interface {
	GenerateRequestURL(params Params) (string, error)
}
