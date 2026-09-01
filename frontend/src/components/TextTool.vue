<template>
  <div class="ace-tool-container">
    <div ref="aceEditor" class="ace-editor">
      <!-- 空内容时的快捷键提示 -->
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

export default {
  name: 'TextTool',
  data() {
    return {
      isEmpty: true
    }
  },
  mounted() {
    this.editor = ace.edit(this.$refs.aceEditor, {
      mode: 'ace/mode/text',
      theme: 'ace/theme/github',
      wrap: true,
      printMargin: false,
      fontSize: 14,
      showPrintMargin: false,
      tabSize: 4,
      useSoftTabs: true
    });
    this.editor.setValue('', -1);

    // 内容变化时切换占位提示
    this.editor.session.on('change', () => {
      this.isEmpty = this.editor.getValue().length === 0;
    });

    this.observeSearchBox();
  },
  beforeUnmount() {
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
}
</script>

<style scoped>
.ace-tool-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 40px);
  box-sizing: border-box;
  gap: 8px;
}

.ace-editor {
  position: relative;
  flex: 1;
  border: 1px solid #dadce0;
  border-radius: 8px;
  min-height: 0;
}

/* 空内容占位提示：覆盖在编辑器上方，不挡点击 */
.editor-placeholder {
  position: absolute;
  top: 12px;
  left: 50px; /* 避开行号槽 */
  pointer-events: none;
  font-size: 14px;
  line-height: 1.9;
  color: #9aa0a6;
  z-index: 5;
}

/* 底部提示栏：快捷键面板，kbd 键帽样式，新同事也能一眼看懂 */
.hint-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 4px 20px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 8px;
  font-size: 13px;
  color: #5f6368;
  flex-shrink: 0;
}

.hint-title {
  font-weight: 600;
  color: #3c4043;
}

.hint-bar kbd {
  display: inline-block;
  padding: 1px 7px;
  background: #fff;
  border: 1px solid #dadce0;
  border-bottom-width: 2px; /* 模拟键帽立体感 */
  border-radius: 4px;
  font-family: Consolas, monospace;
  font-size: 12px;
  color: #3c4043;
  margin-right: 4px;
}

.hint-bar code {
  padding: 0 4px;
  background: #e8eaed;
  border-radius: 3px;
  font-size: 12px;
}

/* ace 查找面板：加宽 + 内部输入框自适应拉伸 */
.ace-tool-container :deep(.ace_search) {
  max-width: none !important;
  border-radius: 8px !important;
  cursor: move;
}

/* 输入框和按钮恢复各自光标 */
.ace-tool-container :deep(.ace_search) input {
  cursor: auto;
}

.ace-tool-container :deep(.ace_button),
.ace-tool-container :deep(.ace_searchbtn),
.ace-tool-container :deep(.ace_searchbtn_close) {
  cursor: pointer;
}

/* 查找/替换表单改 flex，输入框占满剩余宽度（ace 默认 min-width 17em 固定短宽度） */
.ace-tool-container :deep(.ace_search_form),
.ace-tool-container :deep(.ace_replace_form) {
  display: flex;
  align-items: center;
  margin: 4px 36px 4px 12px; /* 右侧留白避开右上角关闭按钮 */
}

.ace-tool-container :deep(.ace_search_field) {
  flex: 1;
  min-width: 0 !important;
  width: auto !important;
}

.ace-tool-container :deep(.ace_searchbtn) {
  flex-shrink: 0;
}
</style>
