package actions

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop"
)

//List 返回列表对象
func List(data interface{}, paginator *pop.Paginator) render.Renderer {
	return r.JSON(ListResponse{
		Data:      data,
		Paginator: paginator,
	})
}

//ListResponse 列表返回消息
type ListResponse struct {
	*pop.Paginator
	Data interface{} `json:"data"`
}
