package productcatlog

type CreateCatalogProductDTO struct {
	Name        string `json:"name" binding:"required"`
	Brand       string `json:"brand"`
	Category    string `json:"category" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"image_url"`
}
