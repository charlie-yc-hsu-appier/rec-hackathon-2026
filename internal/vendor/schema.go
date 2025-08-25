package vendor

type Request struct {
	UserID    string `form:"user_id" binding:"required"`
	ClickID   string `form:"click_id" binding:"required"`
	ImgWidth  int    `form:"w" binding:"required"`
	ImgHeight int    `form:"h" binding:"required"`
}

type Response struct {
	ProductIDs   []string                `json:"product_ids"`
	ProductPatch map[string]ProductPatch `json:"product_patch"`
}

type ProductPatch struct {
	Url   string `json:"url"`
	Image string `json:"image"`
}
