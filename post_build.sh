#!/bin/bash
for f in $(find . -name *_no); do mv $f "${f%_no}_test.go"; done
