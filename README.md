Bloggit
=======

Look, this was a wacky experiment.
If you want to use this...why?

The idea is that you write blog posts as git commits.
Want to edit a post?
Amend the commit.
Simple.

The Process
-----------

Bloggit reads the commit history from the git repository that it is run from.

For each commit that has a "title" and "body" in the format:
```text
This is a title

This is the body.
It can be many lines.
```
The "body" gets parsed as Markdown into HTML.
Then the following fields get passed to a template file (which MUST be called `commit.tmpl`):
- `{{.Title}}`
- `{{.Body}}`
- `{{.Date}}` (this is the commit date)
- `{{.Filename}}` (see below)

Filenames of output files are the commit hash, with the `.html` extension appended.

Once the commits have all be parsed and the output files generated an index file is created from a separate template file (which MUST be called `index.tmpl`).
This file is passed a slice of objects that follow the commit template style outlined above.
