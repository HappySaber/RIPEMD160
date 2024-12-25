package ripemd160

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Ripemd160(message string) string {
	binaryText := textToBit(message)
	l_bit := len(binaryText)
	binaryText = textToAddBit(binaryText, l_bit)
	binary_l_bit := fmt.Sprintf("%064b", l_bit)
	binaryText += binary_l_bit
	var T []string
	result := ""
	h := []int{0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476, 0xC3D2E1F0}

	for i := 0; i < len(binaryText); i += 512 { // делим строку на список по 512 бит
		end := i + 512
		if end > len(binaryText) {
			end = len(binaryText)
		}
		T = append(T, binaryText[i:end])
	}

	hConst := make([]int, len(h))
	copy(hConst, h)
	for _, text512 := range T {
		var text16 []string
		for i := 0; i < len(text512); i += 32 { // делим строки на 16 слов по 32 бита
			end := i + 32
			if end > len(text512) {
				end = len(text512)
			}
			text16 = append(text16, text512[i:end])
		}
		result = mainloop(text16, hConst) // функция 80 итераций алгоритма
	}
	return result
}

func textToAddBit(binaryText string, lBit int) string {
	if lBit%512 != 448 {
		binaryText += "1"
		n := (448 - lBit%512) % 512
		for i := 0; i < n-1; i++ {
			binaryText += "0"
		}
	}
	return binaryText
}

func textToBit(text string) string {
	var binaryText strings.Builder
	for _, char := range text {
		binaryText.WriteString(fmt.Sprintf("%08b", char))
	}
	return binaryText.String()
}

func mainloop(text16 []string, hConst []int) string {
	r := []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		7, 4, 13, 1, 10, 6, 15, 3, 12, 0, 9, 5, 2, 14, 11, 8,
		3, 10, 14, 4, 9, 15, 8, 1, 2, 7, 0, 6, 13, 11, 5, 12,
		1, 9, 11, 10, 0, 8, 12, 4, 13, 3, 7, 15, 14, 5, 6, 2,
		4, 0, 5, 9, 7, 12, 2, 10, 14, 1, 3, 8, 11, 6, 15, 13,
	}
	r0 := []int{
		5, 14, 7, 0, 9, 2, 11, 4, 13, 6, 15, 8, 1, 10, 3, 12,
		6, 11, 3, 7, 0, 13, 5, 10, 14, 15, 8, 12, 4, 9, 1, 2,
		15, 5, 1, 3, 7, 14, 6, 9, 11, 8, 12, 2, 10, 0, 4, 13,
		8, 6, 4, 1, 3, 11, 15, 0, 5, 12, 2, 13, 9, 7, 10, 14,
		12, 15, 10, 4, 1, 5, 8, 7, 6, 2, 13, 14, 0, 3, 9, 11,
	}

	s := []int{
		11, 14, 15, 12, 5, 8, 7, 9, 11, 13, 14, 15, 6, 7, 9, 8,
		7, 6, 8, 13, 11, 9, 7, 15, 7, 12, 15, 9, 11, 7, 13, 12,
		11, 13, 6, 7, 14, 9, 13, 15, 14, 8, 13, 6, 5, 12, 7, 5,
		11, 12, 14, 15, 14, 15, 9, 8, 9, 14, 5, 6, 8, 6, 5, 12,
		9, 15, 5, 11, 6, 8, 13, 12, 5, 12, 13, 14, 11, 8, 5, 6,
	}

	s0 := []int{
		8, 9, 9, 11, 13, 15, 15, 5, 7, 7, 8, 11, 14, 14, 12, 6,
		9, 13, 15, 7, 12, 8, 9, 11, 7, 7, 12, 7, 6, 15, 13, 11,
		9, 7, 15, 11, 8, 6, 6, 14, 12, 13, 5, 14, 13, 13, 7, 5,
		15, 5, 8, 11, 14, 14, 6, 14, 6, 9, 12, 9, 12, 5, 15, 8,
		8, 5, 12, 9, 12, 5, 14, 6, 8, 13, 6, 5, 15, 13, 11, 11,
	}

	var T int
	A, B, C, D, E := hConst[0], hConst[1], hConst[2], hConst[3], hConst[4]

	A1, B1, C1, D1, E1 := hConst[0], hConst[1], hConst[2], hConst[3], hConst[4]

	for j := 0; j < 80; j++ {
		tmp, _ := strconv.Atoi(text16[r[j]])
		T = rol((A+function(j, B, C, D)+tmp+PickConstK(j))%int(math.Pow(2, 32)), s[j]) + E%int(math.Pow(2, 32))
		A = E
		E = D
		D = rol(C, 10)
		C = B
		B = T

		tmp, _ = strconv.Atoi(text16[r0[j]])

		T = rol((A1+function(79-j, B1, C1, D1)+tmp+PickConstK2(j))%int(math.Pow(2, 32)), s0[j]) + E1%int(math.Pow(2, 32))
		A1 = E1
		E1 = D1
		D1 = rol(C1, 10)
		C1 = B1
		B1 = T
	}
	T = (hConst[1] + C + D1) % int(math.Pow(2, 32))
	hConst[1] = (hConst[2] + D + E1) % int(math.Pow(2, 32))
	hConst[2] = (hConst[3] + E + A1) % int(math.Pow(2, 32))
	hConst[3] = (hConst[4] + A + B1) % int(math.Pow(2, 32))
	hConst[4] = (hConst[0] + B + C1) % int(math.Pow(2, 32))
	hConst[0] = T
	// Вывод хэша
	hash := fmt.Sprintf("%08x%08x%08x%08x%08x", hConst[0], hConst[1], hConst[2], hConst[3], hConst[4])
	fmt.Println("Hash:", hash)
	// fmt.Println("LenHash:", len(hash))
	return hash
}

func function(j, x, y, z int) int {
	if 0 <= j && j <= 15 {
		return x ^ y ^ z
	}
	if 16 <= j && j <= 31 {
		return (x & y) | (^x & z)
	}
	if 32 <= j && j <= 47 {
		return (x | ^y) ^ z
	}
	if 48 <= j && j <= 63 {
		return (x & z) | (y & ^z)
	}
	if 64 <= j && j <= 79 {
		return x ^ (y | ^z)
	}
	return -1
}

func PickConstK(j int) int {
	if 0 <= j && j <= 15 {
		return 0x00000000
	}
	if 16 <= j && j <= 31 {
		return 0x5A827999
	}
	if 32 <= j && j <= 47 {
		return 0x6ED9EBA1
	}
	if 48 <= j && j <= 63 {
		return 0x8F1BBCDC
	}
	if 64 <= j && j <= 79 {
		return 0xA953FD4E
	}
	return -1
}

func PickConstK2(j int) int {
	if 0 <= j && j <= 15 {
		return 0x50A28BE6
	}
	if 16 <= j && j <= 31 {
		return 0x5C4DD124
	}
	if 32 <= j && j <= 47 {
		return 0x6D703EF3
	}
	if 48 <= j && j <= 63 {
		return 0x7A6D76E9
	}
	if 64 <= j && j <= 79 {
		return 0x00000000
	}
	return -1
}

// rol выполняет битовый поворот влево
func rol(value int, shift int) int {
	shift = shift % 32
	//fmt.Println(value, shift)
	return (value << shift) | (value >> (32 - shift))
}
