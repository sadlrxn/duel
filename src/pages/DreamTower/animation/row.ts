import gsap from "gsap";

export function setRowBackground({
  target,
  left = false,
  right = false,
  color = "#73218F",
  middleColor = "#1A293D",
  borderColor = "#e8caee",
  duration = 0,
}: {
  target: gsap.TweenTarget;
  left?: boolean;
  right?: boolean;
  duration?: number;
  color?: string;
  middleColor?: string;
  borderColor?: string;
}) {
  const background = `linear-gradient(90deg, ${
    color + (left ? "a0" : "00")
  } 0, ${middleColor}00 50%, ${color + (right ? "a0" : "00")} 100%)`;

  target &&
    gsap.to(target, {
      background,
      borderLeftColor: left ? borderColor : "transparent",
      borderRightColor: right ? borderColor : "transparent",
      duration,
    });
}

export function rowFlowAnims({
  targets,
  repeat = 1,
  duration = 0.2,
}: {
  targets: gsap.TweenTarget[];
  repeat?: number;
  duration?: number;
}) {
  const limit = targets.length;

  const loop = {
    value: 1,
  };

  return gsap.timeline().fromTo(
    loop,
    {
      value: 0.501 - 1,
    },
    {
      value: limit - 0.501,
      roundProps: "value",
      ease: "none",
      duration: duration * limit,
      repeat,
      onUpdate: () => {
        let i = 0;
        for (; i <= limit; i++) {
          let left = false,
            right = false;
          if (i === loop.value || i === (loop.value + 2) % limit) left = true;
          if (
            i === limit - loop.value - 1 ||
            i === limit - ((loop.value + 2) % limit) - 1
          )
            right = true;
          setRowBackground({ target: targets[i], left, right });
        }
      },
    }
  );
}
