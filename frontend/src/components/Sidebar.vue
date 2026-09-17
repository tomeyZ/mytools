<template>
  <aside class="sidebar" :class="{ collapsed }">
    <nav class="sidebar-nav">
      <template v-for="group in groups" :key="group.label">
        <div class="nav-group-title">{{ group.label }}</div>
        <!-- 折叠态下只剩图标，工具名靠 title 兜底 -->
        <div
            v-for="item in group.items"
            :key="item.id"
            class="nav-item"
            :class="{ active: item.id === activeId }"
            :title="item.name"
            @click="switchTool(item.id)"
        >
          <span class="nav-icon" v-html="item.icon"></span>
          <span class="nav-label">{{ item.name }}</span>
        </div>
      </template>
    </nav>

    <!--
      底部「检查更新」：必须放在 .sidebar-nav 外面。
      nav 是 overflow-y:auto 的滚动容器，塞进去会跟着菜单一起滚走，
      而它应该永远钉在左下角。
    -->
    <div class="sidebar-foot">
      <button
          class="upd-row"
          :disabled="checking"
          :title="footerTip"
          @click="$emit('update-click')"
      >
        <span class="upd-icon">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 4v11"/><path d="m7.5 10.5 4.5 4.5 4.5-4.5"/><path d="M5 19h14"/>
          </svg>
          <!-- 红点挂图标右上：挂行尾的话，折叠后版本号隐藏、红点会跟着一起消失 -->
          <span v-if="hasUpdate" class="upd-dot" aria-label="有新版本"></span>
        </span>
        <span class="upd-lbl">检查更新</span>
        <span class="upd-ver" :class="{ has: hasUpdate }">{{ rightText }}</span>
      </button>
    </div>
  </aside>
</template>

<script>
import { TOOL_GROUPS } from '../tools.js';

export default {
  name: 'Sidebar',
  props: {
    collapsed: { type: Boolean, default: false },
    // 当前选中项由 MainLayout 持有（顶栏搜索也会改它），这里只做展示
    activeId: { type: String, default: '' },
    version: { type: String, default: '' },
    // 新版本号，非空即表示检测到了新版本（由 MainLayout 的 updateInfo 派生）
    updateVersion: { type: String, default: '' },
    checking: { type: Boolean, default: false }
  },
  emits: ['tool-change', 'update-click'],
  data() {
    return {
      // 菜单结构统一由 tools.js 提供，顶栏搜索用的是同一份
      groups: TOOL_GROUPS
    };
  },
  computed: {
    hasUpdate() {
      return !!this.updateVersion;
    },
    versionText() {
      return this.version ? 'v' + this.version : '—';
    },
    rightText() {
      return this.checking ? '检查中…' : this.versionText;
    },
    // 折叠态只剩图标，靠悬停提示补语义（与菜单项折叠后同一策略）
    footerTip() {
      if (this.hasUpdate) return '有新版本 ' + this.updateVersion + '，点击查看';
      return '检查更新 · 当前 ' + this.versionText;
    }
  },
  methods: {
    switchTool(id) {
      this.$emit('tool-change', id);
    }
  }
};
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  background: var(--surface);
  border-right: 1px solid var(--border);
  overflow: hidden;
  transition: background var(--tr), border-color var(--tr);
}

.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 12px 10px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border-radius: var(--radius-sm);
  color: var(--text-2);
  font-size: 13.5px;
  line-height: 1.2;
  cursor: pointer;
  white-space: nowrap;
  transition: background var(--tr), color var(--tr);
}

.nav-item:hover {
  background: var(--surface-3);
  color: var(--text);
}

.nav-item.active {
  background: var(--accent-soft);
  color: var(--accent-strong);
  font-weight: 600;
}

.nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 18px;
  height: 18px;
}

.nav-label {
  overflow: hidden;
  text-overflow: ellipsis;
}

.nav-group-title {
  padding: 14px 10px 6px;
  font-size: 10.5px;
  font-weight: 600;
  letter-spacing: .1em;
  color: var(--text-3);
  transition: color var(--tr);
}

/* ---------------- 底部：检查更新 ---------------- */
/* 不留分隔线：靠 nav 的 flex:1 撑出的留白分隔。折叠态下分组标题已退化成
   一条 22px 短线，底部再加一条容易和它混 */
.sidebar-foot {
  flex-shrink: 0;
  padding: 8px 10px 10px;
}

/* 行样式 1:1 复制 .nav-item —— 它就是一个菜单项，不该长得像别的东西 */
.upd-row {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 9px 10px;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--text-2);
  font-family: inherit;
  font-size: 13.5px;
  line-height: 1.2;
  text-align: left;
  cursor: pointer;
  white-space: nowrap;
  transition: background var(--tr), color var(--tr), opacity var(--tr);
}

.upd-row:hover {
  background: var(--surface-3);
  color: var(--text);
}

.upd-row:disabled {
  opacity: .55;
  cursor: default;
}

.upd-icon {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 18px;
  height: 18px;
}

.upd-dot {
  position: absolute;
  top: -3px;
  right: -3px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--danger);
  box-shadow: 0 0 0 2px var(--surface);
}

/* hover 时挖边色跟着底色走，否则会在色块上留一圈白边 */
.upd-row:hover .upd-dot {
  box-shadow: 0 0 0 2px var(--surface-3);
}

.upd-lbl {
  overflow: hidden;
  text-overflow: ellipsis;
}

.upd-ver {
  margin-left: auto;
  font-family: var(--font-mono);
  font-size: 11.5px;
  color: var(--text-3);
}

/* 有新版本：只用红点 + 版本号转强调色，不铺整行底色。
   侧栏是常驻区域，铺满会一直吵 */
.upd-ver.has {
  color: var(--accent-strong);
  font-weight: 600;
}

/* 检查中：状态文字顶替版本号的位置，mono 字体对中文没意义 */
.upd-row:disabled .upd-ver {
  font-family: inherit;
  font-size: 12px;
}

/* ---- 折叠态 ---- */
.collapsed .nav-item {
  justify-content: center;
  padding: 10px 0;
}

.collapsed .nav-label {
  display: none;
}

/* 分组标题在折叠态退化成一条分隔线，保留分组语义又不占地方 */
.collapsed .nav-group-title {
  font-size: 0;
  line-height: 0;
  padding: 8px 0;
}

.collapsed .nav-group-title::after {
  content: '';
  display: block;
  width: 22px;
  height: 1px;
  margin: 0 auto;
  background: var(--border);
}

.collapsed .upd-row {
  justify-content: center;
  padding: 10px 0;
}

.collapsed .upd-lbl,
.collapsed .upd-ver {
  display: none;
}
</style>
