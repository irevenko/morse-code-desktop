package morse

import (
	"io"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/oto"
)

const (
	frequency       = 800.0
	sampleRate      = 20000
	channelNum      = 2
	bitDepthInBytes = 2
)

type SineWave struct {
	freq   float64
	length int64
	pos    int64

	remaining []byte
}

func NewSineWave(freq float64, duration time.Duration) *SineWave {
	l := int64(channelNum) * int64(bitDepthInBytes) * int64(sampleRate) * int64(duration) / int64(time.Second)
	l = l / 4 * 4
	return &SineWave{
		freq:   freq,
		length: l,
	}
}

func (s *SineWave) Read(buf []byte) (int, error) {
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		s.remaining = s.remaining[n:]
		return n, nil
	}

	if s.pos == s.length {
		return 0, io.EOF
	}

	eof := false
	if s.pos+int64(len(buf)) > s.length {
		buf = buf[:s.length-s.pos]
		eof = true
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	length := float64(sampleRate) / float64(s.freq)

	num := (bitDepthInBytes) * (channelNum)
	p := s.pos / int64(num)
	switch bitDepthInBytes {
	case 1:
		for i := 0; i < len(buf)/num; i++ {
			const max = 127
			b := int(math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)
			for ch := 0; ch < channelNum; ch++ {
				buf[num*i+ch] = byte(b + 128)
			}
			p++
		}
	case 2:
		for i := 0; i < len(buf)/num; i++ {
			const max = 32767
			b := int16(math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)
			for ch := 0; ch < channelNum; ch++ {
				buf[num*i+2*ch] = byte(b)
				buf[num*i+1+2*ch] = byte(b >> 8)
			}
			p++
		}
	}

	s.pos += int64(len(buf))

	n := len(buf)
	if origBuf != nil {
		n = copy(origBuf, buf)
		s.remaining = buf[n:]
	}

	if eof {
		return n, io.EOF
	}
	return n, nil
}

func Play(p *oto.Player, freq float64, duration time.Duration) error {
	s := NewSineWave(freq, duration)
	if _, err := io.Copy(p, s); err != nil {
		return err
	}

	return nil
}

func RunShort(c *oto.Player) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := Play(c, frequency, 150*time.Millisecond); err != nil {
			panic(err)
		}
	}()

	wg.Wait()
	return nil
}

func RunLong(c *oto.Player) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := Play(c, frequency, 300*time.Millisecond); err != nil {
			panic(err)
		}
	}()

	wg.Wait()
	return nil
}

func SleepShort() {
	time.Sleep(100 * time.Millisecond)
}

func SleepLong() {
	time.Sleep(200 * time.Millisecond)
}

func MorseToSound(morseSequence string, c *oto.Context, p *oto.Player) {
	morseSlice := strings.Split(morseSequence, "")

	for _, v := range morseSlice {
		if v == "." {
			RunShort(p)
			SleepShort()
		}

		if v == "-" {
			RunLong(p)
			SleepShort()
		}

		if v == " " {
			SleepLong()
		}

		if v == "/" {
			SleepLong()
		}
	}
}

