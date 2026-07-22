# Pathfinder<br>

## Overview<br>

A CLI (/localhost) tool that calculates the minimum number of turns required to move the given amount of trains from the starting position to the goal.

It then prints out the train movements, turn by turn.

## Setup and Usage <br>

`git clone https://gitea.kood.tech/jyrkikarhunen/pathfinder`

`cd pathfinder`

Then either

`go run . [INPUTFILEPATH] [STARTSTATION] [ENDSTATION] [TRAINNUMBER]` for the default CLI mode.  

or

`go run . -s` to launch the server mode. 

The server mode is accessible at `http://localhost:8081/`  
- Select an input file from the list (can manually add in your own files into networks/)  
- Select the starting and ending locations from the next two lists  
- Type in the number of trains to use (capped at a million)  
- Select one or more pathfinding algorithms to use
- Click run

On server mode, the actual turn printouts are still available on the usual terminal, if required.

## Additional features<br>
### Proposed bonuses:<br>
Advanced error handling:  Should not have any generic errors left, they all tell what went wrong.

Super-advanced error handling:  Errors also try to be as specific as they can, such as printing out the line index of the input that failed, or printing out the offending station name(s) in case of duplication.

Turn-efficiency: Matches the example turn limits (impossible to go under them), finds better solutions than a standard max-flow algorithm (like Edmonds-Karp) would.

Test suite: Running `go test` in the root folder runs a test file that tests all of the mandatory examples, and lots of extra stuff that we found useful to test.

Run-efficiency: Biggest (very likely unreasonably large and adversarially designed, 10 thousand stations, 1.5 million connections, 100 thousand trains) inputs can take up to 15 seconds to run (largely dominated by input reading, and printing out to terminal, pathfinding is only 0.5 - 3 seconds, depending on algorithm), anything below that runs more or less instantly.  

### Extra bonuses:<br>
Local server implementation.

Visualization (including turn-by-turn animated trains) shown on the server (so long as the input file is small enough, up to 25 thousand stations + connections combined).

Benchmarking different algorithms by both runtime and by length of solution in turns.  

## Algorithms used<br>
- "Greedy" BFS: Find shortest path, remove used stations, repeat.
- "Max-Flow" Edmonds-Karp and Dinic: Find shortest path, try to augment by finding shortest-by-edges augmenting paths from start to end, thus increasing the flow.  
Dinic being more clever in the way it does this (repeated DFS), resulting in faster runtime.
- "Min-Cost Max-Flow" Bellman-Ford and Dijkstra + Potentials: Same as above, but shortest-by-length instead, often resulting in better set of routes found.  
Bellman-Ford falls off fast on bigger inputs, potentially taking minutes to get a result, Dijkstra remains fast.