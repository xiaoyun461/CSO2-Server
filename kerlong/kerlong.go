package kerlong

import "fmt"

//IntAbs 绝对值
func IntAbs(num int) int {
	ans, ok := Ternary(num > 0, num, -num).(int)
	if ok {
		return ans
	}
	return 0
}

//Ternary 三目运算符
func Ternary(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

func IsAllNumber(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] > '9' || str[i] < '0' {
			return false
		}
	}
	return true
}

//ScanLine 得到一行
func ScanLine() string {
	var c byte
	var err error
	var b []byte
	for err == nil {
		_, err = fmt.Scanf("%c", &c)
		if c != '\n' {
			b = append(b, c)
		} else {
			break
		}
	}
	return string(b)
}
