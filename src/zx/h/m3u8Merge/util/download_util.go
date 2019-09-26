package util

import (
	"github.com/parnurzeal/gorequest"
	"time"
	"net/http"
)

/**
下载相关
 */
var BaseRequest =gorequest.New().
	Retry(10,0 *time.Second,http.StatusBadRequest, http.StatusInternalServerError)

