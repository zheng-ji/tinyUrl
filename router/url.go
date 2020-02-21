package router

import (
	"tinyurl/global"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"path"
	"time"
)

/**
 * 将s_短地址 -> l_长地址
 */
func updateCache(surlPostfix string, longUrl string) bool {

	conn := global.GRedisPool.Get()

	if nil == conn {
		return false
	}
	defer conn.Close()

	var err error

	FmtKey := fmt.Sprintf("%s_%s", FLAG, surlPostfix)

	_, err = conn.Do("HSET", FmtKey, LURL, longUrl)
	if err != nil {
		global.Runlogger.Errorf("[Redis][HSET] error: %v", err)
		return false
	}

	_, err = conn.Do("HSET", FmtKey, POSTFIX, surlPostfix)
	if err != nil {
		global.Runlogger.Errorf("[Redis][HSET] error: %v", err)
		return false
	}

	_, err = conn.Do("HINCRBY", FmtKey, PV, 1)
	if err != nil {
		global.Runlogger.Errorf("[Redis][HSET] error: %v", err)
		return false
	}

	_, err = conn.Do("HINCRBY", FmtKey, UV, 1)
	if err != nil {
		global.Runlogger.Errorf("[Redis][HINCRBY] error: %v", err)
		return false
	}
	return true
}

func storeToDB(surlPostfix string, longUrl string, datetime string, comment string) bool {
	surl := path.Join(global.GlobalConfig.Hostname, surlPostfix)
	stmtStr := "INSERT INTO `tinyurl` (`surl_postfix`, `surl`, `longurl`, `datetime`, `comment`) VALUES (?, ?, ?, ?, ?);"
	stmt, err := global.GDB.Prepare(stmtStr)

	if err != nil {
		global.Runlogger.Errorf("[Mysql][Insert] error: %v", err)
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(surlPostfix, surl, longUrl, datetime, comment, isStat)
	if err != nil {
		global.Runlogger.Errorf("[Mysql][Insert] stmt Exec error: %v", err)
		return false
	}

	if !updateCache(surlPostfix, longUrl) {
		global.Runlogger.Errorf("[Redis][Insert] error: %v", err)
		return false
	}
	return true
}

/**
 * 获取短URL
 */
func getInfoByShort(surlPostfix string) (urlinfo UrlInfo) {
	conn := global.GRedisPool.Get()

	if conn == nil {
		return
	}
	defer conn.Close()

	shortUrl := fmt.Sprintf("%s/%s", global.GlobalConfig.Hostname, surlPostfix)

	FmtKey := fmt.Sprintf("%s_%s", global.FLAG, surlPostfix)

	isExist, err := redis.Int(conn.Do("HEXISTS", FmtKey, LURL))

	if 1 == isExist {

		longUrl, err := redis.String(conn.Do("HGET", FmtKey, LURL))
		if err != nil {
			global.Runlogger.Errorf("[Redis][HGET] error: %v", err)
			return
		}

		urlinfo.shortUrl = shortUrl
		urlinfo.longUrl = longUrl
		urlinfo.postfix = surlPostfix

	} else {

		var longUrl string

		stmtStr := "SELECT `longurl` from `tinyurl` where surl_postfix=?"

		err = globalDb.QueryRow(stmtStr, surlPostfix).Scan(&longUrl)
		if err != nil {
			global.Runlogger.Errorf("[Mysql] Query  error: %v", err)
			return
		}

		if !updateCache(surlPostfix, longUrl) {
			global.Runlogger.Errorf("[Redis][Insert] error: %v", err)
			return
		}
		urlinfo.shortUrl = shortUrl
		urlinfo.longUrl = longUrl
		urlinfo.postfix = surlPostfix
	}

	return
}

/*
* 修改 短链->长链 映射关系
 */
func updateResolve(surlPostfix string) bool {
	conn := global.GRedisPool.Get()

	if nil == conn {
		return false
	}
	defer conn.Close()

	var err error
	var longurl string
	var isStat int

	stmtStr := "SELECT `longurl`, `isStat` from `tinyurl` where surl_postfix=?"
	err = globalDb.QueryRow(stmtStr, surlPostfix).Scan(&longurl, &isStat)

	if err != nil && err != sql.ErrNoRows {
		global.Runlogger.Errorf("[Mysql] Query  error: %v", err)
		return false
	}

	FmtKey := fmt.Sprintf("%s_%s", FLAG, surlPostfix)

	_, err = conn.Do("HSET", FmtKey, LURL, longurl)
	if err != nil {
		global.Runlogger.Errorf("[Redis][HSET] error: %v", err)
		return false
	}

	_, err = conn.Do("HSET", FmtKey, POSTFIX, surlPostfix)
	if err != nil {
		global.Runlogger.Errorf("[Redis][HSET] error: %v", err)
		return false
	}

	return true
}

/**
 * 重定向并统计
 */
func redirectAndStat(r *http.Request, w http.ResponseWriter, surlPostfix string, urlinfo UrlInfo) {

	var err error
    conn := global.GRedisPool.Get()
	if conn == nil {
		return
	}
	defer conn.Close()

	cookie_name := fmt.Sprintf("%s_%s", COOKIE_PRE, surlPostfix)
	_, err = r.Cookie(cookie_name)

	//没有cookie存在，计算UV
	if err == http.ErrNoCookie {
		cookie := http.Cookie{
			Name:    cookie_name,
			Value:   md5Encode(surlPostfix + salt),
			Path:    "/",
			Expires: time.Now().AddDate(0, 0, 1),
		}
		http.SetCookie(w, &cookie)

		//TODO incr short_uv, log
		FmtKey := fmt.Sprintf("%s_%s", FLAG, surlPostfix)
		_, err = conn.Do("HINCRBY", FmtKey, UV, 1)
		if err != nil {
			global.Runlogger.Errorf("[Redis][HINCRYBY] error: %v", err)
			return
		}
	}
	//异步跳转
	go http.Redirect(w, r, urlinfo.longUrl, http.StatusMovedPermanently)

	FmtKey := fmt.Sprintf("%s_%s", FLAG, surlPostfix)
	_, err = conn.Do("HINCRBY", FmtKey, PV, 1)
	if err != nil {
		global.Runlogger.Errorf("error incr short_pv")
	}
	doLog(r, urlinfo)
	return
}

func isValidUrl(rawUrl string) (u *url.URL, err error) {

	if 0 == len(rawUrl) {
		return nil, errors.New("empty url")
	}

	if !strings.HasPrefix(rawUrl, HTTP) {
		rawUrl = fmt.Sprintf("%s://%s", HTTP, rawurl)
	}

	return url.Parse(rawUrl)
}
