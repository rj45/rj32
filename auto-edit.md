
audio normalize step:

ffmpeg-normalize rj32-006b.mp4 -o rj32-006n.mp4 -t -27 -lrt 11 -tp -4.0 \
	-c:a aac -b:a 192k -ar 48000 -f


auto-editor z-2021-01-10__05-19-44.mkv \
	-o test2.mkv --edit_based_on audio_or_motion \
	--ignore 0-6 --motion_threshold 0.003 --min_clip_length 1 -m 5 \
	motionOps -d 5 -b 5

The motion ops:
- b blurs the output -- higher to ignore more motion
- d is a "dilation" -- I think this exaggerates motion
- motion_threshold is the number of change pixels / total pixels in a frame
- min_clip_length will keep short clips, like typing or button clicking

then a speed up pass:

auto-editor test2.mkv -o test3.mkv -s 3 --ignore 0-6

this speeds up any non-speaking left in because there was motion from the first pass
