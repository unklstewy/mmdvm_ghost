#!/usr/bin/env python3
"""
DMR Packet Capture and Log Correlation Analysis
Correlates tcpdump pcap data with MMDVM_Host logs for comprehensive session analysis
"""

import subprocess
import re
import datetime
import sys
import os
from collections import defaultdict

class DMRAnalyzer:
    def __init__(self, pcap_file, log_file, output_file):
        self.pcap_file = pcap_file
        self.log_file = log_file
        self.output_file = output_file
        self.packets = []
        self.log_entries = []
        self.correlated_events = []
        
    def parse_pcap(self):
        """Extract packet data from pcap file"""
        print("Parsing packet capture...")
        try:
            cmd = f"sudo tcpdump -r {self.pcap_file} -tttt -n"
            result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
            
            for line in result.stdout.split('\n'):
                if 'localhost.62031' in line or 'localhost.62032' in line:
                    # Parse timestamp and packet info
                    match = re.match(r'(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+).*UDP, length (\d+)', line)
                    if match:
                        timestamp_str, length = match.groups()
                        timestamp = datetime.datetime.strptime(timestamp_str[:19], '%Y-%m-%d %H:%M:%S')
                        microsec = int(timestamp_str.split('.')[1])
                        timestamp = timestamp.replace(microsecond=microsec)
                        
                        direction = "TO_MMDVM" if "62031 > localhost.62032" in line else "FROM_MMDVM"
                        
                        self.packets.append({
                            'timestamp': timestamp,
                            'direction': direction,
                            'length': int(length),
                            'raw_line': line
                        })
        except Exception as e:
            print(f"Error parsing pcap: {e}")
    
    def parse_mmdvm_log(self):
        """Extract log entries from MMDVM_Host log"""
        print("Parsing MMDVM_Host log...")
        try:
            with open(self.log_file, 'r') as f:
                for line in f:
                    # Parse MMDVM log format: Level: YYYY-MM-DD HH:MM:SS.mmm Message
                    match = re.match(r'([IME]): (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+)\s+(.*)', line.strip())
                    if match:
                        level, timestamp_str, message = match.groups()
                        timestamp = datetime.datetime.strptime(timestamp_str[:19], '%Y-%m-%d %H:%M:%S')
                        microsec = int(timestamp_str.split('.')[1]) * 1000  # Convert to microseconds
                        timestamp = timestamp.replace(microsecond=microsec)
                        
                        self.log_entries.append({
                            'timestamp': timestamp,
                            'level': level,
                            'message': message,
                            'raw_line': line.strip()
                        })
        except Exception as e:
            print(f"Error parsing log: {e}")
    
    def correlate_events(self):
        """Correlate packet events with log entries"""
        print("Correlating events...")
        
        # Group packets by time windows (1 second)
        packet_groups = defaultdict(list)
        for packet in self.packets:
            time_key = packet['timestamp'].replace(microsecond=0)
            packet_groups[time_key].append(packet)
        
        # Correlate with log entries
        for log_entry in self.log_entries:
            log_time = log_entry['timestamp'].replace(microsecond=0)
            
            # Find packets within ±2 seconds of log entry
            related_packets = []
            for i in range(-2, 3):
                check_time = log_time + datetime.timedelta(seconds=i)
                if check_time in packet_groups:
                    related_packets.extend(packet_groups[check_time])
            
            # Create correlation entry
            correlation = {
                'timestamp': log_entry['timestamp'],
                'log_entry': log_entry,
                'related_packets': related_packets,
                'packet_count': len(related_packets)
            }
            
            self.correlated_events.append(correlation)
    
    def analyze_patterns(self):
        """Analyze communication patterns"""
        print("Analyzing patterns...")
        
        analysis = {
            'total_packets': len(self.packets),
            'total_log_entries': len(self.log_entries),
            'voice_transmissions': 0,
            'network_registrations': 0,
            'data_volume': sum(p['length'] for p in self.packets),
            'timespan': None,
            'average_packet_rate': 0,
            'voice_sessions': [],
            'errors': []
        }
        
        if self.packets:
            start_time = min(p['timestamp'] for p in self.packets)
            end_time = max(p['timestamp'] for p in self.packets)
            analysis['timespan'] = end_time - start_time
            analysis['average_packet_rate'] = len(self.packets) / analysis['timespan'].total_seconds()
        
        # Analyze log entries for specific events
        current_voice_session = None
        for log_entry in self.log_entries:
            message = log_entry['message']
            
            if 'voice header' in message:
                analysis['voice_transmissions'] += 1
                # Extract callsign and TG
                match = re.search(r'from (\w+) to TG (\d+)', message)
                if match:
                    callsign, tg = match.groups()
                    current_voice_session = {
                        'start_time': log_entry['timestamp'],
                        'callsign': callsign,
                        'talkgroup': tg,
                        'slot': 'Slot 1' if 'Slot 1' in message else 'Slot 2'
                    }
            
            elif 'end of voice transmission' in message and current_voice_session:
                # Extract duration and quality metrics
                match = re.search(r'(\d+\.\d+) seconds.*?(\d+)% packet loss.*?BER: (\d+\.\d+)%', message)
                if match:
                    duration, packet_loss, ber = match.groups()
                    current_voice_session.update({
                        'end_time': log_entry['timestamp'],
                        'duration': float(duration),
                        'packet_loss': int(packet_loss),
                        'ber': float(ber)
                    })
                    analysis['voice_sessions'].append(current_voice_session)
                current_voice_session = None
            
            elif 'Started' in message or 'Opening' in message:
                analysis['network_registrations'] += 1
            
            elif log_entry['level'] == 'E':  # Error
                analysis['errors'].append({
                    'timestamp': log_entry['timestamp'],
                    'message': message
                })
        
        return analysis
    
    def generate_report(self):
        """Generate comprehensive correlation report"""
        print("Generating correlation report...")
        
        analysis = self.analyze_patterns()
        
        with open(self.output_file, 'w') as f:
            f.write("=" * 80 + "\n")
            f.write("DMR PACKET CAPTURE AND LOG CORRELATION ANALYSIS\n")
            f.write("=" * 80 + "\n")
            f.write(f"Generated: {datetime.datetime.now()}\n")
            f.write(f"PCAP File: {self.pcap_file}\n")
            f.write(f"Log File: {self.log_file}\n")
            f.write("=" * 80 + "\n\n")
            
            # Executive Summary
            f.write("EXECUTIVE SUMMARY\n")
            f.write("-" * 40 + "\n")
            f.write(f"Total Packets Captured: {analysis['total_packets']}\n")
            f.write(f"Total Log Entries: {analysis['total_log_entries']}\n")
            f.write(f"Session Duration: {analysis['timespan']}\n")
            f.write(f"Average Packet Rate: {analysis['average_packet_rate']:.2f} packets/second\n")
            f.write(f"Total Data Volume: {analysis['data_volume']} bytes ({analysis['data_volume']/1024:.1f} KB)\n")
            f.write(f"Voice Transmissions: {analysis['voice_transmissions']}\n")
            f.write(f"Network Events: {analysis['network_registrations']}\n")
            f.write(f"Errors Detected: {len(analysis['errors'])}\n\n")
            
            # Voice Session Analysis
            f.write("VOICE SESSION ANALYSIS\n")
            f.write("-" * 40 + "\n")
            if analysis['voice_sessions']:
                f.write(f"{'Callsign':<12} {'TG':<6} {'Duration':<10} {'Loss%':<8} {'BER%':<8} {'Start Time'}\n")
                f.write("-" * 70 + "\n")
                for session in analysis['voice_sessions']:
                    f.write(f"{session['callsign']:<12} {session['talkgroup']:<6} "
                           f"{session['duration']:<10.1f} {session['packet_loss']:<8} "
                           f"{session['ber']:<8.1f} {session['start_time'].strftime('%H:%M:%S')}\n")
            else:
                f.write("No voice sessions detected in log timeframe.\n")
            f.write("\n")
            
            # Error Analysis
            if analysis['errors']:
                f.write("ERROR ANALYSIS\n")
                f.write("-" * 40 + "\n")
                for error in analysis['errors']:
                    f.write(f"{error['timestamp'].strftime('%H:%M:%S')} - {error['message']}\n")
                f.write("\n")
            
            # Detailed Timeline Correlation
            f.write("DETAILED TIMELINE CORRELATION\n")
            f.write("-" * 40 + "\n")
            f.write("Format: [TIME] LOG_LEVEL: Message | Packets: COUNT\n\n")
            
            for correlation in sorted(self.correlated_events, key=lambda x: x['timestamp']):
                log = correlation['log_entry']
                packet_count = correlation['packet_count']
                
                f.write(f"[{log['timestamp'].strftime('%H:%M:%S.%f')[:-3]}] "
                       f"{log['level']}: {log['message']}")
                
                if packet_count > 0:
                    f.write(f" | Packets: {packet_count}")
                    
                    # Show packet details for significant events
                    if any(keyword in log['message'].lower() for keyword in 
                          ['voice header', 'end of voice', 'started', 'opening']):
                        f.write("\n    Packet Activity:")
                        for packet in correlation['related_packets'][:5]:  # Show first 5
                            f.write(f"\n      {packet['timestamp'].strftime('%H:%M:%S.%f')[:-3]} "
                                   f"{packet['direction']} {packet['length']}B")
                        if len(correlation['related_packets']) > 5:
                            f.write(f"\n      ... and {len(correlation['related_packets'])-5} more packets")
                
                f.write("\n")
            
            # Packet Flow Statistics
            f.write("\nPACKET FLOW STATISTICS\n")
            f.write("-" * 40 + "\n")
            to_mmdvm = sum(1 for p in self.packets if p['direction'] == 'TO_MMDVM')
            from_mmdvm = sum(1 for p in self.packets if p['direction'] == 'FROM_MMDVM')
            f.write(f"Packets TO MMDVM (from network): {to_mmdvm}\n")
            f.write(f"Packets FROM MMDVM (to network): {from_mmdvm}\n")
            f.write(f"Traffic Ratio (Network->MMDVM): {to_mmdvm/max(from_mmdvm,1):.2f}:1\n\n")
            
            # Raw Data Sections
            f.write("RAW PACKET SAMPLE (First 10)\n")
            f.write("-" * 40 + "\n")
            for i, packet in enumerate(self.packets[:10]):
                f.write(f"{packet['timestamp'].strftime('%H:%M:%S.%f')[:-3]} "
                       f"{packet['direction']} {packet['length']}B\n")
            f.write("\n")
        
        print(f"Analysis complete! Report saved to: {self.output_file}")

def main():
    if len(sys.argv) != 4:
        print("Usage: python3 dmr_correlation_analysis.py <pcap_file> <log_file> <output_file>")
        sys.exit(1)
    
    pcap_file, log_file, output_file = sys.argv[1:4]
    
    if not os.path.exists(pcap_file):
        print(f"Error: PCAP file not found: {pcap_file}")
        sys.exit(1)
    
    if not os.path.exists(log_file):
        print(f"Error: Log file not found: {log_file}")
        sys.exit(1)
    
    analyzer = DMRAnalyzer(pcap_file, log_file, output_file)
    analyzer.parse_pcap()
    analyzer.parse_mmdvm_log()
    analyzer.correlate_events()
    analyzer.generate_report()

if __name__ == "__main__":
    main()
