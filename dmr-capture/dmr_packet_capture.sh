#!/bin/bash

# DMR Packet Capture Script for W3MSG WPSD
# Captures DMR network traffic on ports 62031/62032

CAPTURE_DIR="/home/pi-star/dmr_captures"
DATE_STAMP=$(date '+%Y%m%d_%H%M%S')
PCAP_FILE="${CAPTURE_DIR}/dmr_${DATE_STAMP}.pcap"

# Create capture directory if it doesn't exist
mkdir -p "$CAPTURE_DIR"

# Function to display usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -t, --time SECONDS    Capture for specified time (default: 300 seconds)"
    echo "  -c, --count PACKETS   Capture specified number of packets (default: unlimited)"
    echo "  -f, --file FILENAME   Output filename (default: auto-generated)"
    echo "  -h, --help           Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 -t 60             # Capture for 1 minute"
    echo "  $0 -c 1000           # Capture 1000 packets"
    echo "  $0 -f my_capture.pcap # Save to specific file"
}

# Default values
DURATION=""
PACKET_COUNT=""
CUSTOM_FILE=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--time)
            DURATION="$2"
            shift 2
            ;;
        -c|--count)
            PACKET_COUNT="$2"
            shift 2
            ;;
        -f|--file)
            CUSTOM_FILE="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Set output file
if [[ -n "$CUSTOM_FILE" ]]; then
    PCAP_FILE="${CAPTURE_DIR}/${CUSTOM_FILE}"
fi

# Build tcpdump command
TCPDUMP_CMD="sudo tcpdump -i lo -s 0"
TCPDUMP_CMD+=" -w $PCAP_FILE"
TCPDUMP_CMD+=" 'port 62031 or port 62032'"

# Add time limit if specified
if [[ -n "$DURATION" ]]; then
    TCPDUMP_CMD+=" -G $DURATION -W 1"
fi

# Add packet count if specified
if [[ -n "$PACKET_COUNT" ]]; then
    TCPDUMP_CMD+=" -c $PACKET_COUNT"
fi

echo "Starting DMR packet capture..."
echo "Output file: $PCAP_FILE"
echo "Command: $TCPDUMP_CMD"
echo "Press Ctrl+C to stop capture"
echo ""

# Execute capture
eval $TCPDUMP_CMD

echo ""
echo "Capture completed!"
echo "File saved: $PCAP_FILE"
echo "File size: $(ls -lh $PCAP_FILE | awk '{print $5}')"
echo ""
echo "To analyze the capture:"
echo "  tcpdump -r $PCAP_FILE"
echo "  wireshark $PCAP_FILE (if Wireshark is installed)"
