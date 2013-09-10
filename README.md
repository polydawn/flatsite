flatsite
========

`flatsite` generates a flat site.  Template files go in, static files come out.  You can't explain that.

flatsite is a zero-dependencies go program.  The single file `flatsite.go` plus the golang compiler and stdlib is sufficient to get shit done.  No `go get`, no submodules, nothing.  Don't even set your $GOPATH, flatsite doesn't care.  Vendor it into things, if you like.  `cp flatsite.go my_doggie_blog_generator.go` is fine by me.

flatsite is appropriate for use generating html files for sites, as well as css or js as desired.  It's equally appropriate for templating ini config files.  It's just generally appropriate.


Usage
-----

flatsite takes files from a input directory of templates (defaults to `$PWD/tmpl/`), compiles all the templates that are intended for output (defaults to files under `$PWD/tmpl/output/`), and outputs them in another directory (defaults to `$PWD/www/`).

flatsite doesn't have an args parser; configuration is largely done by the templates themselves, or else by environment variables.  Consult the source for your options.  Seriously, it's hovering around like 100 lines.  It would take me longer to explain it.


