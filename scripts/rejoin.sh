#!/bin/bash
awk '{print "---" >> ("complete.yaml"); system("cat "$0" >> complete.yaml")}' files.txt
