#!/usr/bin/env python

import argparse
import sys

args = argparse.ArgumentParser()

args.add_argument(
    '--x',
    dest='x',
    type=int,
    required=True
)
args.add_argument(
    '--y',
    dest='y',
    type=float,
    required=True
)

config = args.parse_args(sys.argv[1:])
print(-config.y**2 - config.x**2)
