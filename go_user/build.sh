#!/bin/sh

#####################################################################
# usage:
# sh build.sh 构建默认的windows 32位程序
# sh build.sh darwin(或linux), 构建指定平台的64位程序

# examples:
# sh build.sh darwin amd64 构建MacOS版本的64位程序
# sh build.sh linux amd64 构建linux版本的64位程序
#####################################################################

source /etc/profile


#
#function getAll(){
#    go get github.com/garyburd/redigo/redis
#    go get github.com/Terry-Mao/goconf
#}

OS="$1"
ARCH="$2"

if [ -n "$OS" ];then
    echo "use defined GOOS: "${OS}
else
    echo "use default GOOS: linux"
    OS=linux
    ARCH=amd64
fi

echo "start building with GOOS: "${OS}", PROFILE: "${PROFILE}

export GOOS=${OS}
export GOARCH=${ARCH}


release_dir="release"
revision=`git describe --long --dirty`


mkdir -p ./${release_dir}
rm -rf ./${release_dir}/*
mkdir -p ./${release_dir}/utils


#flags="-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse HEAD`"
flags="-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git describe --long --dirty --abbrev=14`"
echo ${flags}
go build -ldflags "$flags" -x -o ${release_dir}/import-illegal-mobile import_illegal_mobile.go
go build -ldflags "$flags" -x -o ${release_dir}/import-para-define import_para_define.go
go build -ldflags "$flags" -x -o ${release_dir}/import-super-user import_super_user.go
go build -ldflags "$flags" -x -o ${release_dir}/import-user import_user.go
go build -ldflags "$flags" -x -o ${release_dir}/import-user-bind import_user_bind.go
go build -ldflags "$flags" -x -o ${release_dir}/import-user-number import_user_number.go
go build -ldflags "$flags" -x -o ${release_dir}/inspect-error inspect_error.go
go build -ldflags "$flags" -x -o ${release_dir}/utils/import-user-null utils/import_user_null.go


cp ./config.conf ./${release_dir}/config.conf
cp ./README.md ./${release_dir}/
cp -r data ./${release_dir}/
mkdir -p ./${release_dir}/log

echo "finish building with GOOS: "${OS}", PROFILE: "${PROFILE}

rm -rf release-*.tar.gz
tar zcvf release-${revision}.tar.gz release
