package trans

import (
	"testing"
)

var b BingTranslator
var langMap map[string]string

func init() {
	langMap = map[string]string{
//		"zh-CHS": "你好",
		"en":     "How are you doing",
		"ja":     "こんにちは",
//		"fr":     "Comment vas-tu",
//		"ru":     "Как дела",
//		"ar":     "كيف حالك",
		"ko":     "어떻게 지내니",
	}
	var err error
	b, err = NewBingTransator(5000)
	if err != nil {
		panic(err)
	}
}

func TestBingTranslator_TDetect(t *testing.T) {
	for k, v := range langMap {
		from, err := b.TDetect(v)
		if err != nil {
			t.Errorf("detect %v error:%v\n", v, err.Error())
			continue
		}
		if from != k {
			t.Errorf("detect %v failed. actual:%v, expected:%v\n", v, from, k)
			continue
		}
	}
}

func TestBingTranslator_Translate(t *testing.T) {
	for k, v := range langMap {
		result, err := b.Translate("zh-CHS", k, "你好")
		if err != nil {
			t.Errorf("translate to %v error:%v\n", k, err.Error())
			continue
		}
		if result != v {
			t.Errorf("translate to %v failed. actual:%v, expected:%v\n", k, result, v)
			continue
		}
	}
}
