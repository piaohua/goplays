/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2018-01-17 01:19:11
 * Filename      : niu.go
 * Description   : 玩牌算法
 * *******************************************************/
package algo

import "testing"

func TestAlgo(t *testing.T) {
	a := []uint32{2, 3, 6}
	b := []uint32{6, 4, 1}
	t.Log(Compare(a, b))
	a = []uint32{2, 3, 6, 10, 7}
	b = []uint32{6, 4, 1, 8, 9}
	t.Log(Point(a))
	t.Log(Point(b))
	t.Log(Compare(a, b))
}
