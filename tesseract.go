package tesseract

// #cgo LDFLAGS: -ltesseract
// #include "tesseract/capi.h"
// #include <stdlib.h>
// #include <stdio.h>
import "C"
import (
	"errors"
	"github.com/GeertJohan/go.leptonica"
	"unsafe"
)

const version = "0.1"

// typedef struct TessBaseAPI TessBaseAPI;

// Tess represents a tesseract instance
type Tess struct {
	tba *C.TessBaseAPI
}

// const char* TessVersion();

// Version returns both go.tesseract's version as well as the version from the tesseract lib (>3.02.02)
func Version() string {
	libTessVersion := C.TessVersion()
	return "go.tesseract:" + version + " tesseract lib:" + C.GoString(libTessVersion)
}

// TessBaseAPI* TessBaseAPICreate();
// int TessBaseAPIInit3(TessBaseAPI* handle, const char* datapath, const char* language);

// NewTess creates and returns a new tesseract instance.
func NewTess(datapath string, language string) (*Tess, error) {
	tba := C.TessBaseAPICreate()

	cDatapath := C.CString(datapath)
	defer C.free(unsafe.Pointer(cDatapath))

	cLanguage := C.CString(language)
	defer C.free(unsafe.Pointer(cLanguage))

	res := C.TessBaseAPIInit3(tba, cDatapath, cLanguage)
	if res != 0 {
		return nil, errors.New("could not initiate new Tess instance")
	}

	tess := &Tess{
		tba: tba,
	}
	return tess, nil
}

// void TessBaseAPISetInputName( TessBaseAPI* handle, const char* name);

// SetInputName sets the Tess to read from given input filename
func (t *Tess) SetInputName(filename string) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	C.TessBaseAPISetInputName(t.tba, cFilename)
}

// void TessBaseAPISetImage2(TessBaseAPI* handle, const PIX* pix);

// SetImagePix sets the input image using a leptonica Pix
func (t *Tess) SetImagePix(pix *leptonica.Pix) {
	C.TessBaseAPISetImage2(t.tba, (*[0]byte)(unsafe.Pointer(pix.CPIX)))
}

// char* TessBaseAPIGetUTF8Text(TessBaseAPI* handle);

// GetText returns text after analysing the image(s)
func (t *Tess) GetText() string {
	cText := C.TessBaseAPIGetUTF8Text(t.tba)
	defer C.free(unsafe.Pointer(cText))
	text := C.GoString(cText)
	return text
}

// void TessBaseAPIPrintVariables( const TessBaseAPI* handle, FILE* fp);

// DumpVariables dumps the variables set on a Tess to stdout
func (t *Tess) DumpVariables() {
	C.TessBaseAPIPrintVariables(t.tba, (*C.FILE)(C.stdout))
}

// typedef struct TessPageIterator TessPageIterator;
// typedef struct TessResultIterator TessResultIterator;
// typedef struct TessMutableIterator TessMutableIterator;
// typedef enum TessOcrEngineMode { OEM_TESSERACT_ONLY, OEM_CUBE_ONLY, OEM_TESSERACT_CUBE_COMBINED, OEM_DEFAULT } TessOcrEngineMode;
// typedef enum TessPageSegMode { PSM_OSD_ONLY, PSM_AUTO_OSD, PSM_AUTO_ONLY, PSM_AUTO, PSM_SINGLE_COLUMN, PSM_SINGLE_BLOCK_VERT_TEXT, PSM_SINGLE_BLOCK, PSM_SINGLE_LINE, PSM_SINGLE_WORD, PSM_CIRCLE_WORD, PSM_SINGLE_CHAR, PSM_COUNT } TessPageSegMode;
// typedef enum TessPageIteratorLevel { RIL_BLOCK, RIL_PARA, RIL_TEXTLINE, RIL_WORD, RIL_SYMBOL} TessPageIteratorLevel;
// typedef enum TessPolyBlockType { PT_UNKNOWN, PT_FLOWING_TEXT, PT_HEADING_TEXT, PT_PULLOUT_TEXT, PT_TABLE, PT_VERTICAL_TEXT, PT_CAPTION_TEXT, PT_FLOWING_IMAGE, PT_HEADING_IMAGE, PT_PULLOUT_IMAGE, PT_HORZ_LINE, PT_VERT_LINE, PT_NOISE, PT_COUNT } TessPolyBlockType;
// typedef enum TessOrientation { ORIENTATION_PAGE_UP, ORIENTATION_PAGE_RIGHT, ORIENTATION_PAGE_DOWN, ORIENTATION_PAGE_LEFT } TessOrientation;
// typedef enum TessWritingDirection { WRITING_DIRECTION_LEFT_TO_RIGHT, WRITING_DIRECTION_RIGHT_TO_LEFT, WRITING_DIRECTION_TOP_TO_BOTTOM } TessWritingDirection;
// typedef enum TessTextlineOrder { TEXTLINE_ORDER_LEFT_TO_RIGHT, TEXTLINE_ORDER_RIGHT_TO_LEFT, TEXTLINE_ORDER_TOP_TO_BOTTOM } TessTextlineOrder;
// typedef struct ETEXT_DESC ETEXT_DESC;
// typedef struct Pix PIX;
// typedef struct Boxa BOXA;
// typedef struct Pixa PIXA;

// void TessDeleteText(char* text);
// void TessDeleteTextArray(char** arr);
// void TessDeleteIntArray(int* arr);

///* Base API */

// void TessBaseAPIDelete(TessBaseAPI* handle);

// void TessBaseAPISetOutputName(TessBaseAPI* handle, const char* name);

// BOOL TessBaseAPISetVariable(TessBaseAPI* handle, const char* name, const char* value);
// BOOL TessBaseAPISetDebugVariable(TessBaseAPI* handle, const char* name, const char* value);

// BOOL TessBaseAPIGetIntVariable( const TessBaseAPI* handle, const char* name, int* value);
// BOOL TessBaseAPIGetBoolVariable( const TessBaseAPI* handle, const char* name, BOOL* value);
// BOOL TessBaseAPIGetDoubleVariable(const TessBaseAPI* handle, const char* name, double* value);
// const char* TessBaseAPIGetStringVariable(const TessBaseAPI* handle, const char* name);

// void TessBaseAPIPrintVariables( const TessBaseAPI* handle, FILE* fp);
// BOOL TessBaseAPIPrintVariablesToFile(const TessBaseAPI* handle, const char* filename);

// int TessBaseAPIInit1(TessBaseAPI* handle, const char* datapath, const char* language, TessOcrEngineMode oem, char** configs, int configs_size);
// int TessBaseAPIInit2(TessBaseAPI* handle, const char* datapath, const char* language, TessOcrEngineMode oem);

// const char* TessBaseAPIGetInitLanguagesAsString(const TessBaseAPI* handle);
// char** TessBaseAPIGetLoadedLanguagesAsVector(const TessBaseAPI* handle);
// char** TessBaseAPIGetAvailableLanguagesAsVector(const TessBaseAPI* handle);

// int TessBaseAPIInitLangMod(TessBaseAPI* handle, const char* datapath, const char* language);
// void TessBaseAPIInitForAnalysePage(TessBaseAPI* handle);

// void TessBaseAPIReadConfigFile(TessBaseAPI* handle, const char* filename);
// void TessBaseAPIReadDebugConfigFile(TessBaseAPI* handle, const char* filename);

// void TessBaseAPISetPageSegMode(TessBaseAPI* handle, TessPageSegMode mode);
// TessPageSegMode TessBaseAPIGetPageSegMode(const TessBaseAPI* handle);

// char* TessBaseAPIRect(TessBaseAPI* handle, const unsigned char* imagedata, int bytes_per_pixel, int bytes_per_line, int left, int top, int width, int height);

// void TessBaseAPIClearAdaptiveClassifier(TessBaseAPI* handle);

// void TessBaseAPISetImage(TessBaseAPI* handle, const unsigned char* imagedata, int width, int height, int bytes_per_pixel, int bytes_per_line);

// void TessBaseAPISetSourceResolution(TessBaseAPI* handle, int ppi);

// void TessBaseAPISetRectangle(TessBaseAPI* handle, int left, int top, int width, int height);

// PIX* TessBaseAPIGetThresholdedImage( TessBaseAPI* handle);
// BOXA* TessBaseAPIGetRegions( TessBaseAPI* handle, PIXA** pixa);
// BOXA* TessBaseAPIGetTextlines( TessBaseAPI* handle, PIXA** pixa, int** blockids);
// BOXA* TessBaseAPIGetStrips( TessBaseAPI* handle, PIXA** pixa, int** blockids);
// BOXA* TessBaseAPIGetWords( TessBaseAPI* handle, PIXA** pixa);
// BOXA* TessBaseAPIGetConnectedComponents(TessBaseAPI* handle, PIXA** cc);
// BOXA* TessBaseAPIGetComponentImages( TessBaseAPI* handle, TessPageIteratorLevel level, BOOL text_only, PIXA** pixa, int** blockids);

// int TessBaseAPIGetThresholdedImageScaleFactor(const TessBaseAPI* handle);

// void TessBaseAPIDumpPGM(TessBaseAPI* handle, const char* filename);

// TessPageIterator* TessBaseAPIAnalyseLayout(TessBaseAPI* handle);

// int TessBaseAPIRecognize(TessBaseAPI* handle, ETEXT_DESC* monitor);
// int TessBaseAPIRecognizeForChopTest(TessBaseAPI* handle, ETEXT_DESC* monitor);
// char* TessBaseAPIProcessPages(TessBaseAPI* handle, const char* filename, const char* retry_config, int timeout_millisec);
// char* TessBaseAPIProcessPage(TessBaseAPI* handle, PIX* pix, int page_index, const char* filename, const char* retry_config, int timeout_millisec);

// TessResultIterator* TessBaseAPIGetIterator(TessBaseAPI* handle);
// TessMutableIterator* TessBaseAPIGetMutableIterator(TessBaseAPI* handle);

// char* TessBaseAPIGetHOCRText(TessBaseAPI* handle, int page_number);
// char* TessBaseAPIGetBoxText(TessBaseAPI* handle, int page_number);
// char* TessBaseAPIGetUNLVText(TessBaseAPI* handle);
// int TessBaseAPIMeanTextConf(TessBaseAPI* handle);
// int* TessBaseAPIAllWordConfidences(TessBaseAPI* handle);
// BOOL TessBaseAPIAdaptToWordStr(TessBaseAPI* handle, TessPageSegMode mode, const char* wordstr);

// void TessBaseAPIClear(TessBaseAPI* handle);
// void TessBaseAPIEnd(TessBaseAPI* handle);

// int TessBaseAPIIsValidWord(TessBaseAPI* handle, const char *word);
// BOOL TessBaseAPIGetTextDirection(TessBaseAPI* handle, int* out_offset, float* out_slope);

// const char* TessBaseAPIGetUnichar(TessBaseAPI* handle, int unichar_id);
// void TessBaseAPISetMinOrientationMargin(TessBaseAPI* handle, double margin);

// /* Page iterator */
// void TessPageIteratorDelete(TessPageIterator* handle);
// TessPageIterator* TessPageIteratorCopy(const TessPageIterator* handle);
// void TessPageIteratorBegin(TessPageIterator* handle);
// BOOL TessPageIteratorNext(TessPageIterator* handle, TessPageIteratorLevel level);
// BOOL TessPageIteratorIsAtBeginningOf(const TessPageIterator* handle, TessPageIteratorLevel level);
// BOOL TessPageIteratorIsAtFinalElement(const TessPageIterator* handle, TessPageIteratorLevel level,
// TessPageIteratorLevel element);
// BOOL TessPageIteratorBoundingBox(const TessPageIterator* handle, TessPageIteratorLevel level,
// int* left, int* top, int* right, int* bottom);
// TessPolyBlockType TessPageIteratorBlockType(const TessPageIterator* handle);
// PIX* TessPageIteratorGetBinaryImage(const TessPageIterator* handle, TessPageIteratorLevel level);
// PIX* TessPageIteratorGetImage(const TessPageIterator* handle, TessPageIteratorLevel level, int padding, int* left, int* top);
// BOOL TessPageIteratorBaseline(const TessPageIterator* handle, TessPageIteratorLevel level, int* x1, int* y1, int* x2, int* y2);
// void TessPageIteratorOrientation(TessPageIterator* handle, TessOrientation *orientation, TessWritingDirection *writing_direction, TessTextlineOrder *textline_order, float *deskew_angle);

// /* Result iterator */
// void TessResultIteratorDelete(TessResultIterator* handle);
// TessResultIterator* TessResultIteratorCopy(const TessResultIterator* handle);
// TessPageIterator* TessResultIteratorGetPageIterator(TessResultIterator* handle);
// const TessPageIterator* TessResultIteratorGetPageIteratorConst(const TessResultIterator* handle);
// char* TessResultIteratorGetUTF8Text(const TessResultIterator* handle, TessPageIteratorLevel level);
// float TessResultIteratorConfidence(const TessResultIterator* handle, TessPageIteratorLevel level);
// const char* TessResultIteratorWordFontAttributes(const TessResultIterator* handle, BOOL* is_bold, BOOL* is_italic, BOOL* is_underlined, BOOL* is_monospace, BOOL* is_serif, BOOL* is_smallcaps, int* pointsize, int* font_id);
// BOOL TessResultIteratorWordIsFromDictionary(const TessResultIterator* handle);
// BOOL TessResultIteratorWordIsNumeric(const TessResultIterator* handle);
// BOOL TessResultIteratorSymbolIsSuperscript(const TessResultIterator* handle);
// BOOL TessResultIteratorSymbolIsSubscript(const TessResultIterator* handle);
// BOOL TessResultIteratorSymbolIsDropcap(const TessResultIterator* handle);
