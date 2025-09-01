package main

import (
	"bufio"
	"fmt"
	"wthr/internal/weather"
	"os"
	"strings"
)

func main() {

	client, err := weather.GetTLSClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	weatherInfo, _ := weather.GetWeather("", client)
	weather.Show(weatherInfo)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("[!] q! - exit")
		fmt.Println("[!] c! - credits")
		fmt.Println("[!] Enter the name of the city: ")
		fmt.Print("> ")
		city, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		city = strings.TrimSpace(city)

		if city == "q!" {
			fmt.Println("[!] Have a nice day!")
			return
		} else if city == "c!" {
			fmt.Println()
			fmt.Println("[!] wthr by: 0x3 \t[github.com/trashplusplus]")
			fmt.Println("[!] wttr.in by \t[github.com/chubin]")
			fmt.Println()
			continue
		}

		weatherInfo, _ = weather.GetWeather(city, client)
		weather.Show(weatherInfo)
	}

}
