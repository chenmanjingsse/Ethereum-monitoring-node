import pandas as pd
import numpy as np
import sys
import os
from pathlib import Path

#save to file
import matplotlib as mpl
mpl.use('Agg')

import matplotlib.pyplot as plt

if len(sys.argv) != 2:
    sys.exit(sys.argv[0], ": expecting 1 parameter - txs-propagation-times.log.")

TXS_LOG = sys.argv[1] #"txs-propagation-times.log"
if not os.path.isfile(TXS_LOG):
    sys.exit(TXS_LOG, ": does not exists!")

dtypes = {
        'Hash'          : 'object',
        'ValidityErr'   : 'object',
        'AngainorTimeStamp' : 'object',
        'FalconTimeStamp'   : 'object',
        'PositiveDif'       : 'float',
        'AngainorMinusFalcon' : 'float',
        }

txs = pd.read_csv(TXS_LOG, 
    names=['Hash','ValidityErr','AngainorTimeStamp','FalconTimeStamp',
        'PositiveDif','AngainorMinusFalcon'],
    dtype=dtypes)

#basic info
print(len(txs))
print("delay AngainorMinusFalcon max:", txs['AngainorMinusFalcon'].max())
print("delay AngainorMinusFalcon min:", txs['AngainorMinusFalcon'].min())
print("delay PositiveDif max:", txs['PositiveDif'].max())

## this takes some time (it sorts...)
#print(" mean", txs['PositiveDif'].mean(), "median", txs['PositiveDif'].median())

## DROP txs here
# Drop  txs received after network outage  BASED on   TIME range ..
#txs = all_txs.assign(NeedDrop = np.nan)

## loop through all txs
#for i in txs.index:
#    #CapturedLocally is False?
#    if txs.at[i,'PositiveDif'] >= 20:
#        print(txs.at[i,'Number'],txs.at[i,'AngainorTimeStamp'], txs.at[i,'PositiveDif'])    

#TMP
#exit(0)

max_delay = 0.5 #txs['PositiveDif'].max()

bin_seq = list(np.arange(0,max_delay,0.005))    # (0,  MAX PositiveDif,  step size) 

fig, ax = plt.subplots()
counts, bin_edges = np.histogram (txs['PositiveDif'], bins=bin_seq)

plt.xlabel('Time since first transaction observation [s]')

ax.bar (bin_edges[:-1], counts, width=0.005)
#ax.plot (bin_edges[:-1], counts)

#this helps to set   tyicks down
print(counts)
print(bin_edges)

# be careful here!! 
num_of_y_dots = 10 #  range 20 -> total = 0.05..   num_dots = 5 ...    0.05 /5 = 0.0§  one dot
print("counts",len(counts), type(counts), "max", counts.max(),counts.sum() )
plt.yticks(np.arange(0, counts.max() + counts.max()/num_of_y_dots, counts.max()/num_of_y_dots ),['0','0.02','0.04','0.06','0.08','0.10','0.12','0.15','0.17','0.19','0.21'])   

plt.xscale('symlog')
ax.set_xlim(left=0)
ax.set_xlim(right=max_delay)

nums = [0,0.1,0.2,0.3,0.5]
labels = ['0','0.1','0.2','0.3','0.5']

#nums = [0,1,2,5,10,100,500,1000,1500,2637]
#labels = ['0','1','2','5','10','100','500','1000','1500','2637']

plt.xticks(nums, labels)


#LOCAL show
#plt.show()
#save to file
plt.savefig('5-2-tx-propagation-time.pdf')



