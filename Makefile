all:
	go run pattern.go

test:
	cd develop/dev03 && go build task.go
	cd develop/dev05 && go build task.go
	cd develop/dev06 && go build task.go