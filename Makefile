all:
	go run pattern.go

test:
	cd develop/dev03 && go build task.go
	cd develop/dev05 && go build task.go
	cd develop/dev06 && go build task.go

clean:
	cd develop/dev03 && rm task
	cd develop/dev05 && rm task
	cd develop/dev06 && rm task