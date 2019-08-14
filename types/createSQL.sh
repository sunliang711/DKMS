#!/bin/bash
rpath="$(readlink ${BASH_SOURCE})"
if [ -z "$rpath" ];then
    rpath=${BASH_SOURCE}
fi
root="$(cd $(dirname $rpath) && pwd)"
cd "$root"
shellHeaderLink='https://pic711.oss-cn-shanghai.aliyuncs.com/sh/shell-header.sh'
if [ -e /etc/shell-header.sh ];then
    source /etc/shell-header.sh
else
    (cd /tmp && wget -q "$shellHeaderLink") && source /tmp/shell-header.sh
fi
# write your code below

DB=dkms

usePass=0
while getopts ":hp" opt;do
    case $opt in
        p)
            usePass=1
            echo "Use password."
            ;;
        h)
            echo "Usage: $(basename $0) [-p] \"use password\""
            exit 1
            ;;
    esac
done
if [ $usePass -eq 1 ];then
    mysqladmin -u root -p drop $DB
    mysqladmin -u root -p create $DB
    mysql -u root -p $DB < create.sql
else
    mysqladmin -u root drop $DB
    mysqladmin -u root create $DB
    mysql -u root $DB < create.sql
fi
