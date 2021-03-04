GOCC = go

.PHONY: all install test install-systemd

all:
	$(GOCC) build -o ./bin/multi ./cmd/multi

install:
	$(GOCC) install ./cmd/multi

install-systemd: install
	mkdir -p $(HOME)/.config/systemd/user/
	cp ./init/multiverse-user.service $(HOME)/.config/systemd/user

test:
	$(GOCC) test ./... -cover
