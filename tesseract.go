package tesseract

// #cgo LDFLAGS: -L /usr/local/lib -ltesseract
// #include "tesseract/capi.h"
// #include <stdlib.h>
import "C"

import (
	"bytes"
	"errors"
	"io"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"

	leptonica "gopkg.in/GeertJohan/go.leptonica.v1"
)

const version = "1.0"

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
	// create new empty TessBaseAPI
	tba := C.TessBaseAPICreate()

	// prepare string for C call
	cDatapath := C.CString(datapath)
	defer C.free(unsafe.Pointer(cDatapath))

	// prepare string for C call
	cLanguage := C.CString(language)
	defer C.free(unsafe.Pointer(cLanguage))

	// initialize datapath and language on TessBaseAPI
	res := C.TessBaseAPIInit3(tba, cDatapath, cLanguage)
	if res != 0 {
		return nil, errors.New("could not initiate new Tess instance")
	}

	// create tesseract instance (Tess)
	tess := &Tess{
		tba: tba,
	}

	// set GC finalizer, to be ran in case the user forgets to call Close()
	runtime.SetFinalizer(tess, (*Tess).delete)

	// all done
	return tess, nil
}

// void TessBaseAPIDelete(TessBaseAPI* handle);
// void TessBaseAPIEnd(TessBaseAPI* handle);
func (t *Tess) delete() {
	if t.tba != nil {
		C.TessBaseAPIEnd(t.tba)
		C.TessBaseAPIDelete(t.tba)
	}
}

// Close clears the tesseract instance from memory
func (t *Tess) Close() {
	t.delete()
	t.tba = nil
}

/* void TessBaseAPIClear(TessBaseAPI* handle);

Free up recognition results and any stored image data, without actually
freeing any recognition data that would be time-consuming to reload.
Afterwards, you must call SetImage or TesseractRect before doing
any Recognize or Get* operation.
*/

// Clear frees up recognition results and any stored image data, without actually freeing any recognition data that would be time-consuming to reload.
// Afterwards, you must call SetImagePix before doing any Recognize or Get* operation.
func (t *Tess) Clear() {
	C.TessBaseAPIClear(t.tba)
}

// map t.delete() on t GC as hook/callback in NewXXX() call's

/* void TessBaseAPISetInputName( TessBaseAPI* handle, const char* name);

Set the name of the input file. Needed only for training and
loading a UNLV zone file.
*/

// SetInputName sets the name of the input file. Needed only for training and loading a UNLV zone file.
// ++ TODO: drop this?
func (t *Tess) SetInputName(filename string) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	C.TessBaseAPISetInputName(t.tba, cFilename)
}

// void TessBaseAPISetImage2(TessBaseAPI* handle, const PIX* pix);

// SetImagePix sets the input image using a leptonica Pix
func (t *Tess) SetImagePix(pix *leptonica.Pix) {
	C.TessBaseAPISetImage2(t.tba, (*C.struct_Pix)(unsafe.Pointer(pix.CPIX())))
}

// void TessBaseAPISetSourceResolution(TessBaseAPI* handle, int ppi);

// SetSourceResolution set the resolution of the source image in pixels per inch.
// The must be called after SetInputName or SetImagePix and has the same reslut as the --dpi flag.
func (t *Tess) SetSourceResolution(ppi int) {
	C.TessBaseAPISetSourceResolution(t.tba, C.int(ppi))
}

/* char* TessBaseAPIGetUTF8Text(TessBaseAPI* handle);

Make a text string from the internal data structures.
*/

// Text returns text after analysing the image(s)
func (t *Tess) Text() string {
	cText := C.TessBaseAPIGetUTF8Text(t.tba)
	defer C.free(unsafe.Pointer(cText))
	text := C.GoString(cText)
	return text
}

/* char* TessBaseAPIGetHOCRText(TessBaseAPI* handle, int page_number);

Make a HTML-formatted string with hOCR markup from the internal
data structures.
page_number is 0-based but will appear in the output as 1-based.
Image name/input_file_ can be set by SetInputName before calling
GetHOCRText
STL removed from original patch submission and refactored by rays.
*/

// HOCRText returns the HOCR text for given pagenumber
func (t *Tess) HOCRText(pagenumber int) string {
	cText := C.TessBaseAPIGetHOCRText(t.tba, C.int(pagenumber))
	defer C.free(unsafe.Pointer(cText))
	text := C.GoString(cText)
	return text
}

/* char* TessBaseAPIGetTsvText(TessBaseAPI* handle, int page_number);

Make a TSV-formatted string from the internal data structures.
page_number is 0-based but will appear in the output as 1-based.
Returned string must be freed with the delete [] operator.
*/

// TSVText returns a TSV-formatted string.
func (t *Tess) TSVText(pagenumber int) string {
	cText := C.TessBaseAPIGetTsvText(t.tba, C.int(pagenumber))
	defer C.free(unsafe.Pointer(cText))
	text := C.GoString(cText)
	return text
}

/* char* TessBaseAPIGetBoxText(TessBaseAPI* handle, int page_number);

The recognized text is returned as a char* which is coded
as a UTF8 box file and must be freed with the delete [] operator.
page_number is a 0-base page index that will appear in the box file.
*/

// BoxTextRaw returns the raw box text for given pagenumber
func (t *Tess) BoxTextRaw(pagenumber int) string {
	cText := C.TessBaseAPIGetBoxText(t.tba, C.int(pagenumber))
	defer C.free(unsafe.Pointer(cText))
	text := C.GoString(cText)
	return text
}

// TODO: make this: `type BoxText []BoxCharacter` ?
type BoxText struct {
	Characters []BoxCharacter
}

type BoxCharacter struct {
	Character  rune
	StartX     uint32
	StartY     uint32
	EndX       uint32
	EndY       uint32
	Pagenumber uint32
}

// BoxText returns the output given by BoxTextRaw as BoxText object
func (tess *Tess) BoxText(pagenumber int) (*BoxText, error) {
	text := tess.BoxTextRaw(pagenumber)
	textBuffer := bytes.NewBufferString(text)

	bt := &BoxText{
		Characters: make([]BoxCharacter, 0, len(text)),
	}
	for {
		line, err := textBuffer.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return bt, nil
			}
			return nil, err
		}
		line = strings.TrimRight(line, "\n")
		fields := strings.Split(line, " ")
		if len(fields) != 6 {
			f := strconv.Itoa(len(fields))
			return nil, errors.New("unexpected BoxText format (Length != 6) Length is: " + f)
		}
		if utf8.RuneCountInString(fields[0]) != 1 {
			return nil, errors.New("unexpected BoxText format (RuneCount error): " + fields[0])
		}

		sx, err := strconv.ParseUint(fields[1], 10, 32)
		if err != nil {
			return nil, err
		}
		sy, err := strconv.ParseUint(fields[2], 10, 32)
		if err != nil {
			return nil, err
		}
		ex, err := strconv.ParseUint(fields[3], 10, 32)
		if err != nil {
			return nil, err
		}
		ey, err := strconv.ParseUint(fields[4], 10, 32)
		if err != nil {
			return nil, err
		}
		pgnr, err := strconv.ParseUint(fields[5], 10, 32)
		if err != nil {
			return nil, err
		}
		bt.Characters = append(bt.Characters, BoxCharacter{
			Character:  rune(fields[0][0]),
			StartX:     uint32(sx),
			StartY:     uint32(sy),
			EndX:       uint32(ex),
			EndY:       uint32(ey),
			Pagenumber: uint32(pgnr),
		})
	}
}

// typedef enum TessPageSegMode { PSM_OSD_ONLY, PSM_AUTO_OSD, PSM_AUTO_ONLY, PSM_AUTO, PSM_SINGLE_COLUMN, PSM_SINGLE_BLOCK_VERT_TEXT, PSM_SINGLE_BLOCK, PSM_SINGLE_LINE, PSM_SINGLE_WORD, PSM_CIRCLE_WORD, PSM_SINGLE_CHAR, PSM_COUNT } TessPageSegMode;
type PageSegMode int

const (
	PSM_OSD_ONLY PageSegMode = iota
	PSM_AUTO_OSD
	PSM_AUTO_ONLY
	PSM_AUTO
	PSM_SINGLE_COLUMN
	PSM_SINGLE_BLOCK_VERT_TEXT
	PSM_SINGLE_BLOCK
	PSM_SINGLE_LINE
	PSM_SINGLE_WORD
	PSM_CIRCLE_WORD
	PSM_SINGLE_CHAR
	PSM_COUNT
)

// void TessBaseAPISetPageSegMode(TessBaseAPI* handle, TessPageSegMode mode);
func (tess *Tess) SetPageSegMode(psm PageSegMode) {
	C.TessBaseAPISetPageSegMode(tess.tba, C.TessPageSegMode(psm))
}

/* char* TessBaseAPIGetUNLVText(TessBaseAPI* handle);

The recognized text is returned as a char* which is coded
as UNLV format Latin-1 with specific reject and suspect codes
and must be freed with the delete [] operator.
*/

// UNLVText returns the UNLV text
func (t *Tess) UNLVText() string {
	cText := C.TessBaseAPIGetUNLVText(t.tba)
	defer C.free(unsafe.Pointer(cText))
	text := C.GoString(cText)
	return text
}

/* const char* TessBaseAPIGetInitLanguagesAsString(const TessBaseAPI* handle);

Returns the languages string used in the last valid initialization.
If the last initialization specified "deu+hin" then that will be
returned. If hin loaded eng automatically as well, then that will
not be included in this list. To find the languages actually
loaded use GetLoadedLanguagesAsVector.
The returned string should NOT be deleted.
*/

// InitializedLanguages returns the languages string used in the last valid initialization.
// If the last initialization specified "deu+hin" then that will be returned.
// If hin loaded eng automatically as well, then that will not be included in this list.
// To find the languages actually loaded use (*Tess).LoadedLanguages().
func (t *Tess) InitializedLanguages() string {
	cLang := C.TessBaseAPIGetInitLanguagesAsString(t.tba)
	defer C.free(unsafe.Pointer(cLang))
	return C.GoString(cLang)
}

/* char** TessBaseAPIGetLoadedLanguagesAsVector(const TessBaseAPI* handle);

Returns the loaded languages in the vector of STRINGs.
Includes all languages loaded by the last Init, including those loaded
as dependencies of other loaded languages.
*/

// LoadedLanguages returns the loaded languages in the vector of STRINGs.
// Includes all languages loaded for the given tesseract instance, including those loaded as dependencies of other loaded languages.
func (t *Tess) LoadedLanguages() []string {
	cLangs := C.TessBaseAPIGetLoadedLanguagesAsVector(t.tba)
	defer C.TessDeleteTextArray(cLangs)

	langs := cStringVectorToStringslice(cLangs)
	return langs
}

/* char** TessBaseAPIGetAvailableLanguagesAsVector(const TessBaseAPI* handle);

Returns the available languages in the vector of STRINGs.
*/

// AvailableLanguages returns the languages available to the given tesseract instance.
// To find the languages actually loaded use (*Tess).LoadedLanguages().
func (t *Tess) AvailableLanguages() []string {
	cLangs := C.TessBaseAPIGetAvailableLanguagesAsVector(t.tba)
	defer C.TessDeleteTextArray(cLangs)

	langs := cStringVectorToStringslice(cLangs)
	return langs
}

/* void TessBaseAPIPrintVariables( const TessBaseAPI* handle, FILE* fp);

Print Tesseract parameters to the given file.
*/

// DumpVariables dumps the variables set on a Tess to stdout
func (t *Tess) DumpVariables() {
	C.TessBaseAPIPrintVariables(t.tba, (*C.FILE)(C.stdout))
}

// BOOL TessBaseAPISetVariable(TessBaseAPI* handle, const char* name, const char* value);
func (t *Tess) SetVariable(name, value string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	worked := C.TessBaseAPISetVariable(t.tba, cName, cValue)
	if worked != 1 {
		return errors.New("Unable to set the variable: " + name + " to " + value)
	}
	return nil
}

// void TessBaseAPISetRectangle(TessBaseAPI* handle, int left, int top, int width, int height);
func (t *Tess) SetRectangle(left, top, width, height int) {
	C.TessBaseAPISetRectangle(t.tba, C.int(left), C.int(top), C.int(width), C.int(height))
}

// int TessBaseAPIRecognize(TessBaseAPI* handle, ETEXT_DESC* monitor);
func (t *Tess) Recognize() error {
	ret := C.TessBaseAPIRecognize(t.tba, nil)
	if ret != 0 {
		return errors.New("recognition failed")
	}
	return nil
}

// typedef enum TessPageIteratorLevel { RIL_BLOCK, RIL_PARA, RIL_TEXTLINE, RIL_WORD, RIL_SYMBOL} TessPageIteratorLevel;
type PageIteratorLevel int

const (
	RIL_BLOCK PageIteratorLevel = iota
	RIL_PARA
	RIL_TEXTLINE
	RIL_WORD
	RIL_SYMBOL
)

// TessResultIterator* TessBaseAPIGetIterator(TessBaseAPI* handle);
func (t *Tess) Iterator() (*ResultIterator, error) {
	ri := C.TessBaseAPIGetIterator(t.tba)

	if ri == nil {
		return nil, errors.New("no results")
	}

	resultIterator := &ResultIterator{
		ri: ri,
	}

	runtime.SetFinalizer(resultIterator, (*ResultIterator).delete)
	return resultIterator, nil
}

// typedef struct TessResultIterator TessResultIterator;
type ResultIterator struct {
	ri *C.TessResultIterator
}

// void TessResultIteratorDelete(TessResultIterator* handle);
func (r *ResultIterator) delete() {
	if r.ri != nil {
		C.TessResultIteratorDelete(r.ri)
	}
}

// TESS_API BOOL  TESS_CALL TessResultIteratorNext(TessResultIterator* handle, TessPageIteratorLevel level);
func (r *ResultIterator) Next(level PageIteratorLevel) bool {
	return gobool(C.TessResultIteratorNext(r.ri, C.TessPageIteratorLevel(level)))
}

// char* TessResultIteratorGetUTF8Text(const TessResultIterator* handle, TessPageIteratorLevel level);
func (r *ResultIterator) Text(level PageIteratorLevel) (string, error) {
	cText := C.TessResultIteratorGetUTF8Text(r.ri, C.TessPageIteratorLevel(level))
	if cText == nil {
		return "", errors.New("already at the end")
	}
	defer C.free(unsafe.Pointer(cText))
	text := C.GoString(cText)
	return text, nil
}

// typedef struct TessPageIterator TessPageIterator;
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

// void TessBaseAPISetOutputName(TessBaseAPI* handle, const char* name);

// BOOL TessBaseAPISetDebugVariable(TessBaseAPI* handle, const char* name, const char* value);

// BOOL TessBaseAPIGetIntVariable( const TessBaseAPI* handle, const char* name, int* value);
// BOOL TessBaseAPIGetBoolVariable( const TessBaseAPI* handle, const char* name, BOOL* value);
// BOOL TessBaseAPIGetDoubleVariable(const TessBaseAPI* handle, const char* name, double* value);
// const char* TessBaseAPIGetStringVariable(const TessBaseAPI* handle, const char* name);

// void TessBaseAPIPrintVariables( const TessBaseAPI* handle, FILE* fp);
// BOOL TessBaseAPIPrintVariablesToFile(const TessBaseAPI* handle, const char* filename);

// int TessBaseAPIInit1(TessBaseAPI* handle, const char* datapath, const char* language, TessOcrEngineMode oem, char** configs, int configs_size);
// int TessBaseAPIInit2(TessBaseAPI* handle, const char* datapath, const char* language, TessOcrEngineMode oem);

// int TessBaseAPIInitLangMod(TessBaseAPI* handle, const char* datapath, const char* language);
// void TessBaseAPIInitForAnalysePage(TessBaseAPI* handle);

// void TessBaseAPIReadConfigFile(TessBaseAPI* handle, const char* filename);
// void TessBaseAPIReadDebugConfigFile(TessBaseAPI* handle, const char* filename);

// void TessBaseAPISetPageSegMode(TessBaseAPI* handle, TessPageSegMode mode);
// TessPageSegMode TessBaseAPIGetPageSegMode(const TessBaseAPI* handle);

// char* TessBaseAPIRect(TessBaseAPI* handle, const unsigned char* imagedata, int bytes_per_pixel, int bytes_per_line, int left, int top, int width, int height);

// void TessBaseAPIClearAdaptiveClassifier(TessBaseAPI* handle);

// void TessBaseAPISetImage(TessBaseAPI* handle, const unsigned char* imagedata, int width, int height, int bytes_per_pixel, int bytes_per_line);

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

// int TessBaseAPIRecognizeForChopTest(TessBaseAPI* handle, ETEXT_DESC* monitor);
// char* TessBaseAPIProcessPages(TessBaseAPI* handle, const char* filename, const char* retry_config, int timeout_millisec);
// char* TessBaseAPIProcessPage(TessBaseAPI* handle, PIX* pix, int page_index, const char* filename, const char* retry_config, int timeout_millisec);

// TessMutableIterator* TessBaseAPIGetMutableIterator(TessBaseAPI* handle);

// int TessBaseAPIMeanTextConf(TessBaseAPI* handle);
// int* TessBaseAPIAllWordConfidences(TessBaseAPI* handle);
// BOOL TessBaseAPIAdaptToWordStr(TessBaseAPI* handle, TessPageSegMode mode, const char* wordstr);

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
