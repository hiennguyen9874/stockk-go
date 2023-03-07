package repository

import (
	"github.com/hiennguyen9874/stockk-go/models"
)

type UserRepo struct {
	BaseRepo[models.User]
}

func (r *UserRepo) GetModel(exp *models.UserCreate) *models.User {
	return &models.User{
		Name:        exp.Name,
		Email:       exp.Email,
		Password:    exp.Password,
		IsActive:    true,
		IsSuperUser: false,
	}
}

func (r *UserRepo) GetModelResponse(exp *models.User) *models.UserResponse {
	return &models.UserResponse{
		Id:          exp.Id,
		Name:        exp.Name,
		Email:       exp.Email,
		CreatedAt:   exp.CreatedAt,
		UpdatedAt:   exp.UpdatedAt,
		IsActive:    exp.IsActive,
		IsSuperUser: exp.IsSuperUser,
	}
}
