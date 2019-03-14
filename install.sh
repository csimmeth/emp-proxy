#!/bin/bash

# Install binary to bin
mv emp-proxy /usr/local/bin/.

# Register service
useradd proxyservice -s /sbin/nologin -M
cp emp-proxy.service /lib/systemd/system/.
chmod 755 /lib/systemd/system/emp-proxy.service

# Add syslog configuration
cp emp-proxy.syslog.conf /etc/rsyslog.d/emp-proxy.conf
systemctl restart rsyslog

# Start
systemctl start emp-proxy
systemctl enable emp-proxy
