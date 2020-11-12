package encode

import (
	"bytes"
	"io/ioutil"
	"log"
	"sync"

	iconv "github.com/qiniu/iconv"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	CVtolocal     iconv.Iconv
	CVtoutf8      iconv.Iconv
	localEncode   string
	converterLock sync.Mutex
)

func InitConverter(local string) bool {
	// cv, err := iconv.Open("utf-8", local)
	// if err != nil {
	// 	fmt.Println("Init locale converter failed ! code:1")
	// 	panic(err)
	// }
	// CVtolocal = cv
	// cv, err = iconv.Open(local, "utf-8")
	// if err != nil {
	// 	fmt.Println("Init locale converter failed ! code:2")
	// 	panic(err)
	// }
	// CVtoutf8 = cv
	localEncode = local
	return true
}

//GbkToUtf8 转换GBK编码到UTF-8编码
func GbkToUtf8(str []byte) (b []byte, err error) {
	r := transform.NewReader(bytes.NewReader(str), simplifiedchinese.GBK.NewDecoder())
	b, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}
	return
}

//Utf8ToGbk 转换UTF-8编码到GBK编码
func Utf8ToGbk(str []byte) (b []byte, err error) {
	r := transform.NewReader(bytes.NewReader(str), simplifiedchinese.GBK.NewEncoder())
	b, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}
	return
}

func Utf8ToLocal(str string) (b string, err error) {
	converterLock.Lock()
	cv, err := iconv.Open(localEncode, "utf-8")
	if err != nil {
		log.Println("locale converter failed ! code:1")
		cv.Close()
		converterLock.Unlock()
		return str, err
	}
	buf := cv.ConvString(str)
	cv.Close()
	converterLock.Unlock()
	if len(buf) <= 0 {
		return str, nil
	}
	return buf, nil
}

func LocalToUtf8(str string) (b string, err error) {
	converterLock.Lock()
	cv, err := iconv.Open("utf-8", localEncode)
	if err != nil {
		log.Println("locale converter failed ! code:2")
		cv.Close()
		converterLock.Unlock()
		return str, err
	}
	buf := cv.ConvString(str)
	cv.Close()
	converterLock.Unlock()
	if len(buf) <= 0 {
		return str, nil
	}
	return buf, nil
}
