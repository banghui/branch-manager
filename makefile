build: 
	command -v go
	command -v git
	go build -o bin/bm *.go

link: /usr/local/bin bin/branch-manager
	cp bin/bm /usr/local/bin/bm
		

