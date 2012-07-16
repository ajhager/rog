rog
===
![Rog Screenshot](http://github.com/ajhager/rog/raw/master/data/screenshot.png)

* Depends only on skelterjohn/go.wde.
* 32bit, unicode console with built in font.
* Supports opening multiple windows, each backed by their own console and input.
* Takes inspiration from libtcod (the de facto roguelike library.)

Plans
-----
* Move console to a separate package so other projects can use it.
* Optimizations to console rendering.
* Window scaling.
* User supplied font sets.
* Background color flags.
* String printing.
* Frame time and fps calculation.
* Keyboard and mouse input routines.
* Console to console blitting.
* Console fading.
* Image blit2x.
* Saving screenshots.
* Color constants.
* Field of view and lighting algorithms.
* Documentation and tutorials.

Notes
-----
* On Windows you can build your project with `go build -ldflags -Hwindowsgui` to inhibit the console window that pops up by default.

Thanks
------
* Yamamushi [http://www.theasciiproject.com/]
* libTCOD [http://doryen.eptalys.net/libtcod/]