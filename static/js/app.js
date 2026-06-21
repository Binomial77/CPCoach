/* ============================================================
   CP COACH — shared front-end utilities
   No framework, no build step — just plain JS the Gin server
   can hand out as a static asset.
   ============================================================ */

const CPCoach = (() => {

  /* ---------- toast notifications ---------- */
  function ensureToastStack() {
    let stack = document.querySelector(".toast-stack");
    if (!stack) {
      stack = document.createElement("div");
      stack.className = "toast-stack";
      document.body.appendChild(stack);
    }
    return stack;
  }

  function toast(message, { ok = true, title } = {}) {
    const stack = ensureToastStack();
    const el = document.createElement("div");
    el.className = "toast" + (ok ? "" : " toast--err");
    el.innerHTML = `
      <span class="toast__title">${title || (ok ? "\u2713 Accepted" : "\u2715 Wrong Answer")}</span>
      <span class="toast__body"></span>
      <span class="toast__bar"></span>
    `;
    el.querySelector(".toast__body").textContent = message;
    stack.appendChild(el);

    const remove = () => {
      el.classList.add("is-leaving");
      setTimeout(() => el.remove(), 300);
    };
    const timer = setTimeout(remove, 3600);
    el.addEventListener("click", () => {
      clearTimeout(timer);
      remove();
    });
  }

  /* ---------- count-up number animation ---------- */
  function countUp(el, target, { duration = 900, prefix = "", suffix = "" } = {}) {
    if (!el) return;
    const start = 0;
    const startTime = performance.now();
    const reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    if (reduced || !Number.isFinite(target)) {
      el.textContent = prefix + (Number.isFinite(target) ? target : 0) + suffix;
      return;
    }

    function frame(now) {
      const progress = Math.min((now - startTime) / duration, 1);
      const eased = 1 - Math.pow(1 - progress, 3);
      const value = Math.round(start + (target - start) * eased);
      el.textContent = prefix + value + suffix;
      if (progress < 1) requestAnimationFrame(frame);
    }
    requestAnimationFrame(frame);
  }

  /* ---------- rating tier lookup (codeforces-style ladder) ---------- */
  const TIERS = [
    { max: 1200, key: "tier-newbie", name: "Newbie" },
    { max: 1400, key: "tier-pupil", name: "Pupil" },
    { max: 1600, key: "tier-specialist", name: "Specialist" },
    { max: 1900, key: "tier-expert", name: "Expert" },
    { max: 2100, key: "tier-candidate", name: "Candidate Master" },
    { max: 2300, key: "tier-master", name: "Master" },
    { max: 2400, key: "tier-international-master", name: "International Master" },
    { max: 2600, key: "tier-grandmaster", name: "Grandmaster" },
    { max: 3500, key: "tier-international-grandmaster", name: "International Grandmaster" },
    { max: Infinity, key: "tier-legendary", name: "Legendary Grandmaster" },
  ];

  function ratingTier(rating) {
    const r = Math.max(0, Number.isFinite(rating) ? rating : 0);
    return TIERS.find((t) => r < t.max) || TIERS[TIERS.length - 1];
  }

  /* ---------- typewriter console ---------- */
  // lines: array of strings. Supports simple inline tags: [ok]...[/ok] [warn]...[/warn] [err]...[/err] [dim]...[/dim]
  function typeLines(bodyEl, lines, { speed = 16, lineGap = 260, loop = false, onDone } = {}) {
    if (!bodyEl) return;
    const reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    function render(text) {
      return text
        .replace(/\[ok\](.*?)\[\/ok\]/g, '<span class="ok">$1</span>')
        .replace(/\[warn\](.*?)\[\/warn\]/g, '<span class="warn">$1</span>')
        .replace(/\[err\](.*?)\[\/err\]/g, '<span class="err">$1</span>')
        .replace(/\[dim\](.*?)\[\/dim\]/g, '<span class="dim">$1</span>');
    }

    if (reduced) {
      bodyEl.innerHTML = lines.map((l) => `<div class="console__line">${render(l)}</div>`).join("");
      if (onDone) onDone();
      return;
    }

    bodyEl.innerHTML = "";
    let cancelled = false;
    bodyEl.dataset.typing = "1";

    async function run() {
      while (!cancelled) {
        bodyEl.innerHTML = "";
        for (const raw of lines) {
          if (cancelled) return;
          const plain = raw.replace(/\[(ok|warn|err|dim)\]|\[\/(ok|warn|err|dim)\]/g, "");
          const lineEl = document.createElement("div");
          lineEl.className = "console__line";
          bodyEl.appendChild(lineEl);
          for (let i = 1; i <= plain.length; i++) {
            if (cancelled) return;
            // re-render the visible slice each tick so tag-coloring still applies once complete
            lineEl.textContent = plain.slice(0, i);
            await sleep(speed);
          }
          lineEl.innerHTML = render(raw);
          await sleep(lineGap);
        }
        const cursor = document.createElement("span");
        cursor.className = "console__cursor";
        bodyEl.lastChild?.appendChild(cursor);
        if (onDone) onDone();
        if (!loop) return;
        await sleep(2200);
      }
    }
    run();
    return () => { cancelled = true; };
  }

  function sleep(ms) { return new Promise((r) => setTimeout(r, ms)); }

  /* ---------- staggered reveal for grids of cards ---------- */
  function staggerReveal(selector, { delay = 60 } = {}) {
    document.querySelectorAll(selector).forEach((el, i) => {
      el.style.animationDelay = `${i * delay}ms`;
      el.classList.add("reveal");
    });
  }

  /* ---------- tiny fetch wrapper that understands this backend's JSON shape ---------- */
  async function postJSON(url, body) {
    const res = await fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    if (res.redirected) {
      return { redirected: true, url: res.url };
    }
    let data = {};
    try { data = await res.json(); } catch (_) { /* no body */ }
    return { redirected: false, ok: res.ok, data };
  }

  return { toast, countUp, ratingTier, typeLines, staggerReveal, postJSON };
})();