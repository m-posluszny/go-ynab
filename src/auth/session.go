package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-posluszny/go-ynab/src/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

const userKey = "uid"

func InitAuthSession(authCfg config.AuthConf, redisConf config.RedisConf) gin.HandlerFunc {
	session, err := redis.NewStore(redisConf.Size, "tcp", redisConf.Host, redisConf.Password, []byte(authCfg.Secret))
	if err != nil {
		panic(err)
	}
	return sessions.Sessions("auth_session", session)
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get(userKey)
	if uid == nil {
		c.Redirect(http.StatusConflict, "/login")
		return
	}
	c.Next()
}

func CreateSession(c *gin.Context, userId string) error {
	session := sessions.Default(c)
	session.Set(userKey, userId)
	if err := session.Save(); err != nil {
		return errors.New("cannot login")
	}
	return nil

}

func DeleteSession(c *gin.Context) error {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		return errors.New("session invalid")
	}
	session.Delete(userKey)
	if err := session.Save(); err != nil {
		return errors.New("cannot logout")
	}
	return nil

}

func GetUIDFromSession(c *gin.Context) string {
	session := sessions.Default(c)
	uid := session.Get(userKey)
	fmt.Println("CHUID", uid)
	if uid == nil {
		return ""
	}
	return uid.(string)
}
