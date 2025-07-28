import { cubicOut } from 'svelte/easing';
import type { TransitionConfig } from 'svelte/transition';

export function animateCoordinates(
  targetObject: { x: number; y: number },
  startCoordinates: { x: number; y: number },
  endCoordinates: { x: number; y: number },
  duration = 750,
  easingFunction = null
) {
  // Default easing function (easeInOutCubic)
  const easing =
    easingFunction ||
    function (t: number) {
      return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
    };

  const startTime = performance.now();

  return new Promise<void>((resolve) => {
    function animate(currentTime: number) {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / duration, 1);

      const easedProgress = easing(progress);

      targetObject.x = startCoordinates.x + (endCoordinates.x - startCoordinates.x) * easedProgress;
      targetObject.y = startCoordinates.y + (endCoordinates.y - startCoordinates.y) * easedProgress;

      if (progress < 1) {
        requestAnimationFrame(animate);
      } else {
        resolve();
      }
    }

    requestAnimationFrame(animate);
  });
}

export function scaleBlur(
  node: Element,
  { delay = 0, duration = 125, easing = cubicOut, startScale = 0.95, opacity = 0, blurAmount = 5 }
): TransitionConfig {
  return {
    delay,
    duration,
    easing,
    css: (t, u) => `
        transform: scale(${startScale + t * (1 - startScale)});
        opacity: ${opacity + t * (1 - opacity)};
        filter: blur(${blurAmount * u}px);
      `
  };
}

export function flyBlur(
  node: Element,
  {
    delay = 0,
    duration = 75,
    easing = cubicOut,
    y = 0,
    x = 0,
    opacity = 0,
    blurAmount = 8,
    preserveBackdrop = false
  }
): TransitionConfig {
  const style = getComputedStyle(node);
  const target_opacity = +style.opacity;
  const od = target_opacity * (1 - opacity);
  const [x_value, x_unit] = split_css_unit(x);
  const [y_value, y_unit] = split_css_unit(y);

  return {
    delay,
    duration,
    easing,
    css: (t, u) =>
      preserveBackdrop
        ? `position: relative; top: ${(1 - t) * y_value}${y_unit}; left: ${(1 - t) * x_value}${x_unit}; opacity: ${target_opacity - od * u};`
        : `transform: translate(${(1 - t) * x_value}${x_unit}, ${(1 - t) * y_value}${y_unit}); opacity: ${target_opacity - od * u}; filter: blur(${blurAmount * u}px);`
  };
}

function split_css_unit(value: number | string): [number, string] {
  const split = typeof value === 'string' && value.match(/^\s*(-?[\d.]+)([^\s]*)\s*$/);
  return split ? [parseFloat(split[1]), split[2] || 'px'] : [value as number, 'px'];
}
