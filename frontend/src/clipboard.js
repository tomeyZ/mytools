/* ============================================================================
   复制到剪贴板（统一实现）
   ----------------------------------------------------------------------------
   WebView2 里 navigator.clipboard.writeText 有时不会 resolve（焦点/权限问题），
   直接用会卡住 UI；所以加超时竞速，超时后降级到 execCommand。
   ============================================================================ */

export async function copyText(text) {
  if (text === undefined || text === null || text === '') return false;

  try {
    await Promise.race([
      navigator.clipboard.writeText(text),
      new Promise((_, reject) => setTimeout(() => reject(new Error('clipboard timeout')), 500))
    ]);
    return true;
  } catch (e) {
    /* 落到下面的降级方案 */
  }

  try {
    const ta = document.createElement('textarea');
    ta.value = text;
    ta.style.position = 'fixed';
    ta.style.top = '0';
    ta.style.opacity = '0';
    document.body.appendChild(ta);
    ta.select();
    const ok = document.execCommand('copy');
    document.body.removeChild(ta);
    return ok;
  } catch (e) {
    return false;
  }
}
