default:
	$(info select target - use 'make list' for research)

list:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | grep -E -v -e '^[^[:alnum:]]' -e '^$@$$'

make_linux_amd64:
	GOOS=linux GOARCH=amd64 ./build.sh
make_applesilicone_arm:
	GOOS=darwin GOARCH=arm64 ./build.sh
make_win_amd64:
	GOOS=windows GOARCH=amd64 ./build.sh
make_apple_amd64:
	GOOS=darwin GOARCH=amd64 ./build.sh