package controller

import (
	"BussinessManagement/model"
	"BussinessManagement/repository"
	"BussinessManagement/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/casbin/casbin/v2"

	"github.com/gin-gonic/gin"
)

// UserController : represent the user's controller contract
type UserController interface {
	AddUser(enforcer *casbin.Enforcer) gin.HandlerFunc
	GetUser(*gin.Context)
	GetAllUser(*gin.Context)
	SignInUser(*gin.Context)
	Register(enforcer *casbin.Enforcer) gin.HandlerFunc
	UpdateUser(*gin.Context)
	DeleteUser(*gin.Context)
}

type userController struct {
	userRepo repository.UserRepository
}

//NewUserController -> returns new user controller
func NewUserController(repo repository.UserRepository) UserController {
	return userController{
		userRepo: repo,
	}
}

// ShowAllUser
// @Schemas http
// @Summary Show all user
// @Description Get All User
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} model.User
// @Router /users/ [get]
func (h userController) GetAllUser(ctx *gin.Context) {
	user, err := h.userRepo.GetAllUser()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"code":    -1,
			"message": "something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"user":    user,
	})
}

// ShowUserByID
// @Schemas http
// @Summary Show user with corresponding ID
// @Description Get User by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security Bearer
// @Success 200 {object} model.User
// @Router /users/{id} [get]
func (h userController) GetUser(ctx *gin.Context) {
	id := ctx.Param("user")
	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "can't convert id",
		})

		return
	}
	user, err := h.userRepo.GetUser(intID)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"code":    -1,
			"message": "something went wrong",
		})
		return

	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"user":    user,
	})

}

// Signin
// @Schemas http
// @Summary Logs user into system
// @Description
// @Tags Signin
// @Accept json
// @Produce json
// @Router /signin [post]
// @Param user_info body model.User true "Username and password"
// @Success 200 {object} model.User
func (h userController) SignInUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "missing information",
		})
	}

	dbUser, err := h.userRepo.GetByEmail(user.Email)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"code":    -1,
			"message": "No Such User Found",
		})
		return

	}
	if isTrue := utils.ComparePassword(dbUser.Password, user.Password); isTrue {
		fmt.Println("user before", dbUser.ID)
		token := utils.GenerateToken(dbUser.ID)
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "Successfully Signin", "token": token})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Password not matched",
	})

}

// AddUser
// @Description Add new user into database
// @Schemas http
// @Summary Add user into database
// @Description
// @Tags User
// @Accept json
// @Produce json
// @Router /users [post]
// @Param user body model.User true "User Info"
// @Success 200 {object} model.User
// @Security Bearer
func (h userController) AddUser(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user model.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    -1,
				"message": "missing information",
			})
			return
		}
		utils.HashPassword(&user.Password)
		user, err := h.userRepo.AddUser(user)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"code":    -1,
				"message": "something went wrong",
			})
			return

		}
		enforcer.AddGroupingPolicy(fmt.Sprint(user.ID), user.Role)
		user.Password = ""
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
		})

	}
}

// Register
// @Schemas http
// @Summary Create new account
// @Description
// @Tags Register
// @Accept json
// @Produce json
// @Router /register [post]
// @Param user_info body model.User true "User info"
// @Success 200 {object} model.User
func (h userController) Register(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return h.AddUser(enforcer)
}

// UpdateUser
// @Schemas http
// @Summary Update user info
// @Description
// @Tags User
// @Accept json
// @Produce json
// @Router /users/{id} [put]
// @Param id path string true "User ID need to update"
// @Param user body model.User true "User Info need to update"
// @Success 200 {object} model.User
// @Security Bearer
func (h userController) UpdateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "missing information"})
		return
	}
	id := ctx.Param("user")
	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "user id not exist",
		})
	}
	user.ID = uint(intID)
	user, err = h.userRepo.UpdateUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"code":    -1,
			"message": "something went wrong",
		})
		return

	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"user":    user,
	})

}

// DeleteUser
// @Schemas http
// @Summary Delete user from database
// @Description
// @Tags User
// @Accept json
// @Produce json
// @Router /users/{id} [delete]
// @Param id path string true "User ID need to delete"
// @Success 200 {object} model.User
// @Security Bearer
func (h userController) DeleteUser(ctx *gin.Context) {
	var user model.User
	id := ctx.Param("user")
	intID, _ := strconv.Atoi(id)
	user.ID = uint(intID)
	user, err := h.userRepo.DeleteUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"code":    -1,
			"message": "something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
