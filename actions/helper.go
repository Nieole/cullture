package actions

import (
	"culture/models"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
)

//Monitor  Monitor
var Monitor = flag.String("monitor-url", envy.Get("monitor", "http://localhost:8080"), "this is monitor url")

//Client  Client
var Client = http.Client{
	Timeout: time.Second * 5,
}

//List 返回列表对象
func List(data []Data, paginator *pop.Paginator) render.Renderer {
	return r.JSON(ListResponse{
		Data:      data,
		Paginator: paginator,
	})
}

//Fail Fail
func Fail(message string, a ...interface{}) render.Renderer {
	return r.JSON(map[string]string{"message": fmt.Sprintf(message, a...)})
}

//ListResponse 列表返回消息
type ListResponse struct {
	*pop.Paginator
	Data interface{} `json:"data"`
}

//Data 数据
type Data interface {
	String() string
	FromString(data string) error
}

//Unmarshal Unmarshal
func (l *ListResponse) Unmarshal(d []byte) error {
	return json.Unmarshal(d, l)
}

//Marshal Marshal
func (l *ListResponse) Marshal() ([]byte, error) {
	return json.Marshal(l)
}

//FindByPhone FindByPhone
func FindByPhone(phone string) (*Human, error) {
	if phone == "" {
		return nil, errors.New("phone is blank")
	}
	resp, err := Client.Get(fmt.Sprintf("%s/api/humans/phone/%s", *Monitor, phone))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to check %s", phone)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	h := new(Human)
	err = json.Unmarshal(bytes, h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

//CreateHuman CreateHuman
func CreateHuman(phone string) (*Human, error) {
	if phone == "" {
		return nil, errors.New("phone is blank")
	}
	hq := &Human{
		Name:     RandString(5),
		Role:     "Visitor",
		Sex:      "Other",
		PhoneNum: phone,
	}
	b, err := json.Marshal(hq)
	if err != nil {
		return nil, err
	}
	resp, err := Client.Post(fmt.Sprintf("%s/api/humans", *Monitor), "application/json", strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		all, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(all))
		return nil, fmt.Errorf("failed to add %s", phone)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	h := new(Human)
	err = json.Unmarshal(bytes, h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

//Human Human
type Human struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Sex      string `json:"sex"`
	IDNum    string `json:"idNum"`
	PhoneNum string `json:"phoneNum"`
	OrgName  string `json:"orgName"`
	CarNum   string `json:"carNum"`
	OrgID    string `json:"orgId"`
}

var ran *rand.Rand

func init() {
	ran = rand.New(rand.NewSource(time.Now().Unix()))
}

//RandString 随机字符串
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(ran.Intn(26) + 65)
	}
	return string(bytes)
}

//CurrentUser 获取当前用户
func CurrentUser(c buffalo.Context) (*models.User, error) {
	user, ok := c.Session().Get("current_user").(*models.User)
	if !ok {
		return nil, errors.New("未找到当前用户信息")
	}
	return user, nil
}
