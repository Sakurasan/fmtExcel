#!/bin/bash
#传参测试脚本
echo "tar file Shell script is `basename $0`"
echo "The tar filr name  is : $1"
tar zcvf $1.tar.gz $1

