#!/usr/bin/env python
import argparse
import json


def sort_data(program_item):
    # we sort by program name
    name = program_item['name']
    if type(name) is not str:
        raise ValueError('program name must be a string! please fix before proceeding')
    return name


def run(arguments):
    fd = open('chaos-bugbounty-list.json', 'r')
    content = fd.read()
    fd.close()
    content_json = json.loads(content)
    sorted_data = sorted(content_json['programs'], key=sort_data)
    output_str = json.dumps({"programs": sorted_data}, indent=4, sort_keys=False)
    print(output_str)
    if arguments.fix:
        print('Writing sorted data to file')
        with open('chaos-bugbounty-list.json', 'w') as fd:
            fd.write(output_str)


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--fix', '-f', action='store_true',
                        default=False,
                        help='Fix the data file(this will REPLACE the current chaos data file!')
    args = parser.parse_args()
    run(args)
