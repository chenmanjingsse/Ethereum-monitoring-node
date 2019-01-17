import pandas as pd                     
import numpy as np  

# to make it work on server w/o graphics
import matplotlib as mpl
mpl.use('Agg')

import matplotlib.pyplot as plt
import matplotlib.dates as mdates

PEERS_CSV = "peers.csv"

peers = pd.read_csv(PEERS_CSV, usecols=['LocalTimeStamp', 'Peers', 'InboundPeers'], parse_dates=['LocalTimeStamp'])
peers.rename(columns={'LocalTimeStamp': 'Time'}, inplace=True)
peers.set_index('Time',inplace=True)

peers.plot()

#show now
#plt.show()

#save to pdf
plt.savefig('peers.pdf')
