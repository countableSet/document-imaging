#!/bin/bash

a=1
prefix="scan"
for i in `ls -v ./*.tiff`; do
  newFile=$prefix$(printf "%03d.tiff" "$a")
  mv "$i" "$newFile"
  let a=a+1
done

tiffcp scan???.tiff doc.tiff
tiff2pdf -j -o `date +"%Y%m%d%H%M%S"`.pdf doc.tiff

# Remove the source tiff files
rm *.tiff
