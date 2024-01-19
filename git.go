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
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type log struct {
	iter object.CommitIter
}

type gitCommit struct {
	Hash  string
	Title string
	Body  string
	Date  time.Time
}

func gitLog(repositoryPath string) (*log, error) {
	repo, err := git.PlainOpenWithOptions(repositoryPath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, fmt.Errorf("open git repository: %w", err)
	}
	commitIter, err := repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, fmt.Errorf("log commits: %w", err)
	}
	return &log{iter: commitIter}, nil
}

func (l *log) ForEach(fn func(c *gitCommit) error) error {
	return l.iter.ForEach(func(gc *object.Commit) error {
		var title, body string
		message := strings.SplitN(gc.Message, "\n\n", 2)
		if len(message) > 0 {
			title = message[0]
		}
		if len(message) > 1 {
			body = message[1]
		}
		c := &gitCommit{
			Hash:  gc.Hash.String(),
			Title: title,
			Body:  body,
			Date:  gc.Author.When,
		}
		return fn(c)
	})
}
