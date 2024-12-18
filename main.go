package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Question struct {
	Title   string
	Options []string
	Answer  int
}

type GameState struct {
	PlayerName string
	Points     int
	Questions  []Question
}

func (game *GameState) ProcessCSV() {
	file, error := os.Open("questions.csv")

	if error != nil {
		panic("Error while reading questions file.")
	}

	reader := csv.NewReader(file)
	records, error := reader.ReadAll()

	if error != nil {
		panic("Error while reading questions file.")
	}

	for index, record := range records {
		if index > 0 {
			correctAnswer, _ := toInt(record[5])
			question := Question{
				Title:   record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}

			game.Questions = append(game.Questions, question)
		}
	}

	defer file.Close()
}

func toInt(text string) (int, error) {
	number, error := strconv.Atoi(text)

	if error != nil {
		return 0, errors.New("tried to convert a letter into an integer")
	}

	return number, nil
}

func (game *GameState) Init() {
	fmt.Println("Welcome to WolQuiz!")
	fmt.Println("Please, insert your name: ")

	reader := bufio.NewReader(os.Stdin)
	name, error := reader.ReadString('\n')

	if error != nil {
		panic("Error while reading string.")
	}

	game.PlayerName = name
	game.Points = 0

	fmt.Printf("Hello, %s", game.PlayerName)
}

func (game *GameState) Run() {
	for questionIndex, question := range game.Questions {
		fmt.Printf("\033[36m%d. %s \033[0m\n", (questionIndex + 1), question.Title)

		for optionIndex, option := range question.Options {
			fmt.Printf("[%d] - %s\n", (optionIndex + 1), option)
		}

		fmt.Print("> ")

		var selectedOption int
		var error error

		for {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')
			read = strings.TrimSpace(read)

			selectedOption, error = toInt(read)

			if error != nil {
				fmt.Println(error.Error())
				continue
			}

			break
		}

		if selectedOption == question.Answer {
			fmt.Println("This is correct! +10 points")
			game.Points += 10
		} else {
			fmt.Println("Ops! That's wrong.")
		}

		fmt.Println("--------------")
	}
}

func main() {
	game := &GameState{}
	go game.ProcessCSV()
	game.Init()
	game.Run()

	fmt.Println("Game over! Total points:", game.Points)
}
