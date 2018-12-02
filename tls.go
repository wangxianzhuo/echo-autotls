package autotls

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type AutoTLSManager struct {
	autocert.Manager
	DisableHTTP2 bool
}

func (m AutoTLSManager) StartAutoTLS(address string) *http.Server {
	s := new(http.Server)
	s.Addr = address
	s.TLSConfig = new(tls.Config)
	s.TLSConfig.GetCertificate = m.GetCertificate
	if !m.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}
	s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, acme.ALPNProto)
	return s
}

func DefaultManager(domains ...string) *AutoTLSManager {
	m := AutoTLSManager{}
	m.Prompt = autocert.AcceptTOS
	if len(domains) > 0 {
		m.HostPolicy = autocert.HostWhitelist(domains...)
	}
	dir := cacheDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Printf("warning: autocert.NewListener not using a cache: %v", err)
	} else {
		m.Cache = autocert.DirCache(dir)
	}
	return &m
}

func cacheDir() string {
	const base = "golang-autocert"
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir(), "Library", "Caches", base)
	case "windows":
		for _, ev := range []string{"APPDATA", "CSIDL_APPDATA", "TEMP", "TMP"} {
			if v := os.Getenv(ev); v != "" {
				return filepath.Join(v, base)
			}
		}
		// Worst case:
		return filepath.Join(homeDir(), base)
	}
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, base)
	}
	return filepath.Join(homeDir(), ".cache", base)
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "/"
}
