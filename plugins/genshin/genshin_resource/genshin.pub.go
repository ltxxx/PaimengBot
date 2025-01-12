package genshin_resource

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

func getTodayResourceByGenshinPub(file string) (err error) {
	for i := 0; i < 3; i++ { // 最多尝试3次
		err = tryGetGenshinPubResourceShot(file)
		if err == nil { // 直到成功
			break
		}
	}
	return
}

func tryGetGenshinPubResourceShot(file string) error {
	// 创建 context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Infof),
		chromedp.WithDebugf(log.Debugf),
		chromedp.WithErrorf(log.Errorf),
	)
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// 截图
	var buf []byte
	if err := chromedp.Run(ctx,
		genshinResourceScreenshot(`https://genshin.pub/daily`, `.GSContainer_content_box__1sIXz`, &buf),
	); err != nil {
		log.Warnf("chromedp genshinResourceScreenshot err: %v", err)
		return err
	}
	if err := ioutil.WriteFile(file, buf, 0o644); err != nil {
		log.Warnf("genshinResourceScreenshot write file err: %v", err)
		return err
	}
	return nil
}

// elementScreenshot takes a screenshot of a specific element.
func genshinResourceScreenshot(url, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(1500, 1500),
		chromedp.Navigate(url),
		chromedp.WaitVisible(sel),
		chromedp.Evaluate("document.getElementsByClassName('MewBanner_root__3GKl2')[0].remove();", nil),
		chromedp.Evaluate("document.getElementsByClassName('GSContainer_gs_container__2FbUz')[0].setAttribute(\"style\", \"height:1050px\");", nil),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
