rog
===
![Rog Screenshot](http://github.com/ajhager/rog/raw/master/data/screenshot.png)

* Depends only on skelterjohn/go.wde.
* 32bit, unicode console with built in font.
* Supports opening multiple windows, each backed by their own console and input.
* Takes inspiration from libtcod (the de facto roguelike library.)

Plans
-----
* Background color flags.
* Keyboard and mouse input routines.
* Console to console blitting.
* String printing.
* Field of view and lighting algorithms.
* Image package for scaling up pixel art + window scaling.
* Saving screenshots and videos.
* User supplied font sets.
* Documentation and tutorials.
* Move console to a separate package so other projects can use it.

Notes
-----
* On Windows you can build your project with `go build -ldflags -Hwindowsgui` to inhibit the console window that pops up by default.

Thanks
------
* Yamamushi [http://www.theasciiproject.com/]
* libTCOD [http://doryen.eptalys.net/libtcod/]