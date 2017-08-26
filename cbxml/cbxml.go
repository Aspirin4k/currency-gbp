package cbxml

import (
	"encoding/xml"
	"bytes"

	"github.com/Aspirin4k/currency_gbp/cbquery"

	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
)

// ValCurs ...
// Структуры для парсера
// Список валют
type ValCurs struct {
	XMLName    xml.Name `xml:"ValCurs"`
	ValuteList []Valute `xml:"Valute"`
}

// Valute Валюта
type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	CharCode string
	Value    string
	Name     string
}

func GetParsedXML(valCurs *ValCurs) (error) {
	data, err := cbquery.GetData("http://www.cbr.ru/scripts/XML_daily.asp")

	if err != nil {
		return err
	}

	// Необходимо использовать следующие танцы с бубном
	// т.к. unmarshall работает только если в заголовке
	// xml указана кодировка UTF-8 (здесь windows-1251)
	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader
	decoder.Decode(valCurs)

	return nil
}
