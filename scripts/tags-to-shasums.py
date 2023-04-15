#!/usr/bin/python3
# -*- coding: utf-8 -*-
''' 
Author: Rob Mokkink rob@mokkinksystems.com
Description: Change tags to shasums for oc-mirror usage
'''
#------------------
# Imports
#------------------
import argparse
import os
import glob
import re
import json
import subprocess
import fileinput
from concurrent.futures import ProcessPoolExecutor

#------------------
# Global vars
#------------------

# CSV Files to look for
CSV_FILE = "*clusterserviceversion*"

# Default amount of skopeo processes
SKOPEO_PROCS = 5

#------------------
# Functions
#------------------

def get_ops_csv_files(ops_dir):
    '''
    Get all csv files and get
    all the images
    '''

    file_list = []

    for dir, _, files in os.walk(ops_dir):
        for file in files:
            # Get the clusterservice files
            if glob.fnmatch.fnmatch(file, CSV_FILE):
                # Get the full filepath
                file_path = os.path.join(dir, file)

                file_list.append(file_path)
    # Return a list with all the files
    return file_list

def get_images_locations(filelist):
    '''
    Get all the image locations in the files
    '''

    # Create empty list
    # this will contain all the dictionaries with files
    # image tags
    img_list = []

    # Loop through the files
    for csv_file in filelist:


        # open and read the file
        with open(csv_file, 'rt') as fh:
            get_text = fh.read()

            # Get all the image locations
            get_img_loc = re.findall("image:.*", get_text)

            # Loop through the found image locations and add them
            for img_loc in get_img_loc:

                # Create empty dict
                img_dict = {}

                # Check for vaulty image locations
                # Have seen these a couple of times
                if "@@" in img_loc:
                    print("%s is vaulty" % (img_loc))

                    # Adjust the variable
                    adj_loc = img_loc.replace("@@", "@")
                    print("Adjusted to %s" %(adj_loc))

                    # Assign the variable the proper value
                    img_loc = adj_loc

                # Skip images with shasums
                if "sha256" in img_loc:
                    break

                # Remove example lines
                if ">-" in img_loc:
                    break

                # Skip Empty lines
                if not img_loc:
                    break
                if "''" in img_loc:
                    break

                # Old location
                old_loc = img_loc.split("image:")[-1]
                old_loc = old_loc.strip()

                # Add entry's to the dictionary
                img_dict[csv_file] = old_loc
          
                # Add dict to the list
                img_list.append(img_dict)

    # return the dictionary
    return img_list

def get_shasum_images(csv_file, image):
    '''
    Get the shasum from the image
    '''

    # Create empty dict
    img_dict = {}

    # Assemble command
    command = ['skopeo', 'inspect', 'docker://' + image]

    # Execute skopeo
    proc = subprocess.Popen(command, stdout=subprocess.PIPE)
    proc_out = proc.communicate()

    # Try to find the shasum
    try:
        digest = json.loads(proc_out[0])
        digest = digest['Digest']
    except ValueError:
        digest = "Error"

    # Add stuff to the new dictionary
    split_image_loc = image.split(":")[0]
    new_loc = split_image_loc + "@" + digest

    # Add to the dictionary
    img_dict['csv_file'] =  csv_file
    img_dict['image_tag'] = image
    img_dict['image_shasum'] = new_loc

    # Return the dictionary
    return img_dict


def adj_files(csv_file, image_tag, image_shasum):
    '''
    Adjust the files
    '''
    # Adjust the file
    with fileinput.FileInput(csv_file, inplace=True) as fh:
        for line in fh:
            print(line.replace(image_tag, image_shasum), end='')

def run():
    '''
    Main function
    '''

    parser = argparse.ArgumentParser(description="Tags to Shasums")
    parser.add_argument("-d", "--dir",
                        help="directory containing the operators information",
                        required=True),
    parser.add_argument("-i", "--instances",
                        help="Amount of skopeo instances",
                        required=False)
    args = parser.parse_args()

    # Adjust the default amount of skopeo processes when necessary
    if args.instances is not None:
        skopeo_procs = int(args.instances)
    else:
        skopeo_procs = SKOPEO_PROCS

    # Get all the operator clusterserviceversion files
    get_files = get_ops_csv_files(args.dir)

    # Get all the images inside the files and dump them to json
    get_images = get_images_locations(get_files)

    # Loop through the dictionary that was created and get the shasum
    # Create a list to store the dicts
    shasum_images_list = []

    # Use a processpool to speed up the proces
    with ProcessPoolExecutor(max_workers=skopeo_procs) as executor:
        for entry in get_images:
            for csv_file, image in entry.items():
                get_dict = executor.submit(get_shasum_images, csv_file, image)
                shasum_images_list.append(get_dict)

    # Loop through the results and adjust the files
    # One by one, otherwise we will run into issues with the handlers
    for image_entry in shasum_images_list:
        csv_file = (image_entry.result()['csv_file'])
        image_tag = (image_entry.result()['image_tag'])
        image_shasum = (image_entry.result()['image_shasum'])

        # Adjust the file
        adj_files(csv_file, image_tag, image_shasum)

if __name__ == '__main__':
    run()
