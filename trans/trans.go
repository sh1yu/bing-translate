package trans

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var IGreg = regexp.MustCompile(`IG:"(\w+)"`)
var IID = "translator.5038.1"

type BingTranslator struct {
	ig     string
	cookie string
	//index  int
	client http.Client
}

//func (b *BingTranslator) getIndex() string {
//	b.index ++
//	return strconv.Itoa(b.index)
//}

func NewBingTranslator(timeout int) (BingTranslator, error) {

	t := 5
	if timeout > 0 {
		t = timeout
	}

	b := BingTranslator{
		client: http.Client{Timeout: time.Duration(t) * time.Millisecond},
	}
	req, err := http.NewRequest("GET", "https://cn.bing.com/translator/", nil)
	if err != nil {
		return b, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:65.0) Gecko/20100101 Firefox/65.0")
	req.Header.Set("Host", "cn.bing.com")
	req.Header.Set("Cookie", "MUID=0C680EAF606261FB0E69022164626208; SRCHD=AF=NOFORM; SRCHUID=V=2&GUID=0AF2C07A16DB474DADDEA7FC0729803A&dmnchg=1; SRCHUSR=DOB=20190325&T=1553525221000; _SS=SID=01C6442729D862A92D85490D28F663BD&HV=1553526741; _EDGE_S=SID=01C6442729D862A92D85490D28F663BD; MUIDB=0C680EAF606261FB0E69022164626208; SRCHHPGUSR=WTS=63689122021; MSCC=1; _tarLang=default=es; MSTC=ST=1; _TTSS_OUT=hist=[\"en\"]; btstkn=gKLSDs1tc6ZyZSqLG5V8rl0fKV4e3fotsj1Nfx87QU9IX3rKvd2Pu4bIgnp%252F8l1o")
	//b.index = 0
	res, err := b.client.Do(req)
	if err != nil {
		return b, err
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return b, err
	}
	defer res.Body.Close()

	tokens := IGreg.FindStringSubmatch(string(content))
	if len(tokens) >= 2 {
		b.ig = tokens[1]
		return b, nil
	}
	return b, errors.New("can't find IG")
}

func (b *BingTranslator) TDetect(text string) (string, error) {
	ur := "https://cn.bing.com/tdetect?IG=" + b.ig + "&IID=" + IID
	data := "text=" + url.QueryEscape(text)
	request, err := http.NewRequest("POST", ur, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:65.0) Gecko/20100101 Firefox/65.0")
	request.Header.Set("Host", "cn.bing.com")
	request.Header.Set("Refer", "https://cn.bing.com/translator/")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("TE", "Trailers")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Cookie", "")
	res, err := b.client.Do(request)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", errors.New("invalid status" + res.Status)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	return string(content), nil
}

func (b *BingTranslator) Translate(from, to, text string) (string, error) {
	ur := "https://cn.bing.com/ttranslate?&category=&IG=" + b.ig + "&IID=" + IID
	data := "text=" + url.QueryEscape(text) + "&from=" + from + "&to=" + to
	request, err := http.NewRequest("POST", ur, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:65.0) Gecko/20100101 Firefox/65.0")
	request.Header.Set("Host", "cn.bing.com")
	request.Header.Set("Refer", "https://cn.bing.com/translator/")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("TE", "Trailers")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Cookie", "")
	res, err := b.client.Do(request)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", errors.New("invalid status" + res.Status)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		return "", errors.New("translate result unmarshal error, content:" + string(content))
	}
	if result["statusCode"].(float64) != 200 {
		return "", errors.New("invalid status" + strconv.FormatFloat(result["statusCode"].(float64), 'f', 0, 64))
	}
	return result["translationResponse"].(string), nil
}

func (b *BingTranslator) TSpeak(lang, text string) ([]byte, error) {
	ur := "https://cn.bing.com/tspeak?&format=audio%2Fmp3&language=" + lang + "&IG=" + b.ig +
		"&IID=" + IID + "&options=female&text=" + text
	request, err := http.NewRequest("GET", ur, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:65.0) Gecko/20100101 Firefox/65.0")
	request.Header.Set("Host", "cn.bing.com")
	request.Header.Set("Refer", "https://cn.bing.com/translator/")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("TE", "Trailers")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Cookie", "")
	res, err := b.client.Do(request)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status" + res.Status)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return content, nil
}
