## 1.1Pandoc 简介

Pandoc是使用Haskell语言编写的一款跨平台、自由开源及命令行界面的标记语言转换工具，可实现不同标记语言间的格式转换，堪称该领域中的“瑞士军刀”。

### 1.1.1 Pandoc 支持转换格式

#### from 格式类型：

bibtex (BibTeX bibliography)
biblatex (BibLaTeX bibliography)
bits (BITS XML, alias for jats)
commonmark (CommonMark Markdown)
commonmark_x (CommonMark Markdown with extensions)
creole (Creole 1.0)
csljson (CSL JSON bibliography)
csv (CSV table)
tsv (TSV table)
docbook (DocBook)
docx (Word docx)
dokuwiki (DokuWiki markup)
endnotexml (EndNote XML bibliography)
epub (EPUB)
fb2 (FictionBook2 e-book)
gfm (GitHub-Flavored Markdown), or the deprecated and less accurate markdown_github; use markdown_github only if you
need extensions not supported in gfm.
haddock (Haddock markup)
html (HTML)
ipynb (Jupyter notebook)
jats (JATS XML)
jira (Jira/Confluence wiki markup)
json (JSON version of native AST)
latex (LaTeX)
markdown (Pandoc’s Markdown)
markdown_mmd (MultiMarkdown)
markdown_phpextra (PHP Markdown Extra)
markdown_strict (original unextended Markdown)
mediawiki (MediaWiki markup)
man (roff man)
muse (Muse)
native (native Haskell)
odt (ODT)
opml (OPML)
org (Emacs Org mode)
ris (RIS bibliography)
rtf (Rich Text Format)
rst (reStructuredText)
t2t (txt2tags)
textile (Textile)
tikiwiki (TikiWiki markup)
twiki (TWiki markup)
typst (typst)
vimwiki (Vimwiki)
基于Lua 自定义读取

#### to 格式类型

asciidoc (modern AsciiDoc as interpreted by AsciiDoctor)
asciidoc_legacy (AsciiDoc as interpreted by asciidoc-py).
asciidoctor (deprecated synonym for asciidoc)
beamer (LaTeX beamer slide show)
bibtex (BibTeX bibliography)
biblatex (BibLaTeX bibliography)
chunkedhtml (zip archive of multiple linked HTML files)
commonmark (CommonMark Markdown)
commonmark_x (CommonMark Markdown with extensions)
context (ConTeXt)
csljson (CSL JSON bibliography)
docbook or docbook4 (DocBook 4)
docbook5 (DocBook 5)
docx (Word docx)
dokuwiki (DokuWiki markup)
epub or epub3 (EPUB v3 book)
epub2 (EPUB v2)
fb2 (FictionBook2 e-book)
gfm (GitHub-Flavored Markdown), or the deprecated and less accurate markdown_github; use markdown_github only if you
need extensions not supported in gfm.
haddock (Haddock markup)
html or html5 (HTML, i.e. HTML5/XHTML polyglot markup)
html4 (XHTML 1.0 Transitional)
icml (InDesign ICML)
ipynb (Jupyter notebook)
jats_archiving (JATS XML, Archiving and Interchange Tag Set)
jats_articleauthoring (JATS XML, Article Authoring Tag Set)
jats_publishing (JATS XML, Journal Publishing Tag Set)
jats (alias for jats_archiving)
jira (Jira/Confluence wiki markup)
json (JSON version of native AST)
latex (LaTeX)
man (roff man)
markdown (Pandoc’s Markdown)
markdown_mmd (MultiMarkdown)
markdown_phpextra (PHP Markdown Extra)
markdown_strict (original unextended Markdown)
markua (Markua)
mediawiki (MediaWiki markup)
ms (roff ms)
muse (Muse)
native (native Haskell)
odt (OpenOffice text document)
opml (OPML)
opendocument (OpenDocument)
org (Emacs Org mode)
pdf (PDF)
plain (plain text)
pptx (PowerPoint slide show)
rst (reStructuredText)
rtf (Rich Text Format)
texinfo (GNU Texinfo)
textile (Textile)
slideous (Slideous HTML and JavaScript slide show)
slidy (Slidy HTML and JavaScript slide show)
dzslides (DZSlides HTML5 + JavaScript slide show)
revealjs (reveal.js HTML5 + JavaScript slide show)
s5 (S5 HTML and JavaScript slide show)
tei (TEI Simple)
typst (typst)
xwiki (XWiki markup)
zimwiki (ZimWiki markup)
基于Lua 自定义写入

### 1.1.2 Pandoc 功能拓展

Pandoc 可以基于LaTeX、Groff ms或HTML生成PDF。
Pandoc 针对Markdown 增强语法包含：包括表格、定义列表、元数据块、脚注、引文、数学等语法。
Pandoc 包含模块设计器，它由一组读取器和一组写入器构成。读取器主要用于解析指定文本并产生文档对象。
写入器主要用于将文档对象转换为模板对象。用户基于Lua 实现自定义读取器和写入器的filter来修改AST。
