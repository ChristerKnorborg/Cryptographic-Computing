import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

def makePlot():
    
    file = pd.read_csv("fixed_m_data.csv")
    
    x_axis = (file['m_size'])
    
    OTExtension = file['time_OT_extension']/x_axis
    labels = [ 'Time for OT extension']

    
    plt.scatter(x_axis ,OTExtension, c="blue", label="OT Extension")
    plt.legend(loc="upper left")
    plt.xlabel("length of pattern")
    plt.ylabel("(time ns) / (number of message pairs)")
    
    #plt.ticklabel_format(useOffset=False, style='plain')
    
    plt.show()



makePlot()