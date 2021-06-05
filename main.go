package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func playerSymbol(human bool) (c string) {
	if human {
		return "X" // Human
	} else {
		return "O" // AI
	}
}
func main() {
	game := [3][3]string{{"-", "-", "-"}, {"-", "-", "-"}, {"-", "-", "-"}}
	isRunning := true
	player := true
	var x, y int

	for isRunning {
		gameBoard(game)

		if player {
			x, y = getInput(player)
		} else {
			x, y = getAI(game)
		}
		if game[x][y] != "-" {
			fmt.Printf("\nSquare already taken, try again.")
			x, y = getInput(player)
		}
		game[x][y] = playerSymbol(player)
		winner := decideWinner(game)
		switch winner {
		case "":
			emptycells := getEmptyCells(game)
			if len(emptycells) < 1 {
				isRunning = false
				fmt.Printf("\n Draw\n")
			}
			player = !player

		default:
			isRunning = false
			gameBoard(game)
			fmt.Printf("\n Game Over. Winner Is %s /n", winner)
		}

	}
}

func minimax(state [3][3]string, ij int, depth float64, maximise bool) (max_move int, max_score float64, log_str string) {

	max_move = ij
	nextState := state
	i, j := unCellValue(ij)
	nextState[i][j] = playerSymbol(!maximise)
	winner := decideWinner(nextState)
	log_str += "\n Move #" + strconv.Itoa(ij)

	log_str += nextStateToString(nextState)
	switch winner {
	case playerSymbol(!maximise):
		max_score = 1.0
		return
	case playerSymbol(maximise):
		max_score = -1.0
		return
	case "":
		nextMoves := getEmptyCells(nextState)

		if len(nextMoves) == 0 {
			max_score = 0
			return
		}

		for _, xy := range nextMoves {

			_, score, _ := minimax(nextState, xy, depth+1, !maximise)
			max_score = max_score + score*-1
			log_str += ("\n Score and Max Score for every possiblity # " + fmt.Sprintf("%f", score) + " & " + fmt.Sprintf("%f", max_score) + "\n")
		}

	}
	max_score = max_score * 0.5
	log_str += ("\n Evaulated Max Score for position - " + strconv.Itoa(ij) + " # " + fmt.Sprintf("%f", max_score) + "\n\n")
	log_str += ("\n============================================================================\n")

	return
}
func nextStateToString(state [3][3]string) (log_str string) {
	log_str = ("\n Game State \n")
	log_str += "      " + strings.Join(state[0][0:3], " ") + "\n"
	log_str += "      " + strings.Join(state[1][0:3], " ") + "\n"
	log_str += "      " + strings.Join(state[2][0:3], " ") + "\n"
	return
}
func getInput(player bool) (x, y int) {
	fmt.Printf("\nEnter a square, e.g. 1A. %s's turn: ", playerSymbol(player))
	var xy string

	_, err := fmt.Scan(&xy)
	if err != nil || len(xy) != 2 {
		fmt.Printf("\n Invalid Input \n ")
		return getInput(player)
	}

	x, xerr := strconv.Atoi(string(xy[0]))
	y, yerr := colInt(string(xy[1]))

	if xerr != nil || yerr != nil {
		fmt.Printf("\n Invalid Input \n ")
		return getInput(player)
	}
	return x - 1, y - 1
}
func colInt(char string) (int, error) {
	char = strings.ToLower(char)
	switch char {
	case "a":
		return 1, nil
	case "b":
		return 2, nil
	case "c":
		return 3, nil
	default:
		return -1, errors.New("Invalid column. Please choose A, B, or C.")
	}
}

func getAI(game [3][3]string) (i, j int) {
	nextMoves := getEmptyCells(game)
	max_score := -1000000.0
	var max_move int
	var log_str string
	for _, ij := range nextMoves {
		ij, score, str := minimax(game, ij, 0, true)
		log_str += str
		if score > max_score {
			log_str += "\n Code Line 145, Success! Got Best move,score and previous best max score # " + strconv.Itoa(ij) + " & " + fmt.Sprintf("%f", score) + " &  " + fmt.Sprintf("%f", max_score) + "\n"
			max_move = ij
			max_score = score
		}
	}
	writeLogFile(log_str + "\n")
	i, j = unCellValue(max_move)
	return
}
func unCellValue(ij int) (i, j int) {
	j = int(math.Mod(float64(ij), 10)) - 1
	i = (ij-j)/10 - 1
	return
}
func gameBoard(game [3][3]string) {
	fmt.Printf("\n A B C \n")
	for i, row := range game {
		fmt.Printf("%d|", i+1)
		for _, cell := range row {
			if cell == "-" {
				cell = "_"
			}
			fmt.Printf("%s|", cell)
		}
		fmt.Printf("%d\n", i+1)

	}
}
func decideWinner(game [3][3]string) (winner string) {
	possibleWins := [][3]int{
		{11, 12, 13},
		{21, 22, 23},
		{31, 32, 33},
		{11, 21, 31},
		{12, 22, 32},
		{13, 23, 33},
		{11, 22, 33},
		{13, 22, 31},
	}
	var xs, ys []int
	for i, row := range game {
		for j, cell := range row {
			switch cell {
			case playerSymbol(true):
				xs = append(xs, getCellVal(i, j))

			case playerSymbol(false):
				ys = append(ys, getCellVal(i, j))

			}
		}
	}
	for _, winSet := range possibleWins {
		if subset(winSet, xs) {
			winner = playerSymbol(true)
			return
		} else if subset(winSet, ys) {
			winner = playerSymbol(false)
			return
		} else {
			winner = ""

		}
	}
	return
}
func getCellVal(i, j int) (ij int) {
	ij = (i+1)*10 + (j + 1)
	return
}
func getEmptyCells(game [3][3]string) (cells []int) {
	for i, row := range game {
		for j, cell := range row {
			if cell == "-" {
				cells = append(cells, getCellVal(i, j))
			}
		}
	}
	return
}
func subset(first [3]int, second []int) bool {
	set := make(map[int]int)
	for _, value := range second {
		set[value] += 1
	}

	for _, value := range first {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}

	return true
}
func writeLogFile(str string) (err error) {
	os.Remove("log.txt")
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	t := time.Now()
	f.WriteString("Date Time : " + fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second()) + "\n")
	if _, err = f.WriteString(str + "\n"); err != nil {
		//panic(err)
	}
	return err
}
