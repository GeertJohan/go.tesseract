## go.tesseract
go.tesseract is a wrapper for the tesseract-ocr library.

**go.tesseract is under heavy development and should not be used in a production environment.**

### Installation
First, install tesseract 3.02.02 or later. At time of writing this version of tesseract is not in the ubuntu repository yet. You absolutely need 3.02.02 (or later) as go.tesseract can not and will not compile with earlier versions.

Install dependencies, this will only work if you have a recent version of debian/ubuntu.
```
sudo apt-get install autoconf automake libtool
sudo apt-get install libpng12-dev
sudo apt-get install libjpeg62-dev
sudo apt-get install libtiff4-dev
sudo apt-get install zlib1g-dev
sudo apt-get install libleptonica-dev
```

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
