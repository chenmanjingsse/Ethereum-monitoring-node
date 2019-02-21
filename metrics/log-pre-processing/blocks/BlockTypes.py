#!/usr/bin/python3

import pandas as pd
import numpy as np


BLOCKS_FINAL_LOG = "blocksMergedCheckManuallyNow.log"
#takes preprocessed blocksMerged
blocks = pd.read_csv(BLOCKS_FINAL_LOG, 
    names=['LocalTimeStamp','BlockHash','Number','GasLimit','GasUsed','Difficulty','Time',
    'Coinbase','ParentHash','UncleHash','BlockSize','ListOfTxs','ListOfUncles','BlockType'])

print("Total: ", len(blocks))

#prints number of Main   (check with num of Main in heads)
print("Main: ", len(blocks[blocks.BlockType == "Main"]))

#print num of uncles
print("Uncle: ", len(blocks[blocks.BlockType == "Uncle"]))

#print num of rec uncles
print("Recognized: ", len(blocks[blocks.BlockType == "Recognized"]))

#print num of rec uncles
print("CheckManually: ", len(blocks[blocks.BlockType == "CheckManually"]))

someCheckMan = False
someNan = False

# print all CheckManually
for i, row in blocks.iterrows():
    if row['BlockType'] == "CheckManually":
        print("CheckManually block with num: ", int(row['Number']))
        if not someCheckMan:
            someCheckMan = True

#for i, row in blocks.iterrows():
#    if row['BlockType'] is np.nan: #CheckManually
#        print("BlockType not set in block num: ", int(row['Number']))
#        if not someNan:
#            someNan = True




if not someCheckMan:
    print("there are NO blocks with blockType==CheckManually so nothing to do here")
else:
    print("there ARE some blocks with blockType==CheckManually so fix it")


if not someNan:
    print("there are NO blocks with blockType==NaN so nothing to do here")
else:
    print("there ARE some blocks with blockType==NaN so fix it")