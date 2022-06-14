changelog_args=-o CHANGELOG.md -tag-filter-pattern '^v'

gocognit:
	-gocognit -over 15 .

test_command:
ifeq (, $(shell which richgo))
	$(warning "richgo is not installed, falling back to plain go test")
	$(eval TEST_BIN=go test)
else
	$(eval TEST_BIN=richgo test)
endif
ifdef test_run
	$(eval TEST_ARGS := -run $(test_run))
endif
ifdef verbose
	$(eval TEST_ARGS=$(TEST_ARGS) -v)
endif
	$(eval test_command=$(TEST_BIN) ./... $(TEST_ARGS) --cover)

lint: gocognit
	golangci-lint run --concurrency 4 --print-issued-lines=false --exclude-use-default=false --enable=golint --enable=goimports  --enable=unconvert --enable=unparam

test_only: test_command
	$(test_command) -timeout 60s

test: lint test_only

changelog:
ifdef version
	$(eval changelog_args=--next-tag $(version) $(changelog_args))
endif
	git-chglog $(changelog_args)

.PHONY: gocognit richgo lint test-only test changelog changelog