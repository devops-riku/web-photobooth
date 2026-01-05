<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';

  export let value = '#ffffff';
  const dispatch = createEventDispatcher();

  let canvas: HTMLCanvasElement;
  let container: HTMLDivElement;
  let isDragging = false;
  let brightness = 100;

  // Internal state
  let hue = 0;
  let saturation = 0;
  let hexInput = value;
  let isFocused = false;

  $: {
    // Keep input in sync with value when not focused
    if (!isFocused && value.startsWith('#')) {
      hexInput = value.toUpperCase();
    }
  }

  function handleHexInput(e: Event) {
    const target = e.target as HTMLInputElement;
    let val = target.value;
    if (val && !val.startsWith('#')) val = '#' + val;
    
    // Validate hex
    if (/^#[0-9A-F]{6}$/i.test(val)) {
      value = val;
      const { h, s, l } = hexToHsl(val);
      hue = h;
      saturation = s;
      brightness = l;
      dispatch('change', value);
    }
  }

  $: {
    // Sync external value to internal HSL if it changes from outside
    if (!isDragging) {
      const { h, s, l } = hexToHsl(value);
      hue = h;
      saturation = s;
      brightness = l;
    }
  }

  function hslToHex(h: number, s: number, l: number) {
    l /= 100;
    const a = s * Math.min(l, 1 - l) / 100;
    const f = (n: number) => {
      const k = (n + h / 30) % 12;
      const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1);
      return Math.round(255 * color).toString(16).padStart(2, '0');
    };
    return `#${f(0)}${f(8)}${f(4)}`;
  }

  function hexToHsl(hex: string) {
    let r = parseInt(hex.slice(1, 3), 16) / 255;
    let g = parseInt(hex.slice(3, 5), 16) / 255;
    let b = parseInt(hex.slice(5, 7), 16) / 255;

    let max = Math.max(r, g, b), min = Math.min(r, g, b);
    let h, s, l = (max + min) / 2;

    if (max === min) {
      h = s = 0;
    } else {
      let d = max - min;
      s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
      switch (max) {
        case r: h = (g - b) / d + (g < b ? 6 : 0); break;
        case g: h = (b - r) / d + 2; break;
        case b: h = (r - g) / d + 4; break;
      }
      h /= 6;
    }
    return { h: h * 360, s: s * 100, l: l * 100 };
  }

  function handleMove(e: MouseEvent | TouchEvent) {
    if (!isDragging && e.type !== 'click') return;
    
    const rect = container.getBoundingClientRect();
    const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX;
    const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY;
    
    const x = clientX - rect.left - rect.width / 2;
    const y = clientY - rect.top - rect.height / 2;
    
    const angle = Math.atan2(y, x) * (180 / Math.PI);
    hue = (angle + 360) % 360;
    
    const dist = Math.sqrt(x * x + y * y);
    const maxDist = rect.width / 2;
    saturation = Math.min(100, (dist / maxDist) * 100);
    
    updateValue();
  }

  function updateValue() {
    value = hslToHex(hue, saturation, brightness);
    dispatch('change', value);
  }

  function startDragging(e: MouseEvent | TouchEvent) {
    isDragging = true;
    handleMove(e);
  }

  function stopDragging() {
    isDragging = false;
  }

  $: thumbX = Math.cos(hue * Math.PI / 180) * (saturation * (container?.offsetWidth ?? 0) / 200);
  $: thumbY = Math.sin(hue * Math.PI / 180) * (saturation * (container?.offsetHeight ?? 0) / 200);

</script>

<div class="flex flex-col gap-4 items-center w-full max-w-[240px]">
  <div 
    bind:this={container}
    class="relative w-48 h-48 rounded-full shadow-[inset_0_2px_10px_rgba(0,0,0,0.1)] border-4 border-white cursor-crosshair touch-none"
    style="background: radial-gradient(circle, #fff, transparent), conic-gradient(red, yellow, lime, aqua, blue, magenta, red);"
    on:mousedown={startDragging}
    on:touchstart|preventDefault={startDragging}
    on:mousemove={handleMove}
    on:touchmove|preventDefault={handleMove}
  >
    <!-- Selection Thumb -->
    <div 
      class="absolute w-5 h-5 border-2 border-white rounded-full shadow-md pointer-events-none -translate-x-1/2 -translate-y-1/2 transition-transform duration-75"
      style="left: calc(50% + {thumbX}px); top: calc(50% + {thumbY}px); background-color: {value};"
    ></div>
  </div>

  <!-- Brightness Slider -->
  <div class="w-full flex flex-col gap-2">
    <div class="flex justify-between items-center">
      <span class="text-[9px] font-bold uppercase text-slate-400 tracking-widest">Brightness</span>
      <span class="text-[9px] font-bold text-purple-400">{Math.round(brightness)}%</span>
    </div>
    <div class="relative h-6 flex items-center">
      <input 
        type="range" 
        min="0" 
        max="100" 
        bind:value={brightness}
        on:input={updateValue}
        class="w-full h-1.5 bg-slate-100 rounded-lg appearance-none cursor-pointer accent-purple-500"
        style="background: linear-gradient(to right, #000, {hslToHex(hue, saturation, 50)}, #fff);"
      />
    </div>
  </div>

  <!-- Hex Input -->
  <div class="w-full flex flex-col gap-2">
    <div class="flex justify-between items-center">
      <span class="text-[9px] font-bold uppercase text-slate-400 tracking-widest">Hex Code</span>
    </div>
    <div class="flex items-center gap-2">
      <div class="w-6 h-6 rounded-md border border-slate-200 shrink-0" style="background-color: {value}"></div>
      <input 
        type="text" 
        bind:value={hexInput}
        on:input={handleHexInput}
        on:focus={() => isFocused = true}
        on:blur={() => isFocused = false}
        placeholder="#FFFFFF"
        class="flex-1 bg-slate-50 border-none rounded-lg px-3 py-2 text-[11px] font-mono font-medium text-slate-600 focus:ring-1 focus:ring-purple-200 outline-none uppercase"
        maxlength="7"
      />
    </div>
  </div>
</div>

<svelte:window on:mouseup={stopDragging} on:touchend={stopDragging} />

<style>
  input[type="range"]::-webkit-slider-thumb {
    appearance: none;
    width: 16px;
    height: 16px;
    background: white;
    border: 2px solid #a855f7;
    border-radius: 50%;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  }
</style>
