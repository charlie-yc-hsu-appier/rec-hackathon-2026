from PIL import Image
from PIL import ImageChops
from io import BytesIO
import ffmpeg
import os
import sys


def convert_binary_to_image(string):
    """
        Function to convert the content into pciture
        :param string: binary type, response.content from Http request
        :param path: string type, the path to store the picture
        :return string: dict type, return the basic information of the picture

        Robot Examples:
       | ${picture} | Convert Binary To Image | ${resp.content} | path_to_store
    """
    image = Image.open(BytesIO(string))

    if image.info is not None:
        # Image.open(BytesIO(string)).save(path)
        image_size = len(image.fp.read())
        return [image.mode, image.size, image.format, image_size]
    else:
        return "It's not a binary file"


def differnce_images(path_one, path_two, diff_save_location):
    """
        Function to compare two different pictures
        :param path_one: string type, the path of image 1 (the original image)
        :param path_two: string type, the path of image 2 (the compared image)
        :param diff_save_location: the path to stoe the difference result
        :return string: string type, return the difference result in string,
         Fail: means these pictures are different

        :Robot Examples:
        | ${diff} | Differnce Images | path_pic1 | path_pic2 | path_store_diff
    """
    image_one = Image.open(path_one)
    image_two = Image.open(path_two)

    """ Difference calculated here """
    diff = ImageChops.difference(image_one, image_two)

    """ Method to check if there is no difference """
    if diff.getbbox() is None:
        print("No difference in the images")
        return "Pass"

    else:
        diff.save(diff_save_location)
        print("These pictures are different, see more: " + diff_save_location)
        return "Fail"


def save_video(string, FILE_OUTPUT):

    if os.path.isfile(FILE_OUTPUT):
        os.remove(FILE_OUTPUT)

    out_file = open(FILE_OUTPUT, "wb")
    out_file.write(string)
    out_file.close()


def get_video_info(path):

    try:
        probe = ffmpeg.probe(filename=path)
    except ffmpeg.Error as e:
        print(e.stderr, file=sys.stderr)

    video_stream = next((stream for stream in probe['streams'] if stream['codec_type'] == 'video'), None)

    if video_stream is None:
        print('No video stream found')

    print(video_stream)
    width = int(video_stream['width'])
    height = int(video_stream['height'])
    num_frames = int(video_stream['nb_frames'])
    duration_ts = int(video_stream['duration_ts'])
    codec_long_name = str(video_stream['codec_long_name'])

    return [width, height, num_frames, duration_ts, codec_long_name]
