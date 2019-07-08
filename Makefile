run:
	rm -rf plugin/ &&\
	mkdir plugin &&\
	cd plugin &&\
	git clone git@github.com:af83/$(PLUGIN).git &&\
	cd $(PLUGIN) &&\
	go build -buildmode=plugin &&\
	cd ../.. &&\
	go run scops.go -plugin plugin/$(PLUGIN)/$(PLUGIN).so -debug true