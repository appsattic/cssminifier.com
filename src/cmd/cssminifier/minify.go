package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/chilts/sid"
)

func minifyFile(input string) (io.Reader, error) {
	// create a unique id to use in the filenames
	id := sid.Id()
	filename := path.Join(dir, id+".css")
	minfile := path.Join(dir, id+".min.css")
	outfile := path.Join(dir, id+".out")

	fmt.Printf("id=%s\n", id)
	fmt.Printf("filename=%s\n", filename)
	fmt.Printf("minfile=%s\n", minfile)
	fmt.Printf("outfile=%s\n", outfile)

	// write to both the file and the SHA256 hash at the same time
	fOrig, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	defer fOrig.Close()

	// create the hasher
	hash := sha256.New()

	// might as well write to both the file and the hash at the same time
	multiWriter := io.MultiWriter(fOrig, hash)

	// copy from the input to the file AND hash at the same time
	n, err := io.WriteString(multiWriter, input)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Written %d bytes to original file.\n", n)
	fOrig.Close()

	// get the hashSum
	hashSum := fmt.Sprintf("%x", hash.Sum(nil))
	fmt.Printf("hash=%s\n", hashSum)

	// ToDo: check the status of this hash in Redis

	// run `cleancss` (stdout is the minified file, and stderr has something if there's a problem)
	cmd := exec.Command("./node_modules/.bin/cleancss", filename)

	// read stdout
	stdOutStream, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	// read stderr if the command doesn't exist
	stdErrStream, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	// start the process
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// slurp in stdout
	stdOut, err := ioutil.ReadAll(stdOutStream)
	if err != nil {
		return nil, err
	}
	slurpStdOut := string(stdOut)
	fmt.Printf("stdout : %s\n", slurpStdOut)

	// slurp in stderr
	stdErr, err := ioutil.ReadAll(stdErrStream)
	if err != nil {
		return nil, err
	}
	slurpStdErr := string(stdErr)
	fmt.Printf("stderr : %s\n", slurpStdErr)

	// wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	if len(slurpStdErr) > 0 {
		fmt.Printf("Something appeared on stderr = %d\n", len(slurpStdErr))
		// replace some bits we don't want to show the user (ToDo: generate this from `dir`.)
		r1 := strings.NewReplacer("../../../../var/lib/com-cssminifier/", "", "*/", "", "/*", "")
		slurpStdErr = r1.Replace(slurpStdErr)

		// wrap the error in a CSS comment at the top of the minified CSS so the user will notice it
		slurpStdErr = "/*\n\n" + slurpStdErr + "\n*/\n\n"
		fmt.Printf("stderr : %s\n", slurpStdErr)
	}

	// open the outfile file which has both the stderr AND the minified CSS
	min, err := os.Create(minfile)
	if err != nil {
		return nil, err
	}

	// read both the (modified) stderr and stdout in succession
	r := io.MultiReader(strings.NewReader(slurpStdErr), strings.NewReader(slurpStdOut))

	// ToDo: separate the writing to the file and the writing to the request (ie. don't use the TeeReader below,
	// because I'm not sure if the writing to the file stops if the connection is closed and that's also
	// closed). Dunno. Let's just make it more explicit.

	// finally, write to this file at the same time as the response
	return io.TeeReader(r, min), nil
}
