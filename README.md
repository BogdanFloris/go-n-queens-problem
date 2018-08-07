# go-n-queens-problem
Solution to the N Queens Problem using Go and the asynchronous backtracking algorithm as described in the book **Multiagent Systems: Algorithmic, Game-Theoretic, and Logical Foundations** by *Yoav Shoham and Kevin Leyton-Brown*.
The chess game board is drawn using OpenGL.

## Usage
To run the program without compiling:
```bash
go run main.go
```
To compile and then run the program:
```bash
go build main.go
./main
```
There is one flag, the number N of queens on a NxN chess board. The default is 4. Example:
```bash
./main.go -N 8
```