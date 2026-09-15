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
  </aside>
</template>

<script>
import { TOOL_GROUPS } from '../tools.js';

export default {
  name: 'Sidebar',
  props: {
    collapsed: { type: Boolean, default: false },
    // 当前选中项由 MainLayout 持有（顶栏搜索也会改它），这里只做展示
    activeId: { type: String, default: '' }
  },
  emits: ['tool-change'],
  data() {
    return {
      // 菜单结构统一由 tools.js 提供，顶栏搜索用的是同一份
      groups: TOOL_GROUPS
    };
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
</style>
