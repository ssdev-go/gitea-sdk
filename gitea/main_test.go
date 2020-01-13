// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

func getGiteaURL() string {
	return os.Getenv("GITEA_SDK_TEST_URL")
}

func getGiteaToken() string {
	return os.Getenv("GITEA_SDK_TEST_TOKEN")
}

func getGiteaUsername() string {
	return os.Getenv("GITEA_SDK_TEST_USERNAME")
}

func getGiteaPassword() string {
	return os.Getenv("GITEA_SDK_TEST_PASSWORD")
}

func enableRunGitea() bool {
	r, _ := strconv.ParseBool(os.Getenv("GITEA_SDK_TEST_RUN_GITEA"))
	return r
}

func newTestClient() *Client {
	token := getGiteaToken()
	if token == "" {
		client := NewClientWithHTTP(getGiteaURL(), &http.Client{})
		log.Printf("testing with %v, %v, %v\n", getGiteaURL(), getGiteaUsername(), getGiteaPassword())
		client.SetBasicAuth(getGiteaUsername(), getGiteaPassword())
		return client
	}
	return NewClient(getGiteaURL(), getGiteaToken())
}

func giteaMasterPath() string {
	switch runtime.GOOS {
	case "darwin":
		return fmt.Sprintf("https://dl.gitea.io/gitea/master/gitea-master-darwin-10.6-%s", runtime.GOARCH)
	case "linux":
		return fmt.Sprintf("https://dl.gitea.io/gitea/master/gitea-master-linux-%s", runtime.GOARCH)
	case "windows":
		return fmt.Sprintf("https://dl.gitea.io/gitea/master/gitea-master-windows-4.0-%s.exe", runtime.GOARCH)
	}
	return ""
}

func downGitea() (string, error) {
	for i := 3; i > 0; i-- {
		resp, err := http.Get(giteaMasterPath())
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		f, err := ioutil.TempFile(os.TempDir(), "gitea")
		if err != nil {
			continue
		}
		_, err = io.Copy(f, resp.Body)
		f.Close()
		if err != nil {
			continue
		}

		if err = os.Chmod(f.Name(), 700); err != nil {
			return "", err
		}

		return f.Name(), nil
	}

	return "", fmt.Errorf("Download gitea from %v failed", giteaMasterPath())
}

func runGitea() (*os.Process, error) {
	log.Println("Downloading Gitea from", giteaMasterPath())
	p, err := downGitea()
	if err != nil {
		log.Fatal(err)
	}

	giteaDir := filepath.Dir(p)
	cfgDir := filepath.Join(giteaDir, "custom", "conf")
	os.MkdirAll(cfgDir, os.ModePerm)
	cfg, err := os.Create(filepath.Join(cfgDir, "app.ini"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = cfg.WriteString(`[security]
INTERNAL_TOKEN = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjE1NTg4MzY4ODB9.LoKQyK5TN_0kMJFVHWUW0uDAyoGjDP6Mkup4ps2VJN4
INSTALL_LOCK   = true
SECRET_KEY     = 2crAW4UANgvLipDS6U5obRcFosjSJHQANll6MNfX7P0G3se3fKcCwwK3szPyGcbo
[database]
DB_TYPE  = sqlite3
[log]
MODE = console
LEVEL = Trace
REDIRECT_MACARON_LOG = true
MACARON = ,
ROUTER = ,`)
	cfg.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Run gitea migrate", p)
	err = exec.Command(p, "migrate").Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Run gitea admin", p)
	err = exec.Command(p, "admin", "create-user", "--username=test01", "--password=test01", "--email=test01@gitea.io", "--admin=true", "--must-change-password=false", "--access-token").Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start Gitea", p)
	return os.StartProcess(filepath.Base(p), []string{}, &os.ProcAttr{
		Dir: giteaDir,
	})
}

func TestMain(m *testing.M) {
	if enableRunGitea() {
		p, err := runGitea()
		if err != nil {
			log.Fatal(err)
			return
		}
		defer func() {
			p.Kill()
		}()
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}
