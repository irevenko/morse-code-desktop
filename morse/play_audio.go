package morse

import (
	"flag"
	"github.com/hajimehoshi/oto"
	"io"
	"math"
	"sync"
	"time"
)

//func read() (decoder *minimp3.Decoder, audioData []byte) {
//	file, err := ioutil.ReadFile("beep-07.mp3")
//	if err != nil {
//		log.Fatal(file)
//	}
//
//	var dec *minimp3.Decoder
//	var data []byte
//	dec, data, _ = minimp3.DecodeFull(file)
//
//	return dec, data
//}
//
//func play(data []byte, dec *minimp3.Decoder, p *oto.Player) {
//	p.Write(data)
//
//	<-time.After(time.Second)
//
//	dec.Close()
//}

var (
	sampleRate      = flag.Int("samplerate", 44100, "sample rate")
	channelNum      = flag.Int("channelnum", 2, "number of channel")
	bitDepthInBytes = flag.Int("bitdepthinbytes", 2, "bit depth in bytes")
)

type SineWave struct {
	freq   float64
	length int64
	pos    int64

	remaining []byte
}

func NewSineWave(freq float64, duration time.Duration) *SineWave {
	l := int64(*channelNum) * int64(*bitDepthInBytes) * int64(*sampleRate) * int64(duration) / int64(time.Second)
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

	length := float64(*sampleRate) / float64(s.freq)

	num := (*bitDepthInBytes) * (*channelNum)
	p := s.pos / int64(num)
	switch *bitDepthInBytes {
	case 1:
		for i := 0; i < len(buf)/num; i++ {
			const max = 127
			b := int(math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)
			for ch := 0; ch < *channelNum; ch++ {
				buf[num*i+ch] = byte(b + 128)
			}
			p++
		}
	case 2:
		for i := 0; i < len(buf)/num; i++ {
			const max = 32767
			b := int16(math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)
			for ch := 0; ch < *channelNum; ch++ {
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

func play(context *oto.Context, freq float64, duration time.Duration) error {
	p := context.NewPlayer()
	s := NewSineWave(freq, duration)
	if _, err := io.Copy(p, s); err != nil {
		return err
	}
	if err := p.Close(); err != nil {
		return err
	}
	return nil
}

func run() error {
	const (
		freqC = 523.3
	)

	c, err := oto.NewContext(*sampleRate, *channelNum, *bitDepthInBytes, 4096)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := play(c, freqC, 3*time.Second); err != nil {
			panic(err)
		}
	}()


	wg.Wait()
	c.Close()
	return nil
}
func playBeep() {
	//decoder, data := read()
	//
	//var context *oto.Context
	//context, _= oto.NewContext(decoder.SampleRate, decoder.Channels, 2, 1024)
	//player := context.NewPlayer()
	//
	//
	//play(data, decoder, player)
	//play(data, decoder, player)
}
