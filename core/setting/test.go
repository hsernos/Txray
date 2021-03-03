package setting

import "Txray/log"

func TestSetting() test {
	return setting.Test
}

// 设置超时时间
func SetTimeOut(time uint) {
	defer setting.save()
	setting.Test.TimeOut = time
	log.Infof("设置超时时间为: '%ds'", time)
}

// 设置测试网站
func SetTestUrl(url string) {
	defer setting.save()
	setting.Test.Url = url
	log.Infof("设置测试网站为: '%s'", url)
}
