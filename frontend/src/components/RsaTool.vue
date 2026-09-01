<template>
  <div class="rsa-tool-container">
    <!-- 配置行 -->
    <div class="config-row">
      <div class="config-item inline">
        <label>密钥长度</label>
        <select v-model="keySize" class="modern-select">
          <option :value="1024">1024 bit</option>
          <option :value="2048">2048 bit</option>
          <option :value="4096">4096 bit</option>
        </select>
      </div>
      <div class="config-item inline">
        <label>密钥格式</label>
        <select v-model="keyFormat" class="modern-select">
          <option value="PKCS#1">PKCS#1</option>
          <option value="PKCS#8">PKCS#8</option>
        </select>
      </div>
      <button class="generate-btn" @click="generateKeyPair">
        生成
      </button>
    </div>

    <!-- 密钥结果区域 -->
    <div v-if="keyPair.publicKey || keyPair.privateKey" class="keys-result">
      <div class="key-section public-key-section">
        <div class="key-header">
          <span class="key-label">RSA加密公钥</span>
          <div class="key-actions">
            <button class="copy-btn" @click="copyToClipboard(keyPair.publicKey)">
              <svg viewBox="0 0 24 24" width="14" height="14">
                <path fill="currentColor" d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"/>
              </svg>
              复制公钥
            </button>
            <button class="copy-btn" @click="copySingleLine">
              <svg viewBox="0 0 24 24" width="14" height="14">
                <path fill="currentColor" d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"/>
              </svg>
              复制单行公钥
            </button>
          </div>
        </div>
        <div class="key-section-inner">
          <pre class="key-content">{{ keyPair.publicKey }}</pre>
        </div>
      </div>

      <div class="key-section">
        <div class="key-header">
          <span class="key-label">RSA加密私钥</span>
          <div class="key-actions">
            <button class="copy-btn" @click="copyToClipboard(keyPair.privateKey)">
              <svg viewBox="0 0 24 24" width="14" height="14">
                <path fill="currentColor" d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"/>
              </svg>
              复制私钥
            </button>
            <button class="copy-btn" @click="copyPrivateKeySingleLine">
              <svg viewBox="0 0 24 24" width="14" height="14">
                <path fill="currentColor" d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"/>
              </svg>
              复制单行私钥
            </button>
          </div>
        </div>
        <div class="key-section-inner">
          <pre class="key-content">{{ keyPair.privateKey }}</pre>
        </div>
      </div>
    </div>

    <!-- 提示消息 -->
    <transition name="fade">
      <div v-if="message.show" :class="['message-toast', message.type]">
        <svg v-if="message.type === 'success'" class="icon" viewBox="0 0 1024 1024" width="16" height="16">
          <path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm193.5 301.7l-210.6 292a31.8 31.8 0 0 1-51.7 0L318.5 484.9c-3.8-5.3 0-12.7 6.5-12.7h46.9c10.2 0 19.9 4.9 25.9 13.3l71.2 98.8 157.2-218c6-8.3 15.6-13.3 25.9-13.3H699c6.5 0 10.3 7.4 6.5 12.7z" fill="#fff"/>
        </svg>
        {{ message.text }}
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  data() {
    return {
      keySize: 2048,
      keyFormat: 'PKCS#1',
      keyPair: {
        publicKey: '',
        privateKey: ''
      },
      message: {
        show: false,
        text: '',
        type: 'success'
      }
    }
  },
  methods: {
    showMessage(text, type = 'success') {
      // 清除上一次的定时器，避免旧提示提前关闭新提示
      if (this.messageTimer) {
        clearTimeout(this.messageTimer)
      }
      this.message = { show: true, text, type }
      this.messageTimer = setTimeout(() => {
        this.message.show = false
      }, 3000)
    },
    async copySingleLine() {
      if (!this.keyPair.publicKey) {
        this.showMessage('没有内容可复制', 'error')
        return
      }
      // 去掉 -----BEGIN/END----- 首尾行和所有换行
      const singleLine = this.keyPair.publicKey
        .replace(/-----[A-Z ]*PUBLIC KEY-----/g, '')
        .replace(/\s+/g, '')
      await this.copyToClipboard(singleLine)
    },
    async copyPrivateKeySingleLine() {
      if (!this.keyPair.privateKey) {
        this.showMessage('没有内容可复制', 'error')
        return
      }
      // 去掉 -----BEGIN/END----- 首尾行和所有换行
      const singleLine = this.keyPair.privateKey
        .replace(/-----[A-Z ]*PRIVATE KEY-----/g, '')
        .replace(/\s+/g, '')
      await this.copyToClipboard(singleLine)
    },
    async copyToClipboard(text) {
      if (!text) {
        this.showMessage('没有内容可复制', 'error')
        return
      }
      let ok = false
      // WebView2 中 clipboard API 可能挂起，加超时兜底
      try {
        await Promise.race([
          navigator.clipboard.writeText(text),
          new Promise((_, reject) => setTimeout(() => reject(new Error('timeout')), 500))
        ])
        ok = true
      } catch (err) {
        // 降级：用隐藏 textarea 复制
        try {
          const textarea = document.createElement('textarea')
          textarea.value = text
          textarea.style.position = 'fixed'
          textarea.style.opacity = '0'
          document.body.appendChild(textarea)
          textarea.select()
          ok = document.execCommand('copy')
          document.body.removeChild(textarea)
        } catch (e) {
          ok = false
        }
      }
      this.showMessage(ok ? '复制成功' : '复制失败', ok ? 'success' : 'error')
    },
    async generateKeyPair() {
      try {
        const result = await window.go.handler.RsaHandler.GenerateKeyPair(
          this.keySize,
          this.keyFormat
        )
        if (result.success) {
          this.keyPair.publicKey = result.publicKey
          this.keyPair.privateKey = result.privateKey
        } else {
          this.showMessage(result.message || '生成失败', 'error')
        }
      } catch (err) {
        this.showMessage(`生成失败: ${err.message}`, 'error')
      }
    }
  }
}
</script>

<style scoped>
.rsa-tool-container {
  margin: 0 auto;
  padding: 24px;
  /* 隐藏滚动条 */
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.rsa-tool-container::-webkit-scrollbar {
  display: none;
}

/* 配置行 */
.config-row {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
  flex-shrink: 0;
  background: #f5f7fa;
  border-radius: 12px;
  padding: 10px 16px;
}

.config-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 120px;
}

.config-item.inline {
  flex-direction: row;
  align-items: center;
  gap: 8px;
}

.config-item label {
  font-size: 13px;
  color: #5f6368;
  font-weight: 500;
  white-space: nowrap;
}

.modern-select {
  padding: 5px 10px;
  border: 1px solid #dadce0;
  border-radius: 4px;
  font-size: 13px;
  color: #3c4043;
  background: white;
  cursor: pointer;
  outline: none;
  min-width: 100px;
  transition: all 0.2s;
}

.modern-select:hover,
.modern-select:focus {
  border-color: #1a73e8;
}

.modern-input {
  padding: 8px 12px;
  border: 1px solid #dadce0;
  border-radius: 8px;
  font-size: 13px;
  color: #3c4043;
  outline: none;
  width: 120px;
  transition: all 0.2s;
}

.password-item .modern-input {
  width: 140px;
}

.modern-input:focus {
  border-color: #1a73e8;
}

.modern-input::placeholder {
  color: #9aa0a6;
}

.generate-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 6px 18px;
  margin-left: 12px;
  background: linear-gradient(135deg, #42a5f5 0%, #1e88e5 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.generate-btn:hover {
  background: linear-gradient(135deg, #1e88e5 0%, #1565c0 100%);
}

/* 密钥结果 */
.keys-result {
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: 580px;
  /* 隐藏滚动条 */
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.keys-result::-webkit-scrollbar {
  display: none;
}

.key-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
  border: 1px solid #e8eaed;
  border-radius: 12px;
  padding: 16px;
  overflow: hidden;
}

/* 公钥较短，占较小高度；私钥较长，占较大高度 */
.public-key-section {
  flex: 1.5;
}

.key-section:not(.public-key-section) {
  flex: 2;
}

.key-section-inner {
  background: white;
  border: 1px solid #dadce0;
  border-radius: 8px;
  padding: 12px;
  flex: 1;
  overflow: auto;
  text-align: left;
  /* 隐藏滚动条 */
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.key-section-inner::-webkit-scrollbar {
  display: none;
}

.key-content {
  margin: 0;
  font-family: 'SF Mono', Monaco, 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #3c4043;
  white-space: pre-wrap;
  word-break: break-all;
  text-align: left;
}

.key-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  flex-shrink: 0;
}

.key-actions {
  display: flex;
  gap: 8px;
}

.key-label {
  font-size: 14px;
  font-weight: 600;
  color: #3c4043;
}

.copy-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: white;
  border: 1px solid #dadce0;
  border-radius: 6px;
  font-size: 12px;
  color: #5f6368;
  cursor: pointer;
  transition: all 0.2s;
}

.copy-btn:hover {
  background: #f1f3f4;
  border-color: #5f6368;
}

/* 消息提示 */
.message-toast {
  position: fixed;
  top: 40px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 20px;
  border-radius: 4px;
  font-size: 14px;
  color: white;
  z-index: 1000;
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  animation: slideDown 0.3s ease;
}

.message-toast.success {
  background: #4CAF50;
}

.message-toast.error {
  background: #f44336;
}

.icon {
  margin-right: 4px;
}

.fade-enter-active {
  transition: all 0.4s ease;
}

.fade-leave-active {
  transition: all 0.6s ease;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-12px);
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

/* 响应式 */
@media (max-width: 768px) {
  .config-row {
    flex-wrap: wrap;
  }
  
  .config-item {
    flex: 1;
    min-width: 100px;
  }
  
  .password-item {
    max-width: none;
  }
  
  .generate-btn {
    width: 100%;
    margin-top: 8px;
  }
}
</style>