PKG = github.com/k1LoW/cgrps
COMMIT = $$(git describe --tags --always)
DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)
RELEASE_BUILD_LDFLAGS = -s -w $(BUILD_LDFLAGS)

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

depsdev:
	go get github.com/golang/lint/golint
	go get github.com/motemen/gobump/cmd/gobump
	go get -u github.com/Songmu/goxz/cmd/goxz
	go get -u github.com/tcnksm/ghr
	go get -u github.com/Songmu/ghch/cmd/ghch

build:
	go build -ldflags="$(BUILD_LDFLAGS)"

crossbuild: depsdev
	$(eval ver = v$(shell gobump show -r))
	goxz -pv=$(ver) -os=linux -arch=386,amd64 -build-ldflags="$(RELEASE_BUILD_LDFLAGS)" \
	  -d=./dist/$(ver)

release: crossbuild
	$(eval ver = v$(shell gobump show -r))
	ghr -username k1LoW -replace ${ver} dist/${ver}
