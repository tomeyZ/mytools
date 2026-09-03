<template>
  <div class="json-editor-container">
    <!-- 工具条 -->
    <div class="toolbar">
      <button class="tool-btn" @click="formatJson">
        <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M8 3H7a2 2 0 0 0-2 2v5a2 2 0 0 1-2 2 2 2 0 0 1 2 2v5c0 1.1.9 2 2 2h1"/><path d="M16 21h1a2 2 0 0 0 2-2v-5c0-1.1.9-2 2-2a2 2 0 0 1-2-2V5a2 2 0 0 0-2-2h-1"/></svg>
        格式化
      </button>
      <button class="tool-btn" @click="compactJson">
        <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="4 14 10 14 10 20"/><polyline points="20 10 14 10 14 4"/><line x1="14" y1="10" x2="21" y2="3"/><line x1="3" y1="21" x2="10" y2="14"/></svg>
        压缩
      </button>
      <button class="tool-btn" @click="sortJson">
        <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m3 16 4 4 4-4"/><path d="M7 20V4"/><path d="M11 4h10"/><path d="M11 8h7"/><path d="M11 12h4"/></svg>
        排序
      </button>
      <button class="tool-btn" @click="unescapeAndFormat">
        <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m7 21-4.3-4.3c-1-1-1-2.5 0-3.4l9.6-9.6c1-1 2.5-1 3.4 0l5.6 5.6c1 1 1 2.5 0 3.4L13 21"/><path d="M22 21H7"/><path d="m5 11 9 9"/></svg>
        去转义
      </button>
      <button class="tool-btn" @click="escapeJson">
        <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 3a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2 1 1 0 0 1 1 1v1a2 2 0 0 1-2 2 1 1 0 0 0-1 1v2a1 1 0 0 0 1 1 6 6 0 0 0 6-6V5a2 2 0 0 0-2-2z"/><path d="M5 3a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2 1 1 0 0 1 1 1v1a2 2 0 0 1-2 2 1 1 0 0 0-1 1v2a1 1 0 0 0 1 1 6 6 0 0 0 6-6V5a2 2 0 0 0-2-2z"/></svg>
        转义
      </button>
      <button class="tool-btn" @click="unicodeToChinese">
        <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m5 8 6 6"/><path d="m4 14 6-6 2-3"/><path d="M2 5h12"/><path d="M7 2h1"/><path d="m22 22-5-10-5 10"/><path d="M14 18h6"/></svg>
        Unicode转中文
      </button>
    </div>

    <div ref="jsoneditor" class="jsoneditor"></div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import JSONEditor from 'jsoneditor';
import 'jsoneditor/dist/jsoneditor.min.css';

export default {
  setup() {
    const jsoneditor = ref(null);
    const error = ref('');
    let editor = null;

    // 错误横幅统一入口：设置后 5 秒自动消失，期间再次报错会重置计时。
    // 直接 error.value = xxx 的散落赋值会绕过计时器，禁止再用
    const ERROR_AUTO_DISMISS_MS = 5000;
    let errorTimer = null;
    const showError = (msg) => {
      error.value = msg;
      if (errorTimer) clearTimeout(errorTimer);
      errorTimer = setTimeout(() => {
        error.value = '';
        errorTimer = null;
      }, ERROR_AUTO_DISMISS_MS);
    };
    const clearError = () => {
      error.value = '';
      if (errorTimer) {
        clearTimeout(errorTimer);
        errorTimer = null;
      }
    };

    // 初始化 JSONEditor
    onMounted(() => {
      const container = jsoneditor.value;
      const options = {
        // 仅 code 模式：粘贴/编辑原文。ace 自带 JSON 折叠（鼠标悬停行号区域出现折叠箭头），
        // 原先的树视图与它功能重叠，已移除
        mode: 'code',
        // 关掉自带菜单栏/状态栏：功能由自定义工具条替代。
        // 不能改用 CSS display:none 隐藏——菜单栏对应的 .jsoneditor-outer 上
        // 有 has-main-menu-bar(-35px margin) / has-status-bar(26px padding)，
        // 节点藏了但这些负边距还在，会把 outer 整体上提，底部漏出白条。
        mainMenuBar: false,
        statusBar: false,
        onError: (err) => {
          showError(err.toString());
        },
        onChange: () => {
          try {
            editor.get();
            clearError();
          } catch (e) {
            showError('当前内容不是有效的 JSON');
          }
        }
      };

      editor = new JSONEditor(container, options);
      editor.set();
    });

    // 统一取编辑器文本并收敛成字符串。
    // 防御：库在某些空状态下 getText() 可能返回 undefined 而不是 ''，
    // 直接 .trim() 会抛 TypeError
    const getEditorText = () => String((editor && editor.getText()) ?? '');

    // 销毁 JSONEditor
    onBeforeUnmount(() => {
      if (errorTimer) {
        clearTimeout(errorTimer);
        errorTimer = null;
      }
      if (editor) {
        editor.destroy();
      }
    });

    // 去转义
    //
    // 两类输入必须分开处理，否则会互相破坏：
    //   1. 被转义过的 JSON（{\"a\":1} 或 "{\"a\":1}"）
    //      \n 是 JSON 字符串内部的合法转义序列，属于 JSON 语法本身，必须保持字面量。
    //      提前解码成真实换行会让 JSON 非法 —— JSON 规范禁止字符串里出现裸控制字符
    //      （U+000A 等），JSON.parse 会直接报 "Bad control character"。
    //   2. 被转义过的纯文本（hello\nworld）
    //      这种才期望输出真实换行。
    //
    // 策略：先按 JSON 逐层剥离（期间保护住 \n \t \r），剥完能解析就格式化；
    //       解析不了说明是纯文本，再做完整反转义。

    // 保护符：Unicode 私有区字符，JSON 字符串内合法、纯文本里也不会自然出现
    const GUARD = '\uE000';
    // 位置感知的保护：只把「字符串区域内」的 \n \t \r 换成 GUARD（保持字面量），
    // 区域外的 \n 是格式化换行被转义后的产物，必须让 JSON.parse 还原成真实换行。
    // 同一个 \n 在两种文法里含义相反：
    //   - 转义文法（stringify 产物）：值内换行是 \\n（3字符）、格式换行是 \n（2字符，在字符串外）
    //   - 残缺文法（日志/手抄）：值内换行就是 \n（2字符，在字符串内）——若被还原成真实换行会撕裂 JSON
    // 扫描规则：\\ 成对消费（不能从中间截断）；\" 视为转义后的字符串定界符，切换内外状态
    const guardControls = (s) => {
      let out = '';
      let inString = false;
      for (let i = 0; i < s.length; i++) {
        const c = s[i];
        if (c === '\\' && i + 1 < s.length) {
          const d = s[i + 1];
          if (d === '\\') {
            out += '\\\\';                  // 成对反斜杠：整体保留
          } else if (d === '"') {
            out += '\\"';                   // 转义后的字符串定界符
            inString = !inString;
          } else if (inString && 'ntrbf'.includes(d)) {
            out += GUARD + d;               // 字符串值内的控制序列：保持字面量
          } else {
            out += c + d;                   // 字符串外的 \n（格式换行）、\u 等：原样，交给 parse 还原
          }
          i++;
        } else {
          out += c;
        }
      }
      return out;
    };
    const unguardControls = (s) => s.split(GUARD).join('\\');

    const tryParse = (s) => {
      try {
        return JSON.parse(s);
      } catch (e) {
        return undefined;
      }
    };

    // 完整反转义：仅用于「确定不是 JSON」的纯文本场景
    // 顺序固定：\\ 先占位，最后还原，避免与后续转义相互干扰
    const fullUnescape = (s) => s
      .replace(/\\u([0-9a-fA-F]{4})/g, (_, hex) => String.fromCharCode(parseInt(hex, 16)))
      .replace(/\\\\/g, GUARD)
      .replace(/\\n/g, '\n')
      .replace(/\\t/g, '\t')
      .replace(/\\r/g, '\r')
      .replace(/\\"/g, '"')
      .split(GUARD).join('\\');

    const unescapeAndFormat = () => {
      if (!editor) return;
      const text = getEditorText();
      if (!text.trim()) return; // 空内容：静默不动作，不报错
      clearError();

      // 逐层剥离 JSON 字符串包装，两种形态都支持：
      //   A. 带外层引号  "{\"a\":1}"
      //   B. 裸转义文本  {\"a\":1}      ← 从代码 / 日志 / curl 复制出来时外层引号常被丢掉
      let current = guardControls(text);
      for (let i = 0; i < 10; i++) {
        let next;
        const asJson = tryParse(current);
        if (typeof asJson === 'string') {
          next = guardControls(asJson);
        } else if (current.includes('\\"')) {
          // 补一层引号当作「JSON 字符串的内容」再解析。
          // 必须限定在含 \" 时才走这条路：否则 hello\nworld 这种纯文本会被凭空
          // 补一层再剥掉，看着像剥了一层，实际是原地打转
          const asStringContent = tryParse(`"${current}"`);
          if (typeof asStringContent !== 'string') break;
          next = guardControls(asStringContent);
        } else {
          break;
        }
        // 无进展立即停：否则纯文本会被反复包引号，白白跑满 10 轮
        if (next === current) break;
        current = next;
      }

      // 结构化 / 标量 JSON：格式化展示（stringify 后还原保护符，\n 保持字面量）
      const final = tryParse(current);
      if (final !== null && typeof final === 'object') {
        editor.setText(unguardControls(JSON.stringify(final, null, 2)));
        return;
      }
      if (final !== undefined) {
        editor.setText(unguardControls(String(final)));
        return;
      }

      // 走到这里说明不是合法 JSON，按形态分两种处理
      const plain = unguardControls(current);

      // 形态一：看着像 JSON（以 { 或 [ 开头）但语法有问题 —— 典型是复制时被截断、括号不匹配。
      // 只还原引号，保留 \n \t：它们是 JSON 语法的一部分，解码会让内容更难修复。
      // 同时给出明确提示，避免用户对着一堆断行找不出原因。
      if (/^\s*[[{]/.test(plain)) {
        const unquoted = plain
          .replace(/\\\\/g, GUARD)
          .replace(/\\"/g, '"')
          .split(GUARD).join('\\');
        if (unquoted !== text) {
          editor.setText(unquoted);
        }
        showError('内容不是合法的 JSON，已还原转义。请检查引号 / 括号是否完整（常见于复制时被截断）');
        return;
      }

      // 形态二：纯文本 —— 完整反转义（\n \t \uXXXX 解码为真实字符）
      const decoded = fullUnescape(plain);
      if (decoded !== text) {
        editor.setText(decoded);
        return;
      }

      showError('未检测到可去除的转义层');
    };

    // Unicode 转中文：把 \uXXXX 显示为中文，不改动其他内容
    const unicodeToChinese = () => {
      if (!editor) return;
      const text = getEditorText();
      if (!text.trim()) return; // 空内容：静默不动作，不报错
      clearError();
      // 仅替换 \uXXXX 序列，中文之外的转义（如 \n \"）保持原样
      const converted = text.replace(/\\u([0-9a-fA-F]{4})/g, (_, code) => {
        return String.fromCharCode(parseInt(code, 16));
      });
      if (converted === text) {
        showError('未检测到 Unicode 转义');
        return;
      }
      editor.setText(converted);
    };

    // 格式化当前 JSON
    const formatJson = () => {
      if (!editor) return;
      const text = getEditorText().trim();
      if (!text) return; // 空内容：静默不动作，不报错
      clearError();
      try {
        const obj = JSON.parse(text);
        editor.setText(JSON.stringify(obj, null, 2));
      } catch (e) {
        showError('内容不是有效的 JSON');
      }
    };

    // 压缩 JSON：删除所有多余空格换行
    const compactJson = () => {
      if (!editor) return;
      const text = getEditorText().trim();
      if (!text) return; // 空内容：静默不动作，不报错
      clearError();
      try {
        const obj = JSON.parse(text);
        editor.setText(JSON.stringify(obj));
      } catch (e) {
        showError('内容不是有效的 JSON');
      }
    };

    // 转义：输出 bejson 同款的「裸转义」形态 {"a":1} → {\"a\":1}
    // 不带外层引号——日志/SQL/代码里复制出来的就是这种形态（外层引号属于宿主语言的字符串定界符）。
    // 实现：整体 stringify 后剥掉首尾引号；连点叠加多层，「去转义」可对称还原
    const escapeJson = () => {
      if (!editor) return;
      const text = getEditorText().trim();
      if (!text) return; // 空内容：静默不动作，不报错
      clearError();
      // JSON.stringify 产物必然以 " 开头 " 结尾，slice 安全
      editor.setText(JSON.stringify(text).slice(1, -1));
    };

    // 按 key 排序：所有层级的对象键按字典序排列，便于核对签名参数
    const sortJson = () => {
      if (!editor) return;
      const text = getEditorText().trim();
      if (!text) return; // 空内容：静默不动作，不报错
      clearError();
      try {
        const obj = JSON.parse(text);
        const sorted = sortValue(obj);
        editor.setText(JSON.stringify(sorted, null, 2));
      } catch (e) {
        showError('内容不是有效的 JSON');
      }
    };

    // 递归排序：对象键排序，数组保持顺序
    const sortValue = (value) => {
      if (Array.isArray(value)) {
        return value.map(sortValue);
      }
      if (value !== null && typeof value === 'object') {
        const sorted = {};
        Object.keys(value).sort().forEach(key => {
          sorted[key] = sortValue(value[key]);
        });
        return sorted;
      }
      return value;
    };

    return {
      jsoneditor,
      error,
      unescapeAndFormat,
      unicodeToChinese,
      formatJson,
      compactJson,
      escapeJson,
      sortJson
    };
  }
};
</script>

<style scoped>
.json-editor-container {
  padding: 20px;
  height: calc(100vh - 60px);
  display: flex;
  flex-direction: column;
}

/* 工具条：与 RSA 页配置行一致的卡片风格 */
.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-shrink: 0;
  background: #f5f7fa;
  border-radius: 8px;
  padding: 10px 16px;
}

.tool-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 12px;
  background: white;
  border: 1px solid #dadce0;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 500;
  color: #3c4043;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

/* 内联 SVG 图标：currentColor 跟随文字颜色（hover 蓝 / 激活白），不引入外部图标库 */
.tool-btn .btn-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
}

.tool-btn:hover {
  border-color: #1a73e8;
  color: #1a73e8;
}

.error-message {
  margin-top: 15px;
  padding: 10px;
  background-color: #ffebee;
  color: #f44336;
  border-radius: 4px;
  font-family: monospace;
}

.jsoneditor {
  flex: 1;
  width: 100%;
  /* flex 子项默认 min-height:auto，不置 0 会被内容撑开、导致外层滚动异常 */
  min-height: 0;
  /* 注意：scoped 规则优先级 (0,2,0) 高于全局块的 div.jsoneditor (0,1,1)，
     所以 frame 的 overflow 只能在这里控制，写在全局块会被这条覆盖 */
  overflow: auto;
}
</style>

<style>
/* JSONEditor 自定义样式 */
/* 覆盖源码自带的亮蓝边框（.jsoneditor { border: thin solid #3883fa }），统一成淡灰 */
/* 注意：不要在这里写 overflow——scoped 块的 .jsoneditor[data-v-xxx] 优先级 (0,2,0)
   高于本文件的 div.jsoneditor (0,1,1)，写了也不会生效，反而误导。overflow 统一由 scoped 块控制 */
div.jsoneditor {
  border: 1px solid #dadce0 !important;
  border-radius: 4px !important;
}

.jsoneditor-contextmenu .jsoneditor-menu {
  background-color: white !important;
  display: block !important;
}
</style>
