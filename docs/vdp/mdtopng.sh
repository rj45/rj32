#!/bin/sh

npx markdown-pdf --css-path ~/Downloads/github-markdown.css screen.md
# 8.5 inches wide, 77.6 DPI = 640 px width
pdftoppm screen.pdf video -png
#-r 77.6 -H 770 -y 60 -W 640 -x 1
# pngtopnm video-1.png > video-1.pnm
# pngtopnm video-2.png > video-2.pnm
# pngtopnm video-3.png > video-3.pnm
# pngtopnm video-4.png > video-4.pnm
# pnmcat -tb video-1.pnm video-2.pnm video-3.pnm video-4.pnm | pnmtopng > video.png
# rm video-*.pnm
