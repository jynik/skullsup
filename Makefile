all: go firmware

go:
	$(MAKE) -C go

firmware:
	$(MAKE) -C firmware

prog:
	$(MAKE) -C firmware install

clean:
	$(MAKE) -C go clean
	$(MAKE) -C firmware clean

.PHONY: go firmware prog clean
