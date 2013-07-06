## go.tesseract
go.tesseract is a wrapper for the tesseract-ocr library.

**go.tesseract is under heavy development and should not be used in a production environment.**

### Installation
You are required to install tesseract 3.02.02 or later. At time of writing this version of tesseract is not in the ubuntu repository yet. You absolutely need 3.02.02 (or later) as go.tesseract can not and will not compile with earlier versions.

**Before you continue, make sure you have installed [go.leptonica](//github.com/GeertJohan/go.leptonica). Please follow the directions in it's readme.**

Download and extract tesseract sources
```
svn checkout http://tesseract-ocr.googlecode.com/svn/tags/release-3.02.02 tesseract-ocr-read-only
cd tesseract-ocr-read-only
```

Configure, make and install
```
./autogen.sh
./configure
make
sudo make install
sudo ldconfig
```

Copy language files (do this for any language you require)
```
cp tessdata/eng.* /usr/local/share/tessdata/
```

For more information, view the tesseract [compilation guide](http://code.google.com/p/tesseract-ocr/wiki/Compiling).
