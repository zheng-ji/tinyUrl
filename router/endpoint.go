package router

import (
	"tinyurl/global"
	"context"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
)

// PingHandler for Test
func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"c": 200})
}

func PongHandler(c *gin.Context) {
	c.JSON(200, gin.H{"c": 200})
}


func ShortenHandler(c *gin.Context) {

    conn := global.GRedisPool.Get()
    defer conn.Close()
	if nil == conn {
        c.JSON(http.StatusInternalServerError, gin.H{"c": global.ErrCodeInternal})
        return
	}


	longUrl, _ := c.Get("longurl")
	datetime, _ := c.Get("datetime")
	comment, _ := c.Get("comment")
	if "" == longUrl || "" == datetime {
        c.JSON(http.StatusInternalServerError, gin.H{"c": global.ErrStrMissArg})
        return
	}

	var theUrl *url.URL
	var err error
    theUrl, err = isValidUrl(longUrl)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"c": global.ErrCodeInterval})
    }

    exists, err := redis.Bool(conn.Do("EXISTS", COUNTER))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"c": global.ErrStrInternal})
        return
    }

	if !exists {
		var id string
		stmtStr := "SELECT MAX(id) from `tinyurl`"
        err := global.GDB.QueryRow(stmtStr).Scan(&id)
		if err != nil {
            global.Runlogger.Errorf("%v", err)
			return
		} else {
            conn.Do("SET", global.COUNTER, id)
        }
    }

    ctr, _ := conn.Do("incr", global.COUNTER)
    encoded := global.Encode(ctr.(int64))

	var shortUrl string
	host := global.GlobalConfig.Hostname
    if !strings.HasPrefix(host, global.HTTP) {
        shortUrl = fmt.Sprintf("%s://%s/%s", global.HTTP, host, encoded)
    } else {
        shortUrl = fmt.Sprintf("%s/%s", host, encoded)
    }

    storeToDB(encoded, theUrl.String(), datetime, comment)

    c.JSON(http.StatusInternalServerError, gin.H{"c": global.SUCC_INT, "d":shortUrl})
}

func ResolveHandler(c *gin.Context) {
    index := c.Param("index")

	urlinfo := getInfoByShort(index)
	if "" == urlinfo.longUrl {
		http.Redirect(ctx.ResponseWriter, ctx.Request, ROLL, http.StatusMovedPermanently)
		global.Runlogger.Errorf("[Resolve] longurl is null")
	} else {
		redirectAndStat(c.Request, c.Writer, index, urlinfo)
	}
}
