package mimetype

import (
	"os"
	"strings"
	"sync"

	"github.com/gabriel-vasile/mimetype"
)

var (
	suffixMimeTypes = map[string][]string{
		"3g2":     {"video/3gpp2", "video/3g2", "audio/3gpp2"},
		"3gp":     {"video/3gpp", "video/3gp", "audio/3gpp"},
		"3mf":     {"application/vnd.ms-package.3dmanufacturing-3dmodel+xml"},
		"7z":      {"application/x-7z-compressed"},
		"a":       {"application/x-archive", "application/x-unix-archive"},
		"aac":     {"audio/aac"},
		"aaf":     {"application/octet-stream"},
		"accdb":   {"application/x-msaccess"},
		"aiff":    {"audio/aiff", "audio/x-aiff"},
		"amf":     {"application/x-amf"},
		"amr":     {"audio/amr", "audio/amr-nb"},
		"ape":     {"audio/ape"},
		"apk":     {"application/vnd.android.package-archive"},
		"asf":     {"video/x-ms-asf", "video/asf", "video/x-ms-wmv"},
		"atom":    {"application/atom+xml"},
		"au":      {"audio/basic"},
		"avi":     {"video/x-msvideo", "video/avi", "video/msvideo"},
		"avif":    {"image/avif"},
		"bmp":     {"image/bmp", "image/x-bmp", "image/x-ms-bmp"},
		"bpg":     {"image/bpg"},
		"bz2":     {"application/x-bzip2"},
		"cab":     {"application/x-installshield"},
		"cbor":    {"application/cbor"},
		"class":   {"application/x-java-applet"},
		"cpio":    {"application/x-cpio"},
		"crx":     {"application/x-chrome-extension"},
		"csv":     {"text/csv"},
		"dae":     {"model/vnd.collada+xml"},
		"dbf":     {"application/x-dbf"},
		"dcm":     {"application/dicom"},
		"deb":     {"application/vnd.debian.binary-package"},
		"djvu":    {"image/vnd.djvu"},
		"doc":     {"application/msword", "application/vnd.ms-word"},
		"docx":    {"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		"dvb":     {"video/vnd.dvb.file"},
		"dwg":     {"image/vnd.dwg", "image/x-dwg", "application/acad", "application/x-acad", "application/autocad_dwg", "application/dwg", "application/x-dwg", "application/x-autocad", "drawing/dwg"},
		"eot":     {"application/vnd.ms-fontobject"},
		"epub":    {"application/epub+zip"},
		"exe":     {"application/vnd.microsoft.portable-executable"},
		"fdf":     {"application/vnd.fdf"},
		"fits":    {"application/fits"},
		"flac":    {"audio/flac"},
		"flv":     {"video/x-flv"},
		"gbr":     {"image/x-gimp-gbr"},
		"geojson": {"application/geo+json"},
		"gif":     {"image/gif"},
		"glb":     {"model/gltf-binary"},
		"gml":     {"application/gml+xml"},
		"gpx":     {"application/gpx+xml"},
		"gz":      {"application/gzip", "application/x-gzip", "application/x-gunzip", "application/gzipped", "application/gzip-compressed", "application/x-gzip-compressed", "gzip/document"},
		"har":     {"application/json"},
		"hdr":     {"image/vnd.radiance"},
		"heic":    {"image/heic-sequence"},
		"heif":    {"image/heif-sequence"},
		"html":    {"text/html"},
		"icns":    {"image/x-icns"},
		"ico":     {"image/x-icon"},
		"ics":     {"text/calendar"},
		"jar":     {"application/jar"},
		"jp2":     {"image/jp2"},
		"jpf":     {"image/jpx"},
		"jpg":     {"image/jpeg"},
		"jpm":     {"image/jpm", "video/jpm"},
		"js":      {"text/javascript", "application/x-javascript", "application/javascript"},
		"json":    {"application/json"},
		"jxl":     {"image/jxl"},
		"jxr":     {"image/jxr", "image/vnd.ms-photo"},
		"jxs":     {"image/jxs"},
		"kml":     {"application/vnd.google-earth.kml+xml"},
		"lit":     {"application/x-ms-reader"},
		"lnk":     {"application/x-ms-shortcut"},
		"lua":     {"text/x-lua"},
		"lz":      {"application/lzip", "application/x-lzip"},
		"m3u":     {"application/vnd.apple.mpegurl", "audio/mpegurl"},
		"m4a":     {"audio/x-m4a"},
		"m4v":     {"video/x-m4v"},
		"macho":   {"application/x-mach-binary"},
		"mdb":     {"application/x-msaccess"},
		"midi":    {"audio/midi", "audio/mid", "audio/sp-midi", "audio/x-mid", "audio/x-midi"},
		"mj2":     {"video/mj2"},
		"mkv":     {"video/x-matroska"},
		"mobi":    {"application/x-mobipocket-ebook"},
		"mov":     {"video/quicktime"},
		"mp3":     {"audio/mpeg", "audio/x-mpeg", "audio/mp3"},
		"mp4":     {"audio/mp4", "audio/x-mp4a"},
		"mpc":     {"audio/musepack"},
		"mpeg":    {"video/mpeg"},
		"mqv":     {"video/quicktime"},
		"mrc":     {"application/marc"},
		"msg":     {"application/vnd.ms-outlook"},
		"msi":     {"application/x-ms-installer", "application/x-windows-installer", "application/x-msi"},
		"ndjson":  {"application/x-ndjson"},
		"nes":     {"application/vnd.nintendo.snes.rom"},
		"odc":     {"application/vnd.oasis.opendocument.chart", "application/x-vnd.oasis.opendocument.chart"},
		"odf":     {"application/vnd.oasis.opendocument.formula", "application/x-vnd.oasis.opendocument.formula"},
		"odg":     {"application/vnd.oasis.opendocument.graphics", "application/x-vnd.oasis.opendocument.graphics"},
		"odp":     {"application/vnd.oasis.opendocument.presentation", "application/x-vnd.oasis.opendocument.presentation"},
		"ods":     {"application/vnd.oasis.opendocument.spreadsheet", "application/x-vnd.oasis.opendocument.spreadsheet"},
		"odt":     {"application/vnd.oasis.opendocument.text", "application/x-vnd.oasis.opendocument.text"},
		"oga":     {"audio/ogg"},
		"ogg":     {"application/ogg", "application/x-ogg"},
		"ogv":     {"video/ogg"},
		"otf":     {"font/otf"},
		"otg":     {"application/vnd.oasis.opendocument.graphics-template", "application/x-vnd.oasis.opendocument.graphics-template"},
		"otp":     {"application/vnd.oasis.opendocument.presentation-template", "application/x-vnd.oasis.opendocument.presentation-template"},
		"ots":     {"application/vnd.oasis.opendocument.spreadsheet-template", "application/x-vnd.oasis.opendocument.spreadsheet-template"},
		"ott":     {"application/vnd.oasis.opendocument.text-template", "application/x-vnd.oasis.opendocument.text-template"},
		"owl":     {"application/owl+xml"},
		"p7s":     {"application/pkcs7-signature"},
		"parquet": {"application/vnd.apache.parquet", "application/x-parquet"},
		"pat":     {"image/x-gimp-pat"},
		"pdf":     {"application/pdf", "application/x-pdf"},
		"php":     {"text/x-php"},
		"pl":      {"text/x-perl"},
		"png":     {"image/vnd.mozilla.apng"},
		"ppt":     {"application/vnd.ms-powerpoint", "application/mspowerpoint"},
		"pptx":    {"application/vnd.openxmlformats-officedocument.presentationml.presentation"},
		"ps":      {"application/postscript"},
		"psd":     {"image/vnd.adobe.photoshop", "image/x-psd", "application/photoshop"},
		"pub":     {"application/vnd.ms-publisher"},
		"py":      {"text/x-python", "text/x-script.python", "application/x-python"},
		"qcp":     {"audio/qcelp"},
		"rar":     {"application/x-rar-compressed", "application/x-rar"},
		"rmvb":    {"application/vnd.rn-realmedia-vbr"},
		"rpm":     {"application/x-rpm"},
		"rss":     {"application/rss+xml", "text/rss"},
		"rtf":     {"text/rtf", "application/rtf"},
		"shp":     {"application/vnd.shp"},
		"shx":     {"application/vnd.shx"},
		"so":      {"application/x-sharedlib"},
		"sqlite":  {"application/vnd.sqlite3", "application/x-sqlite3"},
		"srt":     {"application/x-subrip", "application/x-srt", "text/x-srt"},
		"svg":     {"image/svg+xml"},
		"swf":     {"application/x-shockwave-flash"},
		"sxc":     {"application/vnd.sun.xml.calc"},
		"tar":     {"application/x-tar"},
		"tcl":     {"text/x-tcl", "application/x-tcl"},
		"tcx":     {"application/vnd.garmin.tcx+xml"},
		"tiff":    {"image/tiff"},
		"torrent": {"application/x-bittorrent"},
		"tsv":     {"text/tab-separated-values"},
		"ttc":     {"font/collection"},
		"ttf":     {"font/ttf", "font/sfnt", "application/x-font-ttf", "application/font-sfnt"},
		"txt":     {"text/plain"},
		"vcf":     {"text/vcard"},
		"voc":     {"audio/x-unknown"},
		"vtt":     {"text/vtt"},
		"warc":    {"application/warc"},
		"wasm":    {"application/wasm"},
		"wav":     {"audio/wav", "audio/x-wav", "audio/vnd.wave", "audio/wave"},
		"webm":    {"video/webm", "audio/webm"},
		"webp":    {"image/webp"},
		"woff":    {"font/woff"},
		"woff2":   {"font/woff2"},
		"x3d":     {"model/x3d+xml"},
		"xar":     {"application/x-xar"},
		"xcf":     {"image/x-xcf"},
		"xfdf":    {"application/vnd.adobe.xfdf"},
		"xlf":     {"application/x-xliff+xml"},
		"xls":     {"application/vnd.ms-excel", "application/msexcel"},
		"xlsx":    {"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		"xml":     {"text/xml", "application/xml"},
		"xpm":     {"image/x-xpixmap"},
		"xz":      {"application/x-xz"},
		"zip":     {"application/zip", "application/x-zip", "application/x-zip-compressed"},
		"zst":     {"application/zstd"},
	}
)

var mu sync.Mutex
var mimeTypeSuffixMap map[string][]string

// Get 返回各种 suffix 对应的 mimetype
func Get() map[string][]string {
	return suffixMimeTypes
}

// GetSuffixes 返回各种 mimetype 对应的 suffix
func GetSuffixes() map[string][]string {
	if mimeTypeSuffixMap != nil {
		return mimeTypeSuffixMap
	}

	mu.Lock()
	defer mu.Unlock()

	if mimeTypeSuffixMap == nil {
		mimeTypeSuffixMap = make(map[string][]string)
	}

	kvm := map[string][]string{}
	for suffix, mimetypes := range Get() {
		for _, mimetype := range mimetypes {
			if exits, ok := kvm[mimetype]; ok {
				kvm[mimetype] = append(exits, suffix)
			} else {
				kvm[mimetype] = []string{suffix}
			}
		}
	}
	mimeTypeSuffixMap = kvm

	return mimeTypeSuffixMap
}

func Detect(data []byte) string {
	if m := mimetype.Detect(data); m != nil {
		if slice := strings.Split(m.String(), ";"); len(slice) > 0 {
			return strings.TrimSpace(slice[0])
		}
	}
	return ""
}

func DetectFile(file string) string {
	if file == "" {
		return ""
	}

	fi, err := os.Stat(file)
	if err != nil || fi.IsDir() {
		return ""
	}

	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer func() {
		_ = f.Close()
	}()

	buf := make([]byte, 3072)
	n, err := f.Read(buf)
	if err != nil {
		return ""
	}
	return Detect(buf[:n])
}

func Valid(data []byte, mimetypes []string) (string, bool) {
	typ := Detect(data)
	for _, m := range mimetypes {
		if strings.EqualFold(typ, m) {
			return typ, true
		}
	}
	return typ, false
}

func Contains(data []byte, mimetypes []string) (string, bool) {
	typ := Detect(data)
	for _, m := range mimetypes {
		if strings.HasPrefix(typ, m) {
			return typ, true
		}
	}
	return typ, false
}
