#!/bin/bash

set -ex

MU1=/mnt/isilon/Production/raw/2023/10/07/MIEI_20231007_1900_PGM_MU1.mxf
MU2=/mnt/isilon/Production/raw/2023/10/07/3249/MIEI_20231007_1900_PGM_MU2.mxf

FILENAMEMU1=$(basename $MU1)
FILENAMEMU2=$(basename $MU2)

DESTINATION=/mnt/isilon/Fileshare/Matjaz/tmp

mkdir -pv $DESTINATION
cd $DESTINATION


TCMU1=$(ffprobe -v error -show_entries format_tags=timecode -of default=noprint_wrappers=1:nokey=1 $MU1)
TCMU2=$(ffprobe -v error -show_entries format_tags=timecode -of default=noprint_wrappers=1:nokey=1 $MU2)

DURATION=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 $MU2)

echo $TCMU1
echo $TCMU2

IFS=":" read -ra t1 <<< "$TCMU1"
IFS=":" read -ra t2 <<< "$TCMU2"

seconds1=$(echo "scale=3; (${t1[0]} * 3600) + (${t1[1]} * 60) + ${t1[2]} + (${t1[3]} / 10)" | bc)
seconds2=$(echo "scale=3; (${t2[0]} * 3600) + (${t2[1]} * 60) + ${t2[2]} + (${t2[3]} / 10)" | bc)

difference=$(echo "scale=3; $seconds2 - $seconds1" | bc)

echo $difference

parallel 'ffmpeg -i '$MU1' -ss '$difference' -t '$DURATION' -map 0:a:{} -c:a copy -y '$FILENAMEMU1.'{}.wav' ::: {0..15}
parallel 'ffmpeg -i '$MU2' -map 0:a:{} -t '$DURATION' -c:a copy -y '$FILENAMEMU2.'{}.wav' ::: {0..15}

ffmpeg -i $FILENAMEMU1.0.wav -i $FILENAMEMU1.1.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.NOR.wav
ffmpeg -i $FILENAMEMU1.2.wav -i $FILENAMEMU1.3.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.DEU.wav
ffmpeg -i $FILENAMEMU1.4.wav -i $FILENAMEMU1.5.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.NLD.wav
ffmpeg -i $FILENAMEMU1.6.wav -i $FILENAMEMU1.7.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.ENG.wav
ffmpeg -i $FILENAMEMU1.8.wav -i $FILENAMEMU1.8.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.FRA.wav
ffmpeg -i $FILENAMEMU1.9.wav -i $FILENAMEMU1.9.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.SPA.wav
ffmpeg -i $FILENAMEMU1.10.wav -i $FILENAMEMU1.10.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.FIN.wav
ffmpeg -i $FILENAMEMU1.11.wav -i $FILENAMEMU1.11.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.RUS.wav
ffmpeg -i $FILENAMEMU1.12.wav -i $FILENAMEMU1.12.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.POR.wav
ffmpeg -i $FILENAMEMU1.13.wav -i $FILENAMEMU1.13.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.RON.wav
ffmpeg -i $FILENAMEMU1.14.wav -i $FILENAMEMU1.14.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.TUR.wav
ffmpeg -i $FILENAMEMU1.15.wav -i $FILENAMEMU1.15.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.POL.wav

ffmpeg -i $FILENAMEMU2.2.wav -i $FILENAMEMU2.2.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.BUL.wav
ffmpeg -i $FILENAMEMU2.3.wav -i $FILENAMEMU2.3.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.HUN.wav
ffmpeg -i $FILENAMEMU2.4.wav -i $FILENAMEMU2.4.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.ITA.wav
ffmpeg -i $FILENAMEMU2.5.wav -i $FILENAMEMU2.5.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.SLV.wav
ffmpeg -i $FILENAMEMU2.6.wav -i $FILENAMEMU2.6.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.CMN.wav
ffmpeg -i $FILENAMEMU2.7.wav -i $FILENAMEMU2.7.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.HRV.wav
ffmpeg -i $FILENAMEMU2.8.wav -i $FILENAMEMU2.8.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.TAM.wav
ffmpeg -i $FILENAMEMU2.9.wav -i $FILENAMEMU2.9.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.NOB.wav
ffmpeg -i $FILENAMEMU2.10.wav -i $FILENAMEMU2.10.wav -filter_complex "[0:a][1:a]amerge=inputs=2[aout]" -map "[aout]" -y $FILENAMEMU1.YUE.wav

rm -v ./*.mxf.?.wav
rm -v ./*.mxf.??.wav
