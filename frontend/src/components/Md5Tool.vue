<template>
  <div class="page">
    <div class="grid">
      <section class="card wide">
        <div class="card-head">
          <span class="ci">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
              <path d="M4 9h16M4 15h16"/><path d="M10 3.5 8 20.5M16 3.5l-2 17"/>
            </svg>
          </span>
          <div>
            <h3>输入原文</h3>
            <p>任意文本内容，支持多行</p>
          </div>
        </div>
        <div class="card-body">
          <textarea
              v-model="inputText"
              class="inp mono"
              rows="6"
              placeholder="请输入要加密的内容"
          ></textarea>
          <div class="row">
            <button class="btn primary" @click="handleMd5Encrypt">MD5 加密</button>
            <button class="btn ghost" @click="handleClear">清空</button>
          </div>
        </div>
      </section>

      <section class="card wide" v-if="showResults">
        <div class="card-head">
          <span class="ci ci-2">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3.5" y="3.5" width="17" height="17" rx="2.4"/>
              <path d="M8 9h8M8 13h8M8 17h5"/>
            </svg>
          </span>
          <div>
            <h3>计算结果</h3>
            <p>点击右侧按钮复制</p>
          </div>
        </div>
        <div class="card-body">
          <div v-for="item in rows" :key="item.label" class="out out-ok">
            <span class="tag">{{ item.label }}</span>
            <span class="val mono">{{ item.value }}</span>
            <button class="copy" title="复制" @click="handleCopy(item.value, item.label)">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
                <rect x="9" y="9" width="12" height="12" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
              </svg>
            </button>
          </div>
        </div>
      </section>
    </div>

    <transition name="toast-fade">
      <div v-if="toast.show" :class="['message-toast', toast.type]">{{ toast.text }}</div>
    </transition>
  </div>
</template>

<script>
import md5 from 'js-md5';
import { copyText } from '../clipboard.js';
import { toastMixin } from '../toast.js';

export default {
  name: 'Md5Tool',
  mixins: [toastMixin],
  data() {
    return {
      inputText: '',
      result32Upper: '',
      result32Lower: '',
      result16Upper: '',
      result16Lower: '',
      showResults: false
    };
  },
  computed: {
    rows() {
      return [
        { label: '32 位大写', value: this.result32Upper },
        { label: '32 位小写', value: this.result32Lower },
        { label: '16 位大写', value: this.result16Upper },
        { label: '16 位小写', value: this.result16Lower }
      ];
    }
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
    async handleCopy(text, label) {
      const ok = await copyText(text);
      // 4 个结果并排，文案必须带对象，否则分不清复制了哪个
      this.showToast(ok ? label + '已复制' : '复制失败', ok ? 'success' : 'error');
    }
  }
};
</script>

<style scoped>
.out .tag {
  flex: none;
  width: 76px;
  font-size: 12px;
  font-weight: 500;
  color: var(--text-3);
}
</style>
