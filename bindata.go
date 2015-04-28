package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _www_app_js_swp = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x9a\xcf\x6b\x1b\x47\x18\x86\x47\xbd\x6d\x6b\xd7\x6e\x4b\x0f\xa5\xa5\x4c\xe5\x83\x65\xb0\xb4\xae\xe5\xe2\x52\xaf\x16\x5a\xf5\x07\x86\xba\x18\x53\x9b\xd2\x43\x61\xb4\x3b\x92\xd6\x5a\xed\x2c\xb3\x23\xcb\xc6\x55\x4b\x2f\x3d\xf4\xd6\x5e\x0b\x3d\xb5\x87\x52\x7a\x48\x20\x90\x80\xff\x83\xdc\x72\x48\x08\x39\xf8\x92\x4b\x72\x30\x24\x97\x5c\x72\xc8\xbb\x5a\x39\x11\x8e\x63\x1b\x3b\xc4\x04\xbe\x07\x1e\xcd\xec\xec\xec\xb7\xfb\xce\xae\x40\x42\xaa\xcd\xac\x2d\x2e\xf1\xf9\xd2\x1c\x03\xe3\x8c\x6d\xae\xbb\xab\x6b\x63\x9b\xb9\xf9\x77\x19\x13\xb1\x6a\x89\x70\x2b\x36\x41\x8b\x1d\x45\x24\xbb\xa5\xa1\xb9\x25\x4f\xb5\x0f\x9d\xf7\xd3\xd0\x24\x5b\x8b\xc0\xaf\x29\x63\x77\xbb\x5d\x5b\xc4\x71\x69\x3d\x39\xf2\x1c\x04\x41\x9c\x81\x8e\xa9\x17\x3f\x1e\x65\xe5\xd9\x0f\x67\xd2\xcd\x89\xfc\x07\xfc\xad\x37\x57\xcf\xfb\xaa\x08\x82\x20\x08\x82\x20\x08\x82\x78\x81\x98\x38\xc7\x7e\x46\xfb\xca\x60\xbb\x3c\x68\x73\x07\x5a\x82\x20\x08\x82\x20\x08\x82\x20\x08\x82\x20\x5e\x5e\x84\xcf\xd8\xbd\x57\x19\xbb\xf9\x1a\xeb\xff\xfe\xbf\xff\xfd\xff\xca\x18\x63\x97\xe1\x1f\xf0\x47\xa8\xe1\xb7\xd0\x85\x0f\x5f\x67\x6c\x07\x5e\x84\xff\xc3\x7f\xe1\x3f\xf0\x6f\xf8\x0b\x6c\x41\x01\xab\x70\x06\xee\x8d\x32\x76\x1b\xee\xc2\x5b\xf0\x06\xbc\x0e\x2f\xc0\x3f\xe1\xef\x30\x81\x3f\xc0\x1c\x7c\x30\xc2\xd8\x7d\xb8\x07\xef\xc2\x3b\xf0\x2a\xbc\x04\xff\x83\xbf\xc2\x10\x7e\x07\xbf\x82\x65\x58\x80\xef\xc3\x77\xe0\xdb\xf0\x0d\x38\x3e\x92\xe5\x7a\x0f\xed\xe7\x70\x17\x5e\x83\x3b\xf0\x2f\xf8\x1b\xec\xc1\x18\x7e\x0f\x9d\x91\xf3\xbc\x13\x04\x41\x10\x04\x41\x10\x04\x71\x06\x56\xa4\xf0\x4c\x49\xcb\xc8\x97\xba\xe0\x7c\x1a\xc7\xdc\x76\xa7\xb9\xaf\xbc\x4e\x5b\x46\xa6\xd4\x90\xe6\x8b\x50\xa6\xdd\xcf\xb6\x16\xfd\xc2\xa4\x88\xe3\xc9\xa9\xa9\x05\x7c\x20\xc6\x8b\xd5\x9b\x66\x96\x95\x76\x2c\xcb\xb1\xfd\x60\xc3\x4d\x7b\xc3\x5d\xcb\x59\x92\xed\x9a\xd4\x5f\x07\x89\xb1\x5d\x6b\x30\xb6\x22\x02\x3f\x1b\x19\x0c\x54\x9b\x22\x8a\x64\x38\x34\xe6\xa0\x04\xf7\x42\x91\x24\xdf\x88\xb6\xac\xe4\xb5\xea\xe6\xfb\x7b\x0e\xee\xf0\x54\x64\x44\x10\x49\x5d\xac\x87\x9d\xc0\xef\x4f\xd2\xd2\x74\x74\x54\x60\x56\x16\xeb\x13\x5e\xef\x44\x9e\x09\x54\x54\x98\xe2\xdb\x6c\x43\x68\x9e\xe6\xac\xf0\x2c\xbb\xa7\xa5\x30\xb2\x9a\x96\x2c\x6c\x1f\x1d\x2c\x31\x5a\x45\x0d\x77\x39\x94\x22\x91\x3c\x91\xa1\xf4\x0c\x17\x3c\xfd\x8f\x32\x37\x0a\x03\x92\x9b\xa6\xe4\xed\x7e\x66\x1e\x22\x0e\x17\x91\xcf\x6b\x92\x8b\x5a\x28\xd3\x29\xeb\x2a\x88\xb8\xd2\x3c\x16\xda\x38\xf6\xa0\x5e\x56\xbc\x39\xeb\x66\x8b\x95\x38\x36\xfa\x87\x87\x0d\x8b\x6d\xbf\x38\x77\xc2\x94\x4f\xd6\xfe\x39\x86\xf5\xb2\x7b\x35\x9c\xb7\x9f\x3f\x4d\xfb\x74\xa2\xf4\x56\x1f\x9b\xe7\xa3\x13\xe6\xd9\x7f\x6e\x4e\x93\x46\xf0\xa6\x96\xf5\x4a\x7e\x22\xcf\x55\x54\x0d\x03\xaf\x55\xd9\xce\xce\xf8\xa5\x08\x13\xd9\x73\x7d\x99\x98\x20\xda\x2a\xa6\x61\x8a\x1b\xaa\xe1\xd8\xc2\x75\x6a\xda\x3e\xd5\xf1\x9e\x56\x46\x1c\xa8\x80\x35\x18\x3c\xe7\xc7\x2e\x48\xf9\x84\x0b\x32\xf4\xbe\x79\xd6\x9a\xa4\xd3\x86\xae\x13\xd3\x1e\x57\xe1\x92\xa3\x10\x97\xa5\xc4\xa8\x78\x59\xab\x58\x34\x44\x56\x7e\x61\x70\x08\xaf\xa7\xc7\x2c\xf0\x1e\x7b\x14\x00\x00\xff\xff\x13\x3d\xae\x34\x00\x30\x00\x00")

func www_app_js_swp_bytes() ([]byte, error) {
	return bindata_read(
		_www_app_js_swp,
		"www/.app.js.swp",
	)
}

func www_app_js_swp() (*asset, error) {
	bytes, err := www_app_js_swp_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "www/.app.js.swp", size: 12288, mode: os.FileMode(420), modTime: time.Unix(1430153856, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _www_index_html_swp = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x9a\x4d\x6f\xe4\x34\x18\xc7\x5d\x6e\x43\x59\x16\xc1\x65\x25\x10\x72\xc3\xa5\x95\x36\xf6\xbc\x74\x55\x60\x27\xbd\x2c\x3d\x14\x09\x81\x44\x8b\x40\x55\x85\x3c\xb1\xd3\xf1\x34\xb1\x83\xe3\x99\xcc\xa8\x2d\x5c\xe0\xce\x77\xe0\x00\x9f\x82\x63\xaf\x1c\xe1\x86\xc4\x17\x80\x2b\xe2\xc2\x93\x64\x68\x87\xbe\x09\xca\x4b\xb5\xd2\xf3\x53\xff\xb2\x1f\xbf\xfc\xed\xc7\xb1\x22\x55\x93\x41\xfb\x83\xed\x77\xe8\x06\x5b\x27\xc0\x0b\x84\x3c\x38\xd8\xdc\xfd\xfc\xfe\x74\x69\xe3\x65\x42\x44\x6e\x0f\x45\x3a\xcb\xbd\x3e\x24\x37\x61\x54\xc9\x16\xc6\xb2\xd8\x66\x57\x8e\xfb\x74\x61\x10\x77\x42\xcb\x81\xf5\xbc\x2c\x4b\xae\x8d\x54\x53\x36\xf4\x59\x7a\xe3\x3a\x08\x82\xdc\x92\xb1\x4f\xc2\xd7\xef\x91\x5e\xb7\xd3\xae\xc2\xd7\x82\x15\xfa\xd2\x8b\xbb\x77\xbd\x2b\x04\x41\x10\x04\x41\x10\x04\x41\xfe\x47\x7c\xbe\x44\x3e\x83\xf2\x99\x79\xfc\xca\xbc\x5c\xba\x50\x22\x08\x82\x20\x08\x82\x20\x08\x82\x20\x08\xf2\xf4\x22\x24\x21\x3f\x3c\x4b\x88\x5c\x26\xf5\xef\xff\x7f\xfc\xff\xff\xcb\x7d\x42\x7e\x04\x7d\x0f\xfa\x12\x74\x0c\xca\x41\x1f\x83\x9e\x80\x7e\x7d\x9e\x90\x9f\x41\x3f\x81\xbe\x00\xbd\x0a\x3a\xbd\x47\x88\x05\x3d\x00\xb5\x40\xdf\x3c\x47\xc8\x04\xf4\x18\xf4\x1b\xf8\x7f\x07\x3a\x05\x7d\x0b\xfa\x1a\x34\x06\xa5\xcb\xcd\xda\x5f\x2d\xdf\xe5\x29\x20\x08\x82\x20\x08\x82\x20\xc8\x6d\xe8\x5f\x0d\xaf\x3e\x6d\xdd\x24\xad\x3e\x1f\x58\x39\x83\x4a\xab\x5f\xc4\x4e\xe7\x9e\xfa\x59\xae\xa2\xc0\xab\xa9\xe7\xa3\x62\x1a\xd0\xc2\xc5\x51\xc0\x45\x9e\xb3\x51\x11\x6c\xf6\x79\x33\xac\x9e\x21\xf5\x84\x6a\x19\x05\xd0\x59\xf5\x40\x58\x39\xce\x0d\x61\x09\x25\x64\x3d\x2e\xd5\xe6\x90\x3a\x95\x46\x41\xe1\x67\xa9\x2a\x86\x4a\xf9\x80\x0e\x9d\x4a\xc0\xb9\x6e\x62\x71\x01\xe6\x0b\xbb\x68\x56\xe5\xb1\x95\x8a\x8d\x3e\x19\x2b\x37\xab\xbe\xe1\xe5\x4d\x35\xec\xb2\x0e\xeb\xb1\x4c\x9b\xcb\x7b\xfa\xf3\xfc\x64\xc0\x32\xc5\xdf\x7e\xff\xc3\x1d\x27\x4c\x91\x58\x97\x29\x17\xb6\x59\xa7\xc7\xba\x7f\x6d\xaa\x53\x22\xf6\xd7\xce\xb8\x31\x31\x9e\x89\x69\x2c\x0d\x1b\x58\xeb\x0b\xef\x44\x5e\x05\x55\x12\x67\x0d\xbc\x07\x59\xac\x73\xc8\xfd\xbc\xad\xce\xea\xec\x34\x56\xf6\x94\x91\x3a\xd9\x0f\xc3\x2a\xbc\x74\x3a\xd2\x8c\x0a\x16\xa7\x76\x2c\x93\x54\x38\x55\xbb\x8b\x91\x98\xf2\x54\x0f\x0a\xae\x8a\x47\x61\x31\xd4\x19\x2c\xb3\xce\xda\xf3\x50\x64\x97\xd2\xf8\x37\x6c\xf5\x3f\xb6\x6d\x1e\x2d\xef\xb0\x0e\xfc\xcd\xa3\x2b\x1f\xf1\xdf\xf4\xad\x2e\xfa\x23\xd8\xe0\x04\xf6\xbb\xc1\xba\xe7\x71\x98\x3b\x6d\x7c\x55\xfb\x4f\x97\xb9\xc6\x7c\x31\x68\x9d\xac\xae\xad\x3d\xae\x6b\xad\x52\x1b\x69\x4b\x30\x37\x85\x4d\x15\x8d\xe8\x85\x86\xe3\x63\x7a\x94\xda\x83\x37\x55\xf2\xb0\x14\xce\x54\xa5\x72\xce\xba\xaa\x22\x75\x55\x9c\xcc\x9d\x26\xc2\x51\x95\x80\x43\x32\x36\xb1\xd7\xd6\xac\xae\x1d\xcd\xfb\x56\x17\x9a\x16\x52\x6d\xae\x5c\x18\xee\xe9\x84\xa6\x9e\x6e\x6f\xd1\x37\xf6\xeb\xb6\x4c\x79\x41\x8d\xc8\xe0\xc5\x30\xd1\xaa\xcc\xad\x83\x6b\x0e\x3b\xf2\xca\xf8\x28\x28\xb5\xf4\xc3\x48\xaa\x89\x8e\x55\x58\x07\x0f\xa9\x36\xda\x6b\x91\x86\x45\x2c\x52\x15\x75\xaa\xdb\xdc\x6f\xde\x07\xfd\xea\x64\x68\x2a\xcc\x41\x14\x28\x03\x1d\xfd\x95\xb7\xde\x7d\xb2\xf3\xd1\x7b\x5b\xb4\x79\x27\xfd\x1e\x00\x00\xff\xff\x2d\x57\xa6\xdf\x00\x30\x00\x00")

func www_index_html_swp_bytes() ([]byte, error) {
	return bindata_read(
		_www_index_html_swp,
		"www/.index.html.swp",
	)
}

func www_index_html_swp() (*asset, error) {
	bytes, err := www_index_html_swp_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "www/.index.html.swp", size: 12288, mode: os.FileMode(420), modTime: time.Unix(1430152995, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _www_app_js = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x93\xcf\x8a\xdb\x30\x10\xc6\xcf\xf6\x53\x0c\xe9\x61\x1d\x58\x47\xd0\x3f\x97\x46\x31\xb4\xa1\x85\x42\x5b\x96\x7d\x83\x89\x34\x49\xd4\xca\x92\x19\x29\x29\x4b\xc8\xbb\x57\xb2\xe3\xae\x5b\x4a\x59\xca\xe6\x64\x79\x66\xbe\x99\xf9\x7d\xb2\x8f\xc8\xc0\x14\x0f\xec\x3e\xa2\x0d\x04\x2b\xd8\x1e\x9c\x8a\xc6\xbb\x0a\x08\xe6\x70\x02\x5a\x84\xe8\xbb\x3b\xf6\x1d\xee\xb0\x4f\xcc\x97\x17\x09\x6c\xb3\x66\x09\xe7\xb2\xcc\x7d\xd6\x7b\x74\x8e\xec\x67\x13\x62\xea\x73\x4f\xa8\xe2\x42\x31\x61\xa4\xb5\xc5\x10\xaa\x53\x59\x30\x39\x4d\xfc\xf6\x71\x48\x9a\x50\x16\xc5\xd0\xae\x4a\xa7\x42\x6a\x73\x04\x95\xeb\xbf\x62\x4b\xab\x99\xf2\xb6\x6e\x75\xfd\x6a\xd6\xe4\x6c\x21\xf7\x2f\x9b\xcb\x9c\x20\x45\x7a\x19\xa2\x08\x7b\xa6\xed\x6a\xf6\x62\x06\xde\xad\xad\x51\xdf\x57\xa7\x09\xd6\xb9\xd1\x14\xa2\x71\x0f\x35\xa3\xd1\xb5\x62\x1f\x51\x0a\x6c\xe4\x86\xc5\x7f\x75\x38\xfa\xdd\x6f\x7a\x29\xd2\xda\xf9\x34\x5f\x96\xc5\xf9\xb6\x3c\xa7\x67\xef\xc9\x7d\xaa\xbe\x8a\x21\x6f\x26\x86\xe4\x21\x53\x37\x42\x64\xef\x76\xcd\x9d\x25\x4c\x57\x1a\xc8\x92\x8a\x80\xa0\x06\xdf\x20\xfa\x14\x23\x88\x7b\x82\x0c\x03\x36\xed\x27\xc5\x45\xf4\x4f\x9a\x2f\xd4\x6e\x88\xaf\xc2\xf3\x7a\xc2\x33\x8c\x79\x02\x51\xbf\xfe\x04\xa7\xed\x85\x3d\x10\xa0\xd3\xb0\x21\xc0\x8d\xa5\x5c\xf2\xcd\x1b\x07\x9e\xa1\x43\x7e\x22\xec\xbb\xae\x7b\x4e\x4a\x17\xd1\x38\xe2\x7a\x6b\x0f\x46\x8f\xb0\x7f\x54\xb1\xff\x71\xc9\x14\x72\xf2\x37\x89\x31\x36\x7e\x4d\xbf\x02\x8f\x17\x22\x9a\x62\x68\x39\xe2\xfc\x1d\x6c\xc0\x19\xf6\xaf\x64\x46\x14\xcd\x2d\x68\xaf\x0e\x2d\xb9\xb8\xd8\x51\xfc\x60\x29\x1f\xdf\x3f\x7c\xd2\xd5\x0d\x76\xdd\xcd\x3c\xe9\x7e\x06\x00\x00\xff\xff\x6b\xd3\x4e\xa0\x27\x04\x00\x00")

func www_app_js_bytes() ([]byte, error) {
	return bindata_read(
		_www_app_js,
		"www/app.js",
	)
}

func www_app_js() (*asset, error) {
	bytes, err := www_app_js_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "www/app.js", size: 1063, mode: os.FileMode(420), modTime: time.Unix(1430153848, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _www_index_html = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x54\x4f\x6f\x1b\x2f\x14\x3c\x3b\x9f\x02\x73\xb2\xa5\x00\xb1\xf3\x8b\x7e\x6a\xcb\xfa\xd2\xe6\xd0\x5e\x5a\xa9\x39\xb4\x8a\x72\xc0\xf0\x36\x8b\xcb\xc2\x16\xf0\xae\xad\x24\xdf\xbd\xb0\xbb\x4d\xdc\x24\xaa\x7a\x48\xa4\x95\xde\x9b\x59\x98\x19\xfe\x08\x3e\xfd\xf0\xf9\xfd\xc5\xf7\x2f\xe7\xa8\x8a\xb5\x59\x1d\xf1\x5c\x90\x11\xf6\xba\xc0\x60\xf1\xea\x68\xc2\x2b\x10\x2a\xd5\x09\xaf\x21\x0a\x64\x45\x0d\x05\x6e\x35\x74\x8d\xf3\x11\x23\xe9\x6c\x04\x1b\x0b\xdc\x69\x15\xab\x42\x41\xab\x25\x90\x1e\x1c\x23\x6d\x75\xd4\xc2\x90\x20\x85\x81\x62\x81\x7b\x99\x29\x21\x97\xba\x44\x26\xa2\x8f\xe7\xe8\xcd\x55\xe6\x26\x3c\x48\xaf\x9b\xd8\xf7\x93\x59\xb9\xb5\x32\x6a\x67\x67\xf3\x9b\x9e\x98\xb4\xc2\x23\x28\x51\x81\x0e\xfe\xdc\xbd\x1b\xfe\x75\xda\x2a\xd7\xd1\x94\x23\x38\x03\x69\xcc\x23\xe2\xf6\x16\xdd\x18\x77\xfd\x16\xca\xe3\x4e\x78\x9b\x2b\x78\xef\x7c\x6e\x94\xce\x65\x54\xba\x9b\xcd\xe7\x7d\xc7\xd9\x41\x9a\x31\x19\x0a\x5e\x16\x98\x31\xa9\xec\x26\x50\x69\xdc\x56\x95\x46\x78\x48\x36\x35\x13\x1b\xb1\x63\x46\xaf\x03\xcb\xbb\x77\x16\x2a\xdd\xb2\x53\xfa\x3f\x5d\x3e\x60\x5a\x6b\x4b\x37\x01\xaf\x5e\x43\x9c\x34\x5e\xdb\xf8\x52\x36\x9b\x9f\x5b\xf0\x7b\xb6\xa0\x8b\xf4\x8d\xe8\x25\x74\x21\x9c\x91\x94\xb1\x4e\xe9\xff\xa3\x27\xf7\xf0\x55\x64\xc5\x53\x59\x3e\xbd\x04\xab\x74\x79\x45\x48\x0f\x8d\xb6\x3f\x90\x07\x53\xe0\x10\xf7\x06\x42\x05\x90\xae\x73\xe5\xa1\xcc\xa6\xb5\xd8\x25\x5f\xba\x76\x2e\x86\xe8\x45\x93\x41\xf6\xbd\x27\x92\x5d\x32\x64\x32\x84\x07\xae\xdf\xa5\xc4\x0c\xf7\xfc\xcf\x55\x94\x6b\x5a\x03\xf3\x20\x64\x24\x27\x74\x71\x4a\x97\x4f\x13\x3e\x37\xe3\xd3\xd7\x6f\x17\x5e\xd8\x50\x3a\x5f\x83\xff\xc7\xa9\xd2\x29\xa0\xe3\xc9\xe5\xd4\x43\x4b\x96\x74\x91\x42\x3f\x77\x94\x7f\xdf\x8d\x9e\xfa\xbd\x30\xce\xc6\x07\x81\xaf\x9d\xda\xf7\x93\x95\x6e\x91\x56\x05\x16\x4d\x93\x65\x13\x3c\xcc\x14\xf7\x4d\x7a\x31\x22\xec\x22\xdb\x84\x1d\x1e\x33\xa6\xb1\x8f\x53\x70\x36\x28\x72\x36\xbc\x45\xbf\x02\x00\x00\xff\xff\xaf\x4d\xa8\xde\x9c\x04\x00\x00")

func www_index_html_bytes() ([]byte, error) {
	return bindata_read(
		_www_index_html,
		"www/index.html",
	)
}

func www_index_html() (*asset, error) {
	bytes, err := www_index_html_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "www/index.html", size: 1180, mode: os.FileMode(420), modTime: time.Unix(1430152984, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _www_style_css = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x4a\xca\x4f\xa9\x54\xa8\xe6\x52\x50\x48\x4a\x4c\xce\x4e\x2f\xca\x2f\xcd\x4b\xd1\x4d\xce\xcf\xc9\x2f\xb2\x52\x50\x4e\x33\x04\x41\x6b\x85\x5a\x2e\xae\x0c\x23\x5c\x8a\x92\x8d\x40\xd0\x1a\x28\x59\x90\x98\x92\x92\x99\x97\x6e\xa5\x60\xa0\x67\x64\x9a\x9a\x0b\xd2\x07\x08\x00\x00\xff\xff\x7a\x0a\xd1\x13\x5e\x00\x00\x00")

func www_style_css_bytes() ([]byte, error) {
	return bindata_read(
		_www_style_css,
		"www/style.css",
	)
}

func www_style_css() (*asset, error) {
	bytes, err := www_style_css_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "www/style.css", size: 94, mode: os.FileMode(420), modTime: time.Unix(1430155986, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"www/.app.js.swp": www_app_js_swp,
	"www/.index.html.swp": www_index_html_swp,
	"www/app.js": www_app_js,
	"www/index.html": www_index_html,
	"www/style.css": www_style_css,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"www": &_bintree_t{nil, map[string]*_bintree_t{
		".app.js.swp": &_bintree_t{www_app_js_swp, map[string]*_bintree_t{
		}},
		".index.html.swp": &_bintree_t{www_index_html_swp, map[string]*_bintree_t{
		}},
		"app.js": &_bintree_t{www_app_js, map[string]*_bintree_t{
		}},
		"index.html": &_bintree_t{www_index_html, map[string]*_bintree_t{
		}},
		"style.css": &_bintree_t{www_style_css, map[string]*_bintree_t{
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

