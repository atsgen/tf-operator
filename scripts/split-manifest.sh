#!/bin/bash
awk '
BEGIN {
a = 1
}
file ==""{
out="tmp.yaml"
}
!/^--/{
print >> (out)
}
file !="" && /^--/{
close(out)
out="tmp.yaml"
kind=""
name=""
file=""
}
kind =="" && /^kind:/{
kind=tolower($2)
}
name =="" && /^  name:/{
name=tolower($NF)
file=sprintf("%03d-%s-%s.yaml", a, kind, name)
print file >> ("files.txt")
a=(a+1)
close(out)
system("mv tmp.yaml "file)
out=file
}
' $1

