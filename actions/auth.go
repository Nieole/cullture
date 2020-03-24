package actions

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

//LoginHandler LoginHandler
func LoginHandler(c buffalo.Context) error {
	phone := c.Param("phone")
	human, err := FindByPhone(phone)
	if err != nil {
		log.Println(fmt.Sprintf("%s not found", phone))
		human, err = CreateHuman(phone)
		if err != nil {
			return c.Render(http.StatusBadRequest, Fail(err.Error()))
		}
	}
	if human != nil {
		c.Session().Set("current_user_name", human.Name)
		c.Session().Set("current_user_phone", human.PhoneNum)
		c.Session().Save()
		return c.Render(http.StatusCreated, nil)
	}
	return c.Render(http.StatusBadRequest, Fail("failed login %s", phone))
}

//SignOutHandler SignOutHandler
func SignOutHandler(c buffalo.Context) error {
	c.Session().Clear()
	return c.Render(http.StatusAccepted, nil)
}
