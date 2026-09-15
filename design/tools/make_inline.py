# -*- coding: utf-8 -*-
"""把 dist 的构建产物内联成单个 HTML，用 file:// 打开截图。

为什么不用 http.server：本机沙盒会拦 socket bind（WinError 10013），
且环境里设了 HTTP_PROXY=127.0.0.1:6101，本地请求会被代理吞掉返回 502。
内联后既不需要监听端口，也不经过代理。

用法：python make_inline.py <dist_dir>
产出：<dist_dir>/_inline.html
"""
import os
import re
import sys

DIST = sys.argv[1] if len(sys.argv) > 1 else r'D:\go_workspace\mytools\frontend\dist'

html = open(os.path.join(DIST, 'index.html'), encoding='utf-8').read()

m_css = re.search(r'<link rel="stylesheet" href="(/assets/[^"]+\.css)">', html)
m_js = re.search(r'<script type="module"[^>]*src="(/assets/[^"]+\.js)"></script>', html)
assert m_css and m_js, '没找到构建产物引用，检查 dist/index.html'

css_body = open(os.path.join(DIST, m_css.group(1).lstrip('/')), encoding='utf-8').read()
js_body = open(os.path.join(DIST, m_js.group(1).lstrip('/')), encoding='utf-8').read()

# 1) CSS / JS 里对 /assets/ 的绝对引用改成相对（内联页面就放在 dist 下，相对路径能解析到真实文件）
css_body = css_body.replace('url(/assets/', 'url(assets/').replace('url("/assets/', 'url("assets/')
js_body = js_body.replace('"/assets/', '"assets/')

# 2) JS 内联进 <script> 前必须打断 "</script>"，否则标签会被提前闭合
js_body = js_body.replace('</script', '<\\/script')

html = html.replace(m_css.group(0), '<style>\n' + css_body + '\n</style>')
html = html.replace(m_js.group(0), '<script type="module">\n' + js_body + '\n</script>')

# ---------------- 预览脚手架：顶掉 Wails 后端 ----------------
STUB = '''    <script>
      (function () {
        var q = new URLSearchParams(location.search);
        try {
          localStorage.setItem('mytools:theme', q.get('theme') || 'light');
          localStorage.setItem('mytools:nav-collapsed', q.get('cv') === '1' ? '1' : '0');
        } catch (e) {}
        function mk(path) {
          return new Proxy(function () {}, {
            get(_, p) {
              if (p === 'then' || typeof p === 'symbol') return undefined;
              return mk(path + '.' + String(p));
            },
            apply() { return Promise.resolve(window.__stub(path)); }
          });
        }
        var PUB = '-----BEGIN RSA PUBLIC KEY-----\\nMIIBCgKCAQEAy8Dbv8prpJ/0kKhlGeJYozo2t60EG8L0561g13R29LvMR5hyvGZl';
        var PRIV = '-----BEGIN RSA PRIVATE KEY-----\\nMIIEpAIBAAKCAQEAy8Dbv8prpJ/0kKhlGeJYozo2t60EG8L0561g13R29LvMR5hy';
        window.__stub = function (path) {
          // 版本号：fix=true 时伪造「有新版本」，用来预览红点 + 更新弹窗（?upd=1）
          if (/VersionHandler.CheckUpdate/.test(path)) {
            return window.__FAKE_UPDATE
              ? { status: 'update', message: '发现新版本', latest: {
                  version: '1.1.4', create_date: '2026-09-20',
                  change_log: ['新增启动自检更新，有新版本时版本号右上角挂红点',
                               'IP 查询改为多源并发，海外接口超时也不再报错',
                               '顶栏版本号可点击：已是最新 / 新版本弹窗',
                               'RSA 页下拉与我们统一为自绘组件'],
                  download_url: 'https://github.com/tomeyZ/mytools/releases/latest',
                  release_url: 'https://github.com/tomeyZ/mytools/releases/latest'
                } }
              : { status: 'latest', message: '已是最新版本' };
          }
          if (/GetCurrentVersion/.test(path)) return '1.1.3';
          if (/GetIpInfo/.test(path)) return {
            status: 'success', query: '8.8.8.8', country: '美国', countryCode: 'US',
            regionName: '加利福尼亚州', region: 'CA', city: '山景城',
            lat: 37.386, lon: -122.0838, timezone: 'America/Los_Angeles',
            isp: 'Google LLC', org: 'Google LLC', as: 'AS15169 Google LLC'
          };
          if (/TimestampToDate/.test(path)) return '2026-09-01 22:57:03';
          if (/Encrypt/.test(path)) return { success: true, result: 'U2FsdGVkX1+9KpQ0mZ3sT7Wq1oBvYc8xLd2EhRnUfA0=' };
          if (/Decrypt/.test(path)) return { success: true, result: 'hello world' };
          if (/GenerateKeyPair/.test(path)) return { success: true, publicKey: PUB, privateKey: PRIV };
          return '';
        };
        window.go = mk('go');
        window.__TOOL = q.get('tool') || 'time';
        window.__SEED = q.get('seed') === '1';
        window.__FAKE_UPDATE = q.has('upd');   // ?upd=1 顺便自动点开弹窗，?upd=dot 只看红点
      })();
    </script>
'''

TAIL = '''<script>
  function fill(el, v) {
    if (!el) return;
    el.value = v;
    el.dispatchEvent(new Event('input', { bubbles: true }));
  }
  function clickByText(text) {
    var btns = document.querySelectorAll('button');
    for (var i = 0; i < btns.length; i++) {
      if (btns[i].textContent.trim().indexOf(text) === 0) { btns[i].click(); return true; }
    }
    return false;
  }
  setTimeout(function () {
    var MAP = { md5: 'MD5加密', aes: 'AES加解密', rsa: 'RSA密钥生成',
                json: 'JSON美化', text: '文本处理', ip: 'IP地址查询' };
    if (window.__TOOL !== 'time') {
      var el = document.querySelector('.nav-item[title="' + MAP[window.__TOOL] + '"]');
      if (el) el.click();
    }
    // ?tz=open 自动展开时区下拉（headless 没法手工点）
    if (new URLSearchParams(location.search).get('tz') === 'open') {
      setTimeout(function () {
        var trg = document.querySelector('.tz-trigger');
        if (trg) trg.click();
      }, 300);
    }
    // ?upd=1 点一下版本号，看「新版本弹窗」（触发它的红点由 __FAKE_UPDATE 提供）
    if (new URLSearchParams(location.search).get('upd') === '1') {
      setTimeout(function () {
        var pill = document.querySelector('.ver-btn');
        if (pill) pill.click();
      }, 700);
    }
    if (!window.__SEED) return;
    setTimeout(function () {
      var t = window.__TOOL;
      if (t === 'time') {
        fill(document.querySelectorAll('.inp')[0], '1789228377');
        clickByText('转换');
      } else if (t === 'md5') {
        fill(document.querySelector('textarea.inp'), 'hello world');
        clickByText('MD5 加密');
      } else if (t === 'aes') {
        fill(document.querySelector('input[placeholder="请输入加密密钥"]'), 'mytools-2026');
        fill(document.querySelector('textarea.inp'), 'hello world');
        clickByText('加密');
      } else if (t === 'rsa') {
        // 密钥卡是 v-if 渲染的，必须先生成才出现「复制」按钮
        clickByText('生成密钥对');
      }
    }, 300);
  }, 400);
</script>
</body>'''

html = html.replace('<head>', '<head>\n' + STUB, 1).replace('</body>', TAIL, 1)

out = os.path.join(DIST, '_inline.html')
open(out, 'w', encoding='utf-8', newline='\n').write(html)
print('_inline.html 已生成  %.1f MB' % (len(html) / 1048576))
