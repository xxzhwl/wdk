// Package uhttp 包描述
// Author: wanlizhan
// Date: 2023/4/9
package uhttp

import (
	"bytes"
	"fmt"
	"github.com/xxzhwl/wdk"
	"github.com/xxzhwl/wdk/message"
	"github.com/xxzhwl/wdk/ucontext"
	"github.com/xxzhwl/wdk/ulog"
	"github.com/xxzhwl/wdk/ustr"
	"io"
	"net/http"
	"strconv"
	"time"
)

type HttpClient struct {
	Timeout       time.Duration
	EnableRetry   bool
	EnableAlarm   bool
	LastTimeAlarm bool
	RetryDelay    time.Duration
	RetryTimes    int
}

func NewDefaultClient() HttpClient {
	return HttpClient{}
}

func NewHttpClient(timeOut, retryDelay time.Duration, retryTimes int) HttpClient {
	if retryTimes > 10 {
		retryTimes = 10
	}
	if retryDelay > time.Second*180 {
		retryDelay = time.Second * 180
	}
	return HttpClient{Timeout: timeOut, EnableRetry: true, EnableAlarm: true, RetryTimes: retryTimes, RetryDelay: retryDelay}
}

type PostArg struct {
	Title        string
	Url          string
	Body         []byte
	Header       map[string][]string
	CallbackFunc func(response []byte) error
}

func (h HttpClient) Post(arg PostArg) (resp []byte, err error) {
	return h.request(reqArg{
		Method:       http.MethodPost,
		Title:        arg.Title,
		Url:          arg.Url,
		Body:         arg.Body,
		Header:       arg.Header,
		CallbackFunc: arg.CallbackFunc,
	})
}

type GetArg struct {
	Title        string
	Url          string
	Params       map[string]any
	Body         []byte
	Header       map[string][]string
	CallbackFunc func(response []byte) error
}

func (h HttpClient) Get(arg GetArg) (response []byte, err error) {
	paramUrl := ""
	for k, v := range arg.Params {
		paramUrl = paramUrl + k + "=" + v.(string)
	}
	if len(paramUrl) != 0 {
		arg.Url = fmt.Sprintf("%s?%s", arg.Url, paramUrl)
	}
	return h.request(reqArg{
		Method:       http.MethodGet,
		Title:        arg.Title,
		Url:          arg.Url,
		Body:         arg.Body,
		Header:       arg.Header,
		CallbackFunc: arg.CallbackFunc,
	})
}

type ReqArg struct {
	Method       string
	Title        string
	Url          string
	Params       map[string]any
	Body         []byte
	Header       map[string][]string
	CallbackFunc func(response []byte) error
}

func (h HttpClient) Request(arg ReqArg) (response []byte, err error) {
	paramUrl := ""
	for k, v := range arg.Params {
		paramUrl = paramUrl + k + "=" + v.(string)
	}
	if len(paramUrl) != 0 {
		arg.Url = fmt.Sprintf("%s?%s", arg.Url, paramUrl)
	}
	return h.request(reqArg{
		Method:       arg.Method,
		Title:        arg.Title,
		Url:          arg.Url,
		Body:         arg.Body,
		Header:       arg.Header,
		CallbackFunc: arg.CallbackFunc,
	})
}

type reqArg struct {
	Method       string
	Title        string
	Url          string
	Body         []byte
	Header       map[string][]string
	CallbackFunc func(response []byte) error
}

func (h HttpClient) request(arg reqArg) (response []byte, err error) {
	//1.为header增加traceId
	ctx := ucontext.GetCurrentContext()
	if ctx == nil {
		ctx = ucontext.BuildContext()
	}
	if arg.Header == nil {
		arg.Header = make(map[string][]string)
	}
	if len(ctx.TraceId) == 0 {
		ctx.TraceId = ucontext.NewTraceId()
		arg.Header["x-trace"] = []string{ctx.TraceId}
	}
	reader := bytes.NewReader(arg.Body)
	req, err := http.NewRequest(arg.Method, arg.Url, reader)
	if err != nil {
		return nil, err
	}
	req.Header = arg.Header

	client := http.Client{Timeout: h.Timeout}
	var resp *http.Response
	var errTemp error
	var retryTimes = 0
	logData := ulog.ReqLogData{
		HttpMethod: arg.Method,
		Path:       arg.Url,
		Title:      arg.Title,
		TraceId:    ctx.TraceId,
		LocalId:    ctx.LocalId,
		ReqId:      ctx.RequestId,
		Request:    string(arg.Body),
	}
	for i := 0; i <= h.RetryTimes; i++ {
		s1 := time.Now()
		resp, errTemp = client.Do(req)
		s2 := time.Now()
		logData.StartTime = s1.Format(time.DateTime)
		logData.EndTime = s2.Format(time.DateTime)
		logData.Duration = strconv.FormatInt(s2.Sub(s1).Milliseconds(), 10)
		if errTemp != nil {
			logData.ErrMsg = errTemp.Error()
			h.sendAlarm(arg, errTemp.Error(), retryTimes)
			//允许重试且是超时的时候自动重试
			if h.EnableRetry && ustr.Contains(errTemp.Error(), "Timeout") {
				//当i>=1的时候认为开始重试
				if i >= 1 {
					retryTimes++
				}
				if retryTimes >= h.RetryTimes {
					logData.ErrMsg += fmt.Sprintf("请求失败,重试%d次失败,不再重试", h.RetryTimes)
				} else {
					logData.ErrMsg += fmt.Sprintf("请求失败,尝试%.1f秒以后重试", h.RetryDelay.Seconds())
				}
				ulog.ReqLog(logData)
				time.Sleep(h.RetryDelay)
				continue
			} else {
				logData.ErrMsg += fmt.Sprintf("请求失败")
				ulog.ReqLog(logData)
				break
			}
		}
		var respTemp []byte
		if resp != nil {
			respTemp, err = io.ReadAll(resp.Body)
			h.callBack(arg.CallbackFunc, respTemp)
			logData.Code = strconv.Itoa(resp.StatusCode)
			logData.Status = resp.Status
			logData.Response = string(respTemp)
			if err != nil {
				logData.ErrMsg = err.Error()
			}
			logData.Success = true
			ulog.ReqLog(logData)
			break
		}
	}
	defer func(resp *http.Response) {
		if resp == nil {
			return
		}
		err = resp.Body.Close()
		if err != nil {
			ulog.ErrorF("HttpRequest", "[Method:%s],Title:%s,ReqUrl:%s,Header:%v,Response:%s,Err:%s,"+
				"RetryTimes:%d", arg.Title, arg.Method, arg.Url, req.Header, string(arg.Body), err.Error(), retryTimes)
		}
	}(resp)
	return
}

func (h HttpClient) sendAlarm(arg reqArg, errStr string, retryTimes int) {
	if !h.EnableAlarm {
		return
	}
	ctx := ucontext.GetCurrentContext()
	content := fmt.Sprintf("请求地址：%s\n请求方法：%s\n请求头：%s\n请求体：%s\n日志Id：%s\nTraceId：%s\nReqId:%s\n错误信息:%s",
		arg.Url, arg.Method, arg.Header, string(arg.Body), ctx.LocalId, ctx.TraceId, ctx.RequestId, errStr)
	if !h.EnableRetry || !h.LastTimeAlarm || retryTimes >= h.RetryTimes-1 {
		message.SendAlarmMessage("Http请求失败-"+arg.Title, content)
	}
}

func (h HttpClient) callBack(fn func(resp []byte) error, respTemp []byte) {
	defer func() {
		wdk.CatchPanic()
	}()

	if fn == nil {
		return
	}

	callBackErr := fn(respTemp)
	if callBackErr != nil {
		ulog.Error("HttpRequest", fmt.Sprintf("回调执行失败：%s", callBackErr.Error()))
	}
}
