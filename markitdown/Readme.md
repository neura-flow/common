
# github 地址
https://github.com/microsoft/markitdown

# 安装
pip install 'markitdown[all]'

安装markitdown可能遇到speechrecognition安装不成功情况,单独安装
1.pip install speechrecognition
2.pip install pocketsphinx
3.可能缺少pocketsphinx的wheel ,这可以去https://www.lfd.uci.edu/~gohlke/pythonlibs/下载对应的wheel文件安装，再执行　pip install pocketsphinx

# 针对部分文件安装
pip install markitdown[pdf, docx, pptx, xlsx, xls, outlook, az-doc-intel]