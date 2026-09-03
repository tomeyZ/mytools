# mytools 项目长期笔记

## 项目结构

Wails (Go + Vue) 桌面工具集。`frontend/` 为 Vue 前端，`internal/` 为 Go 后端。

- 前端组件：`frontend/src/components/`
- 构建：`cd frontend && npm run build`（用 managed node：`C:/Users/TRS/.workbuddy/binaries/node/versions/22.22.2-2/npm.cmd`）
- 运行：`wails dev`。**Wails 桌面应用没有浏览器 HMR**，改完前端必须重启/重新构建才能看到效果

## JSON 美化页（JsonFormat.vue）

用 `jsoneditor@^10` 库（不是纯 ace；只有 code 模式底层才是 ace）。

- 默认 `mode: 'tree'`（可折叠/展开，对标 bejson），工具条提供「树视图/代码视图」切换
- `modes: ['tree', 'code']`，去掉了 form/text/view 避免右键菜单误触
- 自带菜单栏/状态栏已关闭（`mainMenuBar: false` / `statusBar: false`），功能由自定义工具条替代

### JSONEditor 集成三条铁律（踩了 5 次坑换来）

1. **库有 options 就用 options 关功能，绝不用 CSS `display: none` 掩盖库节点**
   隐藏节点不会移除容器上的 `has-*` 类，对应的负 margin / padding 继续生效：
   - `has-main-menu-bar` → `margin-top: -35px`（把 outer 上提，底部漏 35px 白条）
   - `has-status-bar` → `padding-bottom: 26px`（底部漏 26px 白条）
   源码：`jsoneditor.js:2760-2888`（menu）、`2980`（statusbar）

2. **不要改库子节点属性**（如 `aceEditor.setOptions({ maxLines })`）。库的默认行为经过兼容性测试，外层 flex + overflow 控制即可

3. **加 `!important` 覆盖库的内联样式前，先查清那段内联样式的用途**
   `JSONEditor` 的 `onChangeHeight` 是 **ErrorTable** 的回调（不是 ace 的），
   用 `padding-bottom: N + margin-bottom: -N` 给校验错误表格腾位置——强行归零会破坏该机制

### 布局调试：先实测再推理

第三方库的 CSS 布局问题，纯读代码推理极易连续误判（本项目曾连错 4 轮）。**优先拿到真实计算值再推理。**

诊断脚本：`frontend/debug-jsoneditor.js`
用法：`wails dev` → 打开页面 → 右键 Inspect → Console 粘贴运行。
输出各层元素高度 / 内联与计算后 padding / overflow、底部命中测试（白条是哪个元素）、
以及 `margin-top` 与「frame 底 - outer 底」空隙并直接给 ✅/❌ 结论。

排查要领：
- **拿到数字要算不要看**。`339 - 304 = 35` 精确命中 CSS 里的 `-35px`，这种吻合就是铁证
- 用 `grep -o "规则{[^}]*}" dist/assets/*.css` 确认自定义 CSS 真进了构建产物
- 注意 **scoped CSS 优先级 (0,2,0) 高于全局块的 `div.xxx` (0,1,1)**，全局块写 overflow 会被静默覆盖
- 本机没装 agent-browser（要下 ~500MB Chromium），写诊断脚本让用户跑是轻量替代

**环境备注**：`frontend/debug-jsoneditor.js` 是临时排查工具，问题解决后可删除。
