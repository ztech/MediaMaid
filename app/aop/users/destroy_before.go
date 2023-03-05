package users

import (
	"nasmaid/app/global/consts"
	"nasmaid/app/global/variable"

	"github.com/gin-gonic/gin"
)

// 模拟Aop 实现对某个控制器函数的前置和后置回调

type DestroyBefore struct{}

// 前置函数必须具有返回值，这样才能控制流程是否继续向下执行
func (d *DestroyBefore) Before(context *gin.Context) bool {
	userId := context.GetFloat64(consts.ValidatorPrefix + "id")
	variable.ZapLog.Sugar().Infof("模拟 Users 删除操作， Before 回调,用户ID：%.f\n", userId)
	if userId > 10 {
		return true
	} else {
		return false
	}
}
