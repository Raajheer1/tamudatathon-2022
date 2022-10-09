import os
import unpuzzler 
for filename in os.listdir('/Users/roman/desktop/train/2103'):
    #print('/Users/roman/desktop/train/0123/' + str(filename))
    unpuzzler.unpuzzle('/Users/roman/desktop/train/2103/' + str(filename))