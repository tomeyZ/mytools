# -*- coding: utf-8 -*-
"""在 _inline.html 里注入交互测试脚本，用 --dump-dom 抓 <title> 读结果。

覆盖：初始态 / 点击展开 / 点选切换 / 键盘 ↑↓ + Enter / Esc / 点击外部关闭。

两个必须注意的坑：
1. Vue 的 DOM 更新是异步的，每次操作后都要等一拍再断言
2. 探针函数返回的是字符串，调用处**不要**再写 `x() ? '开' : '关'`——
   字符串永远 truthy，会把结果全部写成「开」（踩过）
"""
import os
import sys

DIST = sys.argv[1] if len(sys.argv) > 1 else r'D:\go_workspace\mytools\frontend\dist'
src = open(os.path.join(DIST, '_inline.html'), encoding='utf-8').read()

# 对照组：禁用过渡后再跑一遍。
# headless 的 --virtual-time-budget 会让 CSS transition 的结束事件不触发，
# 导致 v-if 元素卡在 leave 态留在 DOM（node 恒为 1）。禁掉过渡后 node 应归 0，
# 以此证明「残留是测试环境假象，不是代码缺陷」。
NOTRANS = '''<style>
  .tz-pop-enter-active, .tz-pop-leave-active,
  .tz-pop-enter-from, .tz-pop-leave-to { transition: none !important; }
</style>
</head>'''

TEST = r'''<script>
(function () {
  var R = [];
  var wait = function (ms) { return new Promise(function (r) { setTimeout(r, ms); }); };
  var q = function (s) { return document.querySelector(s); };

  // 面板状态：直接读 aria-expanded（绑定 tzOpen，反映 Vue 的真实状态），
  // 同时报告 DOM 节点是否存在 / 是否卡在 leave 类上
  var menu = function () {
    var t = q('.tz-trigger');
    var el = q('.tz-menu');
    var exp = t ? t.getAttribute('aria-expanded') : '?';
    var cls = el ? String(el.className).replace('tz-menu', '').trim() : '';
    return 'exp=' + exp + ' node=' + (el ? '1' : '0') + (cls ? '(' + cls + ')' : '');
  };
  var cur = function () {
    return q('.tz-current').textContent.trim() + '|' + q('.tz-offset').textContent.trim();
  };
  var onItem = function () {
    var e = q('.tz-item.on');
    return e ? e.querySelector('.tz-name').textContent.trim() : 'NONE';
  };
  var actItem = function () {
    var e = q('.tz-item.act');
    return e ? e.querySelector('.tz-name').textContent.trim() : 'NONE';
  };
  var key = function (k) {
    q('.tz-trigger').dispatchEvent(new KeyboardEvent('keydown', { key: k, bubbles: true, cancelable: true }));
  };

  (async function () {
    await wait(900);

    R.push('1 初始 | ' + menu() + ' | 触发器=' + cur());

    q('.tz-trigger').click();
    await wait(250);
    var items = document.querySelectorAll('.tz-item');
    R.push('2 点击展开 | ' + menu() + ' | 项数=' + items.length
           + ' 首项=' + items[0].querySelector('.tz-name').textContent.trim()
           + ' 首项偏移=' + items[0].querySelector('.tz-off').textContent.trim()
           + ' | on=' + onItem());

    // 点第 5 项（排序后应为印度标准时间）
    var picked = items[4].querySelector('.tz-name').textContent.trim();
    items[4].click();
    await wait(300);
    R.push('3 点「' + picked + '」 | ' + menu() + ' | 触发器=' + cur());

    // 重新打开：高亮应落在当前选中项
    q('.tz-trigger').click();
    await wait(250);
    R.push('4 重开 | ' + menu() + ' | act=' + actItem() + ' on=' + onItem());

    key('ArrowDown'); await wait(150);
    var afterDown = actItem();
    key('Enter'); await wait(300);
    R.push('5 ↓后高亮=' + afterDown + ' → Enter | ' + menu() + ' | 触发器=' + cur());

    q('.tz-trigger').click(); await wait(250);
    var opened = menu();
    key('Escape'); await wait(300);
    R.push('6 Esc（开=' + opened + '） | ' + menu());

    q('.tz-trigger').click(); await wait(250);
    document.body.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    await wait(300);
    R.push('7 点面板外 | ' + menu());

    q('.tz-trigger').click(); await wait(250);
    q('.tz-menu').dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    await wait(300);
    R.push('8 点面板内(应保持开) | ' + menu());

    document.title = 'RESULT::' + R.join(' ;; ');
  })();
})();
</script>
</body>'''

out = src.replace('</head>', NOTRANS, 1).replace('</body>', TEST, 1)
p = os.path.join(DIST, '_itest.html')
open(p, 'w', encoding='utf-8', newline='\n').write(out)
print('_itest.html 已生成')
