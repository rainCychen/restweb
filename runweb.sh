#!/bin/bash
nohup /usr/local/bin/web -port=80 > /var/log/web/web.INFO 2>&1 &