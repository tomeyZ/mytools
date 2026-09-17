<template>
  <div v-if="visible && info" class="upd-mask" @click.self="$emit('close')">
    <div class="upd" role="dialog" aria-modal="true" aria-labelledby="upd-title">
      <!-- 头部淡色强调块：与正文分色，"发生了件事"的重量落在这里；
           版本号升格为大号主角，标题退成一块小胶囊（仍是 aria-labelledby 指向的标题） -->
      <div class="upd-head">
        <h3 id="upd-title" class="upd-badge">发现新版本</h3>
        <button class="upd-x" aria-label="关闭" @click="$emit('close')">
          <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round">
            <path d="M6 6l12 12M18 6L6 18"/>
          </svg>
        </button>

        <div class="upd-ver">
          <span class="v-old">v{{ currentVersion }}</span>
          <svg class="v-arrow" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M4 12h15"/><path d="m13 6 6 6-6 6"/>
          </svg>
          <span class="v-new">v{{ info.version }}</span>
        </div>

        <!-- 发布日期与版本号同源，跟着版本行走，不再单独占一行 -->
        <p v-if="info.create_date" class="upd-date">发布于 {{ info.create_date }}</p>
      </div>

      <div class="upd-body">
        <template v-if="log.length">
          <p class="upd-sub">更新内容</p>
          <!-- 头部已经分色，正文再铺一层灰底就成了"三明治"，所以日志不带面板底 -->
          <ul class="upd-log">
            <li v-for="(item, i) in log" :key="i">{{ item }}</li>
          </ul>
        </template>
        <!-- 空状态是纯文字，不能套内容框样式（那样看起来像个没填的输入框） -->
        <p v-else class="upd-empty">
          <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" aria-hidden="true">
            <circle cx="12" cy="12" r="9"/><path d="M12 8v5"/><path d="M12 16h.01"/>
          </svg>
          本次发布未填写更新说明
        </p>
      </div>

      <div class="upd-foot">
        <button class="btn ghost" @click="$emit('close')">稍后</button>
        <button class="btn primary" @click="download">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 4v11"/><path d="m7.5 10.5 4.5 4.5 4.5-4.5"/><path d="M5 19h14"/>
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

/* ---------------- 头部：淡色强调块 ---------------- */
.upd-head {
  position: relative;
  padding: 17px 18px 15px;
  background: var(--accent-soft);
  border-bottom: 1px solid var(--accent-border);
  color: var(--accent-strong);
}

/* 标题退成胶囊：主角是版本号，但身份标识不能丢（它仍是 aria-labelledby 的目标） */
.upd-badge {
  display: inline-flex;
  align-items: center;
  height: 22px;
  margin: 0;
  padding: 0 10px;
  border: 1px solid var(--accent-border);
  border-radius: 999px;
  background: var(--surface);
  font-size: 11.5px;
  font-weight: 600;
  letter-spacing: .2px;
}

.upd-x {
  position: absolute;
  top: 13px;
  right: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border: none;
  border-radius: var(--radius-xs);
  background: transparent;
  color: inherit;
  opacity: .65;
  cursor: pointer;
  transition: background var(--tr), opacity var(--tr);
}

.upd-x:hover {
  opacity: 1;
  background: var(--accent-border);
}

.upd-ver {
  display: flex;
  align-items: baseline;
  gap: 11px;
  margin-top: 13px;
  font-family: var(--font-mono);
  font-weight: 700;
  letter-spacing: -.4px;
}

/* 旧版本：删除线 = "已被取代"，比单纯变灰更准确 */
.v-old {
  font-size: 18px;
  text-decoration: line-through;
  text-decoration-thickness: 1.5px;
  opacity: .5;
}

.v-arrow {
  flex: none;
  align-self: center;
  opacity: .6;
}

/* 新版本：这弹窗唯一的主角 */
.v-new {
  font-size: 25px;
}

.upd-date {
  margin: 7px 0 0;
  font-family: var(--font-mono);
  font-size: 11.5px;
  opacity: .72;
}

/* ---------------- 正文 ---------------- */
.upd-body {
  padding: 15px 18px 2px;
}

.upd-sub {
  margin: 0 0 8px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text);
}

.upd-log {
  margin: 0;
  padding: 0 4px 0 0;
  list-style: none;
  max-height: 200px;
  overflow-y: auto;
}

.upd-log li {
  position: relative;
  padding-left: 14px;
  font-size: 12.5px;
  line-height: 1.75;
  color: var(--text-2);
}

.upd-log li + li {
  margin-top: 3px;
}

.upd-log li::before {
  content: "";
  position: absolute;
  top: 8px;
  left: 2px;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--accent);
}

.upd-empty {
  display: flex;
  align-items: center;
  gap: 7px;
  margin: 0;
  font-size: 12.5px;
  color: var(--text-3);
}

.upd-empty svg {
  flex: none;
}

/* ---------------- 底部 ---------------- */
/* 头部已有色块，底栏再铺灰底 + 分隔线就是第三次分层，去掉 */
.upd-foot {
  display: flex;
  justify-content: flex-end;
  gap: 9px;
  padding: 14px 18px 17px;
}

/* 弹窗里的按钮比表单里的矮一号，视觉上才不像在填表单 */
.upd-foot .btn {
  height: 32px;
  padding: 0 15px;
  font-size: 12.5px;
}
</style>
