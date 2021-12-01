build-docker:
	cd build && $(MAKE) clean dbuild
push-docker:
	cd build && $(MAKE) clean dpush	
generate:
	go generate ./...
