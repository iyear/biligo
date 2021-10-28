package biligo

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const table = "fZodR9XQDSUm21yCkr6zBqiveYah8bt4xsWpHnJE7jL5VG3guMTKNPAwcF"

var tr = map[string]int64{}
var s = []int{11, 10, 3, 8, 4, 6}

const xor = 177451812
const add = 8728348608

func init() {
	tableByte := []byte(table)
	for i := 0; i < 58; i++ {
		tr[string(tableByte[i])] = int64(i)
	}
}

// BV2AV 带BV前缀
func BV2AV(bv string) int64 {
	var r int64
	arr := []rune(bv)

	for i := 0; i < 6; i++ {
		r += tr[string(arr[s[i]])] * int64(math.Pow(float64(58), float64(i)))
	}
	return (r - add) ^ xor
}

// AV2BV 带BV前缀
func AV2BV(av int64) string {
	x := (av ^ xor) + add
	r := []string{"B", "V", "1", " ", " ", "4", " ", "1", " ", "7", " ", " "}
	for i := 0; i < 6; i++ {
		r[s[i]] = string(table[int64(math.Floor(float64(x/int64(math.Pow(float64(58), float64(i))))))%58])
	}
	var result string
	for i := 0; i < 12; i++ {
		result += r[i]
	}
	return result
}

// parseDynaAt 由于ctrl的location是字符定位的，而FindAllStringIndex获取的是字节定位，只能遍历一遍拿到字符定位
func parseDynaAt(tp int, content string, at map[string]int64) []*dynaCtrl {
	match := regexp.MustCompile("@.*? ").FindAllStringIndex(content, -1)
	var (
		ctrl []*dynaCtrl
		a    = 0
		c    = 0
	)
	for i, t := range []rune(content) {
		c += len(fmt.Sprintf("%c", t))
		if c == match[a][0] {
			up := strings.TrimPrefix(content[match[a][0]:match[a][1]], "@")
			up = strings.TrimSuffix(up, " ")
			ctrl = append(ctrl, &dynaCtrl{
				Location: i + 1,
				Type:     tp,
				Length:   len([]rune(up)) + 2, // 之前删了@和空格，需要加回来
				Data:     strconv.FormatInt(at[up], 10),
			})
			a++
		}
		if a == len(match) {
			break
		}
	}
	return ctrl
}
