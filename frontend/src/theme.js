/* ============================================================================
   当前主题读取 + 变更订阅
   ----------------------------------------------------------------------------
   主题由 MainLayout 写在 <html data-theme> 上，编辑器类组件（ace / jsoneditor）
   没法靠 CSS 变量驱动内部配色，必须主动跟随切换，所以这里提供统一入口。

   【不要用 <style> 里写 v-if 判断深浅色】7 个工具页要各判一遍，容易漏。
   ============================================================================ */

const listeners = new Set();
let observer = null;

export function currentTheme() {
  return document.documentElement.getAttribute('data-theme') || 'light';
}

export function isDark() {
  return currentTheme() === 'dark';
}

function ensureObserver() {
  if (observer) return;
  observer = new MutationObserver(() => {
    listeners.forEach((fn) => fn(currentTheme()));
  });
  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme']
  });
}

/**
 * 订阅主题变化，返回取消订阅函数（在 onBeforeUnmount 里调用）。
 */
export function onThemeChange(fn) {
  listeners.add(fn);
  ensureObserver();
  return () => listeners.delete(fn);
}

/* ace / jsoneditor 的明暗主题搭配（只两档，取值处已用 ACE_THEME.light 兜底） */
export const ACE_THEME = {
  light: 'ace/theme/github',
  dark: 'ace/theme/github_dark'
};
