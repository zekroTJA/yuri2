APPNAME="yuri"
PACKAGE="github.com/zekroTJA/yuri2"
LDPAKAGE="static"
BINPATH="./bin"

GO="go"
DEP="dep"

##########################################

BIN="${BINPATH}/${APPNAME}"

TAG=$(git describe --tags)
COMMIT=$(git rev-parse HEAD)

##########################################

${DEP} ensure -v

${GO} build  \
		-v -o ${BIN} -ldflags "\
			-X ${PACKAGE}/${LDPAKAGE}.AppVersion=${TAG} \
			-X ${PACKAGE}/${LDPAKAGE}.AppCommit=${COMMIT} \
			-X ${PACKAGE}/${LDPAKAGE}.Release=TRUE" \
		./cmd/${APPNAME}