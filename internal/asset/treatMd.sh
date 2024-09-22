#!/bin/bash

for filename in /markdown/*.md; do
    for ((i=0; i<=3; i++)); do
        ./mdtohtml -toc=true "./markdown/$filename" "./static/articles/$filename.html
    done
done
