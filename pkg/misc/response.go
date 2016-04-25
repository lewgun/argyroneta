package misc

import "github.com/gin-gonic/gin"

func FastResponse(c *gin.Context, any interface{}) {

	var obj interface{}

	const (
		result = "result"
	)

	obj = map[string]interface{}{result: "success"}

	switch param := any.(type) {
	case error:
		if param != nil {
			c.Error(param)
			obj = gin.H{"result": "fail", "faildesc": param.Error()}
		}

	case gin.H:
		param[result] = "success"
		obj = param

	case map[string]interface{}:

		if _, ok := param[result]; !ok {
			param[result] = "success"
		}

		obj = param

	case string:
		obj = any

	default:
		obj = gin.H{
			"result": "success",
			"data":   any,
		}
	}

	c.JSON(200, obj)
}
