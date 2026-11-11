#!/bin/bash

cd extensions/upload
source build.sh

cd ../..
pnpx @neutralinojs/neu build
