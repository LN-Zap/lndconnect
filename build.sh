TAG=`git describe --tags`
PACKAGE=lndconnect

mkdir -p build
cd build

SYS=${LNDBUILDSYS:-"windows-386 windows-amd64 openbsd-386 openbsd-amd64 linux-386 linux-amd64 linux-armv6 linux-armv7 linux-arm64 darwin-386 darwin-amd64 dragonfly-amd64 freebsd-386 freebsd-amd64 freebsd-arm netbsd-386 netbsd-amd64 linux-mips64 linux-mips64le linux-ppc64"}

for i in $SYS; do
    OS=$(echo $i | cut -f1 -d-)
    ARCH=$(echo $i | cut -f2 -d-)
    ARM=

    if [[ $ARCH = "armv6" ]]; then
        ARCH=arm
        ARM=6
    elif [[ $ARCH = "armv7" ]]; then
        ARCH=arm
        ARM=7
    fi

    mkdir $PACKAGE-$i-$TAG
    cd $PACKAGE-$i-$TAG

    echo "Building:" $OS $ARCH $ARM
    env GOOS=$OS GOARCH=$ARCH GOARM=$ARM go build -v github.com/LN-Zap/lndconnect
    cd ..

    if [[ $OS = "windows" ]]; then
        zip -r $PACKAGE-$i-$TAG.zip $PACKAGE-$i-$TAG
    else
        tar -cvzf $PACKAGE-$i-$TAG.tar.gz $PACKAGE-$i-$TAG
    fi

    rm -r $PACKAGE-$i-$TAG
done
