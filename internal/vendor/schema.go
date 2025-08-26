package vendor

type Request struct {
	UserID    string `form:"user_id" binding:"required"`
	ClickID   string `form:"click_id" binding:"required"`
	ImgWidth  int    `form:"w" binding:"required"`
	ImgHeight int    `form:"h" binding:"required"`
}

type ProductInfo struct {
	ProductID string `json:"product_id"`
	Url       string `json:"url"`
	Image     string `json:"image"`
}
