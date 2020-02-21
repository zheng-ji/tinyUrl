/*
 * Copyright 2011-2014 YOUMI.NET
 *
 */

package main

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/hoisie/web"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/**
 * Ping
 */
func pingHandler(ctx *web.Context) string {
	return "ok"
}

/**
 * 生成短链接
 */

func shortenHandler(ctx *web.Context) string {
	code := SUCC_INT

	conn := redisPool.Get()
	if nil == conn {
		code |= ERR_SERVER_INT
	}
	defer conn.Close()

	if !isValid(ctx.Params) {
		code |= ERR_SIGNATURE_INT
	}

	longUrl := ctx.Params["url"]
	datetime := ctx.Params["datetime"]
	comment := ctx.Params["comment"]

	var theUrl *url.URL
	var err error

	if "" == longUrl {
		code |= ERR_EMPTYURL_INT
	} else {
		theUrl, err = isValidUrl(longUrl)
		if err != nil {
			code |= ERR_URL_INT
		}
	}

	if "" == datetime {
		datetime = time.Now().Format("2006-01-02 15:04:05")
	}

	if SUCC_INT != code {
		resp := errBuilderByCode(code)
		return resp.jsonString()
	}

	var shortUrl string
	host := appConfig.Hostname

	if nil == err {
		ctr, _ := conn.Do("incr", COUNTER)
		encoded := Encode(ctr.(int64))

		if !strings.HasPrefix(host, HTTP) {
			shortUrl = fmt.Sprintf("%s://%s/%s", HTTP, host, encoded)
		} else {
			shortUrl = fmt.Sprintf("%s/%s", host, encoded)
		}

		isStat := 0
		store(encoded, theUrl.String(), datetime, comment, isStat)

		resp := Resp{
			Code: SUCC_INT,
			Msg:  SUCC_STR,
			Data: shortUrl,
		}
		return resp.jsonString()

	} else {
		http.Redirect(ctx.ResponseWriter, ctx.Request, ROLL, http.StatusNotFound)
	}

	return ""
}

// 刷新短地址并重定向
func modifyHandler(ctx *web.Context) string {
	code := SUCC_INT

	surlPostfix := ctx.Params["surl_postfix"]

	if "" == surlPostfix {
		code |= ERR_SURL_INT
		resp := Resp{
			Code: ERR_SERVER_INT,
			Msg:  ERR_SERVER_STR,
		}
		return resp.jsonString()
	}

	if !updateResolve(surlPostfix) {
		resp := Resp{
			Code: ERR_SERVER_INT,
			Msg:  ERR_SERVER_STR,
		}
		return resp.jsonString()
	} else {
		resp := Resp{
			Code: SUCC_INT,
			Msg:  SUCC_STR,
		}
		return resp.jsonString()
	}
}

// 翻译短地址并重定向
func resolveHandler(ctx *web.Context, short string) {

	if "" == short {
		seelog.Errorf("[Resolve] short is null")
		return
	}

	urlinfo := getInfoByShort(short)
	if "" == urlinfo.longUrl {
		http.Redirect(ctx.ResponseWriter, ctx.Request, ROLL, http.StatusMovedPermanently)
		seelog.Errorf("[Resolve] longurl is null")
	} else {
		redirectAndStat(ctx.Request, ctx.ResponseWriter, short, urlinfo)
	}

}
