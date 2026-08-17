NAME := wallet-api

.PHONY: build
build:
ifneq ($(tag),)
	docker build -t $(NAME):$(tag) .
else
	docker build -t $(NAME):latest .
endif

.PHONY: push
push:
ifneq ($(tag),)
	docker push $(NAME):$(tag)
else
	docker push $(NAME):latest
endif

.PHONY: rmi
rmi:
ifneq ($(tag),)
	docker image rm $(NAME):$(tag)
else
	docker image rm $(NAME):latest
endif

.PHONY: api-test
api-test:
	GOFLAGS=-mod=vendor go test -tags=api_blackbox -count=1 -v ./service/tests/api

.PHONY: api-test-mutating
api-test-mutating:
	WALLET_API_RUN_MUTATING=1 GOFLAGS=-mod=vendor go test -tags=api_blackbox -count=1 -v ./service/tests/api
