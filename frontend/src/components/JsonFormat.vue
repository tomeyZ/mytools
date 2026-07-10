<template>
  <div class="json-editor-container">
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

    // 初始化 JSONEditor
    onMounted(() => {
      const container = jsoneditor.value;
      const options = {
        mode: 'code',
        modes: ['code', 'tree', 'form', 'text', 'view'],
        onError: (err) => {
          error.value = err.toString();
        },
        onChange: () => {
          try {
            editor.get();
            error.value = '';
          } catch (e) {
            error.value = '当前内容不是有效的 JSON';
          }
        },
        onModeChange: (newMode) => {
          if (newMode === 'code') {
            editor.aceEditor.setOptions({
              maxLines: Infinity
            });
          }
        }
      };

      editor = new JSONEditor(container, options);

      // 设置初始空对象
      editor.set();
    });

    // 销毁 JSONEditor
    onBeforeUnmount(() => {
      if (editor) {
        editor.destroy();
      }
    });

    return {
      jsoneditor,
      error
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
  border: 1px solid #ddd;
  border-radius: 4px;
}
</style>

<style>
/* JSONEditor 自定义样式 */
.jsoneditor-menu {
  background-color: #4CAF50 !important;
  border-bottom: none !important;
}

.jsoneditor-contextmenu .jsoneditor-menu {
  background-color: white !important;
}

.jsoneditor-statusbar {
  display: none !important;
}

/* 代码模式下的编辑器高度 */
.ace_editor {
  min-height: 300px;
}
</style>