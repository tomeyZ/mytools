// JSONEditor 布局诊断脚本（临时排查用，排查完可删）
//
// 用法：
//   1. 用 debug 模式启动应用：wails dev（或 wails build -debug）
//   2. 打开 JSON 美化页，切到你看到白条的那个模式（树视图/代码视图）
//   3. 右键 → Inspect / 检查，打开 DevTools
//   4. 切到 Console 面板，把本文件全部内容粘贴进去回车
//   5. 把输出截图发我
//
// 脚本只读取 DOM 和样式，不修改任何东西，可放心运行。

(() => {
  const q = (s) => document.querySelector(s);

  const info = (name, el) => {
    if (!el) {
      console.log(`${name}: 不存在`);
      return;
    }
    const cs = getComputedStyle(el);
    const r = el.getBoundingClientRect();
    console.log(`%c${name}`, 'font-weight:bold', {
      '高度': Math.round(r.height) + 'px',
      '内联paddingBottom': el.style.paddingBottom || '(无)',
      '计算后paddingBottom': cs.paddingBottom,
      'overflow': cs.overflow,
      'position': cs.position,
    });
  };

  console.log('%c=== JSONEditor 布局诊断 ===', 'color:#1a73e8;font-size:14px');

  const frame = q('.jsoneditor');

  // 1. 当前所处模式
  console.log(
    '当前模式:',
    frame ? frame.className.replace(/.*jsoneditor-mode-(\w+).*/, '$1') : '未知'
  );

  // 2. 各层元素的实际高度与 padding
  info('frame      (.jsoneditor)', frame);
  info('outer      (.jsoneditor-outer)', q('.jsoneditor-outer'));
  info('tree       (.jsoneditor-tree)', q('.jsoneditor-tree'));
  info('ace        (.ace_editor)', q('.ace_editor'));
  info('ace scroller (.ace_scroller)', q('.ace_scroller'));

  // 3. 底部命中测试：从底边往上采样，看白条区域究竟是哪个元素
  if (frame) {
    const r = frame.getBoundingClientRect();
    console.log('%c=== 底部命中测试（白条是什么元素）===', 'color:#d93025');
    for (let dy = 2; dy <= 80; dy += 12) {
      const el = document.elementFromPoint(r.left + r.width / 2, r.bottom - dy);
      console.log(`  底部往上 ${dy}px:`, el ? el.className || el.tagName : 'null');
    }
  }

  // 4. 滚动状态：scrollHeight 远大于 clientHeight 说明内容被撑高了
  if (frame) {
    console.log('%c=== 滚动信息 ===', 'color:#188038');
    console.log({
      'frame.scrollHeight': frame.scrollHeight,
      'frame.clientHeight': frame.clientHeight,
      'frame.scrollTop': frame.scrollTop,
      '底部剩余可滚动px': frame.scrollHeight - frame.clientHeight - frame.scrollTop,
    });
  }

  // 5. outer 的滚动状态
  const outer = q('.jsoneditor-outer');
  if (outer) {
    console.log('%c=== outer 滚动信息 ===', 'color:#188038');
    console.log({
      'outer.scrollHeight': outer.scrollHeight,
      'outer.clientHeight': outer.clientHeight,
    });
  }

  // 6. 白条根因指标：outer 的负 margin-top / padding-top。
  //    JSONEditor 给 outer 加 has-main-menu-bar(-35px) / has-nav-bar(-26px)，
  //    用负 margin 上提去盖住菜单栏。菜单栏用 display:none 藏掉后负 margin 还在，
  //    就会把整个 outer 上提，底部漏出等高的白条。
  //    ✅ 健康值：margin-top 与 padding-top 均为 0px（或两者绝对值相等且菜单栏可见）
  //    ❌ 病态值：margin-top 为负数，且页面上没有对应的菜单栏
  if (outer && frame) {
    const ocs = getComputedStyle(outer);
    const gap = frame.getBoundingClientRect().bottom - outer.getBoundingClientRect().bottom;
    console.log('%c=== 白条根因指标 ===', 'color:#a142f4');
    console.log({
      'outer.className': outer.className,
      'outer.marginTop': ocs.marginTop,
      'outer.paddingTop': ocs.paddingTop,
      '底部空隙px(frame底 - outer底)': Math.round(gap),
      '结论': Math.abs(gap) <= 2
        ? '✅ 无白条：outer 与 frame 底边齐平'
        : `❌ 底部漏白 ${Math.round(gap)}px，检查 outer 上的 ${outer.className.split(' ').filter(c => c.startsWith('has-')).join(' / ') || '(无 has-* 类)'}`,
    });
  }
})();
