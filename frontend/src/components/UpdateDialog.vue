<template>
  <div v-if="visible && info" class="upd-mask" @click.self="$emit('close')">
    <div class="upd" role="dialog" aria-modal="true" aria-labelledby="upd-title">
      <div class="upd-head">
        <span class="upd-badge" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="17" height="17" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 3v11"/><path d="m7.5 9.5 4.5 4.5 4.5-4.5"/><path d="M5 20h14"/>
          </svg>
        </span>
        <h3 id="upd-title">发现新版本</h3>
        <button class="upd-x" aria-label="关闭" @click="$emit('close')">
          <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round">
            <path d="M6 6l12 12M18 6L6 18"/>
          </svg>
        </button>
      </div>

      <div class="upd-body">
        <div class="ver-line">
          <span class="vchip">v{{ currentVersion }}</span>
          <svg class="v-arrow" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M4 12h15"/><path d="m13 6 6 6-6 6"/>
          </svg>
          <span class="vchip new">v{{ info.version }}</span>
        </div>

        <p v-if="info.create_date" class="upd-date">发布于 {{ info.create_date }}</p>

        <ul v-if="log.length" class="upd-log">
          <li v-for="(item, i) in log" :key="i">{{ item }}</li>
        </ul>
        <p v-else class="upd-log empty">本次发布没有填写更新说明</p>
      </div>

      <div class="upd-foot">
        <button class="btn ghost" @click="$emit('close')">稍后</button>
        <button class="btn primary" @click="download">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 3v11"/><path d="m7.5 9.5 4.5 4.5 4.5-4.5"/><path d="M5 20h14"/>
          </svg>
          前往下载
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'UpdateDialog',
  props: {
    visible: Boolean,
    currentVersion: String,
    // 后端 VersionInfo：{version, change_log, create_date, download_url, release_url}
    info: { type: Object, default: null }
  },
  computed: {
    log() {
      return (this.info && Array.isArray(this.info.change_log)) ? this.info.change_log : [];
    }
  },
  methods: {
    download() {
      // 没有匹配当前平台的安装包时退回 Release 页面，让用户自己选
      const url = this.info.download_url || this.info.release_url;
      if (!url) return;
      // 走 Go 侧的 OpenExternal：WebView2 里 window.open 会在应用窗口内打开，没有下载能力
      if (window.go && window.go.main && window.go.main.App) {
        window.go.main.App.OpenExternal(url);
      } else {
        window.open(url, '_blank');
      }
      this.$emit('close');
    }
  }
};
</script>

<style scoped>
.upd-mask {
  position: fixed;
  inset: 0;
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: rgba(10, 14, 20, .42);
  backdrop-filter: blur(2px);
}

.upd {
  width: 420px;
  max-width: 100%;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  /* 深色主题下 --shadow-lg 是 none，必须额外给一层实阴影 */
  box-shadow: var(--shadow-lg), 0 18px 44px rgba(0, 0, 0, .28);
  overflow: hidden;
}

/* ---------------- 头部 ---------------- */
.upd-head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 16px 12px 18px;
}

.upd-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: none;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--accent-soft);
  color: var(--accent-strong);
}

.upd-head h3 {
  flex: 1;
  min-width: 0;
  margin: 0;
  font-size: 14.5px;
  font-weight: 600;
  color: var(--text);
}

.upd-x {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: none;
  width: 26px;
  height: 26px;
  border: none;
  border-radius: var(--radius-xs);
  background: transparent;
  color: var(--text-3);
  cursor: pointer;
  transition: background var(--tr), color var(--tr);
}

.upd-x:hover {
  background: var(--surface-2);
  color: var(--text);
}

/* ---------------- 正文 ---------------- */
.upd-body {
  padding: 0 18px 18px;
}

.ver-line {
  display: flex;
  align-items: center;
  gap: 10px;
}

.vchip {
  padding: 3px 9px;
  border: 1px solid var(--border);
  border-radius: var(--radius-xs);
  background: var(--surface-2);
  color: var(--text-2);
  font-family: var(--font-mono);
  font-size: 12px;
  line-height: 1.4;
}

.vchip.new {
  border-color: var(--accent-border);
  background: var(--accent-soft);
  color: var(--accent-strong);
  font-weight: 600;
}

.v-arrow {
  flex: none;
  color: var(--text-3);
}

.upd-date {
  margin: 10px 0 0;
  font-size: 11.5px;
  color: var(--text-3);
}

.upd-log {
  margin: 12px 0 0;
  padding: 12px 14px;
  list-style: none;
  max-height: 210px;
  overflow-y: auto;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--surface-2);
}

.upd-log li {
  position: relative;
  padding-left: 13px;
  font-size: 12.5px;
  line-height: 1.75;
  color: var(--text-2);
}

.upd-log li::before {
  content: "";
  position: absolute;
  top: 9px;
  left: 2px;
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--accent);
}

.upd-log.empty {
  margin-bottom: 0;
  color: var(--text-3);
  font-size: 12.5px;
}

/* ---------------- 底部 ---------------- */
.upd-foot {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 13px 18px;
  border-top: 1px solid var(--border);
  background: var(--surface-2);
}

/* 弹窗里的按钮比表单里的矮一号，视觉上才不像在填表单 */
.upd-foot .btn {
  height: 32px;
  font-size: 12.5px;
}
</style>
