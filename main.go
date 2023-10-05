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
				//h = h / 360.0 //*2 - 1
				// v = v*2 - 1
				oldFreq := freq
				//freq = h * 1000
				freq = hueToFreq(h)
				loudness = v
				oldCorrection := correction
				// see https://www.desmos.com/calculator/2zsziquhoy
				correction = (oldFreq*float64(currI2) + oldCorrection) - freq*float64(currI2)

				loudness = 1 //- math.Sin(float64(currI2)/100000.0)/2

				//freq = (rand.Float64()) * 1000
				//loudness = rand.Float64()

				fmt.Println(freq, loudness)
			}
			//samples[i][0] = loudness * math.Sin((freq*float64(currI2)+correction)*2*math.Pi/float64(sampleRate))
			//samples[i][0] = loudness * math.Pow(math.Tan(math.Pow(math.Sin((freq*float64(currI2)+correction)*2*math.Pi/float64(sampleRate)), 20)), 1111)
			//samples[i][0] = loudness * math.Tan(math.Pow(math.Sin((freq*float64(currI2)+correction)*2*math.Pi/float64(sampleRate)), 20))
			samples[i][0] = loudness * waveform((freq*float64(currI2)+correction)*2*math.Pi/float64(sampleRate))
			samples[i][1] = samples[i][0]
		}
		// fmt.Println("tick")
		return len(samples), true
	})
}

func hueToFreq(hue float64) float64 {
	//return 345 // sax
	return 527 // flutelo
	//return 0 + hue/360.0*500
	//return 0 + hue/360.0*100
	//
	//return 0 + hue/360.0*1000
	//return 100*math.Pow(2, hue/96.6) + 10*math.Sin(10*hue)
}

//func fft() {
//	period := []float64{0.50588235, 0.47843137, 0.47058824, 0.47843137, 0.49411765, 0.48235294,
//		0.47843137, 0.5372549, 0.63137255, 0.70588235, 0.78039216, 0.79215686,
//		0.61568627, 0.4627451, 0.49803922, 0.54117647, 0.4627451, 0.34509804,
//		0.33333333, 0.25882353, 0.12941176, 0.25490196, 0.41176471, 0.41176471,
//		0.38431373, 0.43529412, 0.50588235, 0.54901961, 0.57647059, 0.61960784,
//		0.64313725, 0.61176471}
//}

func waveform(n float64) float64 {
	//return math.Sin(n)
	result := 0.0
	//coeffs := []float64{0.1195814, 0.13768428, 0.02427883, 0.02013908, 0.05850243, 0.01997491, 0.00365248, 0.0318863, 0.0271606, 0.01164827}

	coeffs := []float64{0.1195814, 0.13768428, 0.02427883, 0.02013908, 0.05850243, 0.01997491, 0.00365248, 0.0318863, 0.0271606, 0.01164827, 0.0067154, 0.00352217, 0.00559091, 0.0007892, 0.00174847}

	// flutelo:
	//coeffs := []float64{0.10279466, 0.19061701} // 0.10375996, 0.06003058, 0.09161703, 0.05117856, 0.01291837, 0.04065544, 0.02173721, 0.01358232}

	// 2d slice of coeffs
	//coeffs := [][]float64{
	//	{0.5395480225988699, 0.0},
	//	{-0.08679545696826813, -0.056581890269043915},
	//	{-0.16162712613460659, 0.05106668394299141},
	//	{0.07535705039284196, -0.006755727435837805},
	//	{0.005130120991161465, -0.03878346746182446},
	//	{-0.031999621772155394, 0.014266652177372679},
	//	{-0.011611304172692427, -0.03565540953535884},
	//	{0.010593220338983038, 0.014154168221658553},
	//	{0.021572893075294454, 0.01036308531392525},
	//	{0.006494467508912326, 0.0005524136231722918},
	//	{0.004586091636206997, 0.0046956525508275},
	//	{-0.0017478542969456608, 0.000662547841690786},
	//	{-0.0007826946355727912, -0.0012972670484564413},
	//	{-0.0002681567516006239, 0.00047783174934292366},
	//	{-0.0005044390637605927, -0.002621142263269094},
	//	{0.0007756007423133226, -0.00010119136180056616},
	//}
	//for i := range coeffs {
	//	result += coeffs[i][0]*math.Cos(float64(i+1)*n) + coeffs[i][1]*math.Sin(float64(i+1)*n)
	//}

	for i, coeff := range coeffs {
		result += coeff * math.Sin(float64(i+1)*n)
	}
	return result
	//return math.Sin(n)
	//return 0.3*math.Sin(n) + 0.2*math.Sin(2*n) + 0.05*math.Sin(3*n) + 0.15*math.Sin(4*n) + 0.07*math.Sin(5*n)
	//
	//result := 0.0
	//for i := 1.0; i < 1000; i++ {
	//	//log2 := math.Log2(i)
	//	//if float64(int(log2)) != log2 {
	//	//	continue
	//	//}
	//	result += math.Sin(i*n) / (i * i) // + math.Cos(i*n)/(i*i)
	//}
	//return result
	//return math.Sin(n) + math.Sin(2*n)/2 + math.Sin(4*n)/4 + math.Sin(8*n)/8 + math.Sin(16*n)/16 + math.Sin(32*n)/32 + math.Sin(64*n)/64

	n = math.Mod(n, 2*math.Pi)
	//pi := math.Pi
	//noise := func() float64 { return rand.Float64() / 50.0 }
	//return n
	//return math.Sin(n) + math.Sin(n/(40))
	//return noise/5 + math.Tan(n)
	//n = math.Mod(n, 200*math.Pi)
	//return n
	//square := func(n float64) float64 {
	//	return math.Sin(n+noise()) + math.Sin(3*n+noise())/(3) + math.Sin(5*n+noise())/(5) + math.Sin(7*n+noise())/(7)
	//}
	////
	//saw := func(n float64) float64 {
	//	return math.Sin(2*n+noise())/(2) + math.Sin(4*n+noise())/(4) + math.Sin(6*n+noise())/(6) + math.Sin(8*n+noise())/(8)
	//}
	////
	////return (math.Sin(2*n)/(n) + math.Sin(4*n)/(n) + math.Sin(6*n)/(n) + math.Sin(8*n)/(n)) / 4.0
	//

	//return (0*square(n) + saw(n)) / 4
	//
	//n = math.Mod(n, 2*math.Pi)
	return math.Sin(n*n*n) / (n)
	//return n/(2*math.Pi) + math.Sin(n*n*n)/(n*n*n)
	//return math.Sin(n*n*n) / (n * n * n) // + math.Pow(math.Sin(n), 11)
	//return math.Sin(.2 * math.Exp(1*math.Mod(n, 2*math.Pi)))
	//return math.Sin(n) / n
	//if n < 0 {
	//	n += 2 * math.Pi
	//}
	//if n < math.Pi/2 {
	//	return 4 * n / math.Pi
	//} else if n < math.Pi {
	//	return 4 - 4*n/math.Pi
	//} else if n < 3*math.Pi/2 {
	//	return -4 + 4*n/math.Pi
	//} else {
	//	return 4 - 4*n/math.Pi
	//}
	//return math.Pow(math.Sin(n), 11)
	//return math.Sin(2*n)/2 + math.Sin(4*n)/4 + math.Sin(6*n)/6 + math.Pow(math.Sin(n), 1001)
	//return math.Sin(n) + 1/3*math.Sin(3*n) + 1/5*math.Sin(5*n)
	//return math.Sin(n)
	//return math.Sin(n) + 1/3*math.Sin(3*n+0.8) + 1/5*math.Sin(5*n+0.8)
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
