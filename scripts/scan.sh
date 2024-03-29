#!/bin/bash

PAGECOUNTER=1
FILES=()
SCANFILES=()

############################

answer="y"
while true
do
  case $answer in
    [yY]*|"" ) text="Scaning page "
            text+=$PAGECOUNTER
            echo $text
            PAGECOUNTER=$[$PAGECOUNTER+1]

            filename=`date +"%Y%m%d%H%M%S"`
            filename+=".tiff"
            FILES+=($filename)

            scanimage --device-name 'genesys:libusb:001:011' --mode Color --resolution 300 -x 215.9 -y 279.4 --format=tiff -p > $filename
            ;;
    [nN]* ) break;;
    * ) ;;
  esac

  read -p "Scan more pages? [Y/n] " answer
done

############################
# read pdf file name from user

DATE_REGEX="^([0-9]{4}(-[0-9]{2}){1,2}).*"

read -e -p "Enter pdf filename, no extension: " pdffilename

while true
do
  # if blank pick default name of date
  if [ -z "$pdffilename" ]; then
    pdffilename=`date +"%Y%m%d%H%M%S"`
  fi

  # check and append date to filename
  if ! [[ $pdffilename =~ $DATE_REGEX ]]; then
    originalfilename=$pdffilename
    pdffilename=`date +"%Y-%m-%d_"`
    pdffilename+=$originalfilename

    text="Filename did not include a date, the new filename is "
    text+=$pdffilename
    echo $text
  fi

  # do not overwrite file break if it does not exist
  if [ ! -f "$pdffilename.pdf" ]; then
    break
  fi

  read -e -p "File already exists, pick another: " pdffilename
done

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
tiff2pdf -j -o "$pdffilename.pdf" doc.tiff

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
