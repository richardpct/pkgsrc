// Package pkgsrc for building macOS packages from source
package pkgsrc

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
)

// Workdir directory for building packages source
const Workdir = "/tmp/build"

// Pkg definition
type Pkg struct {
	name       string
	vers       string
	ext        string
	PkgName    string
	pkgNameExt string
	url        string
	hashType   string
	hash       string
}

// Init function
func (p *Pkg) Init(name, vers, ext, url, hashType, hash string) {
	p.name = name
	p.vers = vers
	p.ext = ext
	p.url = url
	p.hashType = hashType
	p.hash = hash
	p.PkgName = name + "-" + vers
	p.pkgNameExt = p.PkgName + "." + ext
}

// CleanWorkdir function for remove existing workdir
// TODO Delete if destination file exists
func (p *Pkg) CleanWorkdir() {
	wdPkgName := path.Join(Workdir, p.PkgName)

	if _, err := os.Stat(wdPkgName); err == nil {
		if err := os.RemoveAll(wdPkgName); err != nil {
			log.Fatal(err)
		}
	} else if _, err := os.Stat(Workdir); err != nil {
		os.Mkdir(Workdir, 0755)
	}
}

// CheckSum function for checking hash
// TODO Other than sha256 algorithm
func (p *Pkg) CheckSum() bool {
	wdPkgNameExt := path.Join(Workdir, p.pkgNameExt)

	f, err := os.Open(wdPkgNameExt)
	if err != nil {
		return false
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	if hashStr := hex.EncodeToString(h.Sum(nil)); hashStr == p.hash {
		return true
	}

	return false
}

// DownloadPkg function for getting package source
func (p *Pkg) DownloadPkg() {
	wdPkgNameExt := path.Join(Workdir, p.pkgNameExt)
	out, err := os.Create(wdPkgNameExt)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// TODO Using join
	resp, err := http.Get(p.url + "/" + p.pkgNameExt)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		log.Fatal(err)
	}
}

// Unpack function for unpacking source file
func (p *Pkg) Unpack() {
	if err := os.Chdir(Workdir); err != nil {
		log.Fatal(err)
	}
	wdPkgNameExt := path.Join(Workdir, p.pkgNameExt)
	cmd := exec.Command("/usr/bin/tar", "xzvf", wdPkgNameExt)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
