login: login.o
	./login

login.o: go.mod
	go	build	./logintool/login.go

go.mod:
	go	mod	init	logintool
	go	mod	tidy

env:
	pip	install	-r	requirements.txt	-i	https://pypi.tuna.tsinghua.edu.cn/simple

main:
	python	main.py

clean:
	rm	-f	./go.mod
	rm	-f	./go.sum
	rm	-f	./login
