<template>
  <div class="page">
    <div class="grid">
      <!-- 加密配置 -->
      <section class="card wide">
        <div class="card-head">
          <span class="ci">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
              <line x1="4" y1="8" x2="20" y2="8"/><line x1="4" y1="16" x2="20" y2="16"/>
              <circle cx="9" cy="8" r="2.2"/><circle cx="15" cy="16" r="2.2"/>
            </svg>
          </span>
          <h3>加密配置</h3>
        </div>
        <div class="card-body">
          <div class="row calc-row">
            <div class="fld">
              <span class="field-label">加密模式</span>
              <AppSelect v-model="aesMode" :options="aesModeOptions"/>
            </div>
            <div class="fld fld-key">
              <span class="field-label">密钥</span>
              <input type="text" v-model="key" placeholder="请输入加密密钥" class="inp">
            </div>
          </div>
        </div>
      </section>

      <!-- 输入 -->
      <section class="card wide">
        <div class="card-head">
          <span class="ci ci-2">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 3.5H6.5a2 2 0 0 0-2 2v13a2 2 0 0 0 2 2h11a2 2 0 0 0 2-2V9z"/>
              <path d="M12 3.5V9h5.5"/>
            </svg>
          </span>
          <h3>输入内容</h3>
        </div>
        <div class="card-body">
          <textarea
              v-model="inputText"
              rows="4"
              placeholder="请输入要加密或解密的文本内容…"
              class="inp mono"
          ></textarea>
          <div class="row">
            <button class="btn primary" @click="encrypt" :disabled="isLoading">
              <svg viewBox="0 0 24 24" width="15" height="15" fill="currentColor">
                <path d="M12,17A2,2 0 0,0 14,15C14,13.89 13.1,13 12,13A2,2 0 0,0 10,15A2,2 0 0,0 12,17M18,8A2,2 0 0,1 20,10V20A2,2 0 0,1 18,22H6A2,2 0 0,1 4,20V10C4,8.89 4.9,8 6,8H7V6A5,5 0 0,1 12,1A5,5 0 0,1 17,6V8H18M12,3A3,3 0 0,0 9,6V8H15V6A3,3 0 0,0 12,3Z"/>
              </svg>
              {{ isLoading ? '处理中…' : '加密' }}
            </button>
            <button class="btn soft" @click="decrypt" :disabled="isLoading">
              <svg viewBox="0 0 24 24" width="15" height="15" fill="currentColor">
                <path d="M12,17A2,2 0 0,0 14,15C14,13.89 13.1,13 12,13A2,2 0 0,0 10,15A2,2 0 0,0 12,17M6,8V20A2,2 0 0,0 8,22H16A2,2 0 0,0 18,20V8H6M7,6A5,5 0 0,1 12,1A5,5 0 0,1 17,6V8H18A2,2 0 0,1 20,10V20A2,2 0 0,1 18,22H6A2,2 0 0,1 4,20V10A2,2 0 0,1 6,8H7V6Z"/>
              </svg>
              {{ isLoading ? '处理中…' : '解密' }}
            </button>
            <button class="btn ghost" @click="handleClear">
              <svg viewBox="0 0 24 24" width="15" height="15" fill="currentColor">
                <path d="M19,4H15.5L14.5,3H9.5L8.5,4H5V6H19M6,19A2,2 0 0,0 8,21H16A2,2 0 0,0 18,19V7H6V19Z"/>
              </svg>
              清空
            </button>
          </div>
        </div>
      </section>

      <!-- 错误提示 -->
      <section class="card wide" v-if="errorMessage">
        <div class="out out-error"><span class="val">{{ errorMessage }}</span></div>
      </section>

      <!-- 结果 -->
      <section class="card wide" v-if="showResult">
        <div class="card-head">
          <span class="ci ci-3">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
              <path d="M20 6 9 17l-5-5"/>
            </svg>
          </span>
          <h3>{{ resultLabel }}</h3>
          <div class="card-head-actions">
            <button class="copy" title="复制" @click="copyResult">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
                <rect x="9" y="9" width="12" height="12" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
              </svg>
              复制
            </button>
          </div>
        </div>
        <div class="card-body">
          <textarea v-model="result" readonly class="inp mono result-view" rows="6"></textarea>
        </div>
      </section>
    </div>

    <transition name="toast-fade">
      <div v-if="toast.show" :class="['message-toast', toast.type]">{{ toast.text }}</div>
    </transition>
  </div>
</template>

<script>
import { copyText } from '../clipboard.js';
import { toastMixin } from '../toast.js';
import AppSelect from './AppSelect.vue';

export default {
  name: 'AesTool',
  components: { AppSelect },
  mixins: [toastMixin],
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
      errorTimer: null
    };
  },
  computed: {
    aesModeOptions() {
      return [
        { value: 'aes128-cbc', label: 'AES-128-CBC' },
        { value: 'aes256-cbc', label: 'AES-256-CBC' },
        { value: 'aes128-ecb', label: 'AES-128-ECB' }
      ];
    },
    resultLabel() {
      return this.resultType === 'encrypt' ? '加密结果' : '解密结果';
    }
  },
  beforeUnmount() {
    if (this.errorTimer) clearTimeout(this.errorTimer);
  },
  methods: {
    async encrypt() {
      this.clearError();
      if (!this.validateInput()) return;

      this.isLoading = true;
      try {
        const response = await window.go.handler.CryptoHandler.Encrypt(
          this.aesMode,
          this.key,
          this.inputText
        );
        if (response.success) {
          this.result = response.result;
          this.resultType = 'encrypt';
          this.showResult = true;
        } else {
          this.showError(response.message);
        }
      } catch (error) {
        console.error('加密失败:', error);
        this.showError('加密失败，请稍后再试');
      } finally {
        this.isLoading = false;
      }
    },

    async decrypt() {
      this.clearError();
      if (!this.validateInput()) return;

      this.isLoading = true;
      try {
        const response = await window.go.handler.CryptoHandler.Decrypt(
          this.aesMode,
          this.key,
          this.inputText
        );
        if (response.success) {
          this.result = response.result;
          this.resultType = 'decrypt';
          this.showResult = true;
        } else {
          this.showError(response.message);
        }
      } catch (error) {
        console.error('解密失败:', error);
        this.showError('解密失败，请稍后再试');
      } finally {
        this.isLoading = false;
      }
    },

    validateInput() {
      if (!this.inputText.trim()) {
        this.showError('请输入要处理的内容');
        return false;
      }
      if (!this.key.trim()) {
        this.showError('请输入密钥');
        return false;
      }
      return true;
    },

    showError(message) {
      if (this.errorTimer) clearTimeout(this.errorTimer);
      this.errorMessage = message;
      this.errorTimer = setTimeout(() => {
        this.errorMessage = '';
      }, 3000);
    },

    clearError() {
      if (this.errorTimer) clearTimeout(this.errorTimer);
      this.errorMessage = '';
    },

    async copyResult() {
      if (!this.result) return;
      const ok = await copyText(this.result);
      // resultLabel 是「加密结果」/「解密结果」，文案带对象
      this.showToast(ok ? this.resultLabel + '已复制' : '复制失败', ok ? 'success' : 'error');
    },

    handleClear() {
      this.inputText = '';
      this.key = '';
      this.result = '';
      this.resultType = '';
      this.showResult = false;
      this.clearError();
    }
  }
};
</script>

<style scoped>
.fld-key {
  flex: 1;
  min-width: 220px;
}
</style>
