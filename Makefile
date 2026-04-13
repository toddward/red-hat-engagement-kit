.PHONY: test test-bash test-go test-vhs

test: test-bash test-go

test-bash:
	bash tui/tests/run_tests.sh

test-go:
	cd tui/viewer && go test ./...

test-vhs:
	bash tui/tests/vhs/run_vhs_tests.sh
