<template>
  <div class="json-editor-container">
    <!-- 工具条 -->
    <div class="toolbar">
      <button class="tool-btn" @click="formatJson">格式化</button>
      <button class="tool-btn" @click="compactJson">压缩</button>
      <button class="tool-btn" @click="sortJson">排序</button>
      <button class="tool-btn" @click="unescapeAndFormat">去转义</button>
      <button class="tool-btn" @click="escapeJson">转义</button>
      <button class="tool-btn" @click="unicodeToChinese">Unicode转中文</button>
      <span v-if="unescapeInfo" class="unescape-info">{{ unescapeInfo }}</span>
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
    const unescapeInfo = ref('');
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

    // 去转义并格式化：循环 JSON.parse 剥掉字符串化的外层
    const unescapeAndFormat = () => {
      if (!editor) return;
      let text = editor.getText().trim();
      if (!text) {
        unescapeInfo.value = '编辑器内容为空';
        return;
      }

      let layers = 0;
      // 循环解包：parse 成功且结果还是字符串，说明还能再剥一层
      for (let i = 0; i < 10; i++) {
        try {
          const parsed = JSON.parse(text);
          if (typeof parsed === 'string') {
            text = parsed;
            layers++;
          } else {
            // 已是对象/数组，停止
            break;
          }
        } catch (e) {
          // parse 失败：如果已剥过层，当前 text 就是结果；否则说明没有转义可去
          break;
        }
      }

      // 尝试把最终内容格式化展示
      try {
        const obj = JSON.parse(text);
        // stringify 不转义中文，\uXXXX 显示为中文
        editor.setText(JSON.stringify(obj, null, 2));
        unescapeInfo.value = layers > 0 ? `已去除 ${layers} 层转义` : '内容无转义，已格式化';
      } catch (e) {
        if (layers > 0) {
          // 剥完转义后仍不是合法 JSON，按文本展示
          editor.setText(text);
          unescapeInfo.value = `已去除 ${layers} 层转义，但内容不是合法 JSON`;
        } else {
          unescapeInfo.value = '未检测到转义，内容保持不变';
        }
      }
    };

    // Unicode 转中文：把 \uXXXX 显示为中文，不改动其他内容
    const unicodeToChinese = () => {
      if (!editor) return;
      const text = editor.getText();
      if (!text.trim()) {
        unescapeInfo.value = '编辑器内容为空';
        return;
      }
      // 仅替换 \uXXXX 序列，中文之外的转义（如 \n \"）保持原样
      const converted = text.replace(/\\u([0-9a-fA-F]{4})/g, (_, code) => {
        return String.fromCharCode(parseInt(code, 16));
      });
      if (converted === text) {
        unescapeInfo.value = '未检测到 Unicode 转义';
        return;
      }
      editor.setText(converted);
      unescapeInfo.value = '已将 Unicode 转为中文';
    };

    // 格式化当前 JSON
    const formatJson = () => {
      if (!editor) return;
      const text = editor.getText().trim();
      if (!text) {
        unescapeInfo.value = '编辑器内容为空';
        return;
      }
      try {
        const obj = JSON.parse(text);
        editor.setText(JSON.stringify(obj, null, 2));
        unescapeInfo.value = '已格式化';
      } catch (e) {
        unescapeInfo.value = '内容不是有效的 JSON';
      }
    };

    // 压缩 JSON：删除所有多余空格换行
    const compactJson = () => {
      if (!editor) return;
      const text = editor.getText().trim();
      if (!text) {
        unescapeInfo.value = '编辑器内容为空';
        return;
      }
      try {
        const obj = JSON.parse(text);
        editor.setText(JSON.stringify(obj));
        unescapeInfo.value = '已压缩';
      } catch (e) {
        unescapeInfo.value = '内容不是有效的 JSON';
      }
    };

    // 转义：把当前内容变成字符串化形态（与"去转义"相反），便于塞进代码/SQL/curl
    const escapeJson = () => {
      if (!editor) return;
      const text = editor.getText().trim();
      if (!text) {
        unescapeInfo.value = '编辑器内容为空';
        return;
      }
      // 无论当前内容是否合法 JSON，都按原文整体字符串化
      editor.setText(JSON.stringify(text));
      unescapeInfo.value = '已转义（可继续点击叠加多层）';
    };

    // 按 key 排序：所有层级的对象键按字典序排列，便于核对签名参数
    const sortJson = () => {
      if (!editor) return;
      const text = editor.getText().trim();
      if (!text) {
        unescapeInfo.value = '编辑器内容为空';
        return;
      }
      try {
        const obj = JSON.parse(text);
        const sorted = sortValue(obj);
        editor.setText(JSON.stringify(sorted, null, 2));
        unescapeInfo.value = '已按 key 排序';
      } catch (e) {
        unescapeInfo.value = '内容不是有效的 JSON';
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
      unescapeInfo,
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
  padding: 5px 14px;
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

.tool-btn:hover {
  border-color: #1a73e8;
  color: #1a73e8;
}

.unescape-info {
  font-size: 13px;
  color: #5f6368;
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
}
</style>

<style>
/* JSONEditor 自定义样式 */
/* 覆盖源码自带的亮蓝边框（.jsoneditor { border: thin solid #3883fa }），统一成淡灰 */
div.jsoneditor {
  border: 1px solid #dadce0 !important;
  border-radius: 4px !important;
}

/* 隐藏组件自带的顶部绿色菜单栏，功能由自定义工具条替代 */
.jsoneditor-menu {
  display: none !important;
}

.jsoneditor-contextmenu .jsoneditor-menu {
  background-color: white !important;
  display: block !important;
}

.jsoneditor-statusbar {
  display: none !important;
}
</style>
