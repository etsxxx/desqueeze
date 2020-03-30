package jpeg

import (
	"fmt"
	"image/jpeg"
	"os"
	"testing"
)

const outDir = "../../test/out"

func TestMain(m *testing.M) {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		os.Mkdir(outDir, os.ModePerm)

	}
	exitCode := m.Run()
	os.Exit(exitCode)
}

func getSize(path string) (x int, y int, err error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return 0, 0, fmt.Errorf("Cannot read input file: %s", path)
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		return 0, 0, fmt.Errorf("Cannot decode input file: %s", path)
	}
	rct := img.Bounds()
	x = rct.Dx()
	y = rct.Dy()
	return
}

func TestDesqueeze(t *testing.T) {
	type args struct {
		inputPath  string
		outputPath string
		multiply   float64
		quality    int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "standard jpeg",
			args: args{
				inputPath:  "../../test/input.jpg",
				outputPath: outDir + "/output.jpg",
				multiply:   1.33,
				quality:    50,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Desqueeze(tt.args.inputPath, tt.args.outputPath, tt.args.multiply, tt.args.quality); (err != nil) != tt.wantErr {
				t.Errorf("Desqueeze() error = %v, wantErr %v", err, tt.wantErr)
			}
			inputDx, inputDy, err := getSize(tt.args.inputPath)
			if err != nil {
				t.Errorf("Data is invalid: %s", err.Error())
			}
			outputDx, outputDy, err := getSize(tt.args.outputPath)
			if err != nil {
				t.Errorf("Data is invalid: %s", err.Error())
			}
			wantDx, wantDy := 512, 512
			if inputDx == outputDx {
				t.Errorf("Desqueeze() outputs width = %d, want = %d", inputDx, outputDx)
			} else if outputDx != wantDx {
				t.Errorf("Desqueeze() outputs width = %d, want = %d ", outputDx, wantDx)
			} else if inputDy != outputDy {
				t.Errorf("Desqueeze() outputs height = %d, want = %d", inputDy, outputDy)
			} else if outputDy != wantDy {
				t.Errorf("Desqueeze() outputs height = %d, want = %d ", outputDy, wantDy)
			}
		})
	}
}
