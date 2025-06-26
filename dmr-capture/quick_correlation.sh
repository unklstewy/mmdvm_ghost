#!/bin/bash

# Quick DMR Correlation Analysis
PCAP_FILE="/home/pi-star/dmr_captures/full_session_20250626_104623.pcap"
LOG_FILE="/var/log/pi-star/MMDVM-$(date '+%Y-%m-%d').log"
OUTPUT_FILE="/home/pi-star/dmr_captures/quick_correlation_$(date '+%Y%m%d_%H%M%S').txt"

echo "DMR PACKET AND LOG CORRELATION ANALYSIS" > "$OUTPUT_FILE"
echo "=======================================" >> "$OUTPUT_FILE"
echo "Generated: $(date)" >> "$OUTPUT_FILE"
echo "PCAP: $PCAP_FILE" >> "$OUTPUT_FILE"
echo "LOG:  $LOG_FILE" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Get packet statistics
echo "PACKET CAPTURE ANALYSIS" >> "$OUTPUT_FILE"
echo "-----------------------" >> "$OUTPUT_FILE"
TOTAL_PACKETS=$(sudo tcpdump -r "$PCAP_FILE" 2>/dev/null | wc -l)
FIRST_PACKET=$(sudo tcpdump -r "$PCAP_FILE" -tttt 2>/dev/null | head -1 | awk '{print $1, $2}')
LAST_PACKET=$(sudo tcpdump -r "$PCAP_FILE" -tttt 2>/dev/null | tail -1 | awk '{print $1, $2}')
FILE_SIZE=$(ls -lh "$PCAP_FILE" | awk '{print $5}')

echo "Total Packets: $TOTAL_PACKETS" >> "$OUTPUT_FILE"
echo "File Size: $FILE_SIZE" >> "$OUTPUT_FILE"
echo "First Packet: $FIRST_PACKET" >> "$OUTPUT_FILE"
echo "Last Packet: $LAST_PACKET" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Analyze packet flow
echo "PACKET FLOW DIRECTION" >> "$OUTPUT_FILE"
echo "---------------------" >> "$OUTPUT_FILE"
TO_MMDVM=$(sudo tcpdump -r "$PCAP_FILE" 2>/dev/null | grep "62031 > localhost.62032" | wc -l)
FROM_MMDVM=$(sudo tcpdump -r "$PCAP_FILE" 2>/dev/null | grep "62032 > localhost.62031" | wc -l)
echo "Network → MMDVM (port 62031→62032): $TO_MMDVM packets" >> "$OUTPUT_FILE"
echo "MMDVM → Network (port 62032→62031): $FROM_MMDVM packets" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Voice session analysis from logs
echo "VOICE SESSION SUMMARY" >> "$OUTPUT_FILE"
echo "---------------------" >> "$OUTPUT_FILE"
echo "Callsign     TG    Duration  Loss%  BER%   Start Time" >> "$OUTPUT_FILE"
echo "----------------------------------------------------" >> "$OUTPUT_FILE"

# Extract voice sessions from today's log within capture timeframe
grep "end of voice transmission" "$LOG_FILE" | while read line; do
    CALLSIGN=$(echo "$line" | sed -n 's/.*from \([A-Z0-9]*\) to TG.*/\1/p')
    TG=$(echo "$line" | sed -n 's/.*to TG \([0-9]*\).*/\1/p')
    DURATION=$(echo "$line" | sed -n 's/.* \([0-9]*\.[0-9]*\) seconds.*/\1/p')
    LOSS=$(echo "$line" | sed -n 's/.* \([0-9]*\)% packet loss.*/\1/p')
    BER=$(echo "$line" | sed -n 's/.*BER: \([0-9]*\.[0-9]*\)%.*/\1/p')
    TIME=$(echo "$line" | awk '{print $3}')
    
    if [[ -n "$CALLSIGN" && -n "$TG" && -n "$DURATION" ]]; then
        printf "%-12s %-5s %-9s %-6s %-6s %s\n" "$CALLSIGN" "$TG" "$DURATION" "$LOSS" "$BER" "$TIME" >> "$OUTPUT_FILE"
    fi
done

echo "" >> "$OUTPUT_FILE"

# Timeline correlation - show packets around voice events
echo "TIMELINE CORRELATION (Voice Events + Packet Activity)" >> "$OUTPUT_FILE"
echo "=====================================================" >> "$OUTPUT_FILE"

# Get voice headers and correlate with packet timestamps
grep "voice header\|end of voice" "$LOG_FILE" | head -20 | while read line; do
    TIMESTAMP=$(echo "$line" | awk '{print $3}')
    EVENT=$(echo "$line" | cut -d' ' -f4-)
    
    echo "[$TIMESTAMP] $EVENT" >> "$OUTPUT_FILE"
    
    # Find packets around this time (±1 second)
    HOUR=$(echo $TIMESTAMP | cut -d':' -f1)
    MIN=$(echo $TIMESTAMP | cut -d':' -f2)  
    SEC=$(echo $TIMESTAMP | cut -d':' -f3 | cut -d'.' -f1)
    
    # Show packet activity around this time
    sudo tcpdump -r "$PCAP_FILE" -tttt 2>/dev/null | grep "$HOUR:$MIN:$SEC" | head -3 | while read pkt; do
        echo "    PACKET: $pkt" >> "$OUTPUT_FILE"
    done
    echo "" >> "$OUTPUT_FILE"
done

echo "ANALYSIS COMPLETE" >> "$OUTPUT_FILE"
echo "Output saved to: $OUTPUT_FILE"
