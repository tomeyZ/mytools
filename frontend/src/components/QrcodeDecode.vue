<template>
  <!-- 单根：QrcodeTool 用 v-show 切两个面板，多根组件挂不上 v-show 的 display 样式 -->
  <div class="qr-wrap">

    <!-- ============ 上段：图台（340） ============ -->
    <!-- 空态：点击 / 拖拽共用这个区，Ctrl+V 走 document 级监听 -->
    <div
        v-if="!previewUrl"
        class="stage drop qr-drop"
        :class="{ hot: dragOver }"
        @click="pickFile"
        @dragover.prevent="dragOver = true"
        @dragleave.prevent="dragOver = false"
        @drop.prevent="onDrop"
    >
      <!-- 中间那张白卡就是待会儿二维码落的位置，里面放一枚淡化的真二维码当水印 ——
           比画一个抽象图标准得多，一眼就知道「这里是放二维码的」。
           ⚠️ 必须 draggable="false" + CSS pointer-events:none：否则从水印上按下拖动会触发
           浏览器原生「拖图片」，松手正好落回本拖拽区 → 水印被当成用户投的图识别一遍 -->
      <div class="dz-frame"><img :src="watermark" alt="" draggable="false"></div>
      <div class="dz-t">拖拽图片到这里，或 <span class="lnk">点击选择</span></div>
      <div class="dz-s">支持 PNG / JPG / WEBP / BMP · 也可以直接 Ctrl + V 粘贴</div>
    </div>

    <!-- 已选图：二维码 280 居中，图片自己的信息 / 动作收在图台底部说明条 -->
    <div v-else class="stage">
      <div class="stage-body">
        <!-- 识别中用扫描线提示：图大了，这条光带才真的看得见 -->
        <div class="big-qr qr-art" :class="{ scan: status === 'loading' }">
          <img :src="previewUrl" alt="已选择的二维码图片" draggable="false">
        </div>
      </div>
      <div class="stage-cap">
        <span class="cap-nm" :title="fileName">{{ fileName }}</span>
        <span>{{ fileMeta }}</span>
        <span class="end">
          <button class="btn ghost sm" :disabled="status === 'loading'" @click="pickFile">换一张图</button>
          <button class="btn ghost sm" :disabled="status === 'loading'" @click="clear">清空</button>
        </span>
      </div>
    </div>

    <!-- ============ 下段：识别结果（154） ============ -->
    <!-- 四态的框必须同高，否则选完图页面高度会跳一下 -->
    <div v-if="status === 'empty'" class="rbox-plain ph">
      选择一张二维码图片后，识别出的字符串会显示在这里
    </div>

    <div v-else-if="status === 'loading'" class="rbox-plain">
      <span class="spin"></span>正在识别二维码…
    </div>

    <!-- 成功 / 失败共用「主体 + 信息栏」结构：状态、字符数、耗时、类型胶囊、操作
         全部收进信息栏，页面上不再有浮在外面的标题行和按钮 -->
    <div v-else class="rbox">
      <div v-if="status === 'error'" class="rbox-main bad">
        <div class="b1">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="9"/><path d="M12 8v4.5"/><path d="M12 16h.01"/>
          </svg>
          未识别到二维码
        </div>
        <div class="b2">{{ errorMessage }}</div>
      </div>
      <div v-else class="rbox-main">
        <!-- 只读展示：可选中、可滚动；短内容居中、长内容自动顶头 —— margin:auto 一行搞定。
             链接不做任何特殊处理，跟别的字符串一样当普通长串显示 -->
        <div class="rbox-txt" :class="{ big: isShortResult }">{{ result }}</div>
      </div>

      <div class="rbox-foot">
        <template v-if="status === 'ok'">
          <span class="st ok">
            <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"/></svg>
            识别成功
          </span>
          <span class="dot"></span>
          <span>{{ result.length }} 字符</span>
          <span class="dot"></span>
        </template>
        <span>耗时 {{ elapsed }} ms</span>

        <span class="end">
          <span v-if="status === 'ok' && typeTag" class="tag">{{ typeTag }}</span>
          <button v-if="status === 'ok'" class="btn soft sm" @click="copyResult">复制</button>
          <button v-else class="btn ghost sm" @click="pickFile">换一张图</button>
        </span>
      </div>
    </div>

    <!-- 点击选择走原生 input：WebView2 原生支持，不用后端开文件对话框 -->
    <input ref="fileEl" type="file" accept="image/*" hidden @change="onFileChange">

    <!-- 复制提示：图标由 .message-toast::before 统一绘制（见 style.css），这里只放文案 -->
    <transition name="toast-fade">
      <div v-if="toast.show" :class="['message-toast', toast.type]">{{ toast.text }}</div>
    </transition>
  </div>
</template>

<script>
import jsQR from 'jsqr';
import { copyText } from '../clipboard.js';
import { toastMixin } from '../toast.js';
// 空态水印用的真实二维码图（不是矢量自绘 —— 29 模块的 SVG 缩到 106px 会发虚，像「图糊了」）
import watermark from '../assets/images/qr-watermark.png';

// 大图降采样阈值：手机原图 4000×3000 光 getImageData 就要几百毫秒，先缩到 1400 再喂解码器，
// 二维码定位图案仍然足够大，识别率不受影响
const MAX_SIDE = 1400;

// 短内容放大一号字并居中展示，长内容回落常规字号顶头 —— 只是个观感阈值，不影响功能
const SHORT_LEN = 64;

// readAsDataURL / loadImage 的 Promise 化，解码流程里要 await
const readAsDataURL = (file) => new Promise((resolve, reject) => {
  const fr = new FileReader();
  fr.onload = () => resolve(fr.result);
  fr.onerror = reject;
  fr.readAsDataURL(file);
});

const loadImage = (url) => new Promise((resolve, reject) => {
  const img = new Image();
  img.onload = () => resolve(img);
  img.onerror = reject;
  img.src = url;
});

export default {
  name: 'QrcodeDecode',
  mixins: [toastMixin],
  props: {
    // 面板是否在前台。两个面板用 v-show 切换、实例都活着，粘贴监听又挂在 document 上，
    // 不判断的话在「生成」tab 里粘一张图会被这里截走去识别
    active: { type: Boolean, default: true }
  },
  data() {
    return {
      status: 'empty',   // empty | loading | ok | error
      previewUrl: '',
      fileName: '',
      fileSize: 0,
      imgW: 0,
      imgH: 0,
      result: '',
      elapsed: 0,
      errorMessage: '',
      dragOver: false,
      watermark
    };
  },
  computed: {
    dimensionLabel() {
      return this.imgW ? this.imgW + ' × ' + this.imgH : '';
    },
    sizeLabel() {
      if (!this.fileSize) return '';
      const kb = this.fileSize / 1024;
      return kb >= 1024 ? (kb / 1024).toFixed(1) + ' MB' : kb.toFixed(1) + ' KB';
    },
    // 尺寸要等图片加载完才知道，所以两者用数组拼 —— 不能写死「A · B」，
    // 否则加载中会渲染成一个孤零零的分隔符（'· 3.1 KB'）
    fileMeta() {
      return [this.dimensionLabel, this.sizeLabel].filter(Boolean).join(' · ');
    },
    isShortResult() {
      return this.result.length > 0 && this.result.length <= SHORT_LEN;
    },
    // 内容类型只是个提示胶囊，识别错了也不影响结果本身，所以规则从宽
    typeTag() {
      const s = this.result;
      if (!s) return '';
      if (/A000000727/.test(s) || s.startsWith('0002')) return 'EMVCo 支付码';
      if (/^https?:\/\//i.test(s)) return '链接';
      if (/^WIFI:/i.test(s)) return 'WiFi 配置';
      if (/^BEGIN:VCARD/i.test(s)) return '名片';
      if (/^-----BEGIN [A-Z ]+-----/.test(s)) return '密钥 / 证书';
      return '';
    }
  },
  // KeepAlive 下组件不销毁：粘贴监听必须在 activated/deactivated 里挂摘，
  // 否则切到别的工具页后，用户往输入框里粘图也会被这边截走
  mounted() {
    this.bindPaste();
  },
  activated() {
    this.bindPaste();
  },
  deactivated() {
    this.unbindPaste();
  },
  beforeUnmount() {
    this.unbindPaste();
  },
  methods: {
    bindPaste() {
      document.addEventListener('paste', this.onPaste);
    },
    unbindPaste() {
      document.removeEventListener('paste', this.onPaste);
    },

    // ---- 三种入图入口，最后都汇到 handleFile ----
    pickFile() {
      if (this.status === 'loading') return;
      this.$refs.fileEl.click();
    },
    onFileChange(e) {
      const file = e.target.files && e.target.files[0];
      // 清掉 value：同一个文件选第二次也能触发 change
      e.target.value = '';
      if (file) this.handleFile(file);
    },
    onDrop(e) {
      this.dragOver = false;
      const files = e.dataTransfer && e.dataTransfer.files;
      const file = files && files[0];
      // 兜底：只放图片进来。站内拖拽（拖元素 / 拖文字）也可能带出 file，
      // MIME 为空时仍放行，避免误杀某些来源的位图
      if (file && (!file.type || file.type.startsWith('image/'))) this.handleFile(file);
    },
    onPaste(e) {
      if (!this.active) return;
      const items = e.clipboardData && e.clipboardData.items;
      if (!items) return;
      for (let i = 0; i < items.length; i++) {
        const it = items[i];
        if (it.kind === 'file' && it.type.startsWith('image/')) {
          const file = it.getAsFile();
          if (file) {
            e.preventDefault();
            this.handleFile(file);
          }
          return;
        }
      }
    },

    // ---- 主流程：读图 → 预览 → 解码，任何一步失败都走同一个失败态 ----
    async handleFile(file) {
      this.status = 'loading';
      this.previewUrl = '';
      this.result = '';
      this.errorMessage = '';
      this.elapsed = 0;
      this.imgW = 0;
      this.imgH = 0;
      this.fileSize = file.size || 0;
      this.fileName = file.name || '剪贴板图片.png';

      let url = '';
      try {
        url = await readAsDataURL(file);
        // 读完就先亮出预览：大图 loadImage 要几百毫秒，这期间图台不该还是「空拖拽区」，
        // 用户得能立刻确认「进去的是不是这张图」
        this.previewUrl = url;
        const img = await loadImage(url);
        this.imgW = img.naturalWidth;
        this.imgH = img.naturalHeight;
        const { code, ms } = this.decodeImage(img);
        this.elapsed = ms;
        if (code && code.data) {
          this.result = code.data;
          this.status = 'ok';
        } else {
          this.fail('图片里没找到可解析的二维码。换一张更清晰、完整的图试试 —— 二维码的四角要都在画面里，别被裁掉。');
        }
      } catch (err) {
        this.elapsed = this.elapsed || 0;
        this.fail('这张图读不出来，换一张试试');
      }
    },

    // 解码只做一件事：canvas 取像素 → jsQR。inversionAttempts 让深底浅码（暗色截图）也能出
    decodeImage(img) {
      const t0 = performance.now();
      const scale = Math.min(1, MAX_SIDE / Math.max(img.naturalWidth, img.naturalHeight));
      const w = Math.max(1, Math.round(img.naturalWidth * scale));
      const h = Math.max(1, Math.round(img.naturalHeight * scale));
      const canvas = document.createElement('canvas');
      canvas.width = w;
      canvas.height = h;
      const ctx = canvas.getContext('2d', { willReadFrequently: true });
      ctx.drawImage(img, 0, 0, w, h);
      const frame = ctx.getImageData(0, 0, w, h);
      const code = jsQR(frame.data, w, h, { inversionAttempts: 'attemptBoth' });
      return { code, ms: Math.max(1, Math.round(performance.now() - t0)) };
    },

    fail(message) {
      this.errorMessage = message;
      this.status = 'error';
    },

    // 回到空态：没有它，选完图只能「再选一张覆盖」，退不回去
    clear() {
      if (this.status === 'loading') return;
      this.status = 'empty';
      this.previewUrl = '';
      this.fileName = '';
      this.fileSize = 0;
      this.imgW = 0;
      this.imgH = 0;
      this.result = '';
      this.elapsed = 0;
      this.errorMessage = '';
      this.dragOver = false;
      if (this.$refs.fileEl) this.$refs.fileEl.value = '';
    },

    async copyResult() {
      if (!this.result) return;
      const ok = await copyText(this.result);
      this.showToast(ok ? '字符串已复制' : '复制失败', ok ? 'success' : 'error');
    }
  }
};
</script>

<style scoped>
.qr-wrap {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* ============ 上段：图台 ============ */
/* 空态（虚线拖拽区）与已选图（实心图台）同高 340 → 选完图不跳高 */
.stage {
  box-sizing: border-box;
  height: 340px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

.stage-body {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px;
}

/* 二维码 280×280 —— 线上原来只有 96 的缩略图，看不清模块，确认「是不是这张图」都费劲 */
.big-qr {
  position: relative;
  width: 280px;
  height: 280px;
  padding: 12px;
  box-sizing: border-box;
  border-radius: var(--radius-xs);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.big-qr img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: contain;
  /* 同理：预览图也不该被拖出去，否则等于把同一张图又投回来重识别一次 */
  pointer-events: none;
  -webkit-user-drag: none;
  user-select: none;
}

/* 识别中：一条扫过整张图的光带 */
.big-qr.scan::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--accent);
  box-shadow: 0 0 8px var(--accent);
  animation: qr-scan 1.25s ease-in-out infinite;
}

@keyframes qr-scan {
  0%, 100% { top: 6%; }
  50%      { top: 92%; }
}

/* 图台底部说明条：左边是图片自己的信息，右边是图片自己的动作。
   识别状态和耗时不在这里 —— 那是「结果」的事，归结果框的信息栏 */
.stage-cap {
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

.cap-nm {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-2);
  max-width: 360px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stage-cap .end {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* ---- 空态 ---- */
.stage.drop {
  align-items: center;
  justify-content: center;
  gap: 14px;
  border-style: dashed;
  border-color: var(--border-strong);
  cursor: pointer;
  transition: border-color var(--tr), background var(--tr);
}

.stage.drop.hot {
  border-color: var(--accent);
  background: var(--accent-soft);
}

/* 白卡 = 待会儿二维码落的位置。用纯白是因为二维码本身就是白底黑块，
   「像原图」压过「跟随主题」（与下方 .qr-art 同一个例外）。
   pointer-events:none 让整块水印对指针完全透明 —— 点击/拖拽全部穿透到外层拖拽区，
   从根上杜绝「把水印图拖出来又被当成投递的图片」 */
.dz-frame {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 132px;
  height: 132px;
  border-radius: var(--radius-xs);
  background: #fff;
  border: 1px solid var(--border);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  pointer-events: none;
}

.dz-frame img {
  width: 106px;
  height: 106px;
  display: block;
  opacity: .14;
  /* 水印只是装饰，不参与任何指针交互 —— 见模板里的 draggable 注释 */
  pointer-events: none;
  -webkit-user-drag: none;
  user-select: none;
}

.dz-t {
  font-size: 13px;
  color: var(--text-2);
  line-height: 1.5;
}

.lnk {
  color: var(--accent-strong);
  font-weight: 600;
}

.dz-s {
  font-size: 11.5px;
  color: var(--text-3);
  line-height: 1.6;
}

/* ============ 下段：结果框 ============ */
/* 标题行整条取消：状态、字符数、耗时、类型胶囊、复制按钮全收进框底信息栏，
   页面上只剩图台和结果框两个干净块 */
.rbox,
.rbox-plain {
  box-sizing: border-box;
  height: 154px;
  overflow: hidden;
  border: 1px solid var(--border);
  /* 与上方图台同圆角：两个同宽盒子上下叠在一起，圆角不一致会有错位感 */
  border-radius: var(--radius);
  background: var(--surface);
}

.rbox {
  display: flex;
  flex-direction: column;
}

.rbox-main {
  flex: 1;
  display: flex;
  overflow: auto;
  padding: 12px 16px;
}

/* 短内容居中、长内容超出时自动回落顶头可滚 —— 不用按长度写分支 */
.rbox-txt {
  margin: auto;
  width: 100%;
  font-family: var(--font-mono);
  font-size: 12.5px;
  line-height: 1.72;
  color: var(--text);
  word-break: break-all;
  white-space: pre-wrap;
}

.rbox-txt.big {
  font-size: 15px;
  letter-spacing: .3px;
}

/* 失败态也走「主体 + 信息栏」。底色只调淡一档的红 —— 大红块会把
   整个页面的注意力全吸走，失败是个提示，不该是主角 */
.rbox-main.bad {
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  text-align: center;
  padding: 0 24px;
  background: color-mix(in srgb, var(--danger-soft) 40%, var(--surface));
}

.rbox-main.bad .b1 {
  display: flex;
  align-items: center;
  gap: 7px;
  font-size: 13px;
  font-weight: 600;
  color: var(--danger);
}

.rbox-main.bad .b2 {
  font-size: 11.5px;
  line-height: 1.75;
  color: var(--text-2);
  max-width: 460px;
}

.rbox-foot {
  flex: none;
  display: flex;
  align-items: center;
  gap: 7px;
  height: 34px;
  padding: 0 12px 0 14px;
  border-top: 1px solid var(--border);
  background: var(--surface-2);
  font-size: 11.5px;
  color: var(--text-3);
  font-variant-numeric: tabular-nums;
}

.rbox-foot .dot {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--border-strong);
}

.rbox-foot .st {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-weight: 600;
}

.rbox-foot .st.ok {
  color: var(--success);
}

.rbox-foot .end {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 空态 / 识别中没有结论可写，就没有信息栏，用一个等高盒子 */
.rbox-plain {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 12.5px;
  color: var(--text-2);
}

/* 空态刻意再弱一档（浅虚线 + 透明底）：两个等重的空盒子并排会被读成
   「两个都要填的输入口」，视觉重点要让给上面的拖拽区 */
.rbox-plain.ph {
  border: 1px dashed var(--border);
  background: transparent;
  color: var(--text-3);
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

/* ---- 加载圈：类名不能叫 .sp（那名字留给右对齐占位），也别复用全局 spinner ---- */
.spin {
  display: block;
  width: 22px;
  height: 22px;
  flex: none;
  border-radius: 50%;
  border: 2.5px solid var(--accent-soft);
  border-top-color: var(--accent);
  animation: qr-spin .7s linear infinite;
}

@keyframes qr-spin {
  to { transform: rotate(360deg); }
}

/* 扫码件本身就是白底黑块，预览必须还原它、不跟随主题。
   这是全站禁硬编码色值的例外之一，理由：「像原图」压过「统一」 */
.qr-art {
  background: #fff;
}
</style>
