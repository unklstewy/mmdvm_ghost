#!/usr/bin/env python3
import sys, re

def parse_hex_blocks(lines):
    """Parse tcpdump hex output split by '--' and extract UDP payloads as byte arrays."""
    blocks = []
    current = []
    for line in lines:
        if line.strip() == '--':
            if current:
                blocks.append(current)
                current = []
        else:
            current.append(line)
    if current:
        blocks.append(current)
    print(f"DEBUG: Found {len(blocks)} blocks", file=sys.stderr)
    payloads = []
    for idx, block in enumerate(blocks):
        hexlines = [l for l in block if re.match(r'^\s*0x[0-9a-f]+:', l)]
        hexstr = ''.join(''.join(l.split(':',1)[1].split()) for l in hexlines)
        udp_payload = hexstr[-110:] # Take the last 110 hex chars (55 bytes)
        print(f"DEBUG: Block {idx+1} hexstr length: {len(hexstr)}", file=sys.stderr)
        print(f"DEBUG: Block {idx+1} udp_payload length: {len(udp_payload)}", file=sys.stderr)
        if len(udp_payload) == 110:
            payloads.append(bytes.fromhex(udp_payload))
    print(f"DEBUG: Extracted {len(payloads)} payloads", file=sys.stderr)
    return payloads

if __name__ == "__main__":
    with open(sys.argv[1], 'r') as f:
        lines = f.readlines()
    payloads = parse_hex_blocks(lines)
    for i, p in enumerate(payloads):
        arr = ', '.join(f'0x{b:02x}' for b in p)
        print(f"// Frame {i+1}\n[]byte{{{arr}}},")
