rog
===
![Rog Screenshot](http://github.com/ajhager/rog/raw/master/data/screenshot.png)

* Depends only on skelterjohn/go.wde, so little to no dependencies.
* 32bit, unicode console with built in font has been implemented.
* Supports opening multiple windows, each backed by their own console and input.
* For the moment, it takes inspiration from libtcod (the de facto roguelike library.)

Plans
-----
* String and shape drawing.
* Keyboard and mouse input routines.
* Console to console blitting.
* Field of view and lighting algorithms.

Notes
-----
* On Windows you can build your project with `go build -ldflags -Hwindowsgui` to inhibit the console window that pops up by default.

Thanks
------
* Yamamushi [http://www.theasciiproject.com/]
* libTCOD [http://doryen.eptalys.net/libtcod/]