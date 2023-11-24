import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

def makePlot():
    # Load the data from the CSV file
    file = pd.read_csv("fixed_m_data.csv")
    
    # Extract data for the plot
    x_axis = file['m_size']
    OTBasic = file['time_OT_Basic']
    OTExtension = file['time_OT_Extension']

    # Creating the line plot
    plt.plot(x_axis, OTBasic, 'g-', marker='o', label="OT Basic")
    plt.plot(x_axis, OTExtension, 'b-', marker='o', label="OT Extension")
    plt.legend(loc="upper left")

    # Setting the axis to log scale
    plt.xscale('log', base=2)
    plt.yscale('log', base=2)

    plt.xlabel("Amount of message pairs")
    plt.ylabel("Time (seconds)")
    plt.grid(True)

    # Display the plot
    plt.show()

# Run the plot function
makePlot()