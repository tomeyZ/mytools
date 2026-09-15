<template>
  <div class="page fill">
    <!-- 注意：ace.edit() 会重建容器的 innerHTML，提示必须放在容器【外面】的兄弟节点上，
         否则初始化时就被抹掉，且 isEmpty 不变 Vue 也不会重新补回来 -->
    <div class="editor-wrap">
      <div ref="aceEditor" class="ace-host"></div>
      <div class="editor-placeholder" v-if="isEmpty" @mousedown.prevent>
        Ctrl+F 查找 / Ctrl+H 替换（支持正则）<br>
        粘贴或输入文本开始编辑
      </div>
    </div>

    <!-- 底部快捷键提示栏 -->
    <div class="hint-bar">
      <span class="hint-title">快捷键</span>
      <span class="hint-item"><kbd>Ctrl+F</kbd> 查找</span>
      <span class="hint-item"><kbd>Ctrl+H</kbd> 替换（支持正则 <code>.*</code>）</span>
      <span class="hint-item"><kbd>Ctrl+Z</kbd> 撤销</span>
      <span class="hint-item"><kbd>Ctrl+D</kbd> 选中下一个相同词</span>
    </div>
  </div>
</template>

<script>
import ace from 'ace-builds';
import 'ace-builds/src-noconflict/ace';
import 'ace-builds/src-noconflict/mode-text';
import 'ace-builds/src-noconflict/ext-searchbox';
import 'ace-builds/src-noconflict/theme-github';
import 'ace-builds/src-noconflict/theme-github_dark';
import { ACE_THEME, currentTheme, onThemeChange } from '../theme.js';

export default {
  name: 'TextTool',
  data() {
    return {
      isEmpty: true
    };
  },
  mounted() {
    this.editor = ace.edit(this.$refs.aceEditor, {
      mode: 'ace/mode/text',
      // ace 的主题是运行时往文档里插 <style>，不跟 CSS 变量走，只能主动切
      theme: ACE_THEME[currentTheme()] || ACE_THEME.light,
      wrap: true,
      printMargin: false,
      fontSize: 14,
      showPrintMargin: false,
      tabSize: 4,
      useSoftTabs: true
    });
    this.editor.setValue('', -1);

    this.offTheme = onThemeChange((t) => {
      if (this.editor) this.editor.setTheme(ACE_THEME[t] || ACE_THEME.light);
    });

    // 内容变化时切换占位提示
    this.editor.session.on('change', () => {
      this.isEmpty = this.editor.getValue().length === 0;
    });

    this.observeSearchBox();
  },
  beforeUnmount() {
    if (this.offTheme) this.offTheme();
    if (this.observer) this.observer.disconnect();
    if (this.editor) {
      this.editor.destroy();
      this.editor.container.remove();
    }
  },
  methods: {
    // ace 的搜索框是首次 Ctrl+F 时才懒加载创建的，用 MutationObserver 等它出现再改造
    observeSearchBox() {
      const container = this.editor.container;
      this.observer = new MutationObserver(() => {
        const box = container.querySelector('.ace_search');
        if (box && !box.dataset.customized) {
          this.customizeSearchBox(box);
          this.observer.disconnect();
        }
      });
      this.observer.observe(container, { childList: true, subtree: true });
    },
    customizeSearchBox(box) {
      box.dataset.customized = '1';
      const SEARCH_W = 560;
      const rect = this.editor.container.getBoundingClientRect();

      // 加宽 + 默认水平居中、垂直偏上显示
      box.style.right = 'auto';
      box.style.width = SEARCH_W + 'px';
      box.style.left = Math.max(8, (rect.width - SEARCH_W) / 2) + 'px';
      box.style.top = Math.max(8, rect.height * 0.3) + 'px';

      // 整个面板按下即拖动；点在输入框/按钮上则不拖，让原有交互正常工作
      box.addEventListener('mousedown', (e) => {
        if (e.target.closest('input, .ace_button, .ace_searchbtn, .ace_searchbtn_close')) return;
        e.preventDefault();
        const startX = e.clientX, startY = e.clientY;
        const originX = box.offsetLeft, originY = box.offsetTop;
        const onMove = (ev) => {
          box.style.left = originX + ev.clientX - startX + 'px';
          box.style.top = originY + ev.clientY - startY + 'px';
        };
        const onUp = () => {
          document.removeEventListener('mousemove', onMove);
          document.removeEventListener('mouseup', onUp);
        };
        document.addEventListener('mousemove', onMove);
        document.addEventListener('mouseup', onUp);
      });
    }
  }
};
</script>

<style scoped>
/* 编辑器页面要占满内容区高度：外层 .content 有确定高度，这里用 100% 而不是
   calc(100vh - N)，避免顶栏尺寸一变就溢出 */
.fill {
  display: flex;
  flex-direction: column;
  height: 100%;
  box-sizing: border-box;
}

.editor-wrap {
  position: relative;
  display: flex;
  flex: 1;
  min-height: 0;
}

.ace-host {
  flex: 1;
  min-width: 0;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
}

/* 空内容占位提示：覆盖在编辑器上方，不挡点击 */
.editor-placeholder {
  position: absolute;
  top: 12px;
  left: 50px; /* 避开行号槽 */
  pointer-events: none;
  font-size: 14px;
  line-height: 1.9;
  color: var(--text-3);
  z-index: 5;
}

/* 底部提示栏：快捷键面板，kbd 键帽样式，新同事也能一眼看懂 */
.hint-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 4px 20px;
  margin-top: 10px;
  padding: 8px 12px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 12.5px;
  color: var(--text-2);
  flex-shrink: 0;
}

.hint-title {
  font-weight: 600;
  color: var(--text);
}

.hint-bar kbd {
  display: inline-block;
  padding: 1px 7px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-bottom-width: 2px; /* 模拟键帽立体感 */
  border-radius: 4px;
  font-family: var(--font-mono);
  font-size: 11.5px;
  color: var(--text);
  margin-right: 4px;
}

.hint-bar code {
  padding: 0 4px;
  background: var(--surface-3);
  border-radius: 3px;
  font-size: 11.5px;
}

/* ace 查找面板：加宽 + 内部输入框自适应拉伸 */
.page :deep(.ace_search) {
  max-width: none !important;
  border-radius: var(--radius-sm) !important;
  cursor: move;
}

/* 输入框和按钮恢复各自光标 */
.page :deep(.ace_search) input {
  cursor: auto;
}

.page :deep(.ace_button),
.page :deep(.ace_searchbtn),
.page :deep(.ace_searchbtn_close) {
  cursor: pointer;
}

/* 查找/替换表单改 flex，输入框占满剩余宽度（ace 默认 min-width 17em 固定短宽度） */
.page :deep(.ace_search_form),
.page :deep(.ace_replace_form) {
  display: flex;
  align-items: center;
  margin: 4px 36px 4px 12px; /* 右侧留白避开右上角关闭按钮 */
}

.page :deep(.ace_search_field) {
  flex: 1;
  min-width: 0 !important;
  width: auto !important;
}

.page :deep(.ace_searchbtn) {
  flex-shrink: 0;
}
</style>
