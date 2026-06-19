package request

import "mime/multipart"

type ProductCreateRequest struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Price       float64               `form:"price" binding:"required,min=0"`
	Stock       int                   `form:"stock" binding:"required,min=0"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
}
