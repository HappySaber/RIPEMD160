package rsa

import (
	"RIPEMD160/ripemd160"
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"

	"golang.org/x/exp/rand"
)

// Функция eds для создания цифровой подписи
func Eds(text string) {
	h := ripemd160.Ripemd160(text)        // Создание хеш-функции (замените на вашу реализацию)
	keys := rsa()                         // Определение ключей для шифра RSA
	S := Encrypt(h, keys["E"], keys["N"]) // Шифрование хеш-функции

	//Преобразование подписи в строку
	str := ""
	for _, elem := range S {
		str += elem
	}

	fmt.Println("Данные, подписанные подписью")
	fmt.Printf("Введенный текст: %s\n", text)
	fmt.Printf("Хеш-функция: %s\n", h)
	fmt.Printf("Подпись: %s\n", str)
	fmt.Println("Проверка")

	var inputText string
	fmt.Println("Введите текст, который был введен в начале:")
	inputText, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	inputText = strings.TrimSpace(inputText)
	h = ripemd160.Ripemd160(inputText)            // Создание хеш-функции
	S_decrypt := Decrypt(S, keys["D"], keys["N"]) // Расшифрование хеш-функции

	fmt.Println("Расшифрованные данные")
	fmt.Printf("Введенный текст: %s\n", inputText)
	fmt.Printf("Хеш-функция: %s\n", h)
	fmt.Printf("Подпись: %s\n", S_decrypt)
	fmt.Println("Проверка")

	if h == S_decrypt {
		fmt.Println("Подпись верна")
	} else {
		fmt.Println("Подпись не верна")
		fmt.Println(h)
		fmt.Println(S_decrypt)
	}
}

// Функция RSA
func rsa() map[string]int {
	p := prime(64, 256)   // Генерация простого числа p
	q := prime(256, 1024) // Генерация простого числа q
	N := p * q
	f_N := (p - 1) * (q - 1)
	E := DefineEk(f_N)
	D := InverseElement(E, f_N)
	if D < 0 {
		D += max(E, f_N)
	}

	// Используем map для хранения ключей
	keys := map[string]int{
		"P": p,
		"Q": q,
		"N": N,
		"E": E,
		"D": D,
	}

	return keys
}

// Определение простого числа в интервале от lower_bound до upper_bound
func prime(lowerBound, higherBound int) int {
	var array []int
	for i := lowerBound + 1; i < higherBound; i += 2 {
		j := 2

		for j = 2; j < i; j++ {
			if i%j == 0 {
				break
			}
		}
		if j == i {
			array = append(array, i)
		}
	}
	randomNumber := rand.Intn(len(array))
	return array[randomNumber]
}

func Evclide(a, m int) []int {
	array := [][]int{
		{a, -1, 1, 0},
		{m, -1, 0, 1},
	}
	for a%m != 0 {
		newElement := []int{
			a % m,
			a / m,
			array[len(array)-2][2] - (a/m)*array[len(array)-1][2],
			array[len(array)-2][3] - (a/m)*array[len(array)-1][3],
		}
		array = append(array, newElement)

		// Обновление a и m
		a = array[len(array)-2][0]
		m = array[len(array)-1][0]
	}
	return array[len(array)-1]
}

func DefineEk(f_N int) int {
	var array []int
	for i := 2; i < f_N; i++ {
		if Evclide(f_N, i)[0] == 1 {
			array = append(array, i)
		}
	}
	randomNumber := rand.Intn(len(array) - 1)
	fmt.Println("Rand: ", randomNumber)
	return array[randomNumber]
}

// InverseElement функция для нахождения обратного элемента
func InverseElement(e, f_n int) int {
	if e > f_n {
		return Evclide(e, f_n)[2]
	}
	return Evclide(f_n, e)[3]
}

// Функция шифрования
func Encrypt(h string, E, N int) []string {
	result := []string{}
	n := big.NewInt(int64(N)) // Преобразуем N в big.Int

	for _, elem := range h {
		// Преобразуем элемент из шестнадцатеричной строки в целое число
		num := new(big.Int)
		num.SetString(string(elem), 16) // Преобразуем символ в строку и затем в big.Int

		// Выполняем шифрование: (num^E) % N
		encrypted := new(big.Int).Exp(num, big.NewInt(int64(E)), n)

		// Преобразуем результат обратно в шестнадцатеричную строку и добавляем в результат
		result = append(result, fmt.Sprintf("%x", encrypted))
	}

	return result
}

// Функция расшифрования
func Decrypt(S []string, D, N int) string {
	result := ""
	n := big.NewInt(int64(N)) // Преобразуем N в big.Int

	for _, elem := range S {
		// Преобразуем элемент из шестнадцатеричной строки в целое число
		num := new(big.Int)
		num.SetString(elem, 16) // Преобразуем шестнадцатеричную строку в big.Int

		// Выполняем расшифрование: (num^D) % N
		decrypted := new(big.Int).Exp(num, big.NewInt(int64(D)), n)

		// Преобразуем результат обратно в шестнадцатеричную строку и добавляем в результат
		result += fmt.Sprintf("%x", decrypted)
	}

	return result
}
