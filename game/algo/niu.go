/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2018-01-17 01:19:59
 * Filename      : niu.go
 * Description   : 玩牌算法
 * *******************************************************/
package algo

var NIUS [][]int = [][]int{{0, 1, 2}, {0, 1, 3}, {0, 2, 3}, {1, 2, 3}, {0, 1, 4}, {0, 2, 4}, {1, 2, 4}, {0, 3, 4}, {1, 3, 4}, {2, 3, 4}}
var NIUL [][]int = [][]int{{3, 4}, {2, 4}, {1, 4}, {0, 4}, {2, 3}, {1, 3}, {0, 3}, {1, 2}, {0, 2}, {0, 1}}

//10=牛牛,0-9牛
func PointNiu(cs []uint32) (i uint32) {
	for k, v := range NIUS {
		if ((cs[v[0]] + cs[v[1]] + cs[v[2]]) % 10) != 0 {
			continue
		}
		n := (cs[NIUL[k][0]] + cs[NIUL[k][1]]) % 10
		if n == 0 {
			return 10
		}
		if n > i {
			i = n
		}
	}
	return
}

//号码1-10
func Point(cs []uint32) (i uint32) {
	//牛牛点数0-10
	if len(cs) == 5 {
		return PointNiu(cs)
	}
	//三公,牌九点数等于所有牌值和对10取余
	for _, v := range cs {
		i += v
	}
	i = i % 10
	return
}

//取大值
func max(cs []uint32) (i uint32) {
	for _, v := range cs {
		if v > i {
			i = v
		}
	}
	return
}

//第二值
func second(m uint32, cs []uint32) (i uint32) {
	for _, v := range cs {
		if m == v {
			continue
		}
		if v > i {
			i = v
		}
	}
	return
}

//大小比较
func Compare(a, b []uint32) bool {
	i := Point(a)
	j := Point(b)
	if i == j {
		m := max(a)
		n := max(b)
		//三公存在相同牌
		//2,3,6,4,1
		if m == n {
			return second(m, a) > second(n, b)
		}
		//最大牌比较
		return m > n
	}
	//牌九,牛牛不存在相同牌
	return i > j
}

//倍率 TODO 暂时所有倍率为1
func Multiple(n uint32) int64 {
	switch n {
	case 0:
	default:
	}
	return 1
}
