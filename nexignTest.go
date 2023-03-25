package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Структура, которая хранит
//
//	Тип вызова
//	Номер абонента
//	Дату и время начала вызова
//	Дату и время окончания вызова
//	Тариф
//	Ссылку на следующий элемент списка
type ListNode struct {
	callType string

	number string

	date1 time.Time

	date2 time.Time

	tariff string

	Next *ListNode
}

// Создаём список
func CreateListNode() *ListNode {
	head := new(ListNode)
	return head
}

// Добавляем значения в список
func (list *ListNode) AddNode(i int, cT string, n string, dt1, dt2 time.Time, tariff string) bool {
	j := 1
	// 0-й узел
	for list != nil && j < i { // Находим i-1-й узел
		list = list.Next
		j++
	}
	if list == nil || j > i {
		return false // i-1-й узел не существует, значит i-й узел тоже не существует
	}
	s := &ListNode{cT, n, dt1, dt2, tariff, list.Next}
	list.Next = s

	return true
}

// Вывод на печать
func (list *ListNode) Print() error {

	list = list.Next
	p1 := list
	p2 := list

	for p1 != nil { //Обход по списку

		var dur time.Duration //Разница между конечным и начальным временем вызова (для всех звонков)
		var currCost float64  //Разница текущего звонка
		var cost float64      //Общая стоимость на конец тарифного периода

		currentNumber := p1.number //Текущий номер

		file, err := os.Create(currentNumber) //Создаём файл для нового абонента

		if err != nil { //Проверяем возможность создания файла
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		defer file.Close() //Закрываем файл
		p2 = list

		//Форматированный вывод в файл
		fmt.Fprintf(file, "Tariff index: %s\n", p1.tariff)
		fmt.Fprintf(file, "----------------------------------------------------------------------------\n")
		fmt.Fprintf(file, "Report for phone number %s:\n", currentNumber)
		fmt.Fprintf(file, "----------------------------------------------------------------------------\n")
		fmt.Fprintf(file, "| Call Type |   Start Time        |     End Time        | Duration |  Cost  |\n")
		fmt.Fprintf(file, "----------------------------------------------------------------------------\n")

		for p2 != nil { //

			if currentNumber == p2.number {
				currDuration := p2.date2.Sub(p2.date1)
				dur += currDuration
				if p2.tariff == "06" {
					if dur <= time.Minute*300 {
						cost = 100
						currCost = 0
					} else {
						currCost = float64(time.Duration(currDuration.Minutes()))
						cost += currCost
					}
				} else if p2.tariff == "03" {
					currCost = float64(time.Duration(currDuration.Minutes())) * 1.5
					cost += currCost
				} else if p2.tariff == "11" {
					if p2.callType == "02" {
						cost = 0
					} else {
						if dur <= time.Minute*100 {
							currCost = float64(time.Duration(currDuration.Minutes())) * 0.5
							cost += currCost
						}
						if dur > time.Minute*100 {
							currCost = float64(time.Duration(currDuration.Minutes())) * 1.5
							cost += currCost
						}
					}
				}
				fmt.Fprintf(file, "|    %v    | %v | %v | %8v | %6.2f |\n", p2.callType, p2.date1.Format("2006-01-02 15:04:05"), p2.date2.Format("2006-01-02 15:04:05"), currDuration, currCost)

			}
			p2 = p2.Next // Переход к следующему элементу списка

		}

		fmt.Fprintf(file, "----------------------------------------------------------------------------\n")
		fmt.Fprintf(file, "|                                           Total Cost: |     %6.2f rubles |\n", cost)
		fmt.Fprintf(file, "----------------------------------------------------------------------------\n")
		p1 = p1.Next //Переход к следующему элементу списка
	}
	fmt.Println()
	return nil
}

func main() {
	// Открываем файл выгрузки
	file, err := os.Open("cdr.txt")

	//Проверям на ошибки при открытии
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file) //Создаём новый сканер для построчного обхода по файлу выгрузки

	//Проверяем возможность создания сканнера на ошибки
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	list := CreateListNode() //Создаём список
	//
	i := 0
	for fileScanner.Scan() { //Построчно сканируем файл
		var str string

		str = fileScanner.Text()        //Заносим отсканированную строку в переменную str
		vars := strings.Split(str, " ") //Разбиваем строку на подстроки по разделителю " "
		//Тип вызова,тариф и номер
		tp := vars[0]
		tar := vars[4]
		num := vars[1]
		//Дата начала и конца вызова
		dt1 := vars[2][:8]
		dt2 := vars[3][:8]
		//Время начала и конца вызова
		tm1 := vars[2][8:]
		tm2 := vars[3][8:]
		//Разбиение на год, месяц и день
		year1, _ := strconv.ParseInt(dt1[:4], 10, 64)
		year2, _ := strconv.ParseInt(dt2[:4], 10, 64)
		month1, _ := strconv.ParseInt(dt1[4:6], 10, 64)
		month2, _ := strconv.ParseInt(dt2[4:6], 10, 64)
		day1, _ := strconv.ParseInt(dt1[6:], 10, 64)
		day2, _ := strconv.ParseInt(dt2[6:], 10, 64)
		//Разбиение на часы, минуты и секунды
		h1, _ := strconv.ParseInt(tm1[:2], 10, 64)
		h2, _ := strconv.ParseInt(tm2[:2], 10, 64)
		m1, _ := strconv.ParseInt(tm1[2:4], 10, 64)
		m2, _ := strconv.ParseInt(tm2[2:4], 10, 64)
		s1, _ := strconv.ParseInt(tm1[4:6], 10, 64)
		s2, _ := strconv.ParseInt(tm2[4:6], 10, 64)
		l, _ := time.LoadLocation("")
		//Дата и время для добавления в список
		date1 := time.Date(int(year1), time.Month(int(month1)), int(day1), int(h1), int(m1), int(s1), 0, l)
		date2 := time.Date(int(year2), time.Month(int(month2)), int(day2), int(h2), int(m2), int(s2), 0, l)

		list.AddNode(i+1, tp, num, date1, date2, tar) //Заносим данные в список
		i++

	}
	file.Close() // Закрываем файл выгрузки
	list.Print() //Печать файлов

}
