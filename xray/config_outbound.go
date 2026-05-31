package xray

import (
	"encoding/json"
	"strconv"
	"strings"
)

// ========================================================  传输协议设置 ==============================================
func genRawSetting(headerType string, host string, path string) map[string]interface{} {
	if headerType == "http" {
		return map[string]interface{}{
			"header": map[string]interface{}{
				"type": "http",
			},
			"request": map[string]interface{}{
				"version": "1.1",
				"method":  "GET",
				"path":    strings.Split(path, ","),
				"headers": map[string]interface{}{
					"Host":            strings.Split(host, ","),
					"User-Agent":      []string{},
					"Accept-Encoding": []string{"gzip, deflate"},
					"Connection":      []string{"keep-alive"},
					"Pragma":          "no-cache",
				},
			},
		}
	} else {
		return map[string]interface{}{
			"header": map[string]interface{}{
				"type": "none",
			},
		}
	}
}

func genMkcpSetting(mtuStr string) map[string]interface{} {
	mtu, err := strconv.Atoi(mtuStr)
	if err != nil {
		mtu = 1350
	}
	return map[string]interface{}{
		"mtu":              mtu,
		"tti":              50,
		"uplinkCapacity":   12,
		"downlinkCapacity": 100,
		"congestion":       false,
		"readBufferSize":   2,
		"writeBufferSize":  2,
	}
}

func genGrpcSetting(mode string, serviceName string, authority string) map[string]interface{} {
	setting := map[string]interface{}{
		"serviceName":           serviceName,
		"multiMode":             mode == "multi",
		"idle_timeout":          60,
		"health_check_timeout":  20,
		"permit_without_stream": false,
		"initial_windows_size":  0,
	}
	if authority != "" {
		setting["authority"] = authority
	}
	return setting
}

func genWsSetting(host string, path string) map[string]interface{} {
	setting := map[string]interface{}{}
	if host != "" {
		setting["host"] = host
	}
	if path != "" {
		setting["path"] = path
	}
	return setting
}

func genHttpSetting(host string, path string) map[string]interface{} {
	setting := map[string]interface{}{
		"host": strings.Split(host, ","),
		"path": path,
	}
	return setting
}

func genXhttpSetting(host string, path string, mode string, extra string) map[string]interface{} {
	setting := map[string]interface{}{}
	if host != "" {
		setting["host"] = host
	}
	if path != "" {
		setting["path"] = path
	}
	if path != "" {
		setting["mode"] = mode
	}
	if extra != "" {
		var extraMap map[string]interface{}
		err := json.Unmarshal([]byte(extra), &extraMap)
		if err == nil {
			setting["extra"] = extraMap
		}
	}
	return setting
}

func genHttpupgradeSetting(host string, path string) map[string]interface{} {
	setting := map[string]interface{}{}
	if host != "" {
		setting["host"] = host
	}
	if path != "" {
		setting["path"] = path
	}
	return setting
}

//  ========================================================  传输安全设置 ==============================================

func genTlsSetting(sni string, alpn string, fp string, allowInsecure bool) map[string]interface{} {
	setting := map[string]interface{}{
		"allowInsecure": allowInsecure,
	}
	if sni != "" {
		setting["serverName"] = sni
	}
	if alpn != "" {
		setting["alpn"] = strings.Split(alpn, ",")
	}
	if fp != "" {
		setting["fingerprint"] = fp
	}
	return setting
}

func genRealitySetting(sni string, fp string, pbk string, sid string, spx string, pqv string) map[string]interface{} {
	setting := map[string]interface{}{
		"serverName":    sni,
		"show":          false,
		"publicKey":     pbk,
		"shortId":       sid,
		"spiderX":       spx,
		"mldsa65Verify": pqv,
	}
	if fp != "" {
		setting["fingerprint"] = fp
	}
	return setting
}


// ===================================================== 其他 ===================================================
func genFinalmask(fm string) map[string]interface{} {
	if fm != "" {
		var fmMap map[string]interface{}
		err := json.Unmarshal([]byte(fm), &fmMap)
		if err == nil {
			return fmMap
		}
	}
	return map[string]interface{} {}
}