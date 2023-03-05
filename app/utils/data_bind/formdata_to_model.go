package data_bind

import (
	"errors"
	"nasmaid/app/global/consts"
	"reflect"

	"github.com/gin-gonic/gin"
)

const (
	modelStructMustPtr = "modelStruct 必须传递一个指针"
)

// 绑定form表单验证器已经验证完成的参数到 model 结构体,
// mode 结构体支持匿名嵌套
// 数据绑定原则：
// 1.表单参数验证器中的结构体字段 json 标签必须和 model 结构体定义的 json 标签一致
// 2.model 中的数据类型与表单参数验证器数据类型保持一致：
// 	例如：model 中的 user_name 是 string 那么表单参数验证器中的 user_name 也必须是 string，bool 类型同理，日期时间字段在 ginskeleton 中请按照 string 处理
// 3.但是 model 中的字段如果是数字类型（int、int8、int16、int64、float32、float64等）都可以绑定表单参数验证中的 float64 类型，程序会自动将原始的 float64 转换为 model 的定义的数字类型

func ShouldBindFormDataToModel(c *gin.Context, modelStruct interface{}) error {
	mTypeOf := reflect.TypeOf(modelStruct)
	if mTypeOf.Kind() != reflect.Ptr {
		return errors.New(modelStructMustPtr)
	}
	mValueOf := reflect.ValueOf(modelStruct)

	//分析 modelStruct 字段
	mValueOfEle := mValueOf.Elem()
	mtf := mValueOf.Elem().Type()
	fieldNum := mtf.NumField()
	for i := 0; i < fieldNum; i++ {
		if !mtf.Field(i).Anonymous && mtf.Field(i).Type.Kind() != reflect.Struct {
			fieldSetValue(c, mValueOfEle, mtf, i)
		} else if mtf.Field(i).Type.Kind() == reflect.Struct {
			//处理结构体(有名+匿名)
			mValueOfEle.Field(i).Set(analysisAnonymousStruct(c, mValueOfEle.Field(i)))
		}
	}
	return nil
}

// 分析匿名结构体,并且获取匿名结构体的值
func analysisAnonymousStruct(c *gin.Context, value reflect.Value) reflect.Value {

	typeOf := value.Type()
	fieldNum := typeOf.NumField()
	newStruct := reflect.New(typeOf)
	newStructElem := newStruct.Elem()
	for i := 0; i < fieldNum; i++ {
		fieldSetValue(c, newStructElem, typeOf, i)
	}
	return newStructElem
}

// 为结构体字段赋值
func fieldSetValue(c *gin.Context, valueOf reflect.Value, typeOf reflect.Type, colIndex int) {
	relaKey := typeOf.Field(colIndex).Tag.Get("json")
	if relaKey != "-" {
		relaKey = consts.ValidatorPrefix + typeOf.Field(colIndex).Tag.Get("json")
		switch typeOf.Field(colIndex).Type.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			valueOf.Field(colIndex).SetInt(int64(c.GetFloat64(relaKey)))
		case reflect.Float32, reflect.Float64:
			valueOf.Field(colIndex).SetFloat(c.GetFloat64(relaKey))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			valueOf.Field(colIndex).SetUint(uint64(c.GetFloat64(relaKey)))
		case reflect.String:
			valueOf.Field(colIndex).SetString(c.GetString(relaKey))
		case reflect.Bool:
			valueOf.Field(colIndex).SetBool(c.GetBool(relaKey))
		default:
			// model 如果有日期时间字段，请统一设置为字符串即可
		}
	}
}
