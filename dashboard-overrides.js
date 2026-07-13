(function () {
  const PYRAMID_COLORS = {
    dark: {
      unit: "#64748b",
      component: "#3b82f6",
      integration: "#06b6d4",
      api: "#8b5cf6",
      e2e: "#f59e0b",
      manual: "#f97316",
    },
    light: {
      unit: "#94a3b8",
      component: "#2563eb",
      integration: "#0891b2",
      api: "#7c3aed",
      e2e: "#d97706",
      manual: "#ea580c",
    },
  };

  const FALLBACK_LAYER_ORDER = ["manual", "e2e", "api", "integration", "component", "unit"];

  function readStoredTheme(rawValue) {
    if (!rawValue) return null;
    try {
      const parsed = JSON.parse(rawValue);
      return parsed === "dark" || parsed === "light" ? parsed : null;
    } catch {
      return rawValue.replace(/^"|"$/g, "") === "dark" ? "dark" : "light";
    }
  }

  function getTheme() {
    const attr = document.documentElement.getAttribute("data-theme");
    if (attr === "dark" || attr === "light") return attr;
    const stored = readStoredTheme(localStorage.getItem("theme"));
    if (stored) return stored;
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }

  function findPyramidWidget(root) {
    return [...root.querySelectorAll('[class*="styles_widget"]')].find((el) =>
      /testing pyramid|пирамида тестирования/i.test(el.textContent || "")
    );
  }

  function pyramidLayersFromWidget(widget) {
    const layers = [];
    widget.querySelectorAll("text, tspan").forEach((node) => {
      const match = (node.textContent || "").match(/^Layer:\s*([a-z]+)/i);
      if (!match) return;
      const layer = match[1].toLowerCase();
      if (!layers.includes(layer)) layers.push(layer);
    });
    return layers;
  }

  function pyramidShapes(svg) {
    return [...svg.querySelectorAll("path, polygon")]
      .filter((shape) => {
        const d = shape.getAttribute("d") || "";
        const points = shape.getAttribute("points") || "";
        return d.length > 16 || points.length > 8;
      })
      .sort((left, right) => {
        const leftY = left.getBBox?.().y ?? 0;
        const rightY = right.getBBox?.().y ?? 0;
        return leftY - rightY;
      });
  }

  function setShapeFill(shape, color) {
    shape.setAttribute("fill", color);
    shape.style.setProperty("fill", color, "important");
  }

  function paintPyramid(root = document) {
    const widget = findPyramidWidget(root);
    if (!widget) return false;

    const svg = widget.querySelector("svg");
    if (!svg) return false;

    const shapes = pyramidShapes(svg);
    if (!shapes.length) return false;

    const layers = pyramidLayersFromWidget(widget);
    const order =
      layers.length === shapes.length ? layers : FALLBACK_LAYER_ORDER.slice(-shapes.length);

    const palette = PYRAMID_COLORS[getTheme()] || PYRAMID_COLORS.light;
    shapes.forEach((shape, index) => {
      const layer = order[index];
      const color = palette[layer];
      if (!color) return;
      setShapeFill(shape, color);
    });

    return true;
  }

  function schedulePaint() {
    paintPyramid();
    window.setTimeout(paintPyramid, 200);
    window.setTimeout(paintPyramid, 800);
    window.setTimeout(paintPyramid, 2000);
  }

  const observer = new MutationObserver(() => paintPyramid());
  observer.observe(document.documentElement, { childList: true, subtree: true });

  new MutationObserver(() => paintPyramid()).observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["data-theme"],
  });

  window.addEventListener("storage", () => paintPyramid());

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", schedulePaint);
  } else {
    schedulePaint();
  }
})();
