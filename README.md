rog
===
![Rog Screenshot](http://github.com/ajhager/rog/raw/master/data/screenshot.png)

* 24bit color, unicode console with built in font.
* Cross platform windowing via github.com/skelterjohn/go.wde.
* Field of view and pathfinding algorithms.

Notes
-----
* On Windows you can build your project with `go build -ldflags -Hwindowsgui` to inhibit the console window that pops up by default.

Thanks
------
* libTCOD [http://doryen.eptalys.net/libtcod/]
* John Asmuth [http://github.com/skelterjohn/go.wde]
* Yamamushi [http://www.theasciiproject.com/]

Plans
-----
* Documentation and tutorials.
* Pathfinding.
* World generation.
* User supplied font sets and tilemaps.
* Noise generators.
* Image scale and rotation.
* Console to console blitting.
* Fold lighting into the library.
* Test suite.
