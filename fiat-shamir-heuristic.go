package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"strconv"
)

var (
	n = 997
	g = 3
	password = "hello"
)

// Returns a three-tuple (gcd, x, y) such that
//    a * x + b * y == gcd, where gcd is the greatest
//    common divisor of a and b.
// This function implements the extended Euclidean
//    algorithm and runs in O(log b) in the worst case.
func extendedEuclideanAlgorithm(a, b int64) (int64, int64, int64) {
	s, oldS := int64(0), int64(1)
	t, oldT := int64(1), int64(0)
	r, oldR := b, a
	for r != 0 {
		quotient := oldR
		oldR, r = r, oldR - quotient * r
		oldS, s = s, oldS - quotient * s
		oldT, t = t, oldT - quotient * t
	}
	return oldR, oldS, oldT
}

// Returns the multiplicative inverse of
//    n modulo p.
// This function returns an integer m such that
//    (n * m) % p == 1.
func inverseOf(n, p int64) int64 {
	gcd, x, y := extendedEuclideanAlgorithm(n, p)
	// assert
	if (n * x + p * y) % p != gcd {
		panic("Error with inverseOf function.")
	}
	if gcd != 1 {
		panic("Has no multiplicative inverse.")
	} else {
		return x % p
	}
}

func pickg(p int) int {
	var next, rand, exp int
	for x:=1; x<=p; x++ {
		rand = x
		exp = 1
		next = rand % p
		for next != 1 {
			next = (next * rand) % p
			exp = exp + 1
		}
		if exp == p-1 {
			return rand
		}
	}
	return 0
}

func main() {
	var Result *big.Int
	var err error
Loop:
	v := rand.Intn(n)
	if v == 0 {
		goto Loop
	}
Loop1:
	c := rand.Intn(n)
	if c == 0 {
		goto Loop1
	}
	if len(os.Args) > 1 {
		password = os.Args[1]
	}
	if len(os.Args) > 2 {
		v, err = strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
	}
	if len(os.Args) > 3 {
		c, err = strconv.Atoi(os.Args[3])
		if err != nil {
			panic(err)
		}
	}
	if len(os.Args) > 4 {
		n, err = strconv.Atoi(os.Args[4])
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Password:\t", password)
	h := md5.New()
	io.WriteString(h, password)
	md5Value := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(md5Value)
	tmp := []byte{'0','x'}
	tmp = append(tmp, md5Value...)
	// 解决了10进制转换到16进制的问题
	value, err := strconv.ParseInt(string(tmp[:10]), 0, 64)
	if err != nil {
		panic(err)
	}
	fmt.Println("value is", value)
	x := value % int64(n)
	// 截止这里没问题
	g = pickg(n)
	// 需要自己做一个bigint
	y := new(big.Int).Exp(big.NewInt(int64(g)), big.NewInt(x), big.NewInt(int64(n)))
	fmt.Println(reflect.TypeOf(y))
	t := new(big.Int).Exp(big.NewInt(int64(g)), big.NewInt(int64(v)), big.NewInt(int64(n)))
	//t := int(math.Pow(float64(g), float64(v)))
	//t = t % n
	r := v - c * int(x)
	if r < 0 {
		// 求解模的逆元
		// pow(g,-r,n)
		a := new(big.Int).Exp(big.NewInt(int64(g)), big.NewInt(int64(-r)), big.NewInt(int64(n)))
		fmt.Println(a)
		//a := int(math.Pow(float64(g), float64(-r)))
		//a = a % n
		// pow(y,c,n)
		b := new(big.Int).Exp(y, big.NewInt(int64(c)), big.NewInt(int64(n)))
		fmt.Println(b)
		//b := int(math.Pow(float64(y), float64(c)))
		//b = b % n

		// 使用 gcd 方法实现
		gcd := big.NewInt(0)
		c := big.NewInt(0)
		gcd.GCD(c, nil, a, big.NewInt(int64(n)))
		//c.Add(c, big.NewInt(int64(n)))
		fmt.Println("c is:", c)
		d := BigMulti(c.String(), b.String())
		e, _ := strconv.Atoi(d)
		f := big.NewInt(int64(e))
		Result = new(big.Int).Exp(f, big.NewInt(1), big.NewInt(int64(n)))
	} else {
		a := new(big.Int).Exp(big.NewInt(int64(g)), big.NewInt(int64(r)), big.NewInt(int64(n)))
		//a := int(math.Pow(float64(g), float64(r)))
		//a = a % n
		b := new(big.Int).Exp(y, big.NewInt(int64(c)), big.NewInt(int64(n)))
		//b := int(math.Pow(float64(y), float64(c)))
		//b = b % n
		// 需要一个大数乘法
		c := BigMulti(a.String(), b.String())
		d, _ := strconv.Atoi(c)
		e := big.NewInt(int64(d))
		//d := big.NewInt(int64(strconv.Atoi(c)))
		Result =  new(big.Int).Exp(e, big.NewInt(1), big.NewInt(int64(n)))
	}

	fmt.Println("======Agreed parameters============")
	fmt.Println("P=", n, "prime number")
	fmt.Println("G=", g, "generator")

	fmt.Println("======The secret==================")
	fmt.Println("X=", x, "alice's secret")

	fmt.Println("======Random values===============")
	fmt.Println("C=", c)
	fmt.Println("V=", v)

	fmt.Println("======Shared value===============")
	fmt.Println("g^x mod P=", y)
	fmt.Println("r=", r)

	fmt.Println("=========Results===================")
	fmt.Println("t=g**v % n =", t)
	fmt.Println("(g**r) * (y**c) =", Result)

	if t.String() == Result.String() {
		fmt.Println("Alice has proven she knows password")
	} else {
		fmt.Println("Alice has not proven she knows x")
	}
}

// 大数相乘
func BigMulti(a, b string) string {
	if a == "0" || b == "0" {
		return "0"
	}
	// string转换成[]byte，容易取得相应位上的具体值
	bsi := []byte(a)
	bsj := []byte(b)

	temp := make([]int, len(bsi)+len(bsj))
	//两数相乘，结果位数不会超过两乘数位数和，即temp的长度只可能为 len(num1)+len(num2) 或 len(num1)+len(num2)-1
	// 选最大的，免得位数不够
	for i := 0; i < len(bsi); i++ {
		for j := 0; j < len(bsj); j++ {
			// 对应每个位上的乘积，直接累加存入 temp 中相应的位置
			temp[i+j+1] += int(bsi[i]-'0') * int(bsj[j]-'0')
		}
	}

	//统一处理进位
	for i := len(temp) - 1; i > 0; i-- {
		temp[i-1] += temp[i] / 10 //对该结果进位（进到前一位）
		temp[i] = temp[i] % 10    //对个位数保留
	}

	// a 和 b 较小的时候，temp的首位为0
	// 为避免输出结果以0开头，需要去掉temp的0首位
	if temp[0] == 0 {
		temp = temp[1:]
	}
	//转换结果：将[]int类型的temp转成[]byte类型,
	//因为在未处理进位的情况下，temp每位的结果可能超过255(go中，byte类型实为uint8，最大为255),所以temp选用[]int类型
	//但在处理完进位后，不再会出现溢出
	res := make([]byte, len(temp)) //res 存放最终结果的ASCII码

	for i := 0; i < len(temp); i++ {
		res[i] = byte(temp[i] + '0')
	}
	return string(res)
}

