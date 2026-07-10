<template>
  <div class="json-format-container">
    <div class="toolbar">
      <button @click="formatJson" class="format-btn">格式化</button>
      <button @click="compressJson" class="compress-btn">压缩</button>
      <button @click="copyJson" class="copy-btn">复制</button>
      <button @click="clearAll" class="clear-btn">清空</button>
    </div>
    <transition name="fade">
      <div v-if="showCopySuccess" class="copy-success-message">
        <svg class="icon" viewBox="0 0 1024 1024" width="16" height="16">
          <path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm193.5 301.7l-210.6 292a31.8 31.8 0 0 1-51.7 0L318.5 484.9c-3.8-5.3 0-12.7 6.5-12.7h46.9c10.2 0 19.9 4.9 25.9 13.3l71.2 98.8 157.2-218c6-8.3 15.6-13.3 25.9-13.3H699c6.5 0 10.3 7.4 6.5 12.7z" fill="#fff"/>
        </svg>
        复制成功
      </div>
    </transition>

    <div class="editor-wrapper">
      <div class="editor-area">
        <div class="input-area">
          <textarea
              v-model="inputJson"
              class="json-textarea"
              placeholder="请输入或粘贴 JSON 数据..."
              ref="textarea"
              @input="autoFormatJson"
          ></textarea>
        </div>

        <div class="output-area">
          <vue-json-pretty
              v-if="formattedJson"
              :data="formattedJson"
              :deep="3"
              :showLength="true"
              class="json-pretty"
          />
          <div v-else class="placeholder">
            <p>格式化后的 JSON 将显示在这里</p>
          </div>
        </div>
      </div>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';

export default {
  components: {
    VueJsonPretty
  },
  setup() {
    const showCopySuccess = ref(false);
    let copySuccessTimer = null;
    const inputJson = ref('');
    const formattedJson = ref(null);
    const error = ref('');
    const textarea = ref(null);

    // 自动格式化JSON（带防抖）
    const autoFormatJson = debounce(() => {
      try {
        if (!inputJson.value.trim()) {
          formattedJson.value = null;
          error.value = '';
          return;
        }

        const parsed = JSON.parse(inputJson.value);
        formattedJson.value = parsed;
        error.value = '';
      } catch (e) {
        // 自动格式化时不显示错误，留待手动格式化时显示
      }
    }, 500);

    // 手动格式化
    const formatJson = () => {
      try {
        if (!inputJson.value.trim()) {
          throw new Error('请输入 JSON 数据');
        }

        const parsed = JSON.parse(inputJson.value);
        formattedJson.value = parsed;
        error.value = '';
      } catch (e) {
        error.value = `JSON 解析错误: ${e.message}`;
        formattedJson.value = null;
      }
    };

    // 压缩JSON
    const compressJson = () => {
      try {
        if (!inputJson.value.trim()) {
          throw new Error('请输入 JSON 数据');
        }

        const parsed = JSON.parse(inputJson.value);
        inputJson.value = JSON.stringify(parsed);
        formattedJson.value = parsed;
        error.value = '';
      } catch (e) {
        error.value = `JSON 解析错误: ${e.message}`;
        formattedJson.value = null;
      }
    };

    // 复制JSON
    const copyJson = async () => {
      try {
        if (!formattedJson.value) {
          throw new Error('没有可复制的 JSON 数据');
        }

        await navigator.clipboard.writeText(JSON.stringify(formattedJson.value, null, 2));

        // 显示复制成功提示
        showCopySuccess.value = true;

        // 3秒后自动隐藏
        if (copySuccessTimer) clearTimeout(copySuccessTimer);
        copySuccessTimer = setTimeout(() => {
          showCopySuccess.value = false;
        }, 2000);

      } catch (e) {
        error.value = `复制失败: ${e.message}`;
      }
    };

    // 清空
    const clearAll = () => {
      inputJson.value = '';
      formattedJson.value = null;
      error.value = '';
    };

    // 防抖函数
    function debounce(fn, delay) {
      let timer = null;
      return function() {
        if (timer) clearTimeout(timer);
        timer = setTimeout(() => {
          fn.apply(this, arguments);
        }, delay);
      };
    }

    return {
      showCopySuccess,
      inputJson,
      formattedJson,
      error,
      textarea,
      autoFormatJson,
      formatJson,
      compressJson,
      copyJson,
      clearAll
    };
  }
};
</script>

<style scoped>
.copy-success-message {
  position: fixed;
  right: 40px;
  top: 40px;
  background-color: #4CAF50;
  color: white;
  padding: 10px 20px;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  display: flex;
  align-items: center;
  //gap: 2px;
  font-size: 14px;
}

.icon {
  margin-right: 4px;
}

/* 淡入淡出动画 */
.fade-enter-active, .fade-leave-active {
  transition: all 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

.json-format-container {
  padding: 20px;
  max-width: 100%;
  margin: 0 auto;
  font-family: Arial, sans-serif;
  height: calc(100vh - 60px);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

h1 {
  color: #333;
  text-align: center;
  margin-bottom: 20px;
  white-space: nowrap;
}

.toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  justify-content: center;
  flex-wrap: nowrap;
}

.toolbar button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: background-color 0.3s;
  white-space: nowrap;
}

.format-btn {
  background-color: #4CAF50;
  color: white;
}

.compress-btn {
  background-color: #2196F3;
  color: white;
}

.copy-btn {
  background-color: #FF9800;
  color: white;
}

.clear-btn {
  background-color: #f44336;
  color: white;
}

.toolbar button:hover {
  opacity: 0.9;
}

.editor-wrapper {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.editor-area {
  display: flex;
  height: 100%;
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
}

.input-area {
  width: 40%;
  padding-right: 10px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.json-textarea {
  width: 100%;
  //padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: monospace;
  resize: none;
  flex: 1;
  min-height: 200px;
}

.output-area {
  width: 60%;
  //height: 100%;
  display: flex;
  flex-direction: column;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: #f9f9f9;
  overflow: auto;
}

.placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #888;
  padding: 20px;
}

.error-message {
  margin-top: 15px;
  padding: 10px;
  background-color: #ffebee;
  color: #f44336;
  border-radius: 4px;
  font-family: monospace;
  white-space: normal;
}

.json-pretty {
  flex: 1;
  padding: 10px;
  font-family: monospace;
  font-size: 14px;
  line-height: 1.5;
  overflow: auto;
  height: 100%;
}
</style>

<style>
::-webkit-scrollbar {
  width: 0;
  height: 0;
  background: transparent;
}
/* vue-json-pretty 样式优化 */
.vjs-tree {
  min-width: fit-content;
}

.vjs-tree-node {
  white-space: nowrap;
}

.vjs-key {
  color: black;
}

.vjs-value-string {
  color: #0b6202;
}

.vjs-value-number {
  color: #1a01cc;
  font-weight: bold;
}

.vjs-value-boolean {
  color: #3498db;
}

.vjs-value-null {
  color: #f44336;
  font-weight: bold;
}
</style>