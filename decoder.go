package pola

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
	"github.com/google/go-jsonnet"
	"github.com/hjson/hjson-go/v4"
	"github.com/tailscale/hujson"
)

const (
	ExtJson    = ".json"
	ExtHjson   = ".hjson"
	ExtHuJson  = ".hujson"
	ExtJwcc    = ".jwcc"
	ExtYml     = ".yml"
	ExtYaml    = ".yaml"
	ExtToml    = ".toml"
	ExtJsonnet = ".jsonnet"
	ExtXml     = ".xml"
)

var (
	ErrDecoderUnsupportedType = errors.New("Decoder, unsuported type")
)

// Decoder decodes input file/stream/data
// into golang's data type.
type Decoder interface {
	Decode(dest any) error
}

type faDecoder struct {
	fa   []fs.FS
	name string
}

// NewFsDecoder decode given file into object.
// Supported format and corresponding decoders are:
// - json: encoding/json
// - hjson: github.com/hjson/hjson-go/v4
// - hujson, jwcc: github.com/tailscale/hujson
// - yaml, yml: github.com/goccy/go-yaml
// - toml: github.com/BurntSushi/toml
// - jsonnet: github.com/google/go-jsonnet
func NewFsDecoder(name string, fa ...fs.FS) Decoder {
	return &faDecoder{fa: fa, name: name}
}

func (d *faDecoder) Decode(dest any) error {
	fa := d.fa
	if len(fa) == 0 {
		abs, err := filepath.Abs(d.name)
		if err != nil {
			return err
		}
		fa = []fs.FS{os.DirFS(filepath.Dir(abs))}
		d.name = filepath.Base(d.name)
	}

	var errs error
	ext := strings.ToLower(filepath.Ext(d.name))
	for _, f := range fa {
		rdr, err := f.Open(d.name)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		err = NewDecoder(rdr, ext).Decode(dest)
		rdr.Close()

		if err == nil {
			return nil
		}
		errs = errors.Join(errs, err)
	}

	// return last errors
	return errs
}

type rdDecoder struct {
	rdr io.Reader
	ext string
}

// NewBytesDecoder return decoder for given stream
func NewBytesDecoder(data []byte, ext string) Decoder {
	return NewDecoder(bytes.NewReader(data), ext)
}

// NewDecoder return decoder for given rider and ext type
func NewDecoder(r io.Reader, ext string) Decoder {
	return &rdDecoder{rdr: r, ext: ext}
}

func (r *rdDecoder) decodeJsonnet(dest any) error {
	data, err := io.ReadAll(r.rdr)
	if err != nil {
		return err
	}
	vm := jsonnet.MakeVM()
	jsStr, err := vm.EvaluateAnonymousSnippet(r.ext, string(data))
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsStr), dest)
}

func (r *rdDecoder) decodeHuJson(dest any) error {
	data, err := io.ReadAll(r.rdr)
	if err != nil {
		return err
	}

	stddata, err := hujson.Standardize(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(stddata, dest)
}

func (r *rdDecoder) decodeHjson(dest any) error {
	data, err := io.ReadAll(r.rdr)
	if err != nil {
		return err
	}

	return hjson.Unmarshal(data, dest)
}

func (r *rdDecoder) Decode(dest any) error {
	switch r.ext {
	case ExtJson:
		return json.NewDecoder(r.rdr).Decode(dest)
	case ExtHjson:
		return r.decodeHjson(dest)
	case ExtHuJson, ExtJwcc:
		return r.decodeHuJson(dest)
	case ExtYaml, ExtYml:
		return yaml.NewDecoder(r.rdr).Decode(dest)
	case ExtToml:
		_, err := toml.NewDecoder(r.rdr).Decode(dest)
		return err
	case ExtJsonnet:
		return r.decodeJsonnet(dest)
	case ExtXml:
		return xml.NewDecoder(r.rdr).Decode(dest)
	}
	return ErrDecoderUnsupportedType
}

// UnmarshalFs decode content specified as name to dest.
// Arg `fa` is an array of file system, in which if its not specified,
// `name` will be searched from current directory (`.`).
// If multiple fa are specified, UnmarshalFs will finish
// once its succeeded to decode the content.
func UnmarshalFs(dest any, name string, fa ...fs.FS) error {
	return NewFsDecoder(name, fa...).Decode(dest)
}
