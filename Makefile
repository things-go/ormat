.PHONY: init
init:
ifeq ($(shell uname -s),Darwin)
	@grep -r -l go-lib-template * | xargs sed -i "" "s/cicd-go-template/$$(basename `git rev-parse --show-toplevel`)/"
else
	@grep -r -l go-lib-template * | xargs sed -i "s/cicd-go-template/$$(basename `git rev-parse --show-toplevel`)/"
endif