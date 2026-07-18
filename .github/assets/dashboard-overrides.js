/**
 * Remap Allure 3 testing pyramid fills to Palette A (--layer-* in CSS).
 * Allure 3.13 hardcodes active bands to var(--color-intent-primary-bg).
 * Keep hex/vars in sync with allure/pyramid-layer-colors.mjs + tokens.css.
 */
(function () {
  const FUNNEL_TOP_TO_BOTTOM = ["manual", "e2e", "api", "integration", "component", "unit"];
  const PALETTE = {
    light: {
      unit: "#94a3b8",
      component: "#2563eb",
      integration: "#0891b2",
      api: "#7c3aed",
      e2e: "#d97706",
      manual: "#ea580c",
    },
    dark: {
      unit: "#64748b",
      component: "#3b82f6",
      integration: "#06b6d4",
      api: "#8b5cf6",
      e2e: "#f59e0b",
      manual: "#f97316",
    },
  };

  function findPyramidWidget(root) {
    return [...root.querySelectorAll('[class*="styles_widget"]')].find((el) =>
      /testing pyramid|пирамида тестирования/i.test(el.textContent || ""),
    );
  }

  function safeBBoxY(node) {
    try {
      return node.getBBox?.().y ?? 0;
    } catch {
      return 0;
    }
  }

  /**
   * Normalize a raw "Layer: <name>…" fragment to a known layer key.
   * Allure concatenates tspans ("manualNo tests") and layer names contain
   * digits ("e2e"), so a naive [a-z]+ capture is wrong — match by prefix.
   */
  function normalizeLayer(raw) {
    const text = (raw || "").trim().toLowerCase();
    let best = null;
    for (const layer of FUNNEL_TOP_TO_BOTTOM) {
      if (text.startsWith(layer) && (!best || layer.length > best.length)) {
        best = layer;
      }
    }
    return best;
  }

  /** Unique Layer: <name> labels with Y (Allure annotation tspans). */
  function layerLabelsFromWidget(widget) {
    const seen = new Set();
    const labels = [];
    widget.querySelectorAll("text, tspan").forEach((node) => {
      const match = (node.textContent || "").match(/Layer:\s*(.+)/i);
      if (!match) return;
      const layer = normalizeLayer(match[1]);
      if (!layer || seen.has(layer)) return;
      seen.add(layer);
      const textEl = node.closest("text") || node;
      labels.push({ layer, y: safeBBoxY(textEl) });
    });
    return labels;
  }

  // A layer with zero tests renders its callout as "Layer: <name>" + "No tests"
  // (Allure i18n), never a "Number of tests:" line — so match the empty phrase.
  const EMPTY_LAYER_RE = /no tests|нет тестов/i;

  /**
   * Layers whose callout says "No tests" — hidden from the funnel. The callout
   * <text> concatenates its tspans ("manualNo tests"), so read the whole node.
   */
  function emptyLayersFromWidget(widget) {
    const empty = new Set();
    widget.querySelectorAll("text").forEach((node) => {
      const text = node.textContent || "";
      const match = text.match(/Layer:\s*(.+)/i);
      if (!match) return;
      const layer = normalizeLayer(match[1]);
      if (layer && EMPTY_LAYER_RE.test(text)) empty.add(layer);
    });
    return empty;
  }

  function pyramidShapes(svg) {
    return [...svg.querySelectorAll("path, polygon")]
      .filter((shape) => {
        const d = shape.getAttribute("d") || "";
        const points = shape.getAttribute("points") || "";
        return d.length > 16 || points.length > 8;
      })
      .map((shape) => ({ shape, y: safeBBoxY(shape) }))
      .sort((left, right) => left.y - right.y);
  }

  /**
   * Pair shapes → layer by label Y proximity (not by data length / order tricks).
   * Same algorithm as allure/pyramid-layer-colors.mjs#pairShapesToLayers.
   */
  function pairShapesToLayers(shapeEntries, labels) {
    if (!shapeEntries.length) return [];

    // Full pyramid: shapes are sorted top→bottom, funnel order is deterministic.
    if (shapeEntries.length === FUNNEL_TOP_TO_BOTTOM.length) {
      return [...FUNNEL_TOP_TO_BOTTOM];
    }

    if (labels.length === shapeEntries.length) {
      const sorted = [...labels].sort((a, b) => a.y - b.y);
      return sorted.map((entry) => entry.layer);
    }

    if (labels.length > 0) {
      return shapeEntries.map((entry) => {
        let best = null;
        let bestDist = Infinity;
        for (const label of labels) {
          const dist = Math.abs(label.y - entry.y);
          if (dist < bestDist) {
            bestDist = dist;
            best = label.layer;
          }
        }
        return best;
      });
    }

    return shapeEntries.map(() => null);
  }

  function currentTheme() {
    return document.documentElement.getAttribute("data-theme") === "dark" ? "dark" : "light";
  }

  function colorForLayer(layer) {
    if (!layer || !PALETTE.light[layer]) return null;
    const cssVar = getComputedStyle(document.documentElement)
      .getPropertyValue(`--layer-${layer}`)
      .trim();
    return cssVar || PALETTE[currentTheme()][layer] || PALETTE.light[layer];
  }

  function setShapeFill(shape, layer) {
    const color = colorForLayer(layer);
    if (!color) return;
    shape.setAttribute("fill", color);
    shape.style.setProperty("fill", color, "important");
    shape.setAttribute("data-pyramid-layer", layer);
  }

  /** Collapse or restore a band. Uses style so Allure's geometry stays intact. */
  function setShapeHidden(shape, hidden) {
    if (hidden) shape.style.setProperty("display", "none", "important");
    else shape.style.removeProperty("display");
  }

  /** Hide the callout <text> of empty layers, restore the rest. */
  function setEmptyLabelVisibility(widget, emptyLayers) {
    widget.querySelectorAll("text").forEach((node) => {
      const match = (node.textContent || "").match(/Layer:\s*(.+)/i);
      if (!match) return;
      const layer = normalizeLayer(match[1]);
      if (!layer) return;
      node.style.display = emptyLayers.has(layer) ? "none" : "";
    });
  }

  function paintPyramid(root = document) {
    const widget = findPyramidWidget(root);
    if (!widget) return false;

    const svg = widget.querySelector("svg");
    if (!svg) return false;

    const shapeEntries = pyramidShapes(svg);
    if (!shapeEntries.length) return false;

    const labels = layerLabelsFromWidget(widget);
    const layers = pairShapesToLayers(shapeEntries, labels);
    const emptyLayers = emptyLayersFromWidget(widget);
    setEmptyLabelVisibility(widget, emptyLayers);

    shapeEntries.forEach((entry, index) => {
      const layer = layers[index];
      const hidden = layer != null && emptyLayers.has(layer);
      setShapeHidden(entry.shape, hidden);
      if (!layer || hidden) return;
      setShapeFill(entry.shape, layer);
    });

    return true;
  }

  // ---- Unified status indicators (widget header dots) ----
  // Every dashboard widget gets macOS-window dots in its header showing ONLY the
  // status-colour families actually present in that widget's chart, in a fixed
  // priority order. Colours are read "по факту" from rendered SVG fills/strokes
  // and matched to the nearest known palette anchor; noise (white, borders, text)
  // is rejected by a distance threshold.
  const INDICATOR_ORDER = ["red", "orange", "yellow", "purple", "gray", "green", "blue"];
  const INDICATOR_ANCHORS = {
    red: [[244, 63, 59], [255, 90, 80], [255, 100, 100], [153, 0, 24], [220, 38, 38], [244, 99, 134], [192, 57, 43], [255, 87, 68]],
    orange: [[255, 130, 0], [255, 140, 66], [255, 168, 51]],
    yellow: [[255, 216, 51], [255, 206, 87], [255, 224, 74], [201, 159, 0], [112, 93, 0], [255, 208, 80]],
    green: [[59, 201, 93], [148, 202, 102], [0, 111, 45], [0, 185, 151], [86, 214, 111], [137, 190, 62], [144, 187, 56], [163, 177, 37], [105, 167, 85]],
    blue: [[102, 186, 254], [84, 168, 237], [69, 155, 222], [97, 182, 251]],
    purple: [[161, 129, 255], [120, 85, 208], [216, 97, 190], [142, 68, 173]],
    gray: [[165, 183, 209], [170, 170, 170]],
  };
  const INDICATOR_THRESHOLD = 90;

  function parseRgbColor(value) {
    if (!value) return null;
    const match = String(value).match(/rgba?\(([^)]+)\)/i);
    if (!match) return null;
    const parts = match[1].split(",").map((piece) => parseFloat(piece));
    const alpha = parts.length > 3 ? parts[3] : 1;
    if (!(alpha > 0.3)) return null;
    if ([parts[0], parts[1], parts[2]].some((n) => Number.isNaN(n))) return null;
    return [parts[0], parts[1], parts[2]];
  }

  function classifyIndicator(rgb) {
    let best = null;
    let bestDist = Infinity;
    for (const family in INDICATOR_ANCHORS) {
      const anchors = INDICATOR_ANCHORS[family];
      for (let i = 0; i < anchors.length; i++) {
        const anchor = anchors[i];
        const dr = rgb[0] - anchor[0];
        const dg = rgb[1] - anchor[1];
        const db = rgb[2] - anchor[2];
        const dist = Math.sqrt(dr * dr + dg * dg + db * db);
        if (dist < bestDist) {
          bestDist = dist;
          best = family;
        }
      }
    }
    return bestDist <= INDICATOR_THRESHOLD ? best : null;
  }

  function widgetIndicatorFamilies(widget) {
    const present = Object.create(null);
    const shapes = widget.querySelectorAll(
      "svg path, svg polyline, svg polygon, svg rect, svg circle, svg ellipse, svg line",
    );
    shapes.forEach((node) => {
      const cs = getComputedStyle(node);
      [cs.fill, cs.stroke].forEach((raw) => {
        const rgb = parseRgbColor(raw);
        if (!rgb) return;
        const family = classifyIndicator(rgb);
        if (family) present[family] = true;
      });
    });
    return INDICATOR_ORDER.filter((family) => present[family]);
  }

  function widgetHeader(widget) {
    return (
      widget.querySelector('[class*="styles_header__"]') ||
      widget.firstElementChild ||
      widget
    );
  }

  // Idempotent: only mutates when a widget's family set changes, so the childList
  // observer that watches the whole document never enters a repaint loop.
  function paintWidgetIndicators(root) {
    const scope = root && root.querySelectorAll ? root : document;
    scope.querySelectorAll('[class*="styles_widget__"]').forEach((widget) => {
      const header = widgetHeader(widget);
      if (!header) return;
      const families = widgetIndicatorFamilies(widget);
      const key = families.join(",");
      let dots = header.querySelector(":scope > .zds-widget-dots");
      if (!families.length) {
        if (dots) dots.remove();
        return;
      }
      if (dots && dots.getAttribute("data-zds-fams") === key) return;
      if (!dots) {
        dots = document.createElement("span");
        dots.className = "zds-widget-dots";
        header.insertBefore(dots, header.firstChild);
      }
      dots.setAttribute("data-zds-fams", key);
      dots.textContent = "";
      families.forEach((family) => {
        const dot = document.createElement("span");
        dot.className = "zds-widget-dot zds-widget-dot--" + family;
        dots.appendChild(dot);
      });
    });
  }

  function paint() {
    paintPyramid();
    paintWidgetIndicators();
  }

  function schedulePaint() {
    paint();
    window.setTimeout(paint, 200);
    window.setTimeout(paint, 800);
    window.setTimeout(paint, 2000);
  }

  let paintQueued = false;
  function queuePaint() {
    if (paintQueued) return;
    paintQueued = true;
    requestAnimationFrame(() => {
      paintQueued = false;
      paint();
    });
  }

  const observer = new MutationObserver(queuePaint);
  observer.observe(document.documentElement, { childList: true, subtree: true });

  new MutationObserver(queuePaint).observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["data-theme"],
  });

  window.addEventListener("storage", queuePaint);

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", schedulePaint);
  } else {
    schedulePaint();
  }
})();
