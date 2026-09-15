# -*- coding: utf-8 -*-
"""把 design/ 下对比稿引用的外部 style.css 内联掉，产出可单独分享的单文件。

为什么要内联：
对比稿里写的是 <link rel="stylesheet" href="../frontend/src/style.css">，
headless 用 file:// 打开时正常，但经「预览面板」打开时服务根通常是 HTML 所在目录，
../frontend/... 会落到目录之外 → 404 → 页面整个没有样式。

用法：
  python inline_design.py                      # 处理 design/ 下所有 *.html
  python inline_design.py design/ui-x.html     # 只处理指定文件

产出：design/preview/<原名>.standalone.html（原文件保持 <link> 不变，方便改稿）
"""
import glob
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
CSS = os.path.join(ROOT, 'frontend', 'src', 'style.css')
PREVIEW = os.path.join(ROOT, 'design', 'preview')

# 对比稿专用补丁：
#  1) 真实 style.css 的 `html/body { height:100% }` + `body { overflow:hidden }` 是为
#     Wails 原生窗口写的（滚动交给内部 .page），套到对比稿这个「长文档」上会锁死滚动，
#     只能看到第一屏 —— 2026-09-15 用户反馈「滑不到下面」的根因。
#  2) headless 截图那边传的是超大 --window-size，整页一次性拍下，所以截图永远发现不了。
SCROLL_FIX = '''
<style id="design-scroll-fix">
  html { height: auto !important; overflow: visible !important; }
  body { height: auto !important; min-height: 100vh; overflow: visible !important; }
</style>
'''

LINK_RE = re.compile(r'<link\s+rel="stylesheet"\s+href="([^"]+)"\s*>')


def inline(src, css_body):
    html = open(src, encoding='utf-8').read()
    for href in LINK_RE.findall(html):
        # 只替换指向 style.css 的那份；页面自己的 <style> 保持不动
        if not href.endswith('style.css'):
            continue
        html = html.replace(LINK_RE.search(html).group(0),
                            '<style>\n/* === 内联自 frontend/src/style.css === */\n'
                            + css_body + '\n</style>', 1)
    # 保险：仍有剩余 link 就报错，避免静默产出半样式页面
    left = LINK_RE.findall(html)
    if left:
        raise RuntimeError('仍有未内联的 link：%s' % left)
    # 无论有没有外链，都要补滚动补丁（自包含的稿子也可能抄了 body{overflow:hidden}）
    html = html.replace('</head>', SCROLL_FIX + '\n</head>', 1)
    return html


def main():
    targets = sys.argv[1:] or sorted(glob.glob(os.path.join(ROOT, 'design', '*.html')))
    css_body = open(CSS, encoding='utf-8').read()
    os.makedirs(PREVIEW, exist_ok=True)
    for src in targets:
        src = os.path.abspath(src)
        html = inline(src, css_body)
        name = os.path.splitext(os.path.basename(src))[0]
        out = os.path.join(PREVIEW, name + '.standalone.html')
        open(out, 'w', encoding='utf-8', newline='\n').write(html)
        print('%-42s -> preview/%s  (%.0f KB)'
              % (os.path.basename(src), os.path.basename(out), len(html) / 1024))


if __name__ == '__main__':
    main()
