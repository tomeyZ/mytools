"""headless 截图 + 几何探针（设计稿验证用）。

为什么不直接敲 chrome 命令：
  - 本机 bash 缺 mkdir/ls/tail 等命令，用 python 驱动最稳；
  - 设计稿的几何量必须实测，肉眼估会错（踩过多次）。

用法：
    python design/tools/shot.py <页面> [--out 名称] [--size 1122x1990] [--dark] [--probe]
                               [--query "upd=1&tool=aes"]

约定：页面里若在 `?probe=1` 时把测量结果塞进 document.title，
      本脚本加 --probe 就会把 title 内容按 "|" 拆行打出来。

注意：--window-size 与实际视口相差 (22, 98)，要视口 H 得传 H+98。
"""
import argparse
import os
import re
import subprocess
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
CHROME = r"C:\Program Files\Google\Chrome\Application\chrome.exe"
OUT_DIR = os.path.join(ROOT, "design", "preview")


def to_url(page):
    if page.startswith("http"):
        return page
    if not os.path.isabs(page):
        page = os.path.join(ROOT, page)
    return "file:///" + page.replace("\\", "/")


def run(args, capture=False):
    return subprocess.run(
        [CHROME, "--headless=new", "--disable-gpu", "--hide-scrollbars",
         "--force-device-scale-factor=2", "--virtual-time-budget=4000"] + args,
        capture_output=capture, text=True, encoding="utf-8", errors="replace",
    )


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("page", help="页面路径或 URL")
    ap.add_argument("--out", default="shot.png", help="输出文件名（存到 design/preview/）")
    ap.add_argument("--size", default="1122x1990", help="窗口尺寸 WxH，不是视口")
    ap.add_argument("--dark", action="store_true", help="附加 ?theme=dark")
    ap.add_argument("--probe", action="store_true", help="抓 document.title 里的测量结果")
    ap.add_argument("--query", default="", help="附加 URL 参数，如 upd=1")
    a = ap.parse_args()

    url = to_url(a.page)
    if a.query:
        url += ("&" if "?" in url else "?") + a.query.lstrip("?&")
    if a.dark:
        url += ("&" if "?" in url else "?") + "theme=dark"
    w, h = a.size.lower().split("x")

    if a.probe:
        probe_url = url + ("&" if "?" in url else "?") + "probe=1"
        r = run(["--dump-dom", probe_url], capture=True)
        m = re.search(r"<title>(.*?)</title>", r.stdout, re.S)
        print("--- probe ---")
        for part in (m.group(1) if m else "(no title)").split("|"):
            part = part.strip()
            if part:
                print("  " + part)

    os.makedirs(OUT_DIR, exist_ok=True)
    dst = os.path.join(OUT_DIR, a.out)
    if os.path.exists(dst):
        os.remove(dst)
    run([f"--window-size={w},{h}", f"--screenshot={dst}", url])
    ok = os.path.exists(dst)
    print(f"{a.out}: exists={ok} size={os.path.getsize(dst) if ok else 0}")
    return 0 if ok else 1


if __name__ == "__main__":
    sys.exit(main())
