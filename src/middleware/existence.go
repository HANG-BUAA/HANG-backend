package middleware

import (
	"HANG-backend/src/api"
	"HANG-backend/src/global"
	"HANG-backend/src/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ParamLocation int

// ！！！ 该类中间件不应赋予读取请求体 json 格式的数据，因为请求体只可以读取一次！！！
const (
	URI ParamLocation = iota + 1
	QUERY
)

func locationToString(location ParamLocation) string {
	if location == URI {
		return "uri"
	} else {
		return "query"
	}
}

func entityNotFoundErr(c *gin.Context, entity string, id uint) {
	api.Fail(c, api.ResponseJson{
		Code: global.ERR_CODE_ENTITY_NOT_FOUND,
		Msg:  fmt.Sprintf("entity %s with the id %d not found", entity, id),
	})
}

func paramMissErr(c *gin.Context, loc ParamLocation, param string) {
	api.Fail(c, api.ResponseJson{
		Code: global.ERR_CODE_MISSING_PARAM,
		Msg:  fmt.Sprintf("param %s is Required as %s", param, locationToString(loc)),
	})
}

func PostExistence(location ParamLocation) gin.HandlerFunc {
	return func(c *gin.Context) {
		var postID uint
		if location == URI {
			uriPostID := c.Param("post_id")
			tmp, err := strconv.ParseUint(uriPostID, 10, 64)
			if err != nil {
				paramMissErr(c, location, "post_id")
				return
			}
			postID = uint(tmp)
		} else {
			queryPostID := c.Query("post_id")
			tmp, err := strconv.ParseUint(queryPostID, 10, 64)
			if err != nil {
				paramMissErr(c, location, "post_id")
				return
			}
			postID = uint(tmp)
		}

		// 判断 post 是否存在
		var post model.Post
		if err := global.RDB.First(&post, postID).Error; err != nil {
			entityNotFoundErr(c, "post", postID)
		}
		c.Set("post", &post)
		c.Next()
	}
}

func CommentExistence(location ParamLocation) gin.HandlerFunc {
	return func(c *gin.Context) {
		var commentID uint
		if location == URI {
			uriCommentID := c.Param("comment_id")
			tmp, err := strconv.ParseUint(uriCommentID, 10, 64)
			if err != nil {
				paramMissErr(c, location, "comment_id")
				return
			}
			commentID = uint(tmp)
		} else {
			queryCommentID := c.Query("comment_id")
			tmp, err := strconv.ParseUint(queryCommentID, 10, 64)
			if err != nil {
				paramMissErr(c, location, "comment_id")
				return
			}
			commentID = uint(tmp)
		}

		// 判断 comment 是否存在
		var comment model.Comment
		if err := global.RDB.First(&comment, commentID).Error; err != nil {
			entityNotFoundErr(c, "comment", commentID)
		}
		c.Set("comment", &comment)
		c.Next()
	}
}