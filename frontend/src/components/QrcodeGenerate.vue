<template>
  <div class="gen">
    <div class="split">

      <!-- ============ 左：内容 ============ -->
      <div class="left">
        <div class="sec">
          <span class="ic">
            <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
              <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/>
              <path d="M14 2v4a2 2 0 0 0 2 2h4"/><path d="M10 9H8"/><path d="M16 13H8"/><path d="M16 17H8"/>
            </svg>
          </span>
          <h4>内容</h4>
          <!-- 计数是「内容自己的信息」，挂在内容分区标题右侧；字节不是字符，
               中文 1 字 = 3 字节，二维码容量按字节算，这样才和右下角信息条对得上 -->
          <span class="sp"><span class="cnt" :class="{ near: nearLimit }">{{ bytes }} B</span></span>
        </div>

        <textarea
            ref="ta"
            v-model="text"
            class="inp ta"
            spellcheck="false"
            placeholder="输入要生成二维码的内容 —— 网址、文本、Wi-Fi 配置都可以…"></textarea>

        <!-- 超限提示放输入侧，不铺到右侧图台上：空态和报错态挤同一格，
             会让人分不清「还没输入」还是「输错了」 -->
        <div class="hint" :class="{ bad: !!error }">
          <template v-if="error">{{ error }}</template>
          <template v-else>
            支持网址 / 文本 / Wi-Fi 配置 / 名片 · {{ level }} 级上限 <b>{{ capacity }} 字节</b>
          </template>
        </div>
      </div>

      <!-- ============ 右：二维码 ============ -->
      <div class="right">
        <div class="sec">
          <span class="ic">
            <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
              <rect width="5" height="5" x="3" y="3" rx="1"/><rect width="5" height="5" x="16" y="3" rx="1"/>
              <rect width="5" height="5" x="3" y="16" rx="1"/><path d="M21 16h-3a2 2 0 0 0-2 2v3"/>
              <path d="M12 7v3a2 2 0 0 1-2 2H7"/>
            </svg>
          </span>
          <h4>二维码</h4>
          <span class="sp">
            <span v-if="error" class="st bad">生成失败</span>
            <span v-else-if="qr" class="st ok">
              <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"/></svg>
              已生成
            </span>
            <span v-else class="st idle">等待输入</span>
          </span>
        </div>

        <!-- 空态与出码态同高（split 定高 + 这里 flex:1），所以开始输入时图台不会长个儿 -->
        <div class="qstage" :class="{ ph: !qr }">
          <div class="qstage-body">
            <template v-if="qr">
              <div class="qcard qr-art"><canvas ref="cv"></canvas></div>
            </template>
            <template v-else>
              <!-- 与识别页同一枚水印。draggable=false + pointer-events:none：
                   页面内的图默认能被浏览器原生拖走，拖动时甩出一张 ghost 图很难看 -->
              <div class="ph-frame"><img :src="watermark" alt="" draggable="false"></div>
              <div class="ph-t">左侧输入内容后<br>二维码会实时出现在这里</div>
            </template>
          </div>

          <div class="qstage-cap">
            <template v-if="qr">
              <span class="mono-n">{{ moduleLabel }}</span>
              <span class="dot"></span>
              <span>{{ bytes }} B</span>
              <span class="end">
                <!-- 版本高说明模块密，导出小尺寸有扫不出的风险 —— 这是唯一能把这个风险
                     说出口的地方。只提示，不自动改用户的选项，也不弹 toast -->
                <span v-if="adviseBig" class="tag warn">V{{ qr.version }} · 建议 1024</span>
                <span v-else class="tag">V{{ qr.version }} · {{ level }}</span>
              </span>
            </template>
            <span v-else class="mono-n">—</span>
          </div>
        </div>
      </div>
    </div>

    <!-- ============ 底部工具条：左参数、右动作 ============ -->
    <div class="bar">
      <span class="lb">尺寸</span>
      <AppSelect v-model="size" :options="sizeOptions" min-width="124px"/>
      <span class="lb">容错</span>
      <AppSelect v-model="level" :options="levelOptions" min-width="104px"/>
      <span class="end">
        <button class="btn ghost" :disabled="!qr" @click="copyImage">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
            <rect x="9" y="9" width="11" height="11" rx="2"/><path d="M5 15V5a2 2 0 0 1 2-2h8"/>
          </svg>
          复制图片
        </button>
        <button class="btn primary" :disabled="!qr" @click="saveImage">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 3v11"/><path d="M8 11l4 4 4-4"/><path d="M4 19h16"/>
          </svg>
          保存 PNG
        </button>
      </span>
    </div>

    <transition name="toast-fade">
      <div v-if="toast.show" :class="['message-toast', toast.type]">{{ toast.text }}</div>
    </transition>
  </div>
</template>

<script>
import QRCode from 'qrcode';
import AppSelect from './AppSelect.vue';
import { toastMixin } from '../toast.js';
import watermark from '../assets/images/qr-watermark.png';

// 各容错级别的字节上限（byte 模式）。纯数字内容实际能装更多 —— 库自己会选数字/字母数字模式，
// 所以这里的数字只用来「提前变黄提醒」，真正的超限判定交给 create() 抛不抛异常
const CAPACITY = { L: 2953, M: 2331, Q: 1663, H: 1273 };

// 静区（模块数）。规范是 4，屏幕和导出都用它 —— 深色截图里码的外圈会紧贴白卡边，
// 留不够扫码器会认不出定位图案
const MARGIN = 4;

// 白卡内容区边长。canvas 按整数倍缩放铺进去，非整数倍会出抗锯齿灰边，扫码率反而掉
const BOX = 276;

// 输入即生成，但版本高的内容算 RS 纠错要几毫秒，连着敲字会卡手；180ms 是「停手就出」的手感
const DEBOUNCE = 180;

export default {
  name: 'QrcodeGenerate',
  components: { AppSelect },
  mixins: [toastMixin],
  data() {
    return {
      text: '',
      size: 512,        // 导出尺寸，与屏幕显示无关
      level: 'M',
      qr: null,         // { version, size }，null = 没内容或有错
      error: '',
      timer: null,
      sizeOptions: [
        { value: 256, label: '小 · 256 px' },
        { value: 512, label: '中 · 512 px' },
        { value: 1024, label: '大 · 1024 px' }
      ],
      levelOptions: [
        { value: 'L', label: 'L · 7%' },
        { value: 'M', label: 'M · 15%' },
        { value: 'Q', label: 'Q · 25%' },
        { value: 'H', label: 'H · 30%' }
      ],
      watermark
    };
  },
  computed: {
    bytes() {
      return new TextEncoder().encode(this.text).length;
    },
    capacity() {
      return CAPACITY[this.level];
    },
    // 保守预警：byte 模式的上限，纯数字内容会提前变黄但照样能生成
    nearLimit() {
      return !this.error && this.bytes > this.capacity * 0.85;
    },
    moduleLabel() {
      return this.qr ? this.qr.size + ' × ' + this.qr.size : '';
    },
    // 版本 ≥ 10 说明模块已经很密，这时导出 256/512 有可能扫不出
    adviseBig() {
      return !!this.qr && this.qr.version >= 10 && this.size < 1024;
    }
  },
  watch: {
    text() {
      this.schedule();
    },
    // 换容错级别要立刻重算：版本会变，模块数和上限都跟着变，等 180ms 会让人以为没反应
    level() {
      clearTimeout(this.timer);
      this.generate();
    }
  },
  mounted() {
    this.generate();
  },
  beforeUnmount() {
    clearTimeout(this.timer);
  },
  methods: {
    schedule() {
      clearTimeout(this.timer);
      this.timer = setTimeout(this.generate, DEBOUNCE);
    },

    generate() {
      if (!this.text) {
        this.qr = null;
        this.error = '';
        return;
      }
      try {
        const q = QRCode.create(this.text, { errorCorrectionLevel: this.level });
        this.qr = { version: q.version, size: q.modules.size };
        this.error = '';
        // 状态先落地（等待输入 → 已生成），canvas 由 v-if 挂上来之后才画
        this.$nextTick(() => this.paint(q));
      } catch (e) {
        this.qr = null;
        // 库只给英文异常。逐条翻成人话，并且给出可执行的两条路，不写「请重试」
        const tooBig = /too big to be stored/i.test(e.message || '');
        this.error = tooBig
          ? '内容太长，超出了二维码的容量 —— 缩短内容，或把容错降到 L 级（上限 ' + CAPACITY.L + ' 字节）'
          : '生成失败：' + (e.message || e);
      }
    },

    // 自己逐模块 fillRect，而不是调 QRCode.toCanvas：这样能保证「一个模块 = 整数个物理像素」，
    // 屏幕上不会出现灰边，扫起来最稳
    paint(q) {
      const cv = this.$refs.cv;
      if (!cv) return;
      const dpr = window.devicePixelRatio || 1;
      const n = q.modules.size;
      const total = n + MARGIN * 2;
      const scale = Math.max(1, Math.floor((BOX * dpr) / total));
      const px = total * scale;

      cv.width = px;
      cv.height = px;
      cv.style.width = px / dpr + 'px';
      cv.style.height = px / dpr + 'px';

      const ctx = cv.getContext('2d');
      ctx.fillStyle = '#fff';
      ctx.fillRect(0, 0, px, px);
      ctx.fillStyle = '#000';
      const data = q.modules.data;
      for (let r = 0; r < n; r++) {
        for (let c = 0; c < n; c++) {
          if (data[r * n + c]) ctx.fillRect((c + MARGIN) * scale, (r + MARGIN) * scale, scale, scale);
        }
      }
    },

    // 导出是重新渲染一遍，不是把屏幕上的 canvas 放大 —— 放大要插值，出毛边，扫码率会掉
    exportDataURL() {
      return QRCode.toDataURL(this.text, {
        errorCorrectionLevel: this.level,
        width: Number(this.size),
        margin: MARGIN,
        color: { dark: '#000', light: '#fff' }
      });
    },

    async copyImage() {
      if (!this.qr) return;
      try {
        const url = await this.exportDataURL();
        const blob = await (await fetch(url)).blob();
        await navigator.clipboard.write([new ClipboardItem({ 'image/png': blob })]);
        this.showToast('图片已复制', 'success');
      } catch (e) {
        // WebView2 里 ClipboardItem 的支持情况跟版本有关，失败不是异常，是常态路径的一种
        this.showToast('复制失败，请改用「保存 PNG」', 'error');
      }
    },

    async saveImage() {
      if (!this.qr) return;
      try {
        const url = await this.exportDataURL();
        const name = 'qrcode-' + Date.now() + '.png';
        const app = window.go && window.go.main && window.go.main.App;
        if (!app || !app.SaveQRCode) {
          this.showToast('当前环境不支持保存文件', 'error');
          return;
        }
        const path = await app.SaveQRCode(url, name);
        // 空路径 = 用户在对话框里点了取消，不是错误，静默返回
        if (path) this.showToast('图片已保存', 'success');
      } catch (e) {
        this.showToast('保存失败：' + (e.message || e), 'error');
      }
    }
  }
};
</script>

<style scoped>
.gen {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 定高 447：面板总高必须和识别页的 508 对齐（447 + 12 + 49 = 508），
   否则切 tab 时页面会跳一下。内部所有块用 flex 吸收，不写死任何一个的高度 */
.split {
  display: flex;
  gap: 16px;
  align-items: stretch;
  height: 447px;
}

.left {
  flex: 1 1 auto;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 350 = 白卡 300 + 白卡内边距 24 + 画布内边距 24 + 边框 2 */
.right {
  flex: 0 0 350px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* ---- 分区标题行 ---- */
.sec {
  flex: none;
  display: flex;
  align-items: center;
  gap: 7px;
  height: 24px;
}

.sec .ic {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: none;
  width: 22px;
  height: 22px;
  border-radius: var(--radius-xs);
  background: var(--surface-3);
  color: var(--text-2);
}

.sec h4 {
  margin: 0;
  font-size: 12.5px;
  font-weight: 650;
  color: var(--text);
}

.sec .sp {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
}

.cnt {
  font-family: var(--font-mono);
  font-size: 11.5px;
  color: var(--text-3);
  font-variant-numeric: tabular-nums;
}

.cnt.near {
  color: var(--tint3-fg);
}

.st {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11.5px;
  font-weight: 600;
}

.st.ok { color: var(--success); }
.st.bad { color: var(--danger); }
.st.idle { color: var(--text-3); font-weight: 400; }

/* ---- 输入区 ---- */
/* .inp 自带 height:36px，这里必须覆盖掉：在 .left 这个 column flex 里
   flex:1 沿主轴拉伸，输入框自动吃掉标题行和提示行之外的全部高度 */
.ta {
  box-sizing: border-box;
  height: auto;
  flex: 1;
  min-height: 0;
  padding: 10px 12px;
  resize: none;
  font-size: 13px;
  line-height: 1.75;
  font-family: inherit;
}

.hint {
  flex: none;
  font-size: 11.5px;
  line-height: 1.6;
  color: var(--text-3);
}

.hint b {
  font-weight: 650;
  color: var(--text-2);
}

.hint.bad {
  color: var(--danger);
  line-height: 1.5;
}

/* ---- 二维码图台 ---- */
.qstage {
  box-sizing: border-box;
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

/* 点阵底纹：把「这是一块画布」说清楚，白卡浮在上面才有层次 */
.qstage-body {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px;
  background-image: radial-gradient(color-mix(in srgb, var(--border-strong) 45%, transparent) 1px, transparent 1px);
  background-size: 16px 16px;
}

.qstage.ph {
  border-style: dashed;
}

.qstage.ph .qstage-body {
  flex-direction: column;
  gap: 13px;
  background-image: none;
}

/* 二维码本身是白底黑块，白卡必须还原它、不跟随主题（与识别页同一个例外） */
.qcard {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 300px;
  height: 300px;
  padding: 12px;
  border-radius: var(--radius-sm);
  background: #fff;
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.qr-art canvas {
  display: block;
}

.ph-frame {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 132px;
  height: 132px;
  border-radius: var(--radius-sm);
  background: #fff;
  border: 1px solid var(--border);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  pointer-events: none;
}

.ph-frame img {
  width: 106px;
  height: 106px;
  display: block;
  opacity: .14;
  pointer-events: none;
  -webkit-user-drag: none;
  user-select: none;
}

.ph-t {
  font-size: 12.5px;
  line-height: 1.6;
  color: var(--text-3);
  text-align: center;
  padding: 0 18px;
}

.qstage-cap {
  flex: none;
  display: flex;
  align-items: center;
  gap: 8px;
  height: 36px;
  padding: 0 12px 0 14px;
  border-top: 1px solid var(--border);
  font-size: 11.5px;
  color: var(--text-3);
  font-variant-numeric: tabular-nums;
}

.qstage-cap .end {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
}

.qstage-cap .dot {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--border-strong);
}

.mono-n {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-2);
}

.tag {
  display: inline-flex;
  align-items: center;
  height: 19px;
  padding: 0 7px;
  border-radius: 999px;
  background: var(--accent-soft);
  border: 1px solid var(--accent-border);
  color: var(--accent-strong);
  font-size: 11px;
  font-weight: 600;
}

.tag.warn {
  background: var(--tint3-bg);
  border-color: var(--tint3-bd);
  color: var(--tint3-fg);
}

/* ---- 底部工具条 ---- */
.bar {
  flex: none;
  display: flex;
  align-items: center;
  gap: 10px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}

.bar .lb {
  font-size: 11.5px;
  color: var(--text-3);
}

.bar .end {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
