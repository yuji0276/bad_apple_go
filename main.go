package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"golang.org/x/term"
)

const (
	ramp = " .:-=+*#%@"
	fps  = 30
)

func main() {
	fd := os.Stdout.Fd()
	W, H, err := term.GetSize(int(fd))
	if err != nil {
		log.Fatal()
	}
	buf := make([]byte, H*W)
	vf := fmt.Sprintf("fps=%d,scale=%d:%d:flags=area,format=gray", fps, W, H)

	cmd := exec.Command("ffmpeg",
		"-hide_banner", "-loglevel", "error",
		"-i", "bad_apple.mp4",
		"-vf", vf,
		"-f", "rawvideo",
		"-pix_fmt", "gray",
		"-",
	)
	sound := exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", "bad_apple.mp4")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := sound.Start(); err != nil {
		log.Fatal()
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	screen := make([]byte, 0, H*W+H+2)
	count := 0
	//1フレームの長さ
	frameDur := time.Second / fps
	//基準時刻。
	startTime := time.Now()
	//繰り返し1回で1フレーム描画
	for {
		//読み込み
		_, err = io.ReadFull(stdout, buf)
		//リセット
		screen = screen[:0]
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//count枚目を描くべき時刻。startTime と count だけで決まるので誤差が積もらない
		target := startTime.Add(time.Duration(count) * frameDur)
		//画面に描画
		screen = append(screen, "\033[2J\033[H"...)
		for i := range H * W {
			idx := int(buf[i]) * (len(ramp) - 1) / 255
			screen = append(screen, ramp[idx])
			if i == H*W-1 {
				break
			}
			if (i+1)%W == 0 {
				screen = append(screen, '\n')
			}
		}
		count++
		//書き込み
		os.Stdout.Write(screen)
		//目標時刻まで待つ。遅れていれば負の値になり、Sleep は即座に返る
		time.Sleep(time.Until(target))
	}
	if err = cmd.Wait(); err != nil {
		os.Exit(1)
	}
	if err = sound.Wait(); err != nil {
		os.Exit(1)
	}
}
