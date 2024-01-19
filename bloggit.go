// Copyright (c) Alisdair MacLeod <copying@alisdairmacleod.co.uk>
//
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
// REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
// INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
// LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
// OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
// PERFORMANCE OF THIS SOFTWARE.

package main

import (
	"fmt"
	"html/template"
	"os"
	"time"
)

type post struct {
	Filename string
	Title    string
	Body     template.HTML
	Date     time.Time
}

type commitErr struct {
	Hash string
	Err  error
}

func (err *commitErr) Error() string {
	return fmt.Sprintf("%s: %v", err.Hash, err.Err.Error())
}

func (err *commitErr) Unwrap() error {
	return err.Err
}

func main() {
	if err := run(); err != nil {
		_, printErr := fmt.Fprintf(os.Stderr, "Bloggit encountered an error: %v", err)
		if printErr != nil {
			panic(printErr)
		}
	}
}

func run() (err error) {
	templates, err := template.ParseFiles("commit.tmpl", "index.tmpl")
	if err != nil {
		return fmt.Errorf("parse templates: %w", err)
	}
	log, err := gitLog(".")
	if err != nil {
		return fmt.Errorf("git log: %w", err)
	}
	var commits []post
	err = log.ForEach(func(c *gitCommit) (err error) {
		// Ignore commits without a body
		if c.Body == "" {
			return nil
		}
		filename := c.Hash + ".html"
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("open commit post file: %w", &commitErr{
				Hash: c.Hash,
				Err:  err,
			})
		}
		defer func() {
			err = f.Close()
		}()
		htmlBody, err := markdown.Convert([]byte(c.Body))
		if err != nil {
			return fmt.Errorf("convert commit body from markdown: %w", &commitErr{
				Hash: c.Hash,
				Err:  err,
			})
		}
		commitData := post{
			Filename: filename,
			Title:    c.Title,
			Body:     template.HTML(htmlBody),
			Date:     c.Date,
		}
		err = templates.ExecuteTemplate(f, "commit.tmpl", commitData)
		if err != nil {
			return fmt.Errorf("write commit post file: %w", &commitErr{
				Hash: c.Hash,
				Err:  err,
			})
		}
		commits = append(commits, commitData)
		return nil
	})
	if err != nil {
		return fmt.Errorf("iterate git log: %w", err)
	}
	f, err := os.Create("index.html")
	if err != nil {
		return fmt.Errorf("open index file: %w", err)
	}
	defer func() {
		err = f.Close()
	}()
	err = templates.ExecuteTemplate(f, "index.tmpl", commits)
	if err != nil {
		return fmt.Errorf("write index file: %w", err)
	}
	return nil
}
