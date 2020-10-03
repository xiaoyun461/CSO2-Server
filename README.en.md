## Counter-Strike Online 2 Server 

[![Build status](https://ci.appveyor.com/api/projects/status/a4pj1il9li5s08k5?svg=true)](https://ci.appveyor.com/project/KouKouChan/cso2-server)
[![](https://img.shields.io/badge/license-MIT-green)](./LICENSE)
[![](https://img.shields.io/badge/version-v0.3.10-blue)](https://github.com/KouKouChan/CSO2-Server/releases)

### 0x01 Description

Counter-Strike Online 2 Server

DataBase:SQLite

*It is my first Golang project , in order to practice myself.*

***This project is incomplete right now !***

***Based on [cso2-master-server](https://github.com/L-Leite/cso2-master-server) by l-leite.***

### 0x02 Plan

    1. Basic game play √
    2. Refactoring code ...

### 0x03 Play

    1. You should have a game client with korea version.
    2. Download a launcher from l-leite's github page.
    3. Download lastest game server file from  ( https://github.com/KouKouChan/CSO2-Server/releases )
    4. Start game server and use bat file to launch your game.
    5. Have fun!

attention!

- If you want to enable registration , you should modify server.conf file , and set EnableRegister to 1 , and you must set your email smtp server and your email code.then you can open localhost:1314 with your browser.

### 0x04 Build

    1. open the folder
    2. enter "go build" command to build
    3. run it

### 0x05 Build Env

    Go 1.14.2
    Use port:30001-TCP、30002-UDP

***If you want to set up a LAN or Internet Server, please open the port of firewall.***

### 0x06 Screenshots

![Image](./photos/main.png)

![Image](./photos/intro.png)

![Image](./photos/channel.png)

![Image](./photos/ingame.jpg)

![Image](./photos/result.jpg)
