================================================================================
W3MSG DMR REPEATER - COMPLETE PACKET CAPTURE & ANALYSIS ARCHIVE
================================================================================
Generated: Thu 26 Jun 2025
Station: W3MSG (ID: 3141966)
Location: Charlotte, EM95mc
Frequencies: RX 433.450 MHz / TX 438.450 MHz

================================================================================
ARCHIVE CONTENTS
================================================================================

📡 PACKET CAPTURE FILES:
-------------------------
full_session_20250626_104623.pcap
  - Complete 10+ minute DMR session capture (741KB)
  - 6,704 packets from MMDVM restart through voice traffic
  - Captures authentication, registration, and live conversations
  - Format: Standard pcap file (analyze with tcpdump/Wireshark)

dmr_capture_20250626_103944.pcap
  - Initial test capture (100 packets, 11KB)
  - Proof of concept packet capture

📊 ANALYSIS REPORTS:
-------------------
correlation_analysis_20250626_110025.txt
  - Comprehensive packet/log correlation analysis
  - Voice session breakdown with quality metrics
  - International participant analysis
  - Timeline correlations

quick_correlation_20250626_110145.txt
  - Executive summary of session analysis
  - Traffic flow statistics (97.9% inbound, 2.1% outbound)  
  - Voice activity summary for 73+ transmissions
  - Key performance metrics

session_20250626_104623.log
  - Session capture log with timestamps
  - Documents the complete capture process

📝 LOG FILES:
-------------
MMDVM_session_log.log
  - Complete MMDVM_Host log for the session date
  - Contains all voice headers, transmissions, and system events
  - Shows international DMR network activity on TG 91

🛠️ TOOLS & SCRIPTS:
------------------
dmr_packet_capture.sh
  - Primary packet capture tool with options
  - Usage: ./dmr_packet_capture.sh -t 600 (10 minutes)
  - Flexible timing and output options

full_session_capture.sh  
  - Complete session lifecycle capture
  - Automatically restarts MMDVM_Host and captures full session
  - Includes packet capture + system monitoring

dmr_correlation_analysis.py
  - Python script for advanced packet/log correlation
  - Cross-references timestamps between pcap and logs
  - Generates detailed voice session analytics

quick_correlation.sh
  - Bash script for rapid analysis
  - Extracts key metrics and correlations
  - Good for quick session summaries

⚙️ CONFIGURATION:
-----------------
etc/mmdvmhost
  - Complete MMDVM_Host configuration file
  - Shows DMR setup, network settings, and hardware config
  - Documents repeater parameters and network connections

================================================================================
KEY FINDINGS FROM ANALYSIS
================================================================================

🎯 TRAFFIC STATISTICS:
- Total Packets: 6,704 (10+ minutes)
- Network→MMDVM: 6,560 packets (97.9%)
- MMDVM→Network: 144 packets (2.1%) 
- Data Volume: 741KB total

📻 VOICE ACTIVITY:
- Your RF Transmissions: 4 sessions (W3MSG)
- Network Transmissions: 70+ international sessions
- Talk Group: TG 91 (Worldwide)
- Countries: Bangladesh, UAE, Malaysia, Brazil, Scotland, Moldova, etc.

📊 QUALITY METRICS:
- Most sessions: 0% packet loss
- BER typically < 1%
- Excellent network connectivity
- International reach confirmed

🌍 INTERNATIONAL PARTICIPANTS:
- S21JSR (Bangladesh) - Multiple long sessions
- A46BCW (UAE) - Very active participant  
- 9M2SFL (Malaysia) - Regular contributor
- PU2XUP (Brazil) - Extended conversations
- MM7SVI (Scotland) - Consistent participant
- ER1CW (Moldova) - Several sessions

================================================================================
ANALYSIS TOOLS
================================================================================

🔍 VIEWING PACKET CAPTURES:
- tcpdump -r filename.pcap
- tcpdump -r filename.pcap -tttt (with timestamps)
- tcpdump -r filename.pcap -x (hex dump)
- wireshark filename.pcap (if Wireshark available)

📈 UNDERSTANDING THE DATA:
- Port 62031→62032: Network to MMDVM_Host
- Port 62032→62031: MMDVM_Host to Network
- UDP length 55: Standard DMR frame size
- High inbound ratio: Active network repeater role

🕐 TIME CORRELATION:
- Packet timestamps correlate with voice events
- Voice headers trigger packet bursts
- End transmissions show in both logs and packets
- Network keepalives maintain connectivity

================================================================================
TECHNICAL DETAILS
================================================================================

🏗️ CAPTURE METHOD:
- Interface: Loopback (lo)
- Ports: 62031, 62032 (DMR network)
- Protocol: UDP
- Tools: tcpdump, systemctl, grep, awk

📡 DMR CONFIGURATION:
- Color Code: 1
- Repeater ID: 3141966
- Whitelist: 3141966, 3176627
- Network: Local gateway connection
- Modes: DMR enabled, others disabled

🔧 QUALITY ASSURANCE:
- 0 packets dropped during capture
- Complete session lifecycle captured
- Proper correlation between multiple data sources
- Validated against live MMDVM_Host logs

================================================================================
USAGE RECOMMENDATIONS
================================================================================

🎯 FOR ANALYSIS:
1. Start with quick_correlation summary
2. Review voice session details
3. Examine packet flows in pcap files
4. Cross-reference with MMDVM logs

🛠️ FOR FUTURE CAPTURES:
1. Use dmr_packet_capture.sh for routine monitoring
2. Use full_session_capture.sh for complete lifecycle analysis
3. Adjust timing parameters as needed
4. Archive important sessions for historical analysis

📊 FOR TROUBLESHOOTING:
1. Check packet loss percentages in analysis
2. Review BER values for RF quality
3. Monitor traffic ratios for network health
4. Correlate timestamps for timing issues

================================================================================
SUPPORT & REFERENCES
================================================================================

📚 DMR Protocol Information:
- ETSI TS 102 361 standard
- DMR Association resources
- MMDVM documentation

🔧 Tools Used:
- MMDVM_Host (DMR software)
- tcpdump (packet capture)
- WPSD Dashboard (W0CHP variant)
- Custom correlation scripts

📞 Contact:
- Station: W3MSG
- QRZ: https://www.qrz.com/db/W3MSG

================================================================================
END OF DOCUMENTATION
================================================================================
