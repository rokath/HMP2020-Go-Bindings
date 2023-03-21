
# Show help.
# hmp -help

# Show com ports.
# hmp -s

# Switch channel 2 OFF.
# hmp -p com7 -ch 2 -output=false

# Switch channel 2 ON.
# hmp -p com7 -ch 2 -output

# Switch 3 times all active channels 2s ON, 1s OFF and beep when finished.
# hmp -p com7 -ch -1 -msON 2000 -count 3 -beep

# Set channel 1 to 12.345 Volt.
# hmp -p com7 -V 12.345 

# Set channel 2 to 13.2V@0.5A and reduce voltage in steps of 0.1V every 200 ms in 150 steps and beep finally.
# hmp -p com7 -ch 2 -A 0.5 -V 13.2 -stepMs 20 -count 150 -beep -v

# Enable general output.
# hmp -p com7 -ch -1 -output

