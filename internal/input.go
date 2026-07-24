package internal

import (
	"bufio"
	"fmt"
	"os"
)

func GetUserAlgorithm() (string, error) {
	var algorithm string

	fmt.Println("Which algorithm do you want to use? (1/2/3/4/5)")
	fmt.Println("1. Greedy BFS\n2. Edmonds-Karp\n3. Dinic\n4. Bellman-Ford\n5. Dijkstra + Potentials")
	fmt.Println("")

	scanner := bufio.NewScanner(os.Stdin)

	for algorithm != "" {
		input := scanner.Text()

		switch input {
		case "1":
			algorithm = "greedy"
		case "2":
			algorithm = "edmonds"
		case "3":
			algorithm = "dinic"
		case "4":
			algorithm = "bellman"
		case "5":
			algorithm = "dijkstra"
		default:
			fmt.Println("Please enter 1, 2, 3, 4, or 5. Try again :)")
			fmt.Println("")
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("Error scanning input: %v", err.Error())
		}
	}

	return algorithm, nil
}
