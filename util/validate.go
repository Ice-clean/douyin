package util

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"reflect"
)

var validate *validator.Validate = validator.New()

type Element interface{}

type User struct {
	Name  string `validate:"required"`                    //非空
	Age   uint8  `validate:"gte=0,lte=130"`               //  0<=Age<=130
	Email string `validate:"required,email"`              //非空，email格式
	Phone string `validate:"numeric,len=11" msg:"电话号码错误"` //数字类型，长度为11
	//dive关键字代表 进入到嵌套结构体进行判断
	Address []*Address `validate:"dive"` //  可以拥有多个地址
}
type Address struct {
	Province string `validate:"required"` //非空
	City     string `validate:"required"` //非空

}

func ValidateStruct(ele Element) error {
	err := validate.Struct(ele)
	if err != nil {
		//断言为：validator.ValidationErrors，类型为：[]FieldError
		for _, e := range err.(validator.ValidationErrors) {
			if success, msg := GetStructTagName(ele, e.Field(), "msg"); success {
				return errors.New(msg)
			}
		}
	} else {
		return nil
	}
	return nil
}

func main() {
	validateStruct()   //结构体校验
	validateVariable() //变量校验
}
func validateStruct() {
	address := Address{
		Province: "重庆",
		City:     "重庆",
	}
	user := User{
		Name:    "江洲",
		Age:     23,
		Email:   "jz@163.com",
		Address: []*Address{&address},
		Phone:   "13366663333x",
	}
	err := ValidateStruct(user)
	fmt.Println(err.Error())
}

//变量校验
func validateVariable() {
	myEmail := "123@qq.com" //邮箱地址：xx@xx.com
	err := validate.Var(myEmail, "required,email")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("变量校验通过！")
	}
}

// GetStructTagName 获取struct的field的tag内容
func GetStructTagName(i interface{}, fieldName string, tagName string) (bool, string) {
	types := reflect.TypeOf(i)
	name, success := types.FieldByName(fieldName)
	if success == false {
		return false, ""
	} else {
		return true, name.Tag.Get(tagName)
	}

}
