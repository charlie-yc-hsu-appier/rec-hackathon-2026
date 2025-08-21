package service

type Response struct {
	ProductIDs   []string                `json:"product_ids"`
	ProductPatch map[string]ProductPatch `json:"product_patch"`
}

type ProductPatch struct {
	Url   string `json:"url"`
	Image string `json:"image"`
}
