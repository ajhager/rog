rog
===
![Rog Screenshot](http://github.com/ajhager/rog/raw/master/data/screenshot2.png)

* Depends only on skelterjohn/go.wde.
* 32bit, unicode console with built in font.
* Supports opening multiple windows, each backed by their own console and input.
* Takes inspiration from libtcod (the de facto roguelike library.)

Plans
-----
* Documentation and tutorials.
* Image scale and rotation.
* Console to console blitting.
* Saving screenshots.
* User supplied font sets and tilemaps.

Notes
-----
* On Windows you can build your project with `go build -ldflags -Hwindowsgui` to inhibit the console window that pops up by default.

Thanks
------
* Yamamushi [http://www.theasciiproject.com/]
* libTCOD [http://doryen.eptalys.net/libtcod/]
* John Asmuth [http://github.com/skelterjohn/go.wde]
