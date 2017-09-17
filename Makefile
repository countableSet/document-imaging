PROJECT=document-imaging
PACKAGE_VERSION=0.3.0
PACKAGE_BASEDIR=package
PACKAGE_DIR=${PACKAGE_BASEDIR}/${PROJECT}
PACKAGE_DATE:=$(shell date +"%a, %d %b %Y %k:%M:%S %z")
PACKAGE_COMMITS:=$(shell git log --pretty=format:"  * %h %s\n" b46555bcf104e90c8b79e67a27d3cd9a7353c20f..HEAD | tr -d '\n')

COMMIT_HASH=`git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE=`date +%FT%T%z`
VERSION_NUMBER=v${PACKAGE_VERSION}
LDFLAGS=-ldflags "-X main.CommitHash=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE} -X main.VersionNumber=${VERSION_NUMBER}"

default: fmt test build

build:
	go build ${LDFLAGS}

package: fmt test build
	mkdir -p ${PACKAGE_DIR}
	cp -R debian ${PACKAGE_DIR}
	mv ${PACKAGE_DIR}/debian/control.ex ${PACKAGE_DIR}/debian/control
	mv ${PACKAGE_DIR}/debian/changelog.ex ${PACKAGE_DIR}/debian/changelog
	sed -i 's/<name>/${PROJECT}/g' ${PACKAGE_DIR}/debian/control
	sed -i 's/<version>/${PACKAGE_VERSION}/g' ${PACKAGE_DIR}/debian/changelog
	sed -i 's/<changes>/${PACKAGE_COMMITS}/g' ${PACKAGE_DIR}/debian/changelog
	sed -i 's/<date>/${PACKAGE_DATE}/g' ${PACKAGE_DIR}/debian/changelog
	sed -i 's/<deb-version>/${PACKAGE_VERSION}/g' ${PACKAGE_DIR}/debian/${PROJECT}.1
	cp ${PROJECT} ${PACKAGE_DIR}
	cd ${PACKAGE_DIR} && debuild -us -uc -b
	tar -czf ${PACKAGE_BASEDIR}/${PROJECT}-${PACKAGE_VERSION}-linux-amd64.tar.gz ${PROJECT}

fmt:
	go fmt

test:
	go test ./...

clean:
	rm -rf ${PROJECT}
	rm -rf ${PACKAGE_BASEDIR}
