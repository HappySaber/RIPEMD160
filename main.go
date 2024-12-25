package main

import (
	"RIPEMD160/ripemd160"
	"RIPEMD160/rsa"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var text string

	fmt.Println("Что вы ходите сделать, хеширование(1) или ЭЦП(2)?")
	var choice int
	fmt.Scanln(&choice)
	if choice == 1 {
		fmt.Println("Введите текст для хэширования")
		text, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		text = strings.TrimSpace(text)
		ripemd160.Ripemd160(text)
	} else if choice == 2 {
		fmt.Println("Введите текст для ЭЦП")
		text, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		text = strings.TrimSpace(text)
		rsa.Eds(text)
	} else {
		fmt.Println("Вы ввели некорректное значение")
	}

}
