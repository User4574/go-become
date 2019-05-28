become: become.go
	go build


.PHONY: install no-really-install

install: become
	@echo "Don't be an idiot."

no-really-install: become
	@echo "Dumbass."
	go install
