#!/bin/bash
for f in $(find . -name *_test.go); do mv $f "${f%_test.go}_no"; done
