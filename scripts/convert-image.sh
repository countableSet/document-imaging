#!/bin/bash

mkdir -p ./output

for file in ./*.tiff
do
  filename=$(basename "$file")
  filename="${filename%.*}"

  convert $file -deskew 40% -background white -level 10%,70%,1 -blur 2 +dither +repage +matte -compress Group4 -colorspace gray -format tiff output/$filename.tiff
done
