#!/usr/bin/env bash

# <bitbar.title>UTC</bitbar.title>
# <bitbar.version>1.0</bitbar.version>
# <bitbar.author>me</bitbar.author>
# <bitbar.author.github>flyingtimes</bitbar.author.github>
# <bitbar.desc>use baidu stock api to monitor stock price. the price only show up in market time.</bitbar.desc>

export PATH=/usr/local/bin:$PATH

# echo `TZ=UTC gdate +'''%Y-%m-%d %H:%M %p'`
utc24Hour=`TZ=UTC gdate +'''%H:%M %p'`
utc12Hour=`TZ=UTC gdate +'''%l:%M %p'`
echo -e "\033[32mUTC\033[0m:$utc24Hour"
echo "---"
echo -e "\033[32mUTC 24 Hour\033[0m: $utc24Hour"
echo -e "\033[32mUTC 12 Hour\033[0m: $utc12Hour"
echo "Refresh Me| terminal=false refresh=true"

# echo $PATH
# echo hi
# echo bye
