package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/config"
)

// WechatConfig 微信配置
type WechatConfig struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
}

// WechatManager 管理微信相关功能
type WechatManager struct {
	config *WechatConfig
	offAcc *officialaccount.OfficialAccount
}

// NewWechatManager 创建新的微信管理器
func NewWechatManager(cfg *WechatConfig) *WechatManager {
	wc := wechat.NewWechat()
	memoryCache := cache.NewMemory()

	offConfig := &config.Config{
		AppID:          cfg.AppID,
		AppSecret:      cfg.AppSecret,
		Token:          cfg.Token,
		EncodingAESKey: cfg.EncodingAESKey,
		Cache:          memoryCache,
	}

	offAcc := wc.GetOfficialAccount(offConfig)

	return &WechatManager{
		config: cfg,
		offAcc: offAcc,
	}
}

// StartAuthServer 启动微信授权服务器
func (wm *WechatManager) StartAuthServer() {
	http.HandleFunc("/auth", wm.handleAuth)
	http.HandleFunc("/callback", wm.handleCallback)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("启动授权服务器失败：%v\n", err)
		}
	}()
}

// handleAuth 处理微信授权请求
func (wm *WechatManager) handleAuth(w http.ResponseWriter, r *http.Request) {
	oauth := wm.offAcc.GetOauth()
	url, err := oauth.GetRedirectURL("http://localhost:8080/callback", "snsapi_userinfo", "")
	if err != nil {
		fmt.Printf("获取重定向URL失败：%v\n", err)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

// handleCallback 处理微信回调
func (wm *WechatManager) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	oauth := wm.offAcc.GetOauth()

	token, err := oauth.GetUserAccessToken(code)
	if err != nil {
		fmt.Printf("获取用户访问令牌失败：%v\n", err)
		return
	}

	userInfo, err := oauth.GetUserInfo(token.AccessToken, token.OpenID, "")
	if err != nil {
		fmt.Printf("获取用户信息失败：%v\n", err)
		return
	}

	// 返回用户信息
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}
