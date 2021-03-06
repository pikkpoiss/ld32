#!/usr/bin/env bash

GITROOT=`git rev-parse --show-toplevel`

cd $GITROOT

mkdir -p tmp

build_aesprite() {
  aseprite \
    --batch assets/${1}.ase \
    --save-as tmp/${1}_00.png
}

build_aesprite numbered_squares
build_aesprite player
build_aesprite boss
build_aesprite stairs
build_aesprite column_whole
cp assets/*.png tmp/

TexturePacker \
  --format json-array \
  --trim-sprite-names \
  --trim-mode None \
  --size-constraints POT \
  --disable-rotation \
  --data src/resources/spritesheet.json \
  --sheet src/resources/spritesheet.png \
  tmp

rm -rf tmp
cd -
