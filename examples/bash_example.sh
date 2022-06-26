#!/bin/bash
awk -v x="$1" -v y="$2" 'BEGIN{ans=(-x*x-y*y); print ans;}'
