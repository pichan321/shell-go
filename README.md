### Important ###
Instead of running the shell using `go run main.go`, please use `go build` to build executable binary, then proceed to execute the built binary. This is due to `go run main.go` running in a terminal mode that does not support Ctrl+Z signal. Ctrl+Z will always cause the main thread to exit in `go run main.go`.

 ### Description
   A simple and basic shell implemented in Go. A personal extension project to deep dive into how shell implementation could change based on different programming languages. Based on Practicum-2 from COMP-49000 at Ithaca College.
