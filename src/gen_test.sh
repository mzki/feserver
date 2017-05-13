#!/bin/sh

# check command exist.
if `type wget > /dev/null 2>&1`; then
	echo "wget found"
else 
	echo "wget not found. please install that" && exit 1
fi

if `type nkf > /dev/null 2>&1`; then 
	echo "nkf found"
else
	echo "nkf not found. please install that" && exit 1
fi

# get testing html.
wget -O ./test.html -A html http://www.fe-siken.com/kakomon/28_haru/q2.html 
nkf -S -w8 --in-place ./test.html 

wget -O ./y19_spring_q26.html -A html http://www.fe-siken.com/kakomon/19_haru/q26.html 
# nkf -S -w8 --in-place ./y19_spring_q26.html 
