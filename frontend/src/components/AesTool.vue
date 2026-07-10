<template>
  <div class="aes-tool-container">
<!--    <div class="tool-header">-->
<!--      <h2>AES 加解密工具</h2>-->
<!--    </div>-->

    <div class="config-section">
      <div class="form-row">
        <div class="form-group">
          <label for="aesMode">加密模式</label>
          <select id="aesMode" v-model="aesMode" class="form-select">
            <option value="aes128-cbc">AES-128-CBC</option>
            <option value="aes256-cbc">AES-256-CBC</option>
            <option value="aes128-ecb">AES-128-ECB</option>
          </select>
        </div>

        <div class="form-group">
          <label for="key">密钥</label>
          <input
              type="text"
              id="key"
              v-model="key"
              placeholder="请输入加密密钥"
              class="form-input"
          />
        </div>
      </div>
    </div>

    <div class="input-section">
      <div class="form-group">
        <label for="inputText">输入内容</label>
        <textarea
            id="inputText"
            v-model="inputText"
            rows="5"
            placeholder="请输入要加密或解密的文本内容..."
            class="form-textarea"
        ></textarea>
      </div>
    </div>

    <div class="action-buttons">
      <button @click="encrypt" class="encrypt-btn" :disabled="isLoading">
        <svg viewBox="0 0 24 24" width="18" height="18">
          <path fill="currentColor" d="M12,17A2,2 0 0,0 14,15C14,13.89 13.1,13 12,13A2,2 0 0,0 10,15A2,2 0 0,0 12,17M18,8A2,2 0 0,1 20,10V20A2,2 0 0,1 18,22H6A2,2 0 0,1 4,20V10C4,8.89 4.9,8 6,8H7V6A5,5 0 0,1 12,1A5,5 0 0,1 17,6V8H18M12,3A3,3 0 0,0 9,6V8H15V6A3,3 0 0,0 12,3Z"/>
        </svg>
        {{ isLoading ? '处理中...' : '加密' }}
      </button>
      <button @click="decrypt" class="decrypt-btn" :disabled="isLoading">
        <svg viewBox="0 0 24 24" width="18" height="18">
          <path fill="currentColor" d="M12,17A2,2 0 0,0 14,15C14,13.89 13.1,13 12,13A2,2 0 0,0 10,15A2,2 0 0,0 12,17M6,8V20A2,2 0 0,0 8,22H16A2,2 0 0,0 18,20V8H6M7,6A5,5 0 0,1 12,1A5,5 0 0,1 17,6V8H18A2,2 0 0,1 20,10V20A2,2 0 0,1 18,22H6A2,2 0 0,1 4,20V10A2,2 0 0,1 6,8H7V6Z"/>
        </svg>
        {{ isLoading ? '处理中...' : '解密' }}
      </button>
      <button @click="handleClear" class="clear-btn">
        <svg viewBox="0 0 24 24" width="18" height="18">
          <path fill="currentColor" d="M19,4H15.5L14.5,3H9.5L8.5,4H5V6H19M6,19A2,2 0 0,0 8,21H16A2,2 0 0,0 18,19V7H6V19Z"/>
        </svg>
        清空
      </button>
    </div>

    <div class="result-section" v-if="showResult">
      <div class="result-header">
        <span class="result-label">{{ resultLabel }}</span>
        <button class="copy-result-btn" @click="copyResult">
          <svg viewBox="0 0 24 24" width="16" height="16">
            <path fill="currentColor" d="M19,21H8V7H19M19,5H8A2,2 0 0,0 6,7V21A2,2 0 0,0 8,23H19A2,2 0 0,0 21,21V7A2,2 0 0,0 19,5M16,1H4A2,2 0 0,0 2,3V17H4V3H16V1Z"/>
          </svg>
          复制
        </button>
      </div>
      <textarea id="result" v-model="result" readonly class="result-textarea"></textarea>
    </div>

    <!-- 错误提示 -->
    <transition name="fade">
      <div v-if="errorMessage" class="error-message">
        <svg viewBox="0 0 24 24" width="20" height="20">
          <path fill="currentColor" d="M13,13H11V7H13M13,17H11V15H13M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z"/>
        </svg>
        {{ errorMessage }}
      </div>
    </transition>

    <!-- 复制成功提示 -->
    <transition name="fade">
      <div v-if="showCopySuccess" class="copy-success-message">
        <svg viewBox="0 0 24 24" width="18" height="18">
          <path fill="currentColor" d="M21,7L9,19L3.5,13.5L4.91,12.09L9,16.17L19.59,5.59L21,7Z"/>
        </svg>
        复制成功
      </div>
    </transition>
  </div>
</template>

<script>
import { ref, computed } from 'vue'

export default {
  setup() {
    const showCopySuccess = ref(false)
    let copyTimer = null

    const showCopyMessage = () => {
      showCopySuccess.value = true
      if (copyTimer) clearTimeout(copyTimer)
      copyTimer = setTimeout(() => {
        showCopySuccess.value = false
      }, 2000)
    }

    return {
      showCopySuccess,
      showCopyMessage
    }
  },
  data() {
    return {
      aesMode: 'aes128-cbc',
      key: '',
      inputText: '',
      result: '',
      resultType: '',
      isLoading: false,
      errorMessage: '',
      showResult: false,
    }
  },
  computed: {
    resultLabel() {
      return this.resultType === 'encrypt' ? '加密结果' : '解密结果'
    }
  },
  methods: {
    async encrypt() {
      this.clearError()
      if (!this.validateInput()) return

      this.isLoading = true
      try {
        const response = await window.go.handler.CryptoHandler.Encrypt(
          this.aesMode,
          this.key,
          this.inputText
        )

        if (response.success) {
          this.result = response.result
          this.resultType = 'encrypt'
          this.showResult = true
        } else {
          this.showError(response.message)
        }
      } catch (error) {
        console.error('加密失败:', error)
        this.showError('加密失败，请稍后再试')
      } finally {
        this.isLoading = false
      }
    },

    async decrypt() {
      this.clearError()
      if (!this.validateInput()) return

      this.isLoading = true
      try {
        const response = await window.go.handler.CryptoHandler.Decrypt(
          this.aesMode,
          this.key,
          this.inputText
        )

        if (response.success) {
          this.result = response.result
          this.resultType = 'decrypt'
          this.showResult = true
        } else {
          this.showError(response.message)
        }
      } catch (error) {
        console.error('解密失败:', error)
        this.showError('解密失败，请稍后再试')
      } finally {
        this.isLoading = false
      }
    },

    validateInput() {
      if (!this.inputText.trim()) {
        this.showError('请输入要处理的内容')
        return false
      }
      if (!this.key.trim()) {
        this.showError('请输入密钥')
        return false
      }
      return true
    },

    showError(message) {
      this.errorMessage = message
      setTimeout(() => {
        this.errorMessage = ''
      }, 3000)
    },

    clearError() {
      this.errorMessage = ''
    },

    async copyResult() {
      if (!this.result) return
      try {
        await navigator.clipboard.writeText(this.result)
        this.showCopyMessage()
      } catch (err) {
        console.error('复制失败:', err)
        this.showError('复制失败，请手动复制')
      }
    },

    handleClear() {
      this.inputText = ''
      this.key = ''
      this.result = ''
      this.resultType = ''
      this.showResult = false
      this.clearError()
    }
  }
}
</script>

<style scoped>
.aes-tool-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px 24px;
  background: #fff;
  border-radius: 8px;
  font-family: "思源宋体", "Noto Serif SC", serif;
}

.aes-tool-container * {
  font-family: inherit;
}

.tool-header {
  text-align: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 2px solid #f0f0f0;
}

.tool-header h2 {
  margin: 0;
  font-size: 22px;
  font-weight: 700;
  color: #2c3e50;
}

.config-section {
  margin-bottom: 20px;
  background: #f8f9fa;
  padding: 16px;
  border-radius: 6px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 16px;
  align-items: start;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-weight: 600;
  color: #2c3e50;
  font-size: 14px;
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.3s ease;
  background: #fff;
  box-sizing: border-box;
  font-family: inherit;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 120px;
  line-height: 1.6;
}

.input-section {
  margin-bottom: 20px;
}

.action-buttons {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.action-buttons button {
  flex: 1;
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.encrypt-btn {
  background-color: #3498db;
  color: white;
}

.encrypt-btn:hover:not(:disabled) {
  background-color: #2980b9;
  transform: translateY(-1px);
}

.decrypt-btn {
  background-color: #27ae60;
  color: white;
}

.decrypt-btn:hover:not(:disabled) {
  background-color: #219a52;
  transform: translateY(-1px);
}

.clear-btn {
  background-color: #95a5a6;
  color: white;
  flex: 0 0 auto;
}

.clear-btn:hover:not(:disabled) {
  background-color: #7f8c8d;
  transform: translateY(-1px);
}

.action-buttons button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.result-section {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 6px;
  border-left: 4px solid #3498db;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.result-label {
  font-weight: 600;
  color: #2c3e50;
  font-size: 15px;
}

.copy-result-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.copy-result-btn:hover {
  background-color: #2980b9;
}

.result-textarea {
  width: 100%;
  background: #fff;
  border: 2px solid #e9ecef;
  color: #495057;
  border-radius: 6px;
  padding: 12px;
  min-height: 200px;
  max-height: 350px;
  box-sizing: border-box;
  line-height: 1.6;
  resize: vertical;
  font-family: "Consolas", "Monaco", "Courier New", monospace;
  font-size: 13px;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: #e74c3c;
  color: white;
  padding: 12px 16px;
  border-radius: 6px;
  margin-top: 16px;
  font-size: 14px;
  box-shadow: 0 4px 12px rgba(231, 76, 60, 0.3);
  animation: slideDown 0.3s ease-out;
}

.copy-success-message {
  position: fixed;
  right: 40px;
  top: 40px;
  background-color: #27ae60;
  color: white;
  padding: 12px 20px;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.3);
  z-index: 1000;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  animation: slideIn 0.3s ease-out forwards;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 768px) {
  .aes-tool-container {
    padding: 20px 16px;
    margin: 10px;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .action-buttons {
    flex-direction: column;
  }

  .form-input,
  .form-select,
  .form-textarea {
    padding: 10px 12px;
  }

  .tool-header h2 {
    font-size: 20px;
  }
}
</style>
