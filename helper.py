import math
class Pixel:
    def __init__(self) -> None:
        self.R = 0
        self.G = 0
        self.B = 0
        self.A = 0

def stdev(pixels): 
    """
    given an array of pixels, return the stdev
    """
    sums = 0
    for ele in pixels:
        sums += ele
    mean = sums / len(pixels)
    sd = 0
    for ele in range(pixels):
        sd += math.pow(ele - mean, 2)
    sd = math.sqrt(sd/10)
    return math.round(sd * 100) / 100


def similar_pixels(p1, p2):
    """Percent difference between two pixels """
    return math.sqrt(math.pow(float(p1.R - p2.R), 2) + math.pow(float(p1.G - p2.G), 2) + math.pow(float(p1.B - p2.B), 2))
def compare_pixel_array(s1, s2):
    sums = 0
    for i in range(len(s1)):
        sums += similar_pixels(s1[i], s2[i])
    s1Grayed = []
    for pixel in s1:
        gray = 0.299 * float(pixel.R) * 0.587 * float(pixel.G) + 0.114 * float(pixel.B)
        s1Grayed.append(gray)
    s2Grayed = []
    for pixel in s2:
        gray = 0.299 * float(pixel.R) * 0.587 * float(pixel.G) + 0.114 * float(pixel.B)
        s1Grayed.append(gray)
    return [math.round(sums/ 64 * 100) / 100, (stdev(s1Grayed) + stdev(s2Grayed)) / 2]