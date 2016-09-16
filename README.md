##go.tesseract
go.tesseract is a wrapper for the tesseract OCR library (text-recognition from image/pdf).

### Installation and dependencies
go.tesseract has two direct dependencies; `go.leptonica` and `libtesseract`

Make sure you have installed [go.leptonica](//github.com/GeertJohan/go.leptonica). go.leptonica has a C library dependency, please read the [go.leptonica/README.md](//github.com/GeertJohan/go.leptonica/blob/master/README.md).

You are required to install the tesseract library including development headers at version 3.04.00 or later. You absolutely need 3.04.00 (or later) as go.tesseract can not compile with earlier versions of tesseract. At time of writing this version of tesseract is not in the ubuntu/debian stable repository yet.

go.tesseract uses gopkg.in for versioned releases:

`go get gopkg.in/GeertJohan/go.tesseract.v1`

#### Debian testing (stretch) package
`sudo apt-get install -t testing libtesseract3 libtesseract-dev`

#### OSX with Homebrew

Do the following before trying to `go get` this package:

```
$ brew install leptonica
$ brew install tesseract
$ export CGO_LDFLAGS="-L$(brew --prefix leptonica)/lib -L$(brew --prefix tesseract)/lib"
$ export CGO_CFLAGS="-I$(brew --prefix leptonica)/include -I$(brew --prefix tesseract)/include"
```

*Note*: this assumes you are using the standard Brew path of `/usr/local/Cellar`

#### Manual installation
Download, configure, make and install
```
git clone https://github.com/tesseract-ocr/tesseract
cd tesseract
git checkout tags/3.04.00
./autogen.sh
./configure
make
sudo make install
sudo ldconfig
```

#### Language files
If you have installed from debian testing (jessie):
```
sudo apt-get install -t testing tesseract-ocr-YOUR-LANGUAGE-SHORTCODE

# example, this installs dutch and english
sudo apt-get install -t testing tesseract-ocr-nld
sudo apt-get install -t testing tesseract-ocr-eng

```

If you have installed manually; copy language files (do this for any language you require)
```
sudo cp tessdata/YOUR-LANGUAGE-SHORTCODE.* /usr/local/share/tessdata/

# example for english and dutch:
sudo cp tessdata/eng.* /usr/local/share/tessdata/
sudo cp tessdata/nld.* /usr/local/share/tessdata/
```

For more information, view the tesseract [compilation guide](http://code.google.com/p/tesseract-ocr/wiki/Compiling).
