package global

const (

	// SUCC CODE
	SUCC = 0

	ErrCodeLackAppid = -1
	ErrStrLackAppid  = "miss appid"

	ErrCodeBadRt = -2
	ErrStrBadRt  = "rt illegal"

	ErrCodeTimeoutRt = -3
	ErrStrTimeoutRt  = "rt timeout"

	ErrCodeBadSign = -4
	ErrStrBadSign  = "bad sign,auth fail"

	ErrCodeMissArg = -5
	ErrStrMissArg  = "miss argument in url"

	ErrCodeInternal = -6
	ErrStrInternal  = "interanl error" //redis请求失败, json unmarl Fail

	ErrCodeRPC = -7
	ErrStrRPC  = "rpc error"

	ExpireSec = 300

	ReqRecordTopic = "reqtopic"

    HTTP = "http"
    FLAG = "flag"

    COUNTER = "tiny_url_counter"
)
