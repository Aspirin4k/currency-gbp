package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/Aspirin4k/currency_gbp/cbxml"
)

func main() {
	// На входе 2 параметра: валюта и ее количество
	// На выходе должно быть значение в разных валютах
	currency := flag.String("currency", "USD", "Валюта, относительно которой считать")
	value := flag.Int("value", 10, "Количество этой валюты")
	flag.Parse()

	if *value <= 0 {
		fmt.Println("Количество валюты не может быть отрицательным")
		return
	}

	// код валюты и (в будущем) ее курс относительно рубля
	cur := strings.ToUpper(*currency)
	curVal := float64(1)

	// Получаем распарсенный список валют
	currs := cbxml.ValCurs{}
	err := cbxml.GetParsedXML(&currs)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Рубль - особый случай
	if cur != "RUB" {
		noVal := true
		for _, valute := range currs.ValuteList {
			if valute.CharCode == cur {
				// Т.к. в хмл вместо \. используются \, , необходимо заменить и запарсить
				curVal, _ = strconv.ParseFloat(strings.Replace(valute.Value, ",", ".", -1), 64)
				noVal = false
				break
			}
		}

		// Если пользователь ввел какую-нибудь абракадабру, то покидаем приложение
		if noVal {
			fmt.Println("Указанной валюты не существует (либо о ней нет информации)")
			return
		}

		fmt.Printf("\tRUB: %10.2f (Российский рубль)\n", float64(*value)*curVal)
	}

	// Генерация выходного списка
	for _, valute := range currs.ValuteList {
		if valute.CharCode != cur {
			t, _ := strconv.ParseFloat(strings.Replace(valute.Value, ",", ".", -1), 64)
			fmt.Printf("\t%s: %10.2f (%s)\n", valute.CharCode, float64(*value)/(t/curVal), valute.Name)
		}
	}
}
