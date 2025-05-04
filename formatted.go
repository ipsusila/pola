package pola

import "path/filepath"

// FormattedText has Decoder and Stringer interface
type FormattedText interface {
	Decoder
	String() string
	Ext() string
}

type JsonText []byte

func (j JsonText) Decode(dest any) error {
	return NewBytesDecoder(j, ExtJson).Decode(dest)
}
func (j JsonText) String() string {
	return string(j)
}
func (j JsonText) Ext() string {
	return ExtJson
}

type YamlText []byte

func (y YamlText) Decode(dest any) error {
	return NewBytesDecoder(y, ExtYaml).Decode(dest)
}
func (y YamlText) String() string {
	return string(y)
}
func (y YamlText) Ext() string {
	return ExtYaml
}

type HjsonText []byte

func (h HjsonText) Decode(dest any) error {
	return NewBytesDecoder(h, ExtHjson).Decode(dest)
}
func (h HjsonText) String() string {
	return string(h)
}
func (h HjsonText) Ext() string {
	return ExtHjson
}

type JwccText []byte

func (j JwccText) Decode(dest any) error {
	return NewBytesDecoder(j, ExtJwcc).Decode(dest)
}
func (j JwccText) String() string {
	return string(j)
}
func (j JwccText) Ext() string {
	return ExtJwcc
}

type HuJsonText []byte

func (h HuJsonText) Decode(dest any) error {
	return NewBytesDecoder(h, ExtHuJson).Decode(dest)
}
func (h HuJsonText) String() string {
	return string(h)
}
func (h HuJsonText) Ext() string {
	return ExtHuJson
}

type JsonnetText []byte

func (j JsonnetText) Decode(dest any) error {
	return NewBytesDecoder(j, ExtJsonnet).Decode(dest)
}
func (j JsonnetText) String() string {
	return string(j)
}
func (j JsonnetText) Ext() string {
	return ExtJsonnet
}

type TomlText []byte

func (t TomlText) Decode(dest any) error {
	return NewBytesDecoder(t, ExtToml).Decode(dest)
}
func (t TomlText) String() string {
	return string(t)
}
func (t TomlText) Ext() string {
	return ExtToml
}

type XmlText []byte

func (x XmlText) Decode(dest any) error {
	return NewBytesDecoder(x, ExtXml).Decode(dest)
}
func (x XmlText) String() string {
	return string(x)
}
func (x XmlText) Ext() string {
	return ExtXml
}

// FormattedTextFile holds the filename for formatted file content.
// Filename extension is use to determined the content type.
type FormattedTextFile string

func (f FormattedTextFile) Decode(dest any) error {
	return NewFsDecoder(string(f)).Decode(dest)
}
func (f FormattedTextFile) String() string {
	return string(f)
}
func (f FormattedTextFile) Ext() string {
	return filepath.Ext(string(f))
}
