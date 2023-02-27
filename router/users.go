package router

import (
	"net/http"

	"faceit/model"
	"faceit/repository"
	"faceit/webhooks"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type userURI struct {
	Id int `uri:"userId"`
}

func LoadUsers(e *gin.Engine) {
	usersGroup := e.Group("/users")

	usersGroup.POST("", createUser)
	usersGroup.DELETE(":userId", deleteUser)
	usersGroup.PUT(":userId", updateUser)
	usersGroup.GET("", getUsers)
}

type userData struct {
	FirstName string `json:"first_name"    binding:"required"`
	LastName  string `json:"last_name"     binding:"required"`
	Nickname  string `json:"nickname"      binding:"required"`
	Password  string `json:"password"      binding:"required"`
	Email     string `json:"email"         binding:"required"`
	Country   string `json:"country"       binding:"required"`
}

func createUser(c *gin.Context) {
	var data userData
	if err := c.ShouldBindJSON(&data); err != nil {
		// todo: humanize messages
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Nickname:  data.Nickname,
		Password:  data.Password,
		Email:     data.Email,
		Country:   data.Country,
	}

	usersRepository, err := repository.NewUsersRepository()
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error while getting userRepository")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	insertedUser, err := usersRepository.Insert(user)
	if err != nil {
		if e, ok := err.(*mysql.MySQLError); ok {
			if e.Number == 1062 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "email already used"})
				return
			}
		}

		log.Error().
			Err(err).
			Msgf("error while inserting user")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	err = webhooks.DispatchEvent(webhooks.EventUserCreated, insertedUser)
	if err != nil {
		log.Error().
			Err(err).
			Str("event", webhooks.EventUserCreated).
			Msgf("error while dispatching event")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": insertedUser})
}

func deleteUser(c *gin.Context) {
	var userURI userURI
	if err := c.ShouldBindUri(&userURI); err != nil {
		// todo: humanize messages
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usersRepository, err := repository.NewUsersRepository()
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error while getting userRepository")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	err = usersRepository.Delete(userURI.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		
		log.Error().
			Err(err).
			Msgf("error while deleting user")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	// todo: add entire user model to the webhook
	err = webhooks.DispatchEvent(webhooks.EventUserDeleted, model.User{
		ID: uint(userURI.Id),
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("event", webhooks.EventUserDeleted).
			Msgf("error while dispatching event")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func updateUser(c *gin.Context) {
	var userURI userURI
	if err := c.ShouldBindUri(&userURI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data userData
	if err := c.ShouldBindJSON(&data); err != nil {
		// todo: humanize messages
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		ID:        uint(userURI.Id),
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Nickname:  data.Nickname,
		Password:  data.Password,
		Email:     data.Email,
		Country:   data.Email,
	}

	usersRepository, err := repository.NewUsersRepository()
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error while getting userRepository")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	newUser, err := usersRepository.Update(user) // todo bug: created_at is empty
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}

		log.Error().
			Err(err).
			Msgf("error while updating user")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	err = webhooks.DispatchEvent(webhooks.EventUserUpdated, newUser)

	if err != nil {
		log.Error().
			Err(err).
			Str("event", webhooks.EventUserUpdated).
			Msgf("error while dispatching event")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": newUser, // todo: hide created_at, updated_at
	})
}

type userFilters struct {
	Country  string `form:"country"`
	Offset   int    `form:"offset"     binding:"gte=0"`
	PageSize int    `form:"page_size"  binding:"gte=1,lte=100"`
}

func getUsers(c *gin.Context) {
	usersRepository, err := repository.NewUsersRepository()
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error while getting userRepository")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	userFilters := userFilters{}
	if err := c.ShouldBindQuery(&userFilters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userFilters.PageSize == 0 {
		userFilters.PageSize = 10 // todo: move to const
	}

	users, err := usersRepository.GetAll(
		userFilters.Offset,
		userFilters.PageSize,
		model.User{
			Country: userFilters.Country,
		},
	)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("error while getting users")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
