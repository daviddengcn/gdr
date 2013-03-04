package gdrf

import (
	"github.com/daviddengcn/go-villa" // Hello
	"go/parser"
	"go/printer"
	"go/token"
	"go/ast"
	"io"
	"os"
	"fmt"
	"bytes"
	"strconv"
)

func init() {
	_ = os.Stdin
	_ = fmt.Println
	_ = printer.Fprint
}

var standardVars map[string]string = map[string]string{
"archive": "archive.",
"archive/tar": "tar.",
"archive/zip": "zip.",
"bufio": "bufio.",
"builtin": "builtin.",
"bytes": "bytes.",
"compress": "compress.",
"compress/bzip2": "bzip2.",
"compress/flate": "flate.",
"compress/gzip": "gzip.",
"compress/lzw": "lzw.",
"compress/testdata": "testdata.",
"compress/zlib": "zlib.",
"container": "container.",
"container/heap": "heap.",
"container/list": "list.",
"container/ring": "ring.",
"crypto": "crypto.",
"crypto/aes": "aes.",
"crypto/cipher": "cipher.",
"crypto/des": "des.",
"crypto/dsa": "dsa.",
"crypto/ecdsa": "ecdsa.",
"crypto/elliptic": "elliptic.",
"crypto/hmac": "hmac.",
"crypto/md5": "md5.",
"crypto/rand": "rand.",
"crypto/rc4": "rc4.",
"crypto/rsa": "rsa.",
"crypto/sha1": "sha1.",
"crypto/sha256": "sha256.",
"crypto/sha512": "sha512.",
"crypto/subtle": "subtle.",
"crypto/tls": "tls.",
"crypto/x509": "x509.",
"crypto/x509/pkix": "pkix.",
"database": "database.",
"database/sql": "sql.",
"database/sql/driver": "driver.",
"debug": "debug.",
"debug/dwarf": "dwarf.",
"debug/elf": "elf.",
"debug/gosym": "gosym.",
"debug/macho": "macho.",
"debug/pe": "pe.",
"encoding": "encoding.",
"encoding/ascii85": "ascii85.",
"encoding/asn1": "asn1.",
"encoding/base32": "base32.",
"encoding/base64": "base64.",
"encoding/binary": "binary.",
"encoding/csv": "csv.",
"encoding/gob": "gob.",
"encoding/hex": "hex.",
"encoding/json": "json.",
"encoding/pem": "pem.",
"encoding/xml": "xml.",
"errors": "errors.",
"expvar": "expvar.",
"flag": "flag.",
"fmt": "fmt.",
"go": "go.",
"go/ast": "ast.",
"go/build": "build.",
"go/doc": "doc.",
"go/parser": "parser.",
"go/printer": "printer.",
"go/scanner": "scanner.",
"go/token": "token.",
"hash": "hash.",
"hash/adler32": "adler32.",
"hash/crc32": "crc32.",
"hash/crc64": "crc64.",
"hash/fnv": "fnv.",
"html": "html.",
"html/template": "template.",
"image": "image.",
"image/color": "color.",
"image/draw": "draw.",
"image/gif": "gif.",
"image/jpeg": "jpeg.",
"image/png": "png.",
"image/testdata": "testdata.",
"index": "index.",
"index/suffixarray": "suffixarray.",
"io": "io.",
"io/ioutil": "ioutil.",
"log": "log.",
"log/syslog": "syslog.",
"math": "math.",
"math/big": "big.",
"math/cmplx": "cmplx.",
"math/rand": "rand.",
"mime": "mime.",
"mime/multipart": "multipart.",
"net": "net.",
"net/http": "http.",
"net/http/cgi": "cgi.",
"net/http/fcgi": "fcgi.",
"net/http/httptest": "httptest.",
"net/http/httputil": "httputil.",
"net/http/pprof": "pprof.",
"net/mail": "mail.",
"net/rpc": "rpc.",
"net/rpc/jsonrpc": "jsonrpc.",
"net/smtp": "smtp.",
"net/testdata": "testdata.",
"net/textproto": "textproto.",
"net/url": "url.",
"os": "os.",
"os/exec": "exec.",
"os/signal": "signal.",
"os/user": "user.",
"path": "path.",
"path/filepath": "filepath.",
"reflect": "reflect.",
"regexp": "regexp.",
"regexp/syntax": "syntax.",
"regexp/testdata": "testdata.",
"runtime": "runtime.",
"runtime/cgo": "cgo.",
"runtime/debug": "debug.",
"runtime/pprof": "pprof.",
"sort": "sort.",
"strconv": "strconv.",
"strings": "strings.",
"sync": "sync.",
"sync/atomic": "atomic.",
"syscall": "syscall.",
"testing": "testing.",
"testing/iotest": "iotest.",
"testing/quick": "quick.",
"text": "text.",
"text/scanner": "scanner.",
"text/tabwriter": "tabwriter.",
"text/template": "template.",
"text/template/parse": "parse.",
"time": "time.",
"unicode": "unicode.",
"unicode/utf16": "utf16.",
"unicode/utf8": "utf8.",
"unsafe": "unsafe."}

func findVar(path string, comments *ast.CommentGroup) string {
	path, err := strconv.Unquote(path)
	if err != nil {
		return ""
	}
	
	v, ok := standardVars[path]
	if ok {
		return v
	}
	return ""
}

func FilterFile(inFn villa.Path, out io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, inFn.S(), nil, parser.ParseComments)
	if err != nil {
		return err
	} // if
	
	var initFunc bytes.Buffer
	initFunc.WriteString(`
func init() {
`)

	for _, imp := range f.Imports {
		//fmt.Printf("%+v\n", imp)
		fmt.Printf("%s: ", imp.Path.Value)
		v := findVar(imp.Path.Value, imp.Comment)
		if len(v) > 0 {
			initFunc.WriteString("\t_ = " + v + "\n")
		}
		fmt.Println(v)
	} // for imp
	
	initFunc.WriteString(`}
`)
	
//	(&printer.Config{Mode: printer.RawFormat, Tabwidth: 4}).Fprint(out, fset, f)

	out.Write(initFunc.Bytes())
	return nil
}
