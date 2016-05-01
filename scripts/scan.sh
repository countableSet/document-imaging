#!/bin/bash

PAGECOUNTER=1
FILES=()
SCANFILES=()

############################

answer="y"
while true
do
  case $answer in
    [yY]* ) text="Scaning page "
            text+=$PAGECOUNTER
            echo $text
            PAGECOUNTER=$[$PAGECOUNTER+1]

            filename=`date +"%Y%m%d%H%M%S"`
            filename+=".tiff"
            FILES+=($filename)

            scanimage --device-name 'genesys:libusb:001:006' --mode Color --resolution 300 -x 215.9 -y 279.4 --format=tiff -p > $filename
            ;;
    [nN]* ) break;;
    * ) ;;
  esac

  read -p "Scan more pages? [y/n] " answer
done

############################
# read pdf file name from user

read -p "Enter pdf filename, no extension: " answer

if [ -z "$answer" ]; then
  answer=`date +"%Y%m%d%H%M%S"`
fi

############################
# convert tiff files into new tiff files

echo "Converting tiff files into scan files"

SCANCOUNTER=1
PREFIX="scan"
for file in "${FILES[@]}"
do
  outputfilename=$PREFIX$(printf "%03d.tiff" "$SCANCOUNTER")
  SCANFILES+=($outputfilename)
  SCANCOUNTER=$[$SCANCOUNTER+1]

  convert $file -deskew 40% -background white -level 10%,70%,1 -blur 2 +dither +repage +matte -compress Group4 -colorspace gray -format tiff $outputfilename
done

echo "Convert Done"

############################
# create final pdf document

echo "Converting tiff scans to pdf"

tiffcp scan???.tiff doc.tiff
tiff2pdf -j -o "$answer.pdf" doc.tiff

if [ $? != 0 ]; then
  echo "Convert Failed. Exiting."
  exit 1
fi

echo "Convert Done"

############################
# remove intermedite files
for i in "${FILES[@]}"
do
  rm -f $i
done
for i in "${SCANFILES[@]}"
do
  rm -f $i
done
rm -f doc.tiff

echo "Cleanup Done"
