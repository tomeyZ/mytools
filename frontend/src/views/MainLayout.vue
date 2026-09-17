<template>
  <div class="app-shell" :class="{ 'nav-collapsed': collapsed }">
    <Sidebar
        :collapsed="collapsed"
        :active-id="activeTool"
        :version="version"
        :update-version="updateInfo ? updateInfo.version : ''"
        :checking="checking"
        @tool-change="handleToolChange"
        @update-click="onVersionClick"
    />

    <div class="shell-main">
      <header class="topbar">
        <button class="icon-btn" :title="collapsed ? '展开侧边栏' : '折叠侧边栏'" @click="toggleCollapse">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="4.5" width="18" height="15" rx="2.4"/>
            <line x1="9.5" y1="4.5" x2="9.5" y2="19.5"/>
          </svg>
        </button>

        <!-- 顶栏工具搜索：Ctrl+K 聚焦，上下键选择，回车直达 -->
        <div class="search">
          <svg class="search-icon" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="7"/><line x1="21" y1="21" x2="16.2" y2="16.2"/>
          </svg>
          <input
              ref="searchEl"
              v-model="keyword"
              class="search-input"
              placeholder="搜索工具…"
              @focus="dropdownOpen = true"
              @blur="dropdownOpen = false"
              @keydown.down.prevent="moveCursor(1)"
              @keydown.up.prevent="moveCursor(-1)"
              @keydown.enter.prevent="chooseCursor"
              @keydown.esc.prevent="resetSearch"
          >
          <kbd class="search-kbd">Ctrl K</kbd>

          <!-- mousedown.prevent：避免 input 先失焦把下拉关掉，导致点击不生效 -->
          <ul v-if="showDropdown" class="search-dropdown">
            <li
                v-for="(tool, i) in results"
                :key="tool.id"
                :class="{ on: i === cursor }"
                @mouseenter="cursor = i"
                @mousedown.prevent="chooseTool(tool.id)"
            >
              <span class="drop-icon" v-html="tool.icon"></span>
              <span class="drop-name">{{ tool.name }}</span>
              <span class="drop-group">{{ tool.group }}</span>
            </li>
            <li v-if="!results.length" class="drop-empty">没有匹配的工具</li>
          </ul>
        </div>

        <div class="topbar-right">
          <div class="theme-seg" role="group" aria-label="主题">
            <button
                v-for="t in themes"
                :key="t.id"
                class="seg-btn"
                :class="{ on: t.id === theme }"
                :title="t.tip"
                @click="setTheme(t.id)"
            >
              <span class="seg-icon" v-html="t.icon"></span>{{ t.label }}
            </button>
          </div>
        </div>
      </header>

      <main class="content">
        <!-- KeepAlive 缓存组件状态，切换菜单后已填内容不丢失 -->
        <KeepAlive>
          <component :is="currentComponent"/>
        </KeepAlive>
      </main>

      <!-- 检查更新：启动时静默查一次，有新版本只在侧栏底部那行挂红点，不打扰 -->
      <UpdateDialog
          :visible="dialogOpen"
          :current-version="version"
          :info="updateInfo"
          @close="dialogOpen = false"
      />

      <transition name="toast-fade">
        <div v-if="toast.show" :class="['message-toast', toast.type]">{{ toast.text }}</div>
      </transition>
    </div>
  </div>
</template>

<script>
import Sidebar from '../components/Sidebar.vue';
import UpdateDialog from '../components/UpdateDialog.vue';
import { TOOL_COMPONENTS, flatTools } from '../tools.js';
import { toastMixin } from '../toast.js';

const THEME_KEY = 'mytools:theme';
const NAV_KEY = 'mytools:nav-collapsed';

// 主题清单：曾经还有一档毛玻璃（半透明面板 + 背景光晕），已移除。
// 如果用户在旧版本里选过它，localStorage 会残留该值，读取时必须在下面白名单归一，
// 否则 this.theme 拿到未知值 → 两个按钮都不高亮，且 <html> 上的属性没有对应 CSS 规则。
const THEMES = [
  { id: 'light', label: '浅色', tip: '浅色主题' },
  { id: 'dark', label: '深色', tip: '深色主题' }
];

const isKnownTheme = (id) => THEMES.some((t) => t.id === id);

const svg = (body) =>
    '<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">' + body + '</svg>';

THEMES[0].icon = svg('<circle cx="12" cy="12" r="4"/><path d="M12 2v2.2M12 19.8V22M4.9 4.9l1.6 1.6M17.5 17.5l1.6 1.6M2 12h2.2M19.8 12H22M4.9 19.1l1.6-1.6M17.5 6.5l1.6-1.6"/>');
THEMES[1].icon = svg('<path d="M20.5 13.2A8.5 8.5 0 1 1 10.8 3.5a6.5 6.5 0 0 0 9.7 9.7z"/>');

// 读写 localStorage 时兜底：首次运行或存储不可用时不能让应用起不来
const readStore = (key, fallback) => {
  try {
    const v = localStorage.getItem(key);
    return v === null ? fallback : v;
  } catch (e) {
    return fallback;
  }
};
const writeStore = (key, value) => {
  try {
    localStorage.setItem(key, value);
  } catch (e) {
    /* 存不了就算了，退化成「本次会话有效」 */
  }
};

export default {
  name: 'MainLayout',
  components: { Sidebar, UpdateDialog },
  mixins: [toastMixin],
  data() {
    return {
      activeTool: 'time',
      version: '',
      theme: 'light',
      collapsed: false,
      themes: THEMES,
      keyword: '',
      dropdownOpen: false,
      cursor: 0,
      checking: false,        // 正在检查更新（防连点）
      updateInfo: null,       // 后端 VersionInfo，非空即表示有新版本
      dialogOpen: false
    };
  },
  computed: {
    currentComponent() {
      return TOOL_COMPONENTS[this.activeTool] || TOOL_COMPONENTS.time;
    },
    results() {
      const kw = this.keyword.trim().toLowerCase();
      const all = flatTools();
      if (!kw) return all;
      return all.filter((t) => t.name.toLowerCase().includes(kw));
    },
    showDropdown() {
      // 有结果或正在输入（含「没有匹配」的提示）时都展开
      return this.dropdownOpen && (this.results.length > 0 || this.keyword.trim() !== '');
    }
  },
  watch: {
    // 关键字变化时把高亮复位，否则过滤后可能指向越界项
    keyword() {
      this.cursor = 0;
    }
  },
  mounted() {
    const saved = readStore(THEME_KEY, 'light');
    this.theme = isKnownTheme(saved) ? saved : 'light';
    this.collapsed = readStore(NAV_KEY, '0') === '1';
    this.applyTheme();
    document.addEventListener('keydown', this.onGlobalKeydown);
    this.loadVersion();
    // 启动静默检查：成功且有新版本就挂红点，失败一律无声（只是没有红点）
    this.checkUpdate(true);
  },
  beforeUnmount() {
    document.removeEventListener('keydown', this.onGlobalKeydown);
  },
  methods: {
    handleToolChange(toolId) {
      this.activeTool = toolId;
    },
    toggleCollapse() {
      this.collapsed = !this.collapsed;
      writeStore(NAV_KEY, this.collapsed ? '1' : '0');
    },
    setTheme(id) {
      this.theme = id;
      this.applyTheme();
      writeStore(THEME_KEY, id);
    },
    // 主题只切 <html data-theme>，颜色由 style.css 里的变量统一响应
    applyTheme() {
      document.documentElement.setAttribute('data-theme', this.theme);
    },
    async loadVersion() {
      try {
        this.version = await window.go.handler.VersionHandler.GetCurrentVersion();
      } catch (e) {
        this.version = '';
      }
    },

    // 检查更新。silent=true 用于启动自检：不弹任何提示、失败也无声
    async checkUpdate(silent) {
      if (this.checking) return;
      this.checking = true;
      let res = null;
      try {
        res = await window.go.handler.VersionHandler.CheckUpdate();
      } catch (e) {
        // 后端已经把可预期失败转成 status，走到这里说明是绑定层面的问题
        if (!silent) this.showToast('检查更新失败，请稍后重试', 'error');
        return;
      } finally {
        this.checking = false;
      }

      this.updateInfo = (res && res.status === 'update' && res.latest) ? res.latest : null;

      if (silent) return;
      if (this.updateInfo) {
        this.dialogOpen = true;
      } else if (res && res.status === 'error') {
        this.showToast(res.message, 'error');
      } else {
        // latest / none：文案由后端给（「已是最新版本」「暂未发布可用版本」）
        this.showToast((res && res.message) || '已是最新版本');
      }
    },

    // 版本号本身就是按钮：已有新版本直接弹窗，没有就先查一次
    onVersionClick() {
      if (this.checking) return;
      if (this.updateInfo) {
        this.dialogOpen = true;
        return;
      }
      this.checkUpdate(false);
    },
    chooseTool(id) {
      this.activeTool = id;
      this.resetSearch();
      this.$refs.searchEl.blur();
    },
    chooseCursor() {
      const tool = this.results[this.cursor];
      if (tool) this.chooseTool(tool.id);
    },
    moveCursor(step) {
      if (!this.results.length) return;
      this.cursor = (this.cursor + step + this.results.length) % this.results.length;
    },
    resetSearch() {
      this.keyword = '';
      this.dropdownOpen = false;
    },
    onGlobalKeydown(e) {
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
        e.preventDefault();
        this.dropdownOpen = true;
        this.$refs.searchEl.focus();
        this.$refs.searchEl.select();
      }
    }
  }
};
</script>

<style scoped>
.app-shell {
  --nav-w: 216px;
  display: grid;
  grid-template-columns: var(--nav-w) 1fr;
  height: 100vh;
  background-color: var(--bg);
  transition: grid-template-columns .22s ease, background-color var(--tr);
}

.app-shell.nav-collapsed {
  --nav-w: 68px;
}

.shell-main {
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
}

/* ---------------- 顶栏 ---------------- */
.topbar {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
  height: 48px;
  padding: 0 16px;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  transition: background var(--tr), border-color var(--tr);
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  border: 1px solid var(--border);
  border-radius: var(--radius-xs);
  background: transparent;
  color: var(--text-3);
  cursor: pointer;
  transition: color var(--tr), border-color var(--tr);
}

.icon-btn:hover {
  color: var(--accent-strong);
  border-color: var(--accent);
}

/* ---------------- 搜索 ---------------- */
.search {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
  width: 232px;
  height: 30px;
  padding: 0 10px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  transition: border-color var(--tr), background var(--tr);
}

.search:focus-within {
  border-color: var(--accent);
  background: var(--surface);
}

.search-icon {
  flex-shrink: 0;
  color: var(--text-3);
}

.search-input {
  flex: 1;
  min-width: 0;
  border: none;
  outline: none;
  background: transparent;
  font-family: inherit;
  font-size: 12.5px;
  color: var(--text);
}

.search-input::placeholder {
  color: var(--text-3);
}

.search-kbd {
  flex-shrink: 0;
  padding: 1px 5px;
  border: 1px solid var(--border);
  border-radius: 4px;
  font-family: var(--font-sans);
  font-size: 10.5px;
  color: var(--text-3);
}

.search:focus-within .search-kbd {
  display: none;
}

.search-dropdown {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  z-index: 50;
  max-height: 320px;
  overflow-y: auto;
  margin: 0;
  padding: 5px;
  list-style: none;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-lg);
}

.search-dropdown li {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 8px 9px;
  border-radius: var(--radius-xs);
  color: var(--text-2);
  font-size: 13px;
  cursor: pointer;
}

.search-dropdown li.on {
  background: var(--accent-soft);
  color: var(--accent-strong);
}

.drop-icon {
  display: flex;
  flex-shrink: 0;
  width: 16px;
  height: 16px;
}

.drop-name {
  font-weight: 500;
}

.drop-group {
  margin-left: auto;
  font-size: 11px;
  color: var(--text-3);
}

.drop-empty {
  justify-content: center;
  color: var(--text-3);
  cursor: default;
}

/* ---------------- 右侧 ---------------- */
.topbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

.theme-seg {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 2px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-radius: calc(var(--radius-sm) + 2px);
}

.seg-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  height: 24px;
  padding: 0 10px;
  border: none;
  border-radius: var(--radius-xs);
  background: transparent;
  color: var(--text-3);
  font-family: inherit;
  font-size: 12px;
  line-height: 1;
  cursor: pointer;
  transition: background var(--tr), color var(--tr);
}

.seg-btn:hover {
  color: var(--text-2);
}

.seg-btn.on {
  background: var(--surface);
  color: var(--accent-strong);
  font-weight: 600;
}

.seg-icon {
  display: flex;
  flex-shrink: 0;
}

/* ---------------- 内容区 ---------------- */
.content {
  flex: 1;
  min-height: 0;
  overflow: auto;
}
</style>
