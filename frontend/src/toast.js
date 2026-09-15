/* ============================================================================
   顶部提示浮层（mixin）
   ----------------------------------------------------------------------------
   配合 style.css 里的 .message-toast / toast-fade-* 使用：

     <transition name="toast-fade">
       <div v-if="toast.show" :class="['message-toast', toast.type]">{{ toast.text }}</div>
     </transition>

   组件里只要 mixins: [toastMixin]，再调用 this.showToast('已复制')。
   ============================================================================ */

const DURATION = 2600;

export const toastMixin = {
  data() {
    return {
      toast: { show: false, text: '', type: 'success' },
      toastTimer: null
    };
  },
  methods: {
    showToast(text, type = 'success') {
      // 连点时重置计时，否则旧提示会提前关掉新提示
      if (this.toastTimer) clearTimeout(this.toastTimer);
      this.toast = { show: true, text, type };
      this.toastTimer = setTimeout(() => {
        this.toast.show = false;
      }, DURATION);
    }
  },
  beforeUnmount() {
    if (this.toastTimer) clearTimeout(this.toastTimer);
  }
};

export default toastMixin;
