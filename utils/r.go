package utils

import (
	"encoding/json"
	"net/http"
)

type ListData struct {
	Rows  interface{} `json:"rows"`
	Total *int64      `json:"total"`
}

type ResponseData struct {
	Code int         `json:"code"`           //相应状态码
	Msg  string      `json:"msg"`            //提示信息
	Data interface{} `json:"data,omitempty"` //数据
}

func Cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Cache-Control", "max-age=21600") // 设置缓存6小时
		return
	}
}

func Ok(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(ResponseData{Code: http.StatusOK, Msg: "操作成功", Data: data})
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func OkMsg(w http.ResponseWriter, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(ResponseData{Code: http.StatusOK, Msg: msg, Data: data})
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func Fail(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(ResponseData{Code: http.StatusInternalServerError, Msg: "请求失败, 服务端错误"})
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonData)
}

func FailMsg(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(ResponseData{Code: http.StatusInternalServerError, Msg: msg})
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonData)
}
