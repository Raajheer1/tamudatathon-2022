import cv2
import numpy as np
import math



# get pieces together in a single image
def combine_pieces(size_vertical, size_horizontal, cnt_row, cnt_col, cnt_total, chn, plist, single = False):
    if single:
        print(size_vertical, size_horizontal, cnt_row, cnt_col, cnt_total, chn)
        temp = np.zeros((size_vertical * cnt_col, size_horizontal * cnt_row, chn), dtype=np.uint8)
        for pIt in range(cnt_total):
            print(str(pIt) + " out of total is " + str(cnt_total))
            start_row = size_vertical * math.floor(pIt / cnt_row)
            start_col = size_horizontal * (pIt % cnt_row)
            temp[start_row:start_row + size_vertical, start_col:start_col + size_horizontal] =\
                plist.pieceData[:size_vertical, :size_horizontal]
        return temp

    temp = np.zeros((size_vertical * cnt_col, size_horizontal * cnt_row, chn), dtype=np.uint8)
    for pIt in range(cnt_total):
        print(str(pIt) + " out of total is " + str(cnt_total))
        start_row = size_vertical * math.floor(pIt / cnt_row)
        start_col = size_horizontal * (pIt % cnt_row)
        temp[start_row:start_row + size_vertical, start_col:start_col + size_horizontal] =\
            plist[pIt].pieceData[:size_vertical, :size_horizontal]
    return temp


# show result puzzle image
def draw_image(temp, window_name):
    print("starting draw image")
    cv2.imshow(window_name, temp)
    print("waiting on image")
    # cv2.waitKey(0)
    # print("waiting on image")
    # cv2.destroyAllWindows()
