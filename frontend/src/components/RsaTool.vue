<template>
  <div class="page">
    <div class="grid">
      <!-- 配置 -->
      <section class="card wide">
        <!-- 这一卡只放一行控件，没有图标和标题。
             留标题时标题与控件之间会空出 200px+（控件为了和下方密钥卡的按钮对齐而被推到右侧），
             比"干脆不写标题"更难看；标题本身也不承担信息（整页只有这一组配置）。
             卡片本身保留 —— 它和下面两张密钥卡构成"配置 → 结果"的层次 -->
        <div class="cfg-row">
          <div class="fld-inline">
            <span class="field-label">密钥长度</span>
            <AppSelect v-model="keySize" :options="keySizeOptions" min-width="160px"/>
          </div>
          <div class="fld-inline">
            <span class="field-label">密钥格式</span>
            <AppSelect v-model="keyFormat" :options="keyFormatOptions" min-width="160px"/>
          </div>
          <!-- 生成 = 钥匙（Material vpn_key，填充式 15px，与 AES 页按钮图标同体系）；
               清空按钮与 AES 页的清空完全同款（.btn ghost + 垃圾桶），
               清掉密钥后两张密钥卡随 v-if 整体消失，不弹提示（AES 同此） -->
          <button class="btn primary" @click="generateKeyPair">
            <svg viewBox="0 0 24 24" width="15" height="15" fill="currentColor">
              <path d="M12.65 10C11.83 7.67 9.61 6 7 6c-3.31 0-6 2.69-6 6s2.69 6 6 6c2.61 0 4.83-1.67 5.65-4H17v4h4v-4h2v-4H12.65zM7 14c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2z"/>
            </svg>
            生成密钥
          </button>
          <button class="btn ghost" @click="clearKeyPair">
            <svg viewBox="0 0 24 24" width="15" height="15" fill="currentColor">
              <path d="M19,4H15.5L14.5,3H9.5L8.5,4H5V6H19M6,19A2,2 0 0,0 8,21H16A2,2 0 0,0 18,19V7H6V19Z"/>
            </svg>
            清空
          </button>
        </div>
      </section>

      <template v-if="keyPair.publicKey || keyPair.privateKey">
        <!-- 公钥 -->
        <section class="card wide">
          <div class="card-head">
            <span class="ci ci-2">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
                <rect x="4" y="10.5" width="16" height="10.5" rx="2.2"/><path d="M8 10.5V7a4 4 0 0 1 8 0v3.5"/>
              </svg>
            </span>
            <h3>公钥 (Public Key)</h3>
            <div class="card-head-actions">
              <button class="copy" @click="copyToClipboard(keyPair.publicKey, '公钥')">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="9" y="9" width="12" height="12" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                </svg>
                复制公钥
              </button>
              <button class="copy" @click="copySingleLine">复制单行公钥</button>
            </div>
          </div>
          <div class="card-body">
            <pre class="key-view mono">{{ keyPair.publicKey }}</pre>
          </div>
        </section>

        <!-- 私钥 -->
        <section class="card wide">
          <div class="card-head">
            <span class="ci ci-3">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 21.5s7.5-3.6 7.5-9.3V5.6L12 2.8 4.5 5.6v6.6c0 5.7 7.5 9.3 7.5 9.3z"/>
                <path d="M9.2 12l2 2 3.6-3.6"/>
              </svg>
            </span>
            <h3>私钥 (Private Key)</h3>
            <div class="card-head-actions">
              <button class="copy" @click="copyToClipboard(keyPair.privateKey, '私钥')">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="9" y="9" width="12" height="12" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                </svg>
                复制私钥
              </button>
              <button class="copy" @click="copyPrivateKeySingleLine">复制单行私钥</button>
            </div>
          </div>
          <div class="card-body">
            <pre class="key-view mono">{{ keyPair.privateKey }}</pre>
          </div>
        </section>
      </template>
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
  name: 'RsaTool',
  components: { AppSelect },
  mixins: [toastMixin],
  data() {
    return {
      keySize: 2048,
      keyFormat: 'PKCS#1',
      keyPair: {
        publicKey: '',
        privateKey: ''
      }
    };
  },
  computed: {
    // 下拉选项。keySize 的值保持 Number 原样传给 Go 后端
    keySizeOptions() {
      return [
        { value: 1024, label: '1024 bit' },
        { value: 2048, label: '2048 bit' },
        { value: 4096, label: '4096 bit' }
      ];
    },
    keyFormatOptions() {
      return [
        { value: 'PKCS#1', label: 'PKCS#1' },
        { value: 'PKCS#8', label: 'PKCS#8' }
      ];
    }
  },
  methods: {
    async copySingleLine() {
      if (!this.keyPair.publicKey) {
        this.showToast('没有内容可复制', 'error');
        return;
      }
      // 去掉 -----BEGIN/END----- 首尾行和所有换行
      const singleLine = this.keyPair.publicKey
        .replace(/-----[A-Z ]*PUBLIC KEY-----/g, '')
        .replace(/\s+/g, '');
      await this.copyToClipboard(singleLine, '单行公钥');
    },
    async copyPrivateKeySingleLine() {
      if (!this.keyPair.privateKey) {
        this.showToast('没有内容可复制', 'error');
        return;
      }
      // 去掉 -----BEGIN/END----- 首尾行和所有换行
      const singleLine = this.keyPair.privateKey
        .replace(/-----[A-Z ]*PRIVATE KEY-----/g, '')
        .replace(/\s+/g, '');
      await this.copyToClipboard(singleLine, '单行私钥');
    },
    async copyToClipboard(text, label) {
      if (!text) {
        this.showToast('没有内容可复制', 'error');
        return;
      }
      const ok = await copyText(text);
      // 用词统一成「已复制」；这一页有 4 个复制点，靠 label 区分
      this.showToast(ok ? label + '已复制' : '复制失败', ok ? 'success' : 'error');
    },
    // 清空已生成的密钥对：两张密钥卡由 v-if 控制随之消失（与 AES 清空一样不弹提示）
    clearKeyPair() {
      this.keyPair.publicKey = '';
      this.keyPair.privateKey = '';
    },
    async generateKeyPair() {
      try {
        const result = await window.go.handler.RsaHandler.GenerateKeyPair(
          this.keySize,
          this.keyFormat
        );
        if (result.success) {
          this.keyPair.publicKey = result.publicKey;
          this.keyPair.privateKey = result.privateKey;
          this.showToast('密钥对已生成');
        } else {
          this.showToast(result.message || '生成失败', 'error');
        }
      } catch (err) {
        this.showToast(`生成失败: ${err.message}`, 'error');
      }
    }
  }
};
</script>

<style scoped>
/* 「label + 下拉」两个一组横排，再跟按钮并成一行。
   这一行是卡片的唯一内容（没有 head），左对齐就能和下面两张密钥卡的图标落在同一条左竖线上。
   组间距 20、label 与控件间距 8 —— 两个距离必须拉开档次，
   否则「密钥长度 [2048 bit] 密钥格式 [PKCS#1]」会连成一串，分不出谁是谁的标题。
   允许换行：窗口太窄时按钮掉到第二行，不会横向溢出 */
.cfg-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;
}

.fld-inline {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: none;
}

.fld-inline .field-label {
  white-space: nowrap;
}

.key-view {
  margin: 0;
  max-height: 260px;
  overflow: auto;
  padding: 12px;
  font-size: 12.5px;
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-all;
  color: var(--text);
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
}
</style>
