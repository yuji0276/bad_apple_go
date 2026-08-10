package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"golang.org/x/term"
)

const (
	ramp = " .:-=+*#%@"
	fps  = 30
	//文字セルの 幅/高さ。フォント依存なので目視で合わせる
	charAspect = 0.5
)

var errInterrupted = errors.New("interrupted")

func run() error {
	screenWeight, screenHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return err
	}
	buf := make([]byte, screenHeight*screenWeight)
	//縦を charAspect 倍に潰してから、比を保ったまま端末に収め、余りを黒で埋める。
	//pad があるので出力は必ず screenWeight*screenHeight バイト/フレームになる
	vf := fmt.Sprintf(
		"fps=%d,scale=iw:ih*%g,scale=%d:%d:force_original_aspect_ratio=decrease:flags=area,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,format=gray",
		fps, charAspect, screenWeight, screenHeight, screenWeight, screenHeight)

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
		return err
	}
	if err := sound.Start(); err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	screen := make([]byte, 0, screenHeight*screenWeight+screenHeight+2)
	count := 0
	//1フレームの長さ
	frameDur := time.Second / fps
	//カーソルを非表示
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")
	//基準時刻。
	var startTime time.Time
	//繰り返し1回で1フレーム描画
	for {
		select {
		case <-ch:
			return errInterrupted
		default:
		}
		//読み込み
		_, err = io.ReadFull(stdout, buf)
		//リセット
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		screen = screen[:0]
		if startTime.IsZero() {
			startTime = time.Now()
		}
		//count枚目を描くべき時刻。startTime と count だけで決まるので誤差が積もらない
		target := startTime.Add(time.Duration(count) * frameDur)
		//目標時刻まで待つ。遅れていれば負の値になり、continue。
		if time.Until(target) < 0 {
			count++
			continue
		} else {
			time.Sleep(time.Until(target))
		}
		//画面に描画
		screen = append(screen, "\033[H"...)
		for i := range screenHeight * screenWeight {
			idx := int(buf[i]) * (len(ramp) - 1) / 255
			screen = append(screen, ramp[idx])
			if i == screenHeight*screenWeight-1 {
				break
			}
			if (i+1)%screenWeight == 0 {
				screen = append(screen, '\n')
			}
		}
		count++
		//書き込み
		os.Stdout.Write(screen)
	}
	if err = cmd.Wait(); err != nil {
		return err
	}
	if err = sound.Wait(); err != nil {
		return err
	}

	return nil
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
