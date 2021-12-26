# HTML Includer
Develop Sites Faster with HTML Includer!

# How to Install

Install HTML Includer on your machine:

```json
go install github.com/GameWorkstore/html-includer@latest
```

# How to Use

Use the command below:

```json
./html-includer \
    absolute/path/to/source \
    absolute/path/to/destiny \
    absolute/path/to/ignore1 \ 
    absolute/path/to/ignore2 \
    ...
```

Source folder should contain the html files you want to be patched with content.
Destiny folder will be deleted, if exists.
Source folder will be copied to destiny folder recursively.
Folders ignored will be skipped entirely and left as it is.

Lines with:

<script>HtmlInclude();</script>
<script src="scripts/html-include.js"></script>

Are replaced with empty string.