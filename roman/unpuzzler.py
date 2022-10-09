import cv2
import numpy as np
import os
import sys
from pathlib import Path
from checkInput import check_input
import readImage as readImg
from makePieces import get_pieces
import drawPieces as drawP
import Piece
from PIL import Image
import imagehash
import collections
# if __name__ == "__main__":
#     # initialization and read image
#     check_input("unpuzzler.py")
#     filename = sys.argv[1]
#     fileTest = Path(filename)
#     if not fileTest.is_file():
#         print("Image file NOT found.")
#         sys.exit(1)
#     img, imgRow, imgCol, imgChn = readImg.read_image(filename, cv2.IMREAD_COLOR)

#     # get pieces from image
#     pList, pSize_vertical, pSize_horizontal, pCnt_row, pCnt_column, pCnt_total = get_pieces(img, imgRow, imgCol, imgChn)

#     # calculate difference between every piece
#     for i in range(pCnt_total):
#         for j in range(i + 1, pCnt_total):
#             if i == j:
#                 continue
#             Piece.piece_difference(pList[i], pList[j])

#     # determine starting pixel to fill in image
#     startPiece = None
#     for piece in pList:
#         Piece.find_neighbors(piece)
#         if piece.neighbors[0] is None and piece.neighbors[3] is None:
#             startPiece = piece
#     if startPiece is None:
#         print("Could not find starting piece... but here is an attempt.")
#         startPiece = pList[0]

#     # fill in image using neighbor information
#     black = np.zeros((pSize_vertical, pSize_horizontal, imgChn), dtype=np.uint8)
#     blackPiece = Piece.Piece(-1, pSize_vertical, pSize_horizontal, imgChn, (0, 0), black, pCnt_total)
#     temp = [blackPiece for x in range(pCnt_total)]
#     temp[0] = startPiece
#     for i in range(pCnt_total):
#         if i % pCnt_row < pCnt_row - 1 and temp[i].neighbors[1] is not None:
#             temp[i + 1] = pList[temp[i].neighbors[1]]
#         if i / pCnt_row < pCnt_column - 1 and temp[i].neighbors[2] is not None:
#             temp[i + pCnt_row] = pList[temp[i].neighbors[2]]

#     # show and save result image
#     filename = os.path.splitext(filename)[0] + "_solve.png"
#     temp = drawP.combine_pieces(pSize_vertical, pSize_horizontal, pCnt_row, pCnt_column, pCnt_total, imgChn, temp)
#     drawP.draw_image(temp, filename + " - Solved Image")
#     cv2.imwrite(filename, temp)
def unpuzzle(filename):
    
    # initialization and read image
    # check_input("unpuzzler.py")
    
    #----
    # filename = sys.argv[1]
    
    # fileTest = Path(filename)
    # if not fileTest.is_file():
    #     print("Image file NOT found.")
    #     sys.exit(1)
    #----
    print(filename)
    img, imgRow, imgCol, imgChn = readImg.read_image(filename, cv2.IMREAD_COLOR)
    initial_vals = collections.defaultdict(int)
    final_vals = []
    # get pieces from image
    pList, pSize_vertical, pSize_horizontal, pCnt_row, pCnt_column, pCnt_total = get_pieces(img, imgRow, imgCol, imgChn)
    #temp2 = drawP.combine_pieces(pSize_vertical, pSize_horizontal, pCnt_row, pCnt_column, pCnt_total, pList, temp2)
    #-----
    # filename = os.path.splitext(filename)[0] + "_solve.jpeg"
    for count,val in enumerate(pList):
        #new_filename = os.path.splitext("scrambled_" + str(count + 1) + "_" + filename)
        new_filename = filename + "_scrambled_" + str(count + 1)
        temp = drawP.combine_pieces(pSize_vertical, pSize_horizontal, pCnt_row, pCnt_column, pCnt_total, imgChn, val, True)
        drawP.draw_image(temp, new_filename + " - Puzzle Image")
        cv2.imwrite(new_filename, temp)
        initial_vals[(str(imagehash.average_hash(Image.open(new_filename))))] = count
    #------   
    #cv2.imwrite(filename, drawP.draw_image(pList[0], filename + " - Solved Image"))



    # calculate difference between every piece
    for i in range(pCnt_total):
        for j in range(i + 1, pCnt_total):
            if i == j:
                continue
            Piece.piece_difference(pList[i], pList[j])

    # determine starting pixel to fill in image
    startPiece = None
    for piece in pList:
        Piece.find_neighbors(piece)
        if piece.neighbors[0] is None and piece.neighbors[3] is None:
            startPiece = piece
    if startPiece is None:
        print("Could not find starting piece... but here is an attempt.")
        startPiece = pList[0]

    # fill in image using neighbor information
    black = np.zeros((pSize_vertical, pSize_horizontal, imgChn), dtype=np.uint8)
    blackPiece = Piece.Piece(-1, pSize_vertical, pSize_horizontal, imgChn, (0, 0), black, pCnt_total)
    temp = [blackPiece for x in range(pCnt_total)]
    temp[0] = startPiece
    for i in range(pCnt_total):
        if i % pCnt_row < pCnt_row - 1 and temp[i].neighbors[1] is not None:
            temp[i + 1] = pList[temp[i].neighbors[1]]
        if i / pCnt_row < pCnt_column - 1 and temp[i].neighbors[2] is not None:
            temp[i + pCnt_row] = pList[temp[i].neighbors[2]]
    print("executing 1")
    # show and save result image
    filename = os.path.splitext(filename)[0] + "_solver.jpeg"
    print("executing 2")
    temp = drawP.combine_pieces(pSize_vertical, pSize_horizontal, pCnt_row, pCnt_column, pCnt_total, imgChn, temp)
    print("executing 3")
    drawP.draw_image(temp, filename + " - Solved Image")
    print("executing 4")
    cv2.imwrite(filename, temp)
    print("executing 5")
    
    img, imgRow, imgCol, imgChn = readImg.read_image(filename, cv2.IMREAD_COLOR)

        # get pieces from image
    pList, pSize_vertical, pSize_horizontal, pCnt_row, pCnt_column, pCnt_total = get_pieces(img, imgRow, imgCol, imgChn)

        # randomize pieces
        #shuffle(pList)

        # show and save result image
        # print(type(pList))
    #-----
    for count,val in enumerate(pList):
        new_filename2 = filename + '_' + str(count + 1) + "_" + 'final'
        temp = drawP.combine_pieces(pSize_vertical, pSize_horizontal, pCnt_row, pCnt_column, pCnt_total, imgChn, val, True)
        drawP.draw_image(temp, new_filename2 + " - Puzzle Image")
        cv2.imwrite(new_filename2, temp)
        final_vals.append(str(imagehash.average_hash(Image.open(new_filename2))))
        
   
    outfile = open("results.txt", "a")
    for val in final_vals:
        outfile.write(str(initial_vals[val]))
    outfile.write('\n')
     #------

    
    
