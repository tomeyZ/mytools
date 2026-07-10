<template>
  <div class="md5-container">
    <h2>MD5加密工具</h2>
    <div class="input-area">
      <textarea
          v-model="inputText"
          class="input-textarea"
          placeholder="请输入要加密的内容"
          rows="5"
      ></textarea>
    </div>
    <div class="btn-group">
      <button class="md5-btn" @click="handleMd5Encrypt">MD5加密</button>
      <button class="clear-btn" @click="handleClear">清空</button>
    </div>
    <div class="result-area" v-if="showResults">
      <div class="result-item">
        <label>32位大写：</label>
        <span>{{ result32Upper }}</span>
        <button class="copy-btn" @click="copyText(result32Upper)">复制</button>
      </div>
      <div class="result-item">
        <label>32位小写：</label>
        <span>{{ result32Lower }}</span>
        <button class="copy-btn" @click="copyText(result32Lower)">复制</button>
      </div>
      <div class="result-item">
        <label>16位大写：</label>
        <span>{{ result16Upper }}</span>
        <button class="copy-btn" @click="copyText(result16Upper)">复制</button>
      </div>
      <div class="result-item">
        <label>16位小写：</label>
        <span>{{ result16Lower }}</span>
        <button class="copy-btn" @click="copyText(result16Lower)">复制</button>
      </div>
    </div>
    <!-- 复制成功提示 -->
    <transition name="fade">
      <div v-if="showCopySuccess" class="copy-success-message">
        <svg class="icon" viewBox="0 0 1024 1024" width="16" height="16">
          <path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm193.5 301.7l-210.6 292a31.8 31.8 0 0 1-51.7 0L318.5 484.9c-3.8-5.3 0-12.7 6.5-12.7h46.9c10.2 0 19.9 4.9 25.9 13.3l71.2 98.8 157.2-218c6-8.3 15.6-13.3 25.9-13.3H699c6.5 0 10.3 7.4 6.5 12.7z" fill="#fff"/>
        </svg>
        复制成功
      </div>
    </transition>
  </div>
</template>

<script>
import {computed, ref} from 'vue';
import md5 from 'js-md5';

export default {
  setup() {
    const showCopySuccess = ref(false);
    let copyTimer = null;

    const showCopyMessage = () => {
      showCopySuccess.value = true;
      if (copyTimer) clearTimeout(copyTimer);
      copyTimer = setTimeout(() => {
        showCopySuccess.value = false;
      }, 2000);
    };

    return {
      showCopySuccess,
      showCopyMessage
    };
  },
  data() {
    return {
      inputText: '',
      result32Upper: '',
      result32Lower: '',
      result16Upper: '',
      result16Lower: '',
      showResults: false,
    };
  },
  methods: {
    handleMd5Encrypt() {
      if (!this.inputText.trim()) return;
      const md5Result = md5(this.inputText);
      this.result32Upper = md5Result.toUpperCase();
      this.result32Lower = md5Result.toLowerCase();
      this.result16Upper = md5Result.slice(8, 24).toUpperCase();
      this.result16Lower = md5Result.slice(8, 24).toLowerCase();
      this.showResults = true;
    },
    handleClear() {
      this.inputText = '';
      this.result32Upper = '';
      this.result32Lower = '';
      this.result16Upper = '';
      this.result16Lower = '';
      this.showResults = false;
    },
    async copyText(text) {
      try {
        await navigator.clipboard.writeText(text);
        this.showCopyMessage();
      } catch (err) {
        console.error('复制失败:', err);
      }
    }
  }
};
</script>

<style scoped>
.md5-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 24px;
  background-color: #fff;
  border-radius: 8px;
  position: relative;
  font-family: "思源宋体", "Noto Serif SC", serif;
}

.md5-container {
  font-family: inherit;
}

h2 {
  color: #2c3e50;
  margin-bottom: 20px;
}
.input-area {
  margin-bottom: 15px;
}
.input-textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  resize: none;
  box-sizing: border-box;
  font-size: 14px;
}
.btn-group {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}
.md5-btn {
  padding: 8px 16px;
  background-color: #3498db;
  color: #fff;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.2s;
}
.md5-btn:hover {
  background-color: #2980b9;
}
.clear-btn {
  padding: 8px 16px;
  background-color: #eee;
  color: #333;
  border: 1px solid #ddd;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.2s;
}
.clear-btn:hover {
  background-color: #ddd;
}
.result-area {
  border-top: 1px solid #eee;
  padding-top: 15px;
}
.result-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}
.result-item label {
  font-weight: bold;
  color: #555;
}
.result-item span {
  flex: 1;
  margin: 0 10px;
  word-break: break-all;
}
.copy-btn {
  padding: 4px 8px;
  background-color: #3498db;
  color: #fff;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  transition: background-color 0.2s;
}
.copy-btn:hover {
  background-color: #2980b9;
}

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
  font-size: 14px;
  animation: slideIn 0.3s ease-out forwards;
}

.icon {
  margin-right: 4px;
}

.fade-enter-active, .fade-leave-active {
  transition: all 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}
</style>