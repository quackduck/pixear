package main //schlobbin, shlobbers

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"
	"time"

	"github.com/crazy3lf/colorconv"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"

	"github.com/google/hilbert" //should i also add my github link??yeah??
	"golang.org/x/image/draw"
)

var s *hilbert.Hilbert

const size = 16
const sampleRate = 44100

var dst *image.RGBA

func main() {
	src, err := getImageFromFilePath("cat2.jpg")
	// src, err := getImageFromFilePath("/Users/ishan/code/clones/fracmd/src/ok-lsd-fractal.png")
	// src, err := getImageFromFilePath("gradient.png")

	if err != nil {
		panic(err) // panic mode yeah?? pee is stored in the balls.. oi? pee?
	}
	// maxP := src.Bounds().Max
	// maxX := maxP.X
	// maxY := maxP.Y

	dst = image.NewRGBA(image.Rect(0, 0, size, size))

	draw.CatmullRom.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	f, _ := os.Create("scaled.png")
	png.Encode(f, dst)
	// Resize:
	// draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// if maxX != maxY {
	// 	panic("WHAT THE FUCK!!!") // real panic right here, oi
	// }
	s, err = hilbert.NewHilbert(size)
	if err != nil {
		panic(err)
	}

	sr := beep.SampleRate(sampleRate) // 44100
	speaker.Init(sr, sr.N(time.Second/10))
	println("now playing")
	speaker.Play(audio2())
	select {}
}

func audio() beep.Streamer {
	var currI = 0
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		// r := rand.Float64()
		// println(len(samples))
		// if len(samples) != 512 {
		// 	return len(samples), true
		// }
		// for x := 0; x < size; x++ {
		// 	for y := 0; y < size; y++ {
		// 		s.MapInverse(x, y)
		// 	}
		// }
		for i := range samples {
			currI++
			mappedI := currI % (size * size)

			if mappedI == 0 {
				time.Sleep(time.Second * 1)
			}

			x, y, _ := s.Map(mappedI)

			// x := mappedI / size
			// y := mappedI % size

			c := dst.At(x, y)

			// r, g, b, _ := c.RGBA()
			// samples[i][0] = 2*(float64(r)/0xffff) - 1
			// samples[i][1] = 2*(float64(g)/0xffff) - 1

			h, s, _ := colorconv.ColorToHSV(c)
			h = h/360.0*2 - 1
			s = s*2 - 1
			// fmt.Println(h, s, v)

			samples[i][0] = math.Sin(10 * math.Cos(math.Tan(h))) //+ //math.Sin(float64(400*i)/float64(len(samples)))
			samples[i][1] = math.Sin(10 * math.Cos(math.Tan(s))) //+ math.Sin(float64(400*i)/float64(len(samples)))

			// if i%3 == 0 {
			// 	samples[i][0] = 4*r + math.Sin(float64(400*i)/float64(len(samples)))
			// 	samples[i][1] = 4*r + math.Tan(200*float64(i)/float64(len(samples))-math.Pi/2)
			// } else if i%3 == 1 {
			// 	samples[i][0] = 4*r + math.Sin(float64(4000*i)/float64(len(samples)))
			// 	samples[i][1] = 4*r + math.Tan(20*float64(i)/float64(len(samples))-math.Pi/2)
			// } else {
			// 	samples[i][0] = 4*r + math.Sin(float64(40*i)/float64(len(samples)))
			// 	samples[i][1] = 4*r + math.Tan(2000*float64(i)/float64(len(samples))-math.Pi/2)
			// }
		}
		return len(samples), true
	})
}

func audio2() beep.Streamer {
	var currI2 = 0
	var pixIndex = 0
	freq := 0.0
	loudness := 0.0
	correction := 0.0 // how much do we need to correct this run's phase by?

	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		// println(len(samples))
		for i := range samples {
			// if i == 0 {
			// 	time.Sleep(5 * time.Second / 1000)
			// }
			currI2 += 1
			if currI2%10000 == 0 {
				pixIndex++
				mappedI := pixIndex % (size * size)
				if mappedI == 0 {
					time.Sleep(time.Second * 23)
				}
				x, y, _ := s.Map(mappedI)
				// x := mappedI / size
				// y := mappedI % size
				c := dst.At(x, y)
				h, _, v := colorconv.ColorToHSV(c)
				h = h / 360.0 //*2 - 1
				// v = v*2 - 1
				oldFreq := freq
				freq = h * 1000
				loudness = v
				oldCorrection := correction
				// see https://www.desmos.com/calculator/2zsziquhoy
				correction = (oldFreq*float64(currI2) + oldCorrection) - freq*float64(currI2)

				loudness = 1

				//freq = (rand.Float64()) * 1000
				//loudness = rand.Float64()

				fmt.Println(freq, loudness)
			}
			//samples[i][0] = loudness * math.Sin((freq*float64(currI2)+correction)*2*math.Pi/float64(sampleRate))
			// samples[i][0] = loudness * math.Pow(math.Tan(math.Pow(math.Sin((freq*float64(currI2)+correction)*2*math.Pi/float64(sampleRate)), 20)), 1111)
			//samples[i][0] = loudness * math.Tan(math.Pow(math.Sin((freq*float64(currI2)+correction)*2*math.Pi/float64(sampleRate)), 20))
			samples[i][0] = loudness * waveform((freq*float64(currI2)+correction)*2*math.Pi/float64(sampleRate))
			samples[i][1] = samples[i][0]
		}
		// fmt.Println("tick")
		return len(samples), true
	})
}

func waveform(n float64) float64 {
	//return math.Sin(n) + 1/3*math.Sin(3*n) + 1/5*math.Sin(5*n)
	return math.Sin(n) + 1/3*math.Sin(3*n+0.8) + 1/5*math.Sin(5*n+0.8)
}

func audio3() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		val := 0.0
		for i := range samples {
			if i%(sampleRate/440/2) == 0 {
				val = 1 - val
			}
			samples[i][0] = val
			samples[i][1] = val
		}
		return len(samples), true
	})
}

// func pixToSample(r, g, b float64) (float64, float64) {
// 	volume =

// 	return volume, 0
// }

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

// oi yeah aye bruv, you might wanna add some comments so i understand this yea... bruv??
// Im so lost
// I'm still holding on to everything that's dead and gone
// I don't wanna say good bye cuz this one means foreveeerrrrr
// Now you're in the stars and 6 fts never felt so far
// Here I am alone between the heavens and the embers
