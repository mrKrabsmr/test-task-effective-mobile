package persons

import (
	core "github.com/mrKrabsmr/commerce-edu-api/internal/apps"
)

type PersonRequestDTO struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic" validate:""`
}

func (d *PersonRequestDTO) Validate() error {
	v := core.GetValidator()
	return v.Struct(d)
}

type ResponseResult struct {
	Age    int
	Gender string
	Nation string
	Err    error
}
