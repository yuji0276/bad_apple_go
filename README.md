# asciiplay

*English | [日本語](README.ja.md)*

Plays "Bad Apple!!" in your terminal as ASCII art — video made of nothing but characters, with the audio playing alongside it.

Each pixel becomes a character: the brighter the pixel, the denser the character (`@`); the darker it is, the lighter the character (` ` or `.`).

---

## Requirements

| What you need | What it's for |
|---|---|
| **ffmpeg** (both `ffmpeg` and `ffplay`) | Converting the video and playing the audio |
| A terminal | Where it's drawn. Any terminal with a monospaced font will do |

### Installing ffmpeg

**macOS (Homebrew)**

```bash
brew install ffmpeg
```

**Ubuntu / Debian**

```bash
sudo apt install ffmpeg
```

**Windows**

Download it from the [official ffmpeg site](https://ffmpeg.org/download.html), or run `winget install ffmpeg` in PowerShell.

To check that it installed correctly, run the commands below. If they print version information, you're set.

```bash
ffmpeg -version
ffplay -version
```

---

## Playing it

1. Move into the project folder.

   ```bash
   cd bad_apple
   ```

2. Play it.

   ```bash
   ./asciiplay
   ```

   If `asciiplay` doesn't run, or you'd like to rebuild it for your own environment, install [Go](https://go.dev/dl/) and run:

   ```bash
   go run .
   ```

3. Wait for it to finish, or press **Ctrl + C** to stop early.

---

## Tips for the best picture

- **Maximize your terminal before you start.** The output size is fixed at the moment the program launches — the bigger the window, the more detail you get.
- **Shrink the font size.** Smaller characters mean higher resolution and a picture that's easier to read. Use your terminal's settings, or `Ctrl + -` (`⌘ + -` on macOS).
- **A dark-background terminal works best.** Bright areas are drawn with dense characters, so white-on-black looks closest to the original video.
- **Don't resize the window while it's playing.** Resizing mid-playback garbles the output. If that happens, stop with Ctrl + C and run it again.

---

## Troubleshooting

**`ffmpeg: no such file or directory`**

ffmpeg either isn't installed or isn't on your `PATH`. See "Installing ffmpeg" above.

**`bad_apple.mp4: No such file or directory`**

The program has to run from the `bad_apple` folder, where `bad_apple.mp4` lives. Run `cd bad_apple` first.

**`./asciiplay` says `permission denied`**

The file isn't marked executable. Add the permission with:

```bash
chmod +x asciiplay
```

**The video plays but there's no sound**

`ffplay` may not be installed. Check with `ffplay -version` — on some systems it's packaged separately from ffmpeg itself.

**The picture looks stretched vertically or horizontally**

Character cells have different width-to-height ratios depending on the font. Switching fonts or adjusting the window's proportions helps.

**The terminal looks broken after Ctrl + C**

Restore it with:

```bash
reset
```

---

## Playing a different video

The filename is hardcoded to `bad_apple.mp4`. To play something else, rename your video to `bad_apple.mp4` and put it in the same folder (it's a good idea to rename the original file rather than overwrite it).

High-contrast, near-black-and-white footage looks best.

---

## About this project

Built as an exercise in learning Go. Rather than having generative AI write the code, I limited it to asking questions about syntax and how things work when I got stuck — the implementation is my own.
