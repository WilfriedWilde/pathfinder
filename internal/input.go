package internal

import (
	"bufio"
	"fmt"
	"os"
)

func GetUserAlgorithm() (string, error) {
	var algorithm, fullName string

	fmt.Println("Which algorithm do you want to use? (1/2/3/4/5)")
	fmt.Println("1. Greedy BFS\n2. Edmonds-Karp\n3. Dinic\n4. Bellman-Ford\n5. Dijkstra + Potentials")
	fmt.Println("")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()

		switch input {
		case "1":
			algorithm = "greedy"
			fullName = "Greedy BFS"
		case "2":
			algorithm = "edmonds"
			fullName = "Edmonds-Karp"
		case "3":
			algorithm = "dinic"
			fullName = "Dinic"
		case "4":
			algorithm = "bellman"
			fullName = "Bellman-Ford"
		case "5":
			algorithm = "dijkstra"
			fullName = "Dijkstra + Potentials"
		default:
			fmt.Println("Please enter 1, 2, 3, 4, or 5. Try again :)")
			fmt.Println("")
		}

		if algorithm != "" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning input: %s", err.Error())
	}

	if algorithm == "" {
		return "", fmt.Errorf("no algorithm selected (input closed)")
	}

	fmt.Printf("\n%s running...\n\n", fullName)

	return algorithm, nil
}
