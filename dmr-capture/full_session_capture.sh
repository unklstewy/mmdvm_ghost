#!/bin/bash

# Full DMR Session Capture Script
# Captures complete MMDVM_Host session from startup to shutdown

CAPTURE_DIR="/home/pi-star/dmr_captures"
DATE_STAMP=$(date '+%Y%m%d_%H%M%S')
PCAP_FILE="${CAPTURE_DIR}/full_session_${DATE_STAMP}.pcap"
LOG_FILE="${CAPTURE_DIR}/session_${DATE_STAMP}.log"

mkdir -p "$CAPTURE_DIR"

echo "=== DMR FULL SESSION CAPTURE ===" | tee "$LOG_FILE"
echo "Start time: $(date)" | tee -a "$LOG_FILE"
echo "Capture file: $PCAP_FILE" | tee -a "$LOG_FILE"
echo "" | tee -a "$LOG_FILE"

# Step 1: Start packet capture in background
echo "[1/4] Starting packet capture..." | tee -a "$LOG_FILE"
sudo tcpdump -i lo -s 0 -w "$PCAP_FILE" 'port 62031 or port 62032' &
TCPDUMP_PID=$!
echo "Packet capture started (PID: $TCPDUMP_PID)" | tee -a "$LOG_FILE"
sleep 2

# Step 2: Stop MMDVM_Host
echo "[2/4] Stopping MMDVM_Host..." | tee -a "$LOG_FILE"
sudo systemctl stop mmdvmhost
sleep 3

# Step 3: Start MMDVM_Host and capture session
echo "[3/4] Starting MMDVM_Host and capturing session..." | tee -a "$LOG_FILE"
sudo systemctl start mmdvmhost
echo "MMDVM_Host restarted, capturing for 10 minutes..." | tee -a "$LOG_FILE"

# Monitor for 10 minutes
for i in {1..10}; do
    echo "Minute $i/10 - $(date)" | tee -a "$LOG_FILE"
    sleep 60
done

# Step 4: Stop capture cleanly
echo "[4/4] Stopping capture..." | tee -a "$LOG_FILE"
sudo kill -TERM $TCPDUMP_PID
sleep 2

echo "" | tee -a "$LOG_FILE"
echo "=== CAPTURE COMPLETE ===" | tee -a "$LOG_FILE"
echo "End time: $(date)" | tee -a "$LOG_FILE"
echo "File: $PCAP_FILE" | tee -a "$LOG_FILE"
echo "Size: $(ls -lh $PCAP_FILE 2>/dev/null | awk '{print $5}' || echo 'N/A')" | tee -a "$LOG_FILE"
echo "" | tee -a "$LOG_FILE"

# Display summary
if [ -f "$PCAP_FILE" ]; then
    PACKET_COUNT=$(sudo tcpdump -r "$PCAP_FILE" 2>/dev/null | wc -l)
    echo "Packets captured: $PACKET_COUNT" | tee -a "$LOG_FILE"
    echo "" | tee -a "$LOG_FILE"
    echo "Analysis commands:" | tee -a "$LOG_FILE"
    echo "  tcpdump -r $PCAP_FILE" | tee -a "$LOG_FILE"
    echo "  tcpdump -r $PCAP_FILE -tttt" | tee -a "$LOG_FILE"
else
    echo "Warning: Capture file not found!" | tee -a "$LOG_FILE"
fi

echo "Log file: $LOG_FILE"
