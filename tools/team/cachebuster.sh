#!/bin/bash

sed -i -r -e "s@(js/production\\.js\\?|css/style\\.css\\?)[0-9]+@\\1$(date +%s)@g" $PWD/www/index.html
