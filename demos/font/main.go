package main

import (
	"github.com/ajhager/rog"
)

func main() {
	font := rog.Font("terminal1-2.png", 8, 16, "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ ✵웃世界¢¥¤§©¨«¬£ª±²³´¶·¸¹º»¼½¾¿☐☑═║╔╗╚╝╠╣╦╩╬░▒▓☺☻☼♀♂▀▁▂▃▄▅▆▇█ÐÑÒÓÔÕÖÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏØÙÚÛÜÝàáâãäåèéêëìíîïðñòóôõö÷ùúûüýÿ♥♦♣♠♪♬æçø←↑→↓↔↕®‼ꀥ")
	rog.Open(80, 20, 2, false, "font", font)
	for i := 0; i < 20; i++ {
		rog.Set(30, i, rog.Black, rog.White, "12345abcdeABCDE!\"#$%")
	}
	for rog.Running() {
		if rog.Key() == rog.Esc {
			rog.Close()
		}
		rog.Flush()
	}
}
