<template>
  <div class="page">
    <div class="grid">
      <div class="card wide">

        <!-- tab 条。它不是页头 —— 页头那套（大标题 + 副标题）已经全站废弃，
             这里是卡片内部的第一个块，自带一条分隔线把「切换」和「内容」分开 -->
        <div class="q-head">
          <div class="seg">
            <button :class="{ on: tab === 'decode' }" @click="tab = 'decode'">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M3 7V5a2 2 0 0 1 2-2h2"/><path d="M17 3h2a2 2 0 0 1 2 2v2"/>
                <path d="M21 17v2a2 2 0 0 1-2 2h-2"/><path d="M7 21H5a2 2 0 0 1-2-2v-2"/>
                <path d="M4 12h16"/>
              </svg>
              识别
            </button>
            <button :class="{ on: tab === 'gen' }" @click="tab = 'gen'">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect width="5" height="5" x="3" y="3" rx="1"/><rect width="5" height="5" x="16" y="3" rx="1"/>
                <rect width="5" height="5" x="3" y="16" rx="1"/><path d="M21 16h-3a2 2 0 0 0-2 2v3"/>
                <path d="M12 7v3a2 2 0 0 1-2 2H7"/>
              </svg>
              生成
            </button>
          </div>

          <!-- 生成是「输入即出码」，没有生成按钮，这件事不写出来看不出来；
               识别那边不需要提示，也不用凑字数，右侧留空 -->
          <span class="end">
            <template v-if="tab === 'gen'">
              <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M13 2 3 14h9l-1 8 10-12h-9l1-8z"/>
              </svg>
              内容变化后<b>自动重新生成</b>
            </template>
          </span>
        </div>

        <!-- 两个面板都用 v-show 而非 v-if：切 tab 时输入的内容、已选的图、
             已经画好的 canvas 都得留着，重来一遍很烦。
             代价是两个实例同时在 DOM 里，所以识别面板要靠 active 判断自己是不是在前台，
             否则在「生成」页粘一张图会被它截走去识别 -->
        <QrcodeDecode v-show="tab === 'decode'" :active="tab === 'decode'"/>
        <QrcodeGenerate v-show="tab === 'gen'"/>

      </div>
    </div>
  </div>
</template>

<script>
import QrcodeDecode from './QrcodeDecode.vue';
import QrcodeGenerate from './QrcodeGenerate.vue';

export default {
  name: 'QrcodeTool',
  components: { QrcodeDecode, QrcodeGenerate },
  data() {
    return {
      // 默认停在识别：这个工具进侧栏时叫「二维码识别」，用户日常打开就是为了扫图，
      // 加了个新功能不该把原有入口顶掉
      tab: 'decode'
    };
  }
};
</script>

<style scoped>
.q-head {
  flex: none;
  display: flex;
  align-items: center;
  gap: 10px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}

/* 分段切换。项目没有全局 border-box，这里的固定高度全得自己补，
   否则 28px 会被 padding 和 border 撑成 32px */
.seg {
  box-sizing: border-box;
  display: inline-flex;
  align-items: center;
  gap: 2px;
  padding: 3px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
}

.seg button {
  box-sizing: border-box;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 28px;
  padding: 0 14px;
  border: 1px solid transparent;
  border-radius: var(--radius-xs);
  background: transparent;
  font-family: inherit;
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text-2);
  cursor: pointer;
  transition: background var(--tr), color var(--tr), border-color var(--tr);
}

.seg button:hover:not(.on) {
  color: var(--text);
}

.seg button.on {
  background: var(--surface);
  border-color: var(--border);
  color: var(--text);
  box-shadow: var(--shadow-sm);
}

.q-head .end {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11.5px;
  color: var(--text-3);
}

.q-head .end svg {
  opacity: .75;
}

.q-head .end b {
  font-weight: 650;
  color: var(--text-2);
}
</style>
