// Created by Noah Hinger
// Inspired by / based off this video: https://www.youtube.com/watch?v=8VsiYWW9r48 by Matt Godbolt (namesake)
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var accumulator = 0
var indexRegister = 0
var ram [64]int

func main() {

	file, err := os.Open("examples/fibonacci.gb")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	//Load program
	var prgrm []string
	for scanner.Scan() {
		prgrm = append(prgrm, scanner.Text())
	}

	running := true

	programCounter := 0
	for running {

		//ignore comments
		if strings.HasPrefix(prgrm[programCounter], "#") {
			continue
		}
		switch cmd := strings.Split(prgrm[programCounter], " ")[0]; cmd {
		case "inc":
			indexRegister++
		case "dec":
			indexRegister--
		case "load":
			loadValue := ""
			if len(strings.Split(prgrm[programCounter], " ")) == 2 {
				loadValue = strings.Split(prgrm[programCounter], " ")[1]
			} else {
				log.Fatal("Invalid load command: ", prgrm[programCounter])
			}
			if strings.HasPrefix(loadValue, "@") {
				//Load the value at the given address into the accumulator
				ramIndex := strings.Split(loadValue, "@")[1]
				ri, err := strconv.Atoi(ramIndex)
				if err != nil {
					log.Fatal("Invalid ram index: ", ramIndex)
				}
				accumulator = ram[ri]

			} else if strings.HasPrefix(loadValue, "#") {
				//Load the value of the number into the accumulator
				num := strings.Split(loadValue, "#")[1]
				n, err := strconv.Atoi(num)
				if err != nil {
					log.Fatal("Invalid number: ", num)
				}
				accumulator = n

			} else if loadValue == "i" {
				//Load the value of the index register into the accumulator
				accumulator = ram[indexRegister]

			} else {
				log.Fatal("Unknown load command: ", prgrm[programCounter])
			}
		case "store":
			storeValue := ""
			if len(strings.Split(prgrm[programCounter], " ")) == 2 {
				storeValue = strings.Split(prgrm[programCounter], " ")[1]
			} else {
				log.Fatal("Invalid store command: ", prgrm[programCounter])
			}
			if strings.HasPrefix(storeValue, "@") {
				//Store the value of the accumulator at the given address
				ramIndex := strings.Split(storeValue, "@")[1]
				ri, err := strconv.Atoi(ramIndex)
				if err != nil {
					log.Fatal("Invalid ram index: ", ramIndex)
				}
				ram[ri] = accumulator

			} else if storeValue == "i" {
				//Store the value of the accumulator in the index register
				ram[indexRegister] = accumulator

			} else {
				log.Fatal("Unknown store command: ", prgrm[programCounter])
			}

		case "add":
			addValue := ""
			if len(strings.Split(prgrm[programCounter], " ")) == 2 {
				addValue = strings.Split(prgrm[programCounter], " ")[1]
			} else {
				log.Fatal("Invalid add command: ", prgrm[programCounter])
			}
			if strings.HasPrefix(addValue, "@") {
				//Add the value at the given address to the accumulator
				ramIndex := strings.Split(addValue, "@")[1]
				ri, err := strconv.Atoi(ramIndex)
				if err != nil {
					log.Fatal("Invalid ram index: ", ramIndex)
				}
				accumulator += ram[ri]
			} else if strings.HasPrefix(addValue, "#") {
				//Add the value of the number to the accumulator
				num := strings.Split(addValue, "#")[1]
				n, err := strconv.Atoi(num)
				if err != nil {
					log.Fatal("Invalid number: ", num)
				}
				accumulator += n
			} else if addValue == "i" {
				//Add the value of the index register to the accumulator
				accumulator += ram[indexRegister]
			} else {
				log.Fatal("Unknown add command: ", prgrm[programCounter])
			}

		case "jump":
			jumpValue := ""
			if len(strings.Split(prgrm[programCounter], " ")) == 2 {
				jumpValue = strings.Split(prgrm[programCounter], " ")[1]
			} else {
				log.Fatal("Invalid jump command: ", prgrm[programCounter])
			}
			jumpPoint, err := strconv.Atoi(jumpValue)
			if err != nil {
				log.Fatal("Invalid jump point: ", jumpPoint)
			}
			if jumpPoint < 0 || jumpPoint > len(prgrm) {
				log.Fatal("Jump point out of bounds: ", jumpPoint)
			}
			programCounter = jumpPoint

		case "read":
			//Read a value from stdin into the accumulator
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal("Error reading input: ", err)
			}
			input = strings.TrimSuffix(input, "\n")
			n, err := strconv.Atoi(input)
			if err != nil {
				log.Fatal("Invalid input: ", input)
			}
			accumulator = n

		case "write":
			//Write the value of the accumulator to stdout
			println(accumulator)

		default:
			log.Fatal("Unknown command: ", cmd)
		}
		programCounter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
