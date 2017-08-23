package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
	"flag"
	"encoding/xml"
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
	"bytes"
	"strconv"
	"strings"
)

// Структуры для парсера
// Список валют
type ValCurs struct {
	XMLName  	xml.Name 	`xml:"ValCurs"`
	ValuteList 	[]Valute  	`xml:"Valute"`
}

// Валюта
type Valute struct {
	XMLName 	xml.Name 	`xml:"Valute"`
	CharCode 	string
	Value 		string
	Name 		string
}

func getXML(url string) ([]byte, error) {
	var netClient = &http.Client {
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}

	return data, nil
}

func main() {
	// На входе 2 параметра: валюта и ее количество
	// На выходе должно быть значение в разных валютах
	currency := flag.String("currency", "USD", "")
	value := flag.Int("value", 10, "")
	flag.Parse()

	// код валюты и (в будущем) ее курс относительно рубля
	cur := strings.ToUpper(*currency)
	cur_val := float64(1)

	data, err := getXML("http://www.cbr.ru/scripts/XML_daily.asp")
	if (err != nil) {
		fmt.Println(err)
	} else {
		// Необходимо использовать следующие танцы с бубном
		// т.к. unmarshall работает только если в заголовке
		// xml указана кодировка UTF-8 (здесь windows-1251)
		q := ValCurs{}
		reader := bytes.NewReader(data)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReader
		decoder.Decode(&q)

		fmt.Printf("%d %s, в других валютах это:\n",*value, cur)

		if cur != "RUB" {
			for _, valute := range q.ValuteList {
				if valute.CharCode == cur {
					// Т.к. в хмл вместо \. используются \, , необходимо заменить и запарсить
					cur_val, _ = strconv.ParseFloat(strings.Replace(valute.Value,",",".",-1),64)
					break
				}
			}
			fmt.Printf("\tRUB: %10.2f (Российский рубль)\n", float64(*value) * cur_val)
		}

		// Генерация выходного списка
		for _, valute := range q.ValuteList {
			if valute.CharCode != cur {
				t, _ := strconv.ParseFloat(strings.Replace(valute.Value, ",", ".", -1), 64)
				fmt.Printf("\t%s: %10.2f (%s)\n", valute.CharCode, float64(*value)/ (t / cur_val), valute.Name)
			}
		}
	}
}
