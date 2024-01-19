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
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
)

type md struct {
	converter goldmark.Markdown
}

var markdown = &md{converter: goldmark.New()}

func (m *md) Convert(in []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := m.converter.Convert(in, &buf)
	if err != nil {
		return nil, fmt.Errorf("convert from markdown: %w", err)
	}
	return buf.Bytes(), nil
}
