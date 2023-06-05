#! /bin/bash

# Check if the directory exists
if [ ! -d "$1" ]
then
	echo "Directory $1 does not exist."
	exit 1
fi

if [ ! -d "$1" ]
then
	echo "Directory $2 does not exist."
	exit 1
fi

for file in $(find "$1" -maxdepth 1 -type f)
do
	echo "Processing $file"
	cp -vR "$file" "$2"

	# Check in a loop if file exists
	# When the file is removed we can copy over a new file
	while [ -f "$2/$(basename $file)" ]
	do
		echo "Waiting for file to be processed. Seleeping 10 seconds."
		sleep 10
	done
done
