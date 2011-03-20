package fftw_test

import (
  "fftw"
  "gospec"
  . "gospec"
  "math"
)


func Alloc1dSpec(c gospec.Context) {
  d10   := fftw.Alloc1d(10)
  d100  := fftw.Alloc1d(100)
  d1000 := fftw.Alloc1d(1000)
  c.Specify("Allocates the appropriate memory for 1d arrays.", func() {
    c.Expect(len(d10), Equals, 10)
    c.Expect(len(d100), Equals, 100)
    c.Expect(len(d1000), Equals, 1000)
  })
}

func Alloc2dSpec(c gospec.Context) {
  d100x50 := fftw.Alloc2d(100, 50)
  c.Specify("Allocates the appropriate memory for 2d arrays.", func() {
    c.Expect(len(d100x50), Equals, 100)
    for _,v := range d100x50 {
      c.Expect(len(v), Equals, 50)
    }
    counter := 0.0
    for i := range d100x50 {
      for j := range d100x50[i] {
        d100x50[i][j] = complex(counter, 0)
        counter += 1.0
      }
    }
    counter = 0.0
    for i := range d100x50 {
      for j := range d100x50[i] {
        c.Expect(real(d100x50[i][j]), Equals, counter)
        counter += 1.0
      }
    }
  })
}

func Alloc3dSpec(c gospec.Context) {
  d100x20x10 := fftw.Alloc3d(100, 20, 10)
  c.Specify("Allocates the appropriate memory for 3d arrays.", func() {
    c.Expect(len(d100x20x10), Equals, 100)
    for _,v := range d100x20x10 {
      c.Expect(len(v), Equals, 20)
      for _,v := range v {
        c.Expect(len(v), Equals, 10)
      }
    }
    counter := 0.0
    for i := range d100x20x10 {
      for j := range d100x20x10[i] {
        for k := range d100x20x10[i][j] {
          d100x20x10[i][j][k] = complex(counter, 0)
          counter += 1.0
        }
      }
    }
    counter = 0.0
    for i := range d100x20x10 {
      for j := range d100x20x10[i] {
        for k := range d100x20x10[i][j] {
          c.Expect(real(d100x20x10[i][j][k]), Equals, counter)
          counter += 1.0
        }
      }
    }
  })
}

func FFT1dSpec(c gospec.Context) {
  signal := fftw.Alloc1d(16)
  for i := range signal {
    signal[i] = complex(float64(i), float64(-i))
  }
  forward := fftw.PlanDft1d(signal, signal, fftw.Forward, fftw.Estimate)
  c.Specify("Creating a plan doesn't overwrite an existing array if fftw.Estimate is used.", func() {
    for i := range signal {
      c.Expect(signal[i], Equals, complex(float64(i), float64(-i)))
    }
  })

  // A simple real cosine should result in transform with two spikes, one at S[1] and one at S[-1]
  // The spikes should be real and have amplitude equal to len(S)/2 (because fftw doesn't normalize)
  for i := range signal {
    signal[i] = complex(math.Cos(float64(i) / float64(len(signal)) * math.Pi * 2), 0)
  }
  forward.Execute()
  c.Specify("Forward 1d FFT works properly.", func() {
    c.Expect(real(signal[0]), IsWithin(1e-9), 0.0)
    c.Expect(imag(signal[0]), IsWithin(1e-9), 0.0)
    c.Expect(real(signal[1]), IsWithin(1e-9), float64(len(signal))/2)
    c.Expect(imag(signal[1]), IsWithin(1e-9), 0.0)
    for i := 2; i < len(signal) - 1; i++ {
      c.Expect(real(signal[i]), IsWithin(1e-9), 0.0)
      c.Expect(imag(signal[i]), IsWithin(1e-9), 0.0)
    }
    c.Expect(real(signal[len(signal)-1]), IsWithin(1e-9), float64(len(signal))/2)
    c.Expect(imag(signal[len(signal)-1]), IsWithin(1e-9), 0.0)
  })
}

func FFT2dSpec(c gospec.Context) {
  signal := fftw.Alloc2d(64, 8)
  for i := range signal {
    for j := range signal[i] {
      signal[i][j] = complex(float64(i+j), float64(-i-j))
    }
  }
  forward := fftw.PlanDft2d(signal, signal, fftw.Forward, fftw.Estimate)
  c.Specify("Creating a plan doesn't overwrite an existing array if fftw.Estimate is used.", func() {
    for i := range signal {
      for j := range signal[i] {
        c.Expect(signal[i][j], Equals, complex(float64(i+j), float64(-i-j)))
      }
    }
  })

  // As long as fx < dx/2 and fy < dy/2, where dx and dy are the lengths in each dimension,
  // there will be 2^n spikes, where n is the number of dimensions.  Each spike will be
  // real and have magnitude equal to dx*dy / 2^n
  dx := len(signal)
  fx := float64(dx) / 4
  dy := len(signal[0])
  fy := float64(dy) / 4
  for i := range signal {
    for j := range signal[i] {
      cosx := math.Cos(float64(i) / float64(dx) * fx * math.Pi * 2)
      cosy := math.Cos(float64(j) / float64(dy) * fy * math.Pi * 2)
      signal[i][j] = complex(cosx * cosy, 0)
    }
  }
  forward.Execute()
  c.Specify("Forward 2d FFT works properly.", func() {
    for i := range signal {
      for j := range signal[i] {
        if (i == int(fx) || i == dx - int(fx)) &&
           (j == int(fy) || j == dy - int(fy)) {
          c.Expect(real(signal[i][j]), IsWithin(1e-7), float64(dx*dy/4))
          c.Expect(imag(signal[i][j]), IsWithin(1e-7), 0.0)
        } else {
          c.Expect(real(signal[i][j]), IsWithin(1e-7), 0.0)
          c.Expect(imag(signal[i][j]), IsWithin(1e-7), 0.0)
        }
      }
    }
  })
}

func FFT3dSpec(c gospec.Context) {
  signal := fftw.Alloc3d(32, 16, 8)
  for i := range signal {
    for j := range signal[i] {
      for k := range signal[i][j] {
        signal[i][j][k] = complex(float64(i+j+k), float64(-i-j-k))
      }
    }
  }
  forward := fftw.PlanDft3d(signal, signal, fftw.Forward, fftw.Estimate)
  c.Specify("Creating a plan doesn't overwrite an existing array if fftw.Estimate is used.", func() {
    for i := range signal {
      for j := range signal[i] {
        for k := range signal[i][j] {
          c.Expect(signal[i][j][k], Equals, complex(float64(i+j+k), float64(-i-j-k)))
        }
      }
    }
  })

  // As long as fx < dx/2, fy < dy/2, and fz < dz/2, where dx,dy,dz  are the lengths in
  // each dimension, there will be 2^n spikes, where n is the number of dimensions.
  // Each spike will be real and have magnitude equal to dx*dy*dz / 2^n
  dx := len(signal)
  fx := float64(dx) / 4
  dy := len(signal[0])
  fy := float64(dy) / 4
  dz := len(signal[0][0])
  fz := float64(dz) / 4
  for i := range signal {
    for j := range signal[i] {
      for k := range signal[i][j] {
        cosx := math.Cos(float64(i) / float64(dx) * fx * math.Pi * 2)
        cosy := math.Cos(float64(j) / float64(dy) * fy * math.Pi * 2)
        cosz := math.Cos(float64(k) / float64(dz) * fz * math.Pi * 2)
        signal[i][j][k] = complex(cosx * cosy * cosz, 0)
      }
    }
  }
  forward.Execute()
  c.Specify("Forward 3d FFT works properly.", func() {
    for i := range signal {
      for j := range signal[i] {
        for k := range signal[i][j] {
          if (i == int(fx) || i == dx - int(fx)) &&
             (j == int(fy) || j == dy - int(fy)) &&
             (k == int(fz) || k == dz - int(fz)) {
            c.Expect(real(signal[i][j][k]), IsWithin(1e-7), float64(dx*dy*dz/8))
            c.Expect(imag(signal[i][j][k]), IsWithin(1e-7), 0.0)
          } else {
            c.Expect(real(signal[i][j][k]), IsWithin(1e-7), 0.0)
            c.Expect(imag(signal[i][j][k]), IsWithin(1e-7), 0.0)
          }
        }
      }
    }
  })
}

func FFTR2CSpec(c gospec.Context) {
  signal := make([]float64, 16)
  F_signal := make([]complex128, 9)

  for i := range signal {
    signal[i] = math.Sin(float64(i) / float64(len(signal)) * math.Pi * 2)
  }
  forward := fftw.PlanDftR2C1d(signal, F_signal, fftw.Estimate)
  forward.Execute()
  c.Specify("Running a R2C transform doesn't destroy the input.", func() {
    for i := range signal {
      c.Expect(signal[i], Equals, math.Sin(float64(i) / float64(len(signal)) * math.Pi * 2))
    }
  })

  c.Specify("Forward 1d Real to Complex FFT  works properly.", func() {
    c.Expect(real(F_signal[0]), IsWithin(1e-9), 0.0)
    c.Expect(imag(F_signal[0]), IsWithin(1e-9), 0.0)
    c.Expect(real(F_signal[1]), IsWithin(1e-9), 0.0)
    c.Expect(imag(F_signal[1]), IsWithin(1e-9), -float64(len(signal))/2)
    for i := 2; i < len(F_signal) - 1; i++ {
      c.Expect(real(F_signal[i]), IsWithin(1e-9), 0.0)
      c.Expect(imag(F_signal[i]), IsWithin(1e-9), 0.0)
    }
  })
}

func FFTC2RSpec(c gospec.Context) {
  signal := make([]float64, 16)
  F_signal := make([]complex128, 9)

  forward := fftw.PlanDftR2C1d(signal, F_signal, fftw.Estimate)
  backward := fftw.PlanDftC2R1d(F_signal, signal, fftw.Estimate)
  for i := range signal {
    signal[i] = float64(i)
  }
  forward.Execute()
  backward.Execute()
  
  c.Specify("Forward 1d Complx to Real FFT  works properly.", func() {
    for i := 0; i < len(signal); i++ {
      c.Expect(signal[i], IsWithin(1e-7), float64(i * len(signal)))
    }
  })
}
