Valve server queries "A2S_INFO" for 7 Days to Die.
==================================================
https://developer.valvesoftware.com/wiki/Server_queries#A2S_INFO

Example output.
==============
Console
-------
```
Header: -1
Return Code: 73 
Protocol Version: 17
Server Name: xCloud Games Host
World Type: Random Gen
Short Game Name: 7DTD
Long Game Name: 7 Days To Die
MAX Players: 7
Online Players: 1
Server Type: Dedicated server
Operation System: Linux
Protected: Public
Game Version: 00.15.02
```
JSON
----
```
{
  "Header":-1,
  "ReturnCode":73,
  "Protocol":17,
  "ServerName":"xCloud Games Host",
  "World":"Random Gen",
  "DescShort":"7DTD",
  "DescLong":"7 Days To Die",
  "Players":1,
  "PlayersMAX":7,
  "Bots":0,
  "ServerType":"Dedicated server",
  "Environment":"Linux",
  "Visibility":"Public",
  "Version":"00.15.02"
}
```
