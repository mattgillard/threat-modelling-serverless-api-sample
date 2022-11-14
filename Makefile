
$(ARTIFACTS_DIR)/userget: functions/userget/main.go
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o $(ARTIFACTS_DIR)/userget functions/userget/main.go
build-UserGetFunction:$(ARTIFACTS_DIR)/userget


$(ARTIFACTS_DIR)/userput: functions/userput/main.go
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o $(ARTIFACTS_DIR)/userput functions/userput/main.go
build-UserPutFunction: $(ARTIFACTS_DIR)/userput

$(ARTIFACTS_DIR)/userpost: functions/userpost/main.go
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o $(ARTIFACTS_DIR)/userpost functions/userpost/main.go
build-UserPostFunction: $(ARTIFACTS_DIR)/userpost

$(ARTIFACTS_DIR)/userdelete: functions/userdelete/main.go
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o $(ARTIFACTS_DIR)/userdelete functions/userdelete/main.go
build-UserDeleteFunction: $(ARTIFACTS_DIR)/userdelete

.PHONY: build
build:
	sam build 

.PHONY: init
init: build
	sam deploy --guided

.PHONY: deploy
deploy: build
	sam deploy

.PHONY: delete
delete:
	sam delete
